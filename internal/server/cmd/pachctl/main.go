package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/log"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/tracing"
	"github.com/laityjet/mammoth/v0/internal/server/cmd/pachctl/cmd"
	"github.com/spf13/pflag"
)

func main() {
	log.InitPachctlLogger()
	log.SetLevel(log.InfoLevel)
	ctx := pctx.Background("pachctl")

	// Remove kubernetes client flags from the spf13 flag set
	// (we link the kubernetes client, so otherwise they're in 'pachctl --help')
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	tracing.InstallJaegerTracerFromEnv()
	err := func() error {
		defer tracing.CloseAndReportTraces()
		pachctl, err := cmd.PachctlCmd()
		if err != nil {
			return errors.Wrap(err, "could not create pachctl command")
		}
		return errors.EnsureStack(pachctl.ExecuteContext(ctx))
	}()
	if err != nil {
		if errString := strings.TrimSpace(err.Error()); errString != "" {
			fmt.Fprintf(os.Stderr, "%s\n", errString)
		}
		os.Exit(1)
	}
}
