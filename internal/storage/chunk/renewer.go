package chunk

import (
	"context"
	"github.com/laityjet/mammoth/v0/internal/meters"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"time"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/storage/renew"
	"github.com/laityjet/mammoth/v0/internal/storage/track"
)

type Renewer struct {
	ss *renew.StringSet
}

func NewRenewer(ctx context.Context, tr track.Tracker, name string, ttl time.Duration) *Renewer {
	ctx = pctx.Child(ctx, "trackerRenewer", pctx.WithCounter("renewals", 0))
	renewFunc := func(ctx context.Context, x string, ttl time.Duration) error {
		_, err := tr.SetTTL(ctx, x, ttl)
		if err != nil {
			return errors.EnsureStack(err)
		}
		meters.Inc(ctx, "renewals", 1)
		return nil
	}
	composeFunc := renew.NewTmpComposer(tr, name)
	return &Renewer{
		ss: renew.NewStringSet(ctx, ttl, renewFunc, composeFunc),
	}
}

func (r *Renewer) Add(ctx context.Context, id ID) error {
	return r.ss.Add(ctx, id.TrackerID())
}

func (r *Renewer) Close() error {
	return r.ss.Close()
}
