// Package spout needs to be documented.
//
// TODO: document
package spout

import (
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/server/worker/driver"
	"github.com/laityjet/mammoth/v0/internal/server/worker/logs"
)

// Run will run a spout pipeline until the driver is canceled.
func Run(driver driver.Driver, logger logs.TaggedLogger) error {
	logger = logger.WithJob("spout")
	return errors.EnsureStack(driver.RunUserCode(driver.PachClient().Ctx(), logger, nil))
}
