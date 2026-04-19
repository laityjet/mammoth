package clusterstate

import (
	v2_12_0 "github.com/laityjet/mammoth/v0/internal/clusterstate/v2.12.0"
	"github.com/laityjet/mammoth/v0/internal/migrations"
)

var state_2_12_0 migrations.State = v2_12_0.Migrate(state_2_11_0)
