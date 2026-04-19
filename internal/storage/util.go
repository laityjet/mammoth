// Package storage needs to be documented.
//
// TODO: document
package storage

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/grpcutil"
	"github.com/laityjet/mammoth/v0/internal/pachconfig"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/storage"
	"google.golang.org/grpc"
)

func NewTestServer(t testing.TB, db *pachsql.DB) *Server {
	ctx := pctx.TestContext(t)
	b, buckURL := dockertestenv.NewTestBucket(ctx, t)
	t.Log("bucket", buckURL)
	s, err := New(ctx, Env{
		DB:     db,
		Bucket: b,
		Config: pachconfig.StorageConfiguration{},
	})
	require.NoError(t, err)
	return s
}

func NewTestFilesetClient(t testing.TB, s *Server) storage.FilesetClient {
	gc := grpcutil.NewTestClient(t, func(gs *grpc.Server) {
		storage.RegisterFilesetServer(gs, s)
	})
	return storage.NewFilesetClient(gc)
}
