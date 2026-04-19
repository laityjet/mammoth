//go:build !linux

package driver

import "github.com/laityjet/mammoth/v0/internal/server/worker/logs"

func logRunningProcesses(l logs.TaggedLogger, pgid int) {
	l.Logf("warning: listing processes on this OS is not supported; you won't see debug information about which child processes we're about to kill")
}
