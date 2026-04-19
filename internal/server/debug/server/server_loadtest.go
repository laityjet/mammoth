package server

import (
	"context"
	_ "embed"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"sigs.k8s.io/yaml"

	"github.com/laityjet/mammoth/v0/internal/debug"
	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pfsload"
	"github.com/laityjet/mammoth/v0/internal/task"
	"github.com/laityjet/mammoth/v0/internal/uuid"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	srvpfs "github.com/laityjet/mammoth/v0/internal/server/pfs"
)

// RunLoadTest implements the pfs.RunLoadTest RPC
func (a *debugServer) RunPFSLoadTest(ctx context.Context, req *debug.RunPFSLoadTestRequest) (_ *debug.RunPFSLoadTestResponse, retErr error) {
	pachClient := a.env.GetPachClient(ctx)
	taskService := a.env.TaskService
	var project string
	repo := "load_test"
	if req.Branch != nil {
		project = req.Branch.Repo.Project.GetName()
		repo = req.Branch.Repo.Name
	}
	if err := pachClient.CreateRepo(project, repo); err != nil && !srvpfs.IsRepoExistsErr(err) {
		return nil, err
	}
	branch := uuid.New()
	if req.Branch != nil {
		branch = req.Branch.Name
	}
	if err := pachClient.CreateBranch(project, repo, branch, "", "", nil); err != nil {
		return nil, err
	}
	seed := time.Now().UTC().UnixNano()
	if req.Seed > 0 {
		seed = req.Seed
	}
	resp := &debug.RunPFSLoadTestResponse{
		Spec:   req.Spec,
		Branch: client.NewBranch(req.Branch.GetRepo().GetProject().GetName(), repo, branch),
		Seed:   seed,
	}
	start := time.Now()
	var err error
	resp.StateId, err = a.runLoadTest(pachClient, taskService, resp.Branch, req.Spec, seed, req.StateId)
	if err != nil {
		resp.Error = err.Error()
	}
	resp.Duration = durationpb.New(time.Since(start))
	return resp, nil
}

func (a *debugServer) runLoadTest(pachClient *client.APIClient, taskService task.Service, branch *pfs.Branch, specStr string, seed int64, stateID string) (string, error) {
	jsonBytes, err := yaml.YAMLToJSON([]byte(specStr))
	if err != nil {
		return "", errors.EnsureStack(err)
	}
	spec := &pfsload.CommitSpec{}
	if err := protojson.Unmarshal(jsonBytes, spec); err != nil {
		return "", errors.Wrap(err, "unmarshal CommitSpec")
	}
	return pfsload.Commit(pachClient.Ctx(), pachClient.PfsAPIClient, taskService, branch, spec, seed, stateID)
}

func (a *debugServer) RunPFSLoadTestDefault(ctx context.Context, _ *emptypb.Empty) (resp *debug.RunPFSLoadTestResponse, retErr error) {
	for _, spec := range defaultLoadSpecs {
		var err error
		resp, err = a.RunPFSLoadTest(ctx, &debug.RunPFSLoadTestRequest{
			Spec: spec,
		})
		if err != nil {
			return nil, err
		}
		if resp.Error != "" {
			return resp, nil
		}
	}
	return resp, nil
}

var (
	//go:embed load-test-0.yaml
	loadSpec1 string
	//go:embed load-test-1.yaml
	loadSpec2 string
	//go:embed load-test-2.yaml
	loadSpec3 string

	defaultLoadSpecs = []string{loadSpec1, loadSpec2, loadSpec3}
)
