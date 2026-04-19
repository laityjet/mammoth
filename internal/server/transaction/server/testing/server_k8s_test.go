//go:build k8s

package testing

import (
	"bytes"
	"fmt"
	"github.com/laityjet/mammoth/v0/internal/uuid"
	"strings"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/minikubetestenv"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/testutil"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"
)

func TestCreatePipelineTransaction(t *testing.T) {
	c, _ := minikubetestenv.AcquireCluster(t)
	repo := uuid.UniqueString("in")
	pipeline := uuid.UniqueString("pipeline")
	_, err := c.ExecuteInTransaction(func(txnClient *client.APIClient) error {
		require.NoError(t, txnClient.CreateRepo(pfs.DefaultProjectName, repo))
		require.NoError(t, txnClient.CreatePipeline(pfs.DefaultProjectName,
			pipeline,
			"",
			[]string{"bash"},
			[]string{fmt.Sprintf("cp /pfs/%s/* /pfs/out", repo)},
			&pps.ParallelismSpec{Constant: 1},
			client.NewPFSInput(pfs.DefaultProjectName, repo, "/"),
			"master",
			false,
		))
		return nil
	})
	require.NoError(t, err)

	commit := client.NewCommit(pfs.DefaultProjectName, repo, "master", "")
	require.NoError(t, c.PutFile(commit, "foo", strings.NewReader("bar")))

	commitInfo, err := c.WaitCommit(pfs.DefaultProjectName, pipeline, "master", "")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, c.GetFile(commitInfo.Commit, "foo", &buf))
	require.Equal(t, "bar", buf.String())
}

func TestCreateProjectlessPipelineTransaction(t *testing.T) {
	c, _ := minikubetestenv.AcquireCluster(t)
	repo := uuid.UniqueString("in")
	pipeline := uuid.UniqueString("pipeline")
	_, err := c.ExecuteInTransaction(func(txnClient *client.APIClient) error {
		require.NoError(t, txnClient.CreateRepo(pfs.DefaultProjectName, repo))
		_, err := txnClient.PpsAPIClient.CreatePipeline(txnClient.Ctx(),
			&pps.CreatePipelineRequest{
				Pipeline: &pps.Pipeline{Name: pipeline},
				Transform: &pps.Transform{
					Image: testutil.DefaultTransformImage,
					Cmd:   []string{"bash"},
					Stdin: []string{fmt.Sprintf("cp /pfs/%s/* /pfs/out", repo)},
				},
				ParallelismSpec: &pps.ParallelismSpec{Constant: 1},
				Input:           client.NewPFSInput(pfs.DefaultProjectName, repo, "/"),
				OutputBranch:    "master",
			})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	commit := client.NewCommit(pfs.DefaultProjectName, repo, "master", "")
	require.NoError(t, c.PutFile(commit, "foo", strings.NewReader("bar")))

	commitInfo, err := c.WaitCommit(pfs.DefaultProjectName, pipeline, "master", "")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, c.GetFile(commitInfo.Commit, "foo", &buf))
	require.Equal(t, "bar", buf.String())
}
