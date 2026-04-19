package fileset

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pctx"
)

func TestPostgresStore(t *testing.T) {
	StoreTestSuite(t, func(t testing.TB) MetadataStore {
		ctx := pctx.TestContext(t)
		db := dockertestenv.NewTestDB(t)
		return NewTestStore(ctx, t, db)
	})
}
