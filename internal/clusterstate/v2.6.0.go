package clusterstate

import (
	"github.com/laityjet/mammoth/v0/internal/migrations"

	v2_6_0 "github.com/laityjet/mammoth/v0/internal/clusterstate/v2.6.0"
)

var state_2_6_0 migrations.State = v2_6_0.Migrate(state_2_5_4)

// DO NOT MODIFY THIS STATE
// IT HAS ALREADY SHIPPED IN A RELEASE
