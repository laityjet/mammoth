package server

import (
	"github.com/laityjet/mammoth/v0/internal/collection"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/transactionenv"
	txnenv "github.com/laityjet/mammoth/v0/internal/transactionenv"
	"github.com/laityjet/mammoth/v0/internal/transaction"
)

// APIServer represents an api server.
type APIServer interface {
	transaction.APIServer
	txnenv.TransactionServer
}

type Env struct {
	DB         *pachsql.DB
	PGListener collection.PostgresListener
	TxnEnv     *transactionenv.TransactionEnv
}

// NewAPIServer creates an APIServer.
func NewAPIServer(env Env) (APIServer, error) {
	return newAPIServer(env)
}
