package server

import (
	"context"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/laityjet/mammoth/v0/internal/admin"
	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/debug"
	"github.com/laityjet/mammoth/v0/internal/enterprise"
	"github.com/laityjet/mammoth/v0/internal/identity"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/testpachd/realenv"
	tu "github.com/laityjet/mammoth/v0/internal/testutil"
	"github.com/laityjet/mammoth/v0/internal/license"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"
	"github.com/laityjet/mammoth/v0/internal/proxy"
	"github.com/laityjet/mammoth/v0/internal/transaction"
	"github.com/laityjet/mammoth/v0/internal/version/versionpb"
	"github.com/laityjet/mammoth/v0/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var protos = []protoreflect.FileDescriptor{
	admin.File_admin_admin_proto,
	auth.File_auth_auth_proto,
	debug.File_debug_debug_proto,
	enterprise.File_enterprise_enterprise_proto,
	identity.File_identity_identity_proto,
	license.File_license_license_proto,
	pfs.File_pfs_pfs_proto,
	pps.File_pps_pps_proto,
	proxy.File_proxy_proxy_proto,
	transaction.File_transaction_transaction_proto,
	versionpb.File_version_versionpb_version_proto,
	worker.File_worker_worker_proto,
}

func rangeRPCs(f func(fd protoreflect.FileDescriptor, sd protoreflect.ServiceDescriptor, md protoreflect.MethodDescriptor)) {
	for _, fd := range protos {
		svcs := fd.Services()
		for si := 0; si < svcs.Len(); si++ {
			sd := svcs.Get(si)
			methods := sd.Methods()
			for mi := 0; mi < methods.Len(); mi++ {
				md := methods.Get(mi)
				f(fd, sd, md)
			}
		}
	}
}

func testRPC(ctx context.Context, t *testing.T, sd protoreflect.ServiceDescriptor, md protoreflect.MethodDescriptor, cc *grpc.ClientConn, req proto.Message) {
	fullName := "/" + path.Join(string(sd.FullName()), string(md.Name()))
	reply := &emptypb.Empty{}
	if err := cc.Invoke(ctx, fullName, req, reply); err != nil {
		if s, ok := status.FromError(err); ok {
			switch {
			case s.Code() == codes.Aborted && strings.Contains(s.Message(), "panic: "):
				t.Fatal(err)
			case strings.Contains(s.Message(), "this is a bug"):
				t.Fatal(err)
			case strings.Contains(s.Message(), "pachd mock"):
				t.Skip("skipping method that has no mock")
			case strings.Contains(s.Message(), "not activated"):
				t.Fatal("auth not activated?")
			case strings.Contains(s.Message(), "already activated"):
				// This is fine and is a useful test.
				return
			case s.Code() == codes.Unimplemented:
				t.Skip("skipping unimplemented method/service")
			case s.Code() == codes.Unauthenticated:
				t.Fatal("unauthenticated?")
			}
		} else {
			t.Log(err)
		}
	}
}

func protoBytes(msg proto.Message) []byte {
	b, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func FuzzRPCs(f *testing.F) {
	f.Add(protoBytes(&emptypb.Empty{}))
	f.Add(protoBytes(&pfs.InspectRepoRequest{
		Repo: nil,
	}))
	f.Add(protoBytes(&pfs.FindCommitsRequest{
		Start: &pfs.Commit{},
	}))
	f.Add(protoBytes(&pfs.DiffFileRequest{
		OldFile: &pfs.File{
			Commit: &pfs.Commit{
				Repo: &pfs.Repo{
					Name: "test",
					Project: &pfs.Project{
						Name: "default",
					},
				},
			},
		},
	}))
	ctx := pctx.Background("")
	env := realenv.NewRealEnvWithIdentity(ctx, f, dockertestenv.NewTestDBConfig(f).PachConfigOption)
	peerPort := strconv.Itoa(int(env.ServiceEnv.Config().PeerPort))
	tu.ActivateAuthClient(f, env.PachClient, peerPort)

	f.Fuzz(func(t *testing.T, a []byte) {
		rangeRPCs(func(fd protoreflect.FileDescriptor, sd protoreflect.ServiceDescriptor, md protoreflect.MethodDescriptor) {
			name := string(sd.FullName()) + "." + string(md.Name())
			t.Run(name, func(t *testing.T) {
				switch {
				case strings.Contains(name, "RunLoadTest"):
					t.Skip("skipping load tests")
				case strings.Contains(name, "Deactivate"):
					t.Skip("skipping auth deactivation")
				case strings.Contains(name, "RotateRootToken"):
					t.Skip("skipping RotateRootToken")
				}
				ctx, c := context.WithTimeout(pctx.Child(ctx, name), 5*time.Second)
				defer c()
				client := env.PachClient.WithCtx(ctx)

				msg := dynamicpb.NewMessage(md.Input())
				if err := proto.Unmarshal(a, msg); err != nil {
					// Random data was not a proto, don't go any further.
					return
				}
				testRPC(client.Ctx(), t, sd, md, client.ClientConn(), msg)
			})
		})
	})
}
