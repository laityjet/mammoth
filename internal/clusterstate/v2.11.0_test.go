package clusterstate

import (
	"crypto/rand"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/migrations"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/storage/fileset"
	"github.com/laityjet/mammoth/v0/internal/testetcd"
)

func Test_v2_11_0_ClusterState(t *testing.T) {
	ctx := pctx.TestContext(t)
	db := dockertestenv.NewTestDirectDB(t)
	migrationEnv := migrations.Env{EtcdClient: testetcd.NewEnv(ctx, t).EtcdClient}

	// Pre-migration
	// Note that we are applying 2.6 migration here because we need to create collections.repos table
	require.NoError(t, migrations.ApplyMigrations(ctx, db, migrationEnv, state_2_6_0))
	setupTestData(t, ctx, db)

	// Apply migrations up to and including 2.11.0
	require.NoError(t, migrations.ApplyMigrations(ctx, db, migrationEnv, state_2_11_0))
	require.NoError(t, migrations.BlockUntil(ctx, db, state_2_11_0))
}

func newFilesetToken() fileset.Token {
	token := fileset.Token{}
	if _, err := rand.Read(token[:]); err != nil {
		panic(err)
	}
	return token
}
