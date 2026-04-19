package server

import (
	"context"

	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/pfsload"
	"github.com/laityjet/mammoth/v0/internal/task"
	"github.com/laityjet/mammoth/v0/internal/pfs"
)

type WorkerEnv struct {
	PFS         pfs.APIClient
	TaskService task.Service
}

type Worker struct {
	env WorkerEnv
}

func NewWorker(env WorkerEnv) *Worker {
	return &Worker{
		env: env,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	return pfsload.Worker(pctx.Child(ctx, "pfsload"), w.env.PFS, w.env.TaskService)
}
