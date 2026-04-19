package clusterstate

import (
	"context"

	"github.com/laityjet/mammoth/v0/internal/authdb"
	"github.com/laityjet/mammoth/v0/internal/migrations"
)

var state_2_3_0 migrations.State = state_2_1_0.
	Apply("Add internal auth user as a cluster admin", func(ctx context.Context, env migrations.Env) error {
		return authdb.InternalAuthUserPermissions(ctx, env.Tx)
	})
	// DO NOT MODIFY THIS STATE
	// IT HAS ALREADY SHIPPED IN A RELEASE
