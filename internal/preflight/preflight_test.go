package preflight

import (
	"strings"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/dbutil"
	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/testutil"
)

func TestTestMigrations(t *testing.T) {
	ctx := pctx.TestContext(t)
	dbc := dockertestenv.NewTestDBConfig(t)
	opts := dbc.Direct.DBOptions()
	opts = append(opts, dbutil.WithQueryLog(true, "migrations"))
	db := testutil.OpenDB(t, opts...)
	if err := TestMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	row := db.QueryRow(`select count(1) from migrations`)
	var rows int
	err := row.Scan(&row)
	want := `relation "migrations" does not exist`
	if err == nil {
		t.Fatalf("should get an error when querying the migrations table\n got migration row count: %v", rows)
	} else if got := err.Error(); !strings.Contains(got, want) {
		t.Fatalf("unexpected error:\n  got: %v\n want: %v", got, want)
	}

}
