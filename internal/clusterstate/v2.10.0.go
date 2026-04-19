package clusterstate

import (
	v2_10_0 "github.com/laityjet/mammoth/v0/internal/clusterstate/v2.10.0"
	"github.com/laityjet/mammoth/v0/internal/migrations"
)

var state_2_10_0 migrations.State = v2_10_0.Migrate(State_2_8_0)
