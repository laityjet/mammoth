package main

import (
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"go.uber.org/zap"
)

func main() {
	log.InitPachctlLogger()
	log.SetLevel(log.DebugLevel)
	ctx := pctx.Background("")

	sctx, done := log.SpanContext(ctx, "postgres")
	if err := dockertestenv.EnsureDBEnv(sctx); err != nil {
		done(zap.Error(err))
		log.Exit(ctx, "error", zap.Error(err))
	}
	done()

	sctx, done = log.SpanContext(ctx, "minio")
	if _, err := dockertestenv.NewMinioClient(sctx); err != nil {
		done(zap.Error(err))
		log.Exit(ctx, "error", zap.Error(err))
	}
	done()
}
