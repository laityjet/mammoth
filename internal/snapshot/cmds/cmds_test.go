package cmds

import (
	"context"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dbutil"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pachd"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/snapshot"
	"github.com/laityjet/mammoth/v0/internal/testutil"
	tu "github.com/laityjet/mammoth/v0/internal/testutilpachctl"

	"github.com/laityjet/mammoth/v0/internal/pctx"
)

func TestListSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := pctx.TestContext(t)
	c := pachd.NewTestPachd(t)
	if err := tu.PachctlBashCmdCtx(ctx, t, c, `
		pachctl create snapshot
		pachctl create snapshot
		pachctl create snapshot
		pachctl list snapshot
		`,
	).Run(); err != nil {
		t.Fatalf("list snapshot RPC: %v", err)
	}
}

// there are no tests for delete. because delete gRPC implementation has not been
// merged into master

func TestRestore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	var (
		ctx   = pctx.TestContext(t)
		s     *snapshot.Snapshotter
		dbcfg dockertestenv.DBConfig
	)
	c := pachd.NewTestPachd(t, pachd.TestPachdOption{
		OnReady: func(ctx context.Context, p *pachd.Full) error {
			s = p.Snapshotter()
			return nil
		},
		CopyDBConfig: &dbcfg,
	})
	if err := tu.PachctlBashCmdCtx(ctx, t, c, `
		pachctl create repo foo
		pachctl create snapshot
		pachctl delete repo foo
		pachctl create repo bar
		(pachctl inspect repo foo && (echo "repo foo exists" >&2; exit 1)) || true
		`,
	).Run(); err != nil {
		t.Fatalf("mutate, create snapshot & mutate RPC: %v", err)
	}
	require.NoError(t, s.RestoreSnapshot(ctx, 1, snapshot.RestoreSnapshotOptions{}), "restoration must succeed")

	c = pachd.NewTestPachd(t, pachd.TestPachdOption{
		MutateEnv: func(env *pachd.Env) {
			db := testutil.OpenDB(t, dbcfg.PGBouncer.DBOptions()...)
			directDB := testutil.OpenDB(t, dbcfg.Direct.DBOptions()...)
			dbListenerConfig := dbutil.GetDSN(
				ctx,
				dbutil.WithHostPort(dbcfg.Direct.Host, int(dbcfg.Direct.Port)),
				dbutil.WithDBName(dbcfg.Direct.DBName),
				dbutil.WithUserPassword(dbcfg.Direct.User, dbcfg.Direct.Password),
				dbutil.WithSSLMode("disable"))
			env.DB = db
			env.DirectDB = directDB
			env.DBListenerConfig = dbListenerConfig
		},
	})
	// check that foo exists and bar does not
	if err := tu.PachctlBashCmdCtx(ctx, t, c, `
		pachctl list repo
		pachctl inspect repo foo
		(pachctl inspect repo bar && (echo "repo bar exists" >&2; exit 1)) || true
		`,
	).Run(); err != nil {
		t.Fatalf("post-restore check: %v", err)
	}
}
