// Package wrapcheck provides a wrapcheck analyzer.
package wrapcheck

import (
	"github.com/tomarrell/wrapcheck/v2/wrapcheck"
	"golang.org/x/tools/go/analysis"
)

var Analyzer *analysis.Analyzer

func init() {
	cfg := wrapcheck.NewDefaultConfig()
	cfg.IgnoreSigs = []string{
		"github.com/laityjet/mammoth/v0/internal/errors.Errorf",
		"github.com/laityjet/mammoth/v0/internal/errors.New",
		"github.com/laityjet/mammoth/v0/internal/errors.Unwrap",
		"github.com/laityjet/mammoth/v0/internal/errors.EnsureStack",
		"github.com/laityjet/mammoth/v0/internal/errors.Join",
		"github.com/laityjet/mammoth/v0/internal/errors.JoinInto",
		"github.com/laityjet/mammoth/v0/internal/errors.Close",
		"github.com/laityjet/mammoth/v0/internal/errors.Invoke",
		"github.com/laityjet/mammoth/v0/internal/errors.Invoke1",
		"google.golang.org/grpc/status.Error",
		"google.golang.org/grpc/status.Errorf",
		"(*google.golang.org/grpc/internal/status.Status).Err",
		"google.golang.org/protobuf/types/known/anypb.New",
		".Wrap(",
		".Wrapf(",
		".WithMessage(",
		".WithMessagef(",
		".WithStack(",
	}
	cfg.IgnorePackageGlobs = []string{
		"github.com/laityjet/mammoth/v0/*",
	}
	cfg.IgnoreInterfaceRegexps = []string{
		`^fileset\.`,
		`^collection\.`,
		`^track\.`,
	}
	Analyzer = wrapcheck.NewAnalyzer(cfg)
}
