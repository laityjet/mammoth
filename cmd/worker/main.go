package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	debugclient "github.com/laityjet/mammoth/v0/internal/debug"
	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/cmdutil"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/grpcutil"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/middleware/logging"
	"github.com/laityjet/mammoth/v0/internal/pachconfig"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/ppsutil"
	"github.com/laityjet/mammoth/v0/internal/proc"
	"github.com/laityjet/mammoth/v0/internal/profileutil"
	"github.com/laityjet/mammoth/v0/internal/restart"
	"github.com/laityjet/mammoth/v0/internal/serviceenv"
	"github.com/laityjet/mammoth/v0/internal/tracing"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"
	debugserver "github.com/laityjet/mammoth/v0/internal/server/debug/server"
	"github.com/laityjet/mammoth/v0/internal/server/worker"
	workerserver "github.com/laityjet/mammoth/v0/internal/server/worker/server"
	"github.com/laityjet/mammoth/v0/internal/version"
	"github.com/laityjet/mammoth/v0/internal/version/versionpb"
	workerapi "github.com/laityjet/mammoth/v0/internal/worker"
)

func main() {
	log.InitWorkerLogger()
	ctx := pctx.Child(pctx.Background(""), "", pctx.WithFields(pps.WorkerIDField(os.Getenv(client.PPSPodNameEnv))))
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println(runtime.GOARCH)
		fmt.Println(runtime.GOOS)
		fmt.Println(version.PrettyPrintVersion(version.Version))
		os.Exit(0)
	}
	go log.WatchDroppedLogs(ctx, time.Minute)
	go proc.MonitorSelf(ctx)
	log.Debug(ctx, "version info", log.Proto("versionInfo", version.Version))

	// append pachyderm bins to path to allow use of pachctl
	os.Setenv("PATH", os.Getenv("PATH")+":/pach-bin")
	cmdutil.Main(ctx, do, &pachconfig.WorkerFullConfiguration{})
}

func do(ctx context.Context, config *pachconfig.WorkerFullConfiguration) error {
	// must run InstallJaegerTracer before InitWithKube/pach client initialization
	tracing.InstallJaegerTracerFromEnv()
	env := serviceenv.InitWithKube(ctx, pachconfig.NewConfiguration(config))

	// Enable cloud profilers if the configuration allows.
	profileutil.StartCloudProfiler(ctx, "pachyderm-worker", env.Config())

	// Enable restart watcher.
	r, err := restart.New(ctx, env.GetDBClient(), env.GetPostgresListener())
	if err != nil {
		return errors.Wrap(err, "restart.New")
	}
	go func() {
		if err := r.RestartWhenRequired(ctx); err != nil {
			log.Error(ctx, "restart notifier failed", zap.Error(err))
		}
	}()

	// Construct a client that connects to the sidecar.
	pachClient := env.GetPachClient(ctx)
	p := &pps.Pipeline{
		Project: &pfs.Project{Name: env.Config().PPSProjectName},
		Name:    env.Config().PPSPipelineName,
	}
	pipelineInfo, err := ppsutil.GetWorkerPipelineInfo(
		pachClient,
		env.GetDBClient(),
		env.GetPostgresListener(),
		p,
		env.Config().PPSSpecCommitID,
	) // get pipeline creds for pachClient
	if err != nil {
		return errors.Wrapf(err, "worker: get pipelineInfo for %q", p)
	}
	ctx = pachClient.AddMetadata(ctx)

	// Construct worker API server.
	workerInstance, err := worker.NewWorker(pctx.Child(ctx, ""), env, pachClient, pipelineInfo, "/")
	if err != nil {
		return err
	}

	// grpc logger
	logs := logging.NewLoggingInterceptor(ctx)
	logs.Level = log.DebugLevel

	// Start worker api server
	server, err := grpcutil.NewServer(ctx, false,
		grpc.ChainUnaryInterceptor(logs.UnarySetup, logs.UnaryAnnounce),
		grpc.ChainStreamInterceptor(logs.StreamSetup, logs.StreamAnnounce),
	)
	if err != nil {
		return err
	}

	workerapi.RegisterWorkerServer(server.Server, workerInstance.APIServer)
	versionpb.RegisterAPIServer(server.Server, version.NewAPIServer(version.Version, version.APIServerOptions{}))
	debugServer := debugserver.NewDebugServer(debugserver.Env{
		DB:                   env.GetDBClient(),
		SidecarClient:        pachClient,
		GetLokiClient:        env.GetLokiClient,
		Name:                 env.Config().PodName,
		GetPachClient:        pachClient.WithCtx,
		GetKubeClient:        env.GetKubeClient,
		GetDynamicKubeClient: env.GetDynamicKubeClient,
		Config:               *env.Config(),
		TaskService:          env.GetTaskService("debug"),
	})
	debugclient.RegisterDebugServer(server.Server, debugServer)

	// Put our IP address into etcd, so pachd can discover us
	workerRcName := ppsutil.PipelineRcName(pipelineInfo)
	key := path.Join(env.Config().PPSEtcdPrefix, workerserver.WorkerEtcdPrefix, workerRcName, env.Config().PPSWorkerIP)

	// Prepare to write "key" into etcd by creating lease -- if worker dies, our
	// IP will be removed from etcd
	leaseID, err := getETCDLease(ctx, env.GetEtcdClient(), 10*time.Second)
	if err != nil {
		return errors.Wrapf(err, "worker: get etcd lease")
	}

	// keepalive forever
	keepAliveChan, err := env.GetEtcdClient().KeepAlive(ctx, leaseID)
	if err != nil {
		return errors.Wrapf(err, "worker: etcd KeepAlive")
	}
	go func() {
		for {
			_, more := <-keepAliveChan
			if !more {
				log.Error(ctx, "failed to renew worker IP address etcd lease")
				return
			}
		}
	}()

	if err := writeKey(ctx, env.GetEtcdClient(), key, leaseID, 10*time.Second); err != nil {
		return errors.Wrapf(err, "worker: etcd key %s", key)
	}

	// If server ever exits, return error
	if _, err := server.ListenTCP("", env.Config().PPSWorkerPort); err != nil {
		return err
	}
	return server.Wait()
}

func getETCDLease(ctx context.Context, client *etcd.Client, duration time.Duration) (etcd.LeaseID, error) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	sec := int64(duration / time.Second)
	if sec == 0 { // do not aallow durations < 1 second to round down to 0 seconds
		sec = 1
	}
	resp, err := client.Grant(ctx, sec)
	if err != nil {
		return 0, errors.Wrapf(err, "getETCDLease: etcd grant")
	}
	return resp.ID, nil
}

func writeKey(ctx context.Context, client *etcd.Client, key string, id etcd.LeaseID, duration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := client.Put(ctx, key, "", etcd.WithLease(id)); err != nil {
		return errors.Wrapf(err, "writeKey: etcd put")
	}
	return nil
}
