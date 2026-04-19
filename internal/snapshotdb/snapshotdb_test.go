package snapshotdb

import (
	"context"
	"github.com/laityjet/mammoth/v0/internal/dbutil"
	"github.com/laityjet/mammoth/v0/internal/storage/track"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/clusterstate"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/migrations"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/storage/chunk"
	"github.com/laityjet/mammoth/v0/internal/storage/fileset"
	"github.com/laityjet/mammoth/v0/internal/storage/kv"
	"github.com/laityjet/mammoth/v0/internal/testetcd"
)

type dependencies struct {
	ctx context.Context
	db  *pachsql.DB
	tx  *pachsql.Tx
	s   *fileset.Storage
}

func DB(t testing.TB) (context.Context, *pachsql.DB) {
	t.Helper()
	ctx := pctx.Child(pctx.TestContext(t), t.Name())
	db := dockertestenv.NewTestDB(t)
	migrationEnv := migrations.Env{EtcdClient: testetcd.NewEnv(ctx, t).EtcdClient}
	require.NoError(t, migrations.ApplyMigrations(ctx, db, migrationEnv, clusterstate.DesiredClusterState), "should be able to set up tables")
	return ctx, db
}

func FilesetStorage(t testing.TB, db *pachsql.DB) *fileset.Storage {
	t.Helper()
	tracker := track.NewPostgresTracker(db)
	s := fileset.NewStorage(fileset.NewPostgresStore(db), tracker, chunk.NewStorage(kv.NewMemStore(), nil, db, tracker))
	return s
}

func withDependencies(t *testing.T, f func(d dependencies)) {
	ctx, db := DB(t)
	s := FilesetStorage(t, db)
	err := dbutil.WithTx(ctx, db, func(ctx context.Context, sqlTx *pachsql.Tx) error {
		f(dependencies{ctx: ctx, db: db, tx: sqlTx, s: s})
		return nil
	})
	require.NoError(t, err)
}
