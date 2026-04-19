package testing

import (
	"context"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/pfs"

	"github.com/laityjet/mammoth/v0/internal/pachd"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
)

// TestCheckStorage checks that the CheckStorage rpc is wired up correctly.
// An more extensive test lives in the `chunk` package.
func TestCheckStorage(t *testing.T) {
	ctx := pctx.TestContext(t)
	t.Parallel()
	client := newClient(ctx, t)
	res, err := client.CheckStorage(ctx, &pfs.CheckStorageRequest{
		ReadChunkData: false,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func newClient(ctx context.Context, t testing.TB) pfs.APIClient {
	pachClient := pachd.NewTestPachd(t)
	return pachClient.PfsAPIClient
}
