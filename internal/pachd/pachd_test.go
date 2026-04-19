package pachd

import (
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/pjs"

	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestFull(t *testing.T) {
	ctx := pctx.TestContext(t)
	pc := NewTestPachd(t)
	res, err := pc.VersionAPIClient.GetVersion(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	t.Log(res)
}

func TestPJSWorkerAuth(t *testing.T) {
	ctx := pctx.TestContext(t)
	pc := NewTestPachd(t, PJSWorkerAuthOption(auth.HashToken("iampjs")))
	pc.SetAuthToken(auth.HashToken("iampjs"))
	_, err := pc.ListQueue(ctx, &pjs.ListQueueRequest{})
	require.NoError(t, err)
}
