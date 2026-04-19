package pfs

import (
	"testing"

	"google.golang.org/grpc/status"

	"github.com/laityjet/mammoth/v0/internal/pfs"

	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/require"
)

func TestErrorMatching(t *testing.T) {
	c := client.NewCommit(pfs.DefaultProjectName, "foo", "bar", "")
	require.True(t, IsCommitNotFoundErr(ErrCommitNotFound{c}))
	require.False(t, IsCommitNotFoundErr(ErrCommitDeleted{c}))
	require.False(t, IsCommitNotFoundErr(ErrCommitFinished{c}))

	require.False(t, IsCommitDeletedErr(ErrCommitNotFound{c}))
	require.True(t, IsCommitDeletedErr(ErrCommitDeleted{c}))
	require.False(t, IsCommitDeletedErr(ErrCommitFinished{c}))

	require.False(t, IsCommitFinishedErr(ErrCommitNotFound{c}))
	require.False(t, IsCommitFinishedErr(ErrCommitDeleted{c}))
	require.True(t, IsCommitFinishedErr(ErrCommitFinished{c}))
}

type grpcStatus interface{ GRPCStatus() *status.Status }

var _ grpcStatus = ErrCommitNotFinished{}
