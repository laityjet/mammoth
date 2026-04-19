// Command pachhttp runs the Pachyderm HTTP server locally, against your current pachctl context.
package main

import (
	"context"
	"os/signal"

	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/pachctl"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/signals"
	"github.com/laityjet/mammoth/v0/internal/server/http"
	"go.uber.org/zap"
)

func main() {
	log.InitPachctlLogger()
	log.SetLevel(log.DebugLevel)
	ctx, cancel := signal.NotifyContext(pctx.Background(""), signals.TerminationSignals...)
	defer cancel()

	pc := &pachctl.Config{Verbose: true}
	c, err := pc.NewOnUserMachine(ctx, false)
	if err != nil {
		log.Exit(ctx, "problem creating pachyderm client", zap.Error(err))
	}
	s, err := http.New(ctx, 1659, func(ctx context.Context) *client.APIClient { return c.WithCtx(ctx) })
	if err != nil {
		log.Exit(ctx, "problem creating http server", zap.Error(err))
	}
	log.Info(ctx, "starting server on port 1659")
	if err := s.ListenAndServe(ctx); err != nil {
		log.Exit(ctx, "problem running http server", zap.Error(err))
	}
}
