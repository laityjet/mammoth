// Package testing needs to be documented.
//
// TODO: document
package testing

import (
	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	pfsserver "github.com/laityjet/mammoth/v0/internal/server/pfs"
)

func finishProjectCommit(pachClient *client.APIClient, project, repo, branch, id string) error {
	if err := pachClient.FinishCommit(project, repo, branch, id); err != nil {
		if !pfsserver.IsCommitFinishedErr(err) {
			return err
		}
	}
	_, err := pachClient.WaitCommit(project, repo, branch, id)
	return err
}

func finishCommit(pachClient *client.APIClient, repo, branch, id string) error {
	return finishProjectCommit(pachClient, pfs.DefaultProjectName, repo, branch, id)
}
