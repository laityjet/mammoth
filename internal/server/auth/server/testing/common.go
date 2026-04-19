// Package testing needs to be documented.
//
// TODO: document
package testing

import (
	"context"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/config"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/testpachd/realenv"
	tu "github.com/laityjet/mammoth/v0/internal/testutil"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"
)

func EnvWithAuth(t *testing.T) *realenv.RealEnv {
	t.Helper()
	ctx := pctx.TestContext(t)
	env := realenv.NewRealEnv(ctx, t, dockertestenv.NewTestDBConfig(t).PachConfigOption)
	_, err := env.AuthServer.Activate(env.PachClient.Ctx(), &auth.ActivateRequest{RootToken: tu.RootToken})
	require.NoError(t, err, "activate server should work")
	env.PachClient.SetAuthToken(tu.RootToken)
	require.NoError(t, config.WritePachTokenToConfig(tu.RootToken, false))
	client := env.PachClient.WithCtx(context.Background())
	_, err = client.PfsAPIClient.ActivateAuth(client.Ctx(), &pfs.ActivateAuthRequest{})
	require.NoError(t, err, "should be able to activate auth")
	_, err = client.PpsAPIClient.ActivateAuth(client.Ctx(), &pps.ActivateAuthRequest{})
	require.NoError(t, err, "should be able to activate auth")
	return env
}
