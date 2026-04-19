// Package taskprotos needs to be documented.
//
// TODO: document
package taskprotos

import (
	// storage and compaction tasks
	_ "github.com/laityjet/mammoth/v0/internal/server/pfs/server"
	// worker transform tasks
	_ "github.com/laityjet/mammoth/v0/internal/server/worker/pipeline/transform"
	// worker datum tasks
	_ "github.com/laityjet/mammoth/v0/internal/server/worker/datum"
)
