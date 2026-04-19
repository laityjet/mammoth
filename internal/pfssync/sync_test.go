package pfssync_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/pachd"
	"github.com/laityjet/mammoth/v0/internal/pfs"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pfssync"
	"github.com/laityjet/mammoth/v0/internal/randutil"
	"github.com/laityjet/mammoth/v0/internal/require"
	"github.com/laityjet/mammoth/v0/internal/storage/renew"
)

func BenchmarkDownload(b *testing.B) {
	pachClient := pachd.NewTestPachd(b)
	repo := "repo"
	require.NoError(b, pachClient.CreateRepo(pfs.DefaultProjectName, repo))
	commit, err := pachClient.StartCommit(pfs.DefaultProjectName, repo, "master")
	require.NoError(b, err)
	require.NoError(b, pachClient.WithModifyFileClient(commit, func(mf client.ModifyFile) error {
		for i := 0; i < 100; i++ {
			if err := mf.PutFile(fmt.Sprintf("file%d", i), randutil.NewBytesReader(rand.New(rand.NewSource(0)), 500)); err != nil {
				return errors.EnsureStack(err)
			}
		}
		return nil
	}))
	require.NoError(b, pachClient.FinishCommit(pfs.DefaultProjectName, repo, "master", commit.Id))
	fis, err := pachClient.ListFileAll(commit, "")
	require.NoError(b, err)
	require.NoError(b, pachClient.WithRenewer(func(ctx context.Context, renewer *renew.StringSet) error {
		cacheClient := pfssync.NewCacheClient(pachClient, renewer)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dir := b.TempDir()
			require.NoError(b, pfssync.WithDownloader(cacheClient, func(d pfssync.Downloader) error {
				for _, fi := range fis {
					if err := d.Download(dir, fi.File); err != nil {
						return errors.EnsureStack(err)
					}
				}
				return nil
			}))
		}
		return nil
	}))
}
