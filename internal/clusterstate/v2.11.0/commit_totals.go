package v2_11_0

import (
	"context"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/migrations"
	"github.com/laityjet/mammoth/v0/internal/pctx"
)

func normalizeCommitTotals(ctx context.Context, env migrations.Env) error {
	ctx = pctx.Child(ctx, "normalizeCommitTotals")
	tx := env.Tx
	if _, err := tx.ExecContext(ctx, `ALTER TABLE pfs.commits ADD COLUMN total_fileset_id UUID REFERENCES storage.filesets(id)`); err != nil {
		return errors.Wrap(err, "add total_fileset column to pfs.commits")
	}
	log.Info(ctx, "normalizing pfs.commit_totals")
	_, err := tx.ExecContext(ctx, "UPDATE pfs.commits c SET total_fileset_id = fileset_id FROM pfs.commit_totals t WHERE c.commit_id = t.commit_id")
	if err != nil {
		return errors.Wrap(err, "migrate totals")
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE pfs.commit_totals;"); err != nil {
		return errors.Wrap(err, "drop the pfs.commit_totals")
	}
	return nil
}
