// Package server needs to be documented.
//
// TODO: document
package server

import (
	"context"
	"strings"
	"time"

	"github.com/laityjet/mammoth/v0/internal/errors"

	etcd "go.etcd.io/etcd/client/v3"
	"gocloud.dev/blob"

	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"

	pfsserver "github.com/laityjet/mammoth/v0/internal/server/pfs"
	pps_server "github.com/laityjet/mammoth/v0/internal/server/pps"

	"github.com/laityjet/mammoth/v0/internal/client"
	col "github.com/laityjet/mammoth/v0/internal/collection"
	"github.com/laityjet/mammoth/v0/internal/pachconfig"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/task"
	txnenv "github.com/laityjet/mammoth/v0/internal/transactionenv"
	"github.com/laityjet/mammoth/v0/internal/transactionenv/txncontext"
)

const (
	StorageTaskNamespace = "storage"
	fileSetsRepo         = client.FileSetsRepoName
	defaultTTL           = client.DefaultTTL
	maxTTL               = 30 * time.Minute
)

type APIServer = *validatedAPIServer

type PipelineInspector interface {
	InspectPipelineInTransaction(context.Context, *txncontext.TransactionContext, *pps.Pipeline) (*pps.PipelineInfo, error)
}

// PFSAuth contains the auth methods called by PFS.
// It is a subset of what the Auth Service provides.
type PFSAuth interface {
	CheckRepoIsAuthorized(ctx context.Context, repo *pfs.Repo, p ...auth.Permission) error
	WhoAmI(ctx context.Context, req *auth.WhoAmIRequest) (*auth.WhoAmIResponse, error)
	GetPermissions(ctx context.Context, req *auth.GetPermissionsRequest) (*auth.GetPermissionsResponse, error)

	CheckProjectIsAuthorizedInTransaction(ctx context.Context, txnCtx *txncontext.TransactionContext, project *pfs.Project, p ...auth.Permission) error
	CheckRepoIsAuthorizedInTransaction(ctx context.Context, txnCtx *txncontext.TransactionContext, repo *pfs.Repo, p ...auth.Permission) error
	CreateRoleBindingInTransaction(ctx context.Context, txnCtx *txncontext.TransactionContext, principal string, roleSlice []string, resource *auth.Resource) error
	DeleteRoleBindingInTransaction(ctx context.Context, transactionContext *txncontext.TransactionContext, resource *auth.Resource) error
	GetPermissionsInTransaction(ctx context.Context, txnCtx *txncontext.TransactionContext, req *auth.GetPermissionsRequest) (*auth.GetPermissionsResponse, error)
}

// Env is the dependencies needed to run the PFS API server
type Env struct {
	Bucket      *blob.Bucket
	DB          *pachsql.DB
	EtcdPrefix  string
	EtcdClient  *etcd.Client
	TaskService task.Service
	TxnEnv      *txnenv.TransactionEnv
	Listener    col.PostgresListener

	Auth                 PFSAuth
	GetPipelineInspector func() PipelineInspector

	StorageConfig pachconfig.StorageConfiguration
	GetPPSServer  func() pps_server.APIServer
}

// NewAPIServer creates an APIServer.
func NewAPIServer(ctx context.Context, env Env) (pfsserver.APIServer, error) {
	a, err := newAPIServer(ctx, env)
	if err != nil {
		return nil, err
	}
	return newValidatedAPIServer(a, env.Auth), nil
}

// IsPermissionError returns true if a given error is a permission error.
func IsPermissionError(err error) bool {
	return strings.Contains(err.Error(), "has already finished")
}

func (a *apiServer) getPermissionsInTransaction(ctx context.Context, txnCtx *txncontext.TransactionContext, repo *pfs.Repo) ([]auth.Permission, []string, error) {
	resp, err := a.env.Auth.GetPermissionsInTransaction(ctx, txnCtx, &auth.GetPermissionsRequest{Resource: repo.AuthResource()})
	if err != nil {
		return nil, nil, errors.EnsureStack(err)
	}

	return resp.Permissions, resp.Roles, nil
}
