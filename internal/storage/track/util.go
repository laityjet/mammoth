package track

import (
	"context"
	"time"

	"github.com/laityjet/mammoth/v0/internal/dbutil"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
)

// Create creates uses tracker to create the object id.
func Create(ctx context.Context, tr Tracker, id string, pointsTo []string, ttl time.Duration) error {
	return dbutil.WithTx(ctx, tr.DB(), func(ctx context.Context, tx *pachsql.Tx) error {
		return errors.EnsureStack(tr.CreateTx(tx, id, pointsTo, ttl))
	})
}

// Delete deletes id from the tracker
func Delete(ctx context.Context, tr Tracker, id string) error {
	return dbutil.WithTx(ctx, tr.DB(), func(ctx context.Context, tx *pachsql.Tx) error {
		return errors.EnsureStack(tr.DeleteTx(tx, id))
	})
}

// Drop sets the object at id to expire now
func Drop(ctx context.Context, tr Tracker, id string) error {
	_, err := tr.SetTTL(ctx, id, ExpireNow)
	return errors.EnsureStack(err)
}
