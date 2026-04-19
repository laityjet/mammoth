// Package transactiondb contains the database schema that Pachyderm
// transactions use.
package transactiondb

import (
	col "github.com/laityjet/mammoth/v0/internal/collection"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/transaction"
)

const (
	transactionsCollectionName = "transactions"
)

var transactionsIndexes = []*col.Index{}

// Transactions returns a collection of open transactions
func Transactions(db *pachsql.DB, listener col.PostgresListener) col.PostgresCollection {
	return col.NewPostgresCollection(
		transactionsCollectionName,
		db,
		listener,
		&transaction.TransactionInfo{},
		transactionsIndexes,
	)
}

// CollectionsV0 returns a list of all the Transaction API collections for
// postgres-initialization purposes. These collections are not usable for
// querying.
//
// DO NOT MODIFY THIS FUNCTION
// IT HAS BEEN USED IN A RELEASED MIGRATION
func CollectionsV0() []col.PostgresCollection {
	return []col.PostgresCollection{
		col.NewPostgresCollection(transactionsCollectionName, nil, nil, nil, transactionsIndexes),
	}
}
