package server

import (
	"context"

	"github.com/laityjet/mammoth/v0/internal/identity"
	col "github.com/laityjet/mammoth/v0/internal/collection"
	"github.com/laityjet/mammoth/v0/internal/pachconfig"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	txnenv "github.com/laityjet/mammoth/v0/internal/transactionenv"
	"github.com/laityjet/mammoth/v0/internal/server/enterprise"
	"github.com/laityjet/mammoth/v0/internal/server/pfs"
	"github.com/laityjet/mammoth/v0/internal/server/pps"
	etcd "go.etcd.io/etcd/client/v3"
)

// Env is the environment required for an apiServer
type Env struct {
	DB         *pachsql.DB
	EtcdClient *etcd.Client
	Listener   col.PostgresListener
	TxnEnv     *txnenv.TransactionEnv

	// circular dependency
	GetEnterpriseServer func() enterprise.APIServer
	GetIdentityServer   func() identity.APIServer
	GetPfsServer        func() pfs.APIServer
	GetPpsServer        func() pps.APIServer

	BackgroundContext context.Context
	Config            pachconfig.Configuration
}
