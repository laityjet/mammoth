package server

import (
	"context"
	"time"

	"github.com/laityjet/mammoth/v0/internal/backoff"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/task"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/server/worker/pipeline/transform"
	"go.uber.org/zap"
)

type WorkerEnv struct {
	PFS         pfs.APIClient
	TaskService task.Service
}

type Worker struct {
	env WorkerEnv
}

func NewWorker(env WorkerEnv) *Worker {
	return &Worker{env: env}
}

func (w *Worker) Run(ctx context.Context) error {
	return backoff.RetryUntilCancel(ctx, func() error {
		return transform.PreprocessingWorker(ctx, w.env.PFS, w.env.TaskService, nil)
	}, backoff.NewInfiniteBackOff(), func(err error, _ time.Duration) error {
		log.Debug(ctx, "error in pps worker", zap.Error(err))
		return nil
	})
}
