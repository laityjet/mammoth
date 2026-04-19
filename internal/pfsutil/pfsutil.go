// Package pfsutil needs to be documented.
//
// TODO: document
package pfsutil

import (
	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/pfs"
)

func MetaCommit(commit *pfs.Commit) *pfs.Commit {
	branch := ""
	if commit.Branch != nil {
		branch = commit.Branch.Name
	}
	return client.NewSystemRepo(commit.Repo.Project.GetName(), commit.Repo.Name, pfs.MetaRepoType).NewCommit(branch, commit.Id)
}
