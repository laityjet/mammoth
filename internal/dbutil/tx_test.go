package dbutil_test

import (
	"context"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dbutil"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestNestingTransactionsPanics(t *testing.T) {
	ctx := pctx.TestContext(t)
	db := dockertestenv.NewTestDB(t)
	var withTxErr error
	require.YesPanic(t, func() {
		withTxErr = dbutil.WithTx(ctx, db, func(ctx context.Context, _ *pachsql.Tx) error {
			err := dbutil.WithTx(ctx, db, func(_ context.Context, _ *pachsql.Tx) error {
				return nil
			})
			return errors.Wrap(err, "nested WithTx")
		})
	}, "nesting transactions should panic")
	require.NoError(t, withTxErr, "WithTx should not have errored")
}
