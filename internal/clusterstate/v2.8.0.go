package clusterstate

import (
	v2_8_0 "github.com/laityjet/mammoth/v0/internal/clusterstate/v2.8.0"
	"github.com/laityjet/mammoth/v0/internal/migrations"
)

var State_2_8_0 migrations.State = v2_8_0.Migrate(state_2_7_0)
