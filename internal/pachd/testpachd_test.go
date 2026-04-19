package pachd_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/pachd"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/pfs"
)

func TestNewTestPachd(t *testing.T) {
	ctx, done := context.WithTimeout(pctx.TestContext(t), 5*time.Second)
	c := pachd.NewTestPachd(t).WithCtx(ctx)

	defer done()
	repo := &pfs.Repo{
		Name: "test",
		Project: &pfs.Project{
			Name: pfs.DefaultProjectName,
		},
		Type: pfs.UserRepoType,
	}
	commit := &pfs.Commit{
		Repo: repo,
		Branch: &pfs.Branch{
			Repo: repo,
			Name: "master",
		},
	}
	filename := "test.txt"
	want := "hello, world!\n"
	_, err := c.PfsAPIClient.CreateRepo(c.Ctx(), &pfs.CreateRepoRequest{Repo: repo})
	require.NoError(t, err, "should create repo")
	if err := c.PutFile(commit, filename, strings.NewReader(want)); err != nil {
		t.Fatalf("put file: %v", err)
	}
	var buf bytes.Buffer
	if err := c.GetFile(commit, filename, &buf); err != nil {
		t.Fatalf("get file: %v", err)
	}
	require.Equal(t, want, buf.String(), "file content should be equal")
}

func TestNewTestPachd_WithAuth(t *testing.T) {
	c := pachd.NewTestPachd(t, pachd.ActivateAuthOption(""))
	robot, err := c.AuthAPIClient.GetRobotToken(c.Ctx(), &auth.GetRobotTokenRequest{
		Robot: "robotguy42",
	})
	if err != nil {
		t.Fatalf("create robot token: %v", err)
	}
	c.SetAuthToken(robot.Token)
	whoami, err := c.AuthAPIClient.WhoAmI(c.Ctx(), &auth.WhoAmIRequest{})
	if err != nil {
		t.Fatalf("who am i: %v", err)
	}
	require.Equal(t, "robot:robotguy42", whoami.Username, "whoami should tell me my username")
}
