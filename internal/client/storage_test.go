package client_test

import (
	"testing"
	"testing/fstest"

	"github.com/laityjet/mammoth/v0/internal/dockertestenv"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/testpachd/realenv"
)

func TestFileSystemToFileSet(t *testing.T) {
	ctx := pctx.TestContext(t)
	env := realenv.NewRealEnv(ctx, t, dockertestenv.NewTestDBConfig(t).PachConfigOption)
	var testFS = fstest.MapFS{
		"foo":      &fstest.MapFile{Data: []byte("bar")},
		"baz/quux": &fstest.MapFile{Data: []byte("quuux")},
	}
	fs, err := env.PachClient.FileSystemToFileset(ctx, testFS)
	require.NoError(t, err, "must be able to upload filesystem and get fileset ID")
	t.Log("fileset", fs)
}
