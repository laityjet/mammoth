// Package transform needs to be documented.
//
// TODO: document
package transform

import (
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/laityjet/mammoth/v0/internal/backoff"
	"github.com/laityjet/mammoth/v0/internal/errutil"
	"github.com/laityjet/mammoth/v0/internal/pps"
	"github.com/laityjet/mammoth/v0/internal/server/worker/driver"
	"github.com/laityjet/mammoth/v0/internal/server/worker/logs"
)

// Run will run a transform pipeline until the driver is canceled.
func Run(driver driver.Driver, logger logs.TaggedLogger) error {
	reg, err := NewRegistry(driver, logger)
	if err != nil {
		return err
	}
	logger.Logf("transform spawner started")
	// wrap SubscribeJob() in a retry to mitigate database connection flakiness.
	return backoff.RetryUntilCancel(driver.PachClient().Ctx(),
		func() error {
			err := driver.PachClient().SubscribeJob(
				driver.PipelineInfo().Pipeline.Project.GetName(),
				driver.PipelineInfo().Pipeline.Name,
				true,
				func(jobInfo *pps.JobInfo) error {
					if jobInfo.PipelineVersion != driver.PipelineInfo().Version {
						// Skip this job - we should be shut down soon, but don't error out in the meantime
						return nil
					}
					if jobInfo.State == pps.JobState_JOB_FINISHING {
						return nil
					}
					return reg.StartJob(proto.Clone(jobInfo).(*pps.JobInfo))
				},
			)
			if errutil.IsDatabaseDisconnect(err) {
				logger.Logf("retry SubscribeJob() in transform.Run(); err: %v", err)
				return backoff.ErrContinue
			}
			return err
		}, backoff.RetryEvery(time.Second), nil)
}
