package server

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/admin"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/version"

	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestInspectCluster(t *testing.T) {
	var a apiServer
	_, err := a.InspectCluster(pctx.TestContext(t), &admin.InspectClusterRequest{
		ClientVersion:  version.Version,
		CurrentProject: &pfs.Project{Name: "#<does-not-exist>"},
	})
	require.NoError(t, err, "InspectCluster must not err")
}
