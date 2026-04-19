package v2_12_0

import (
	"context"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/migrations"
	"github.com/laityjet/mammoth/v0/internal/pctx"
)

func alterBranchProvenanceTable(ctx context.Context, env migrations.Env) error {
	ctx = pctx.Child(ctx, "alterBranchProvenanceTable")
	_, err := env.Tx.ExecContext(ctx, `
		ALTER TABLE pfs.branch_provenance
		ADD COLUMN never BOOL NOT NULL DEFAULT false;
	`)
	if err != nil {
		return errors.Wrap(err, "alter branch provenance table")
	}
	return nil
}
