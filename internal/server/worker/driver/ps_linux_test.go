//go:build linux

package driver

import (
	"regexp"
	"strings"
	"syscall"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/server/worker/logs"
)

func TestLogRunningProcesses(t *testing.T) {
	ctx := pctx.TestContext(t)
	l := logs.NewTest(ctx)
	logRunningProcesses(l, syscall.Getpgrp())
	var found bool
	finder := regexp.MustCompile(`^note: about to kill.*driver[._]test`)
	for _, log := range l.Logs {
		if finder.MatchString(log) {
			found = true
			break
		}
	}
	if !found {
		logs := new(strings.Builder)
		for i, log := range l.Logs {
			if i != 0 {
				logs.WriteRune('\n')
			}
			logs.WriteString("    ")
			logs.WriteString(log)
		}
		t.Errorf("did not get info about self (driver.test or driver_test); all logs:\n%v", logs.String())
	}
}
