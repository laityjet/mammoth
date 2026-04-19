package server

import (
	"context"
	"time"

	"github.com/laityjet/mammoth/v0/internal/identity"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pachsql"
	"github.com/laityjet/mammoth/v0/internal/protoutil"
)

type dbUser struct {
	Email             string     `db:"email"`
	LastAuthenticated *time.Time `db:"last_authenticated"`
}

func (u dbUser) ToProto() *identity.User {
	return &identity.User{
		Email:             u.Email,
		LastAuthenticated: protoutil.MustTimestampFromPointer(u.LastAuthenticated),
	}
}

func listUsers(ctx context.Context, db *pachsql.DB) ([]*identity.User, error) {
	var users []dbUser
	err := db.SelectContext(ctx, &users, "SELECT email, last_authenticated FROM identity.users WHERE enabled=true;")
	if err != nil {
		return nil, errors.EnsureStack(err)
	}
	result := make([]*identity.User, len(users))
	for i, u := range users {
		result[i] = u.ToProto()
	}
	return result, nil
}

func addUserInTx(ctx context.Context, tx *pachsql.Tx, email string) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO identity.users (email, last_authenticated, enabled) VALUES ($1, now(), true) ON CONFLICT(email) DO UPDATE SET last_authenticated=NOW()`, email)
	return errors.EnsureStack(err)
}
