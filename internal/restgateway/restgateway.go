// Package restgateway implements a REST gateway.
//
// TODO: document REST gateway.  Is it really RESTful, or just a JSON-over-HTML gateway?
package restgateway

import (
	"context"
	"net/http"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/pctx"
	"github.com/laityjet/mammoth/v0/internal/metadata"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"

	"github.com/laityjet/mammoth/v0/internal/admin"
	"github.com/laityjet/mammoth/v0/internal/auth"
	"github.com/laityjet/mammoth/v0/internal/debug"
	"github.com/laityjet/mammoth/v0/internal/enterprise"
	"github.com/laityjet/mammoth/v0/internal/identity"
	"github.com/laityjet/mammoth/v0/internal/license"
	"github.com/laityjet/mammoth/v0/internal/logs"
	"github.com/laityjet/mammoth/v0/internal/proxy"
	"github.com/laityjet/mammoth/v0/internal/transaction"
	"github.com/laityjet/mammoth/v0/internal/version/versionpb"
	"github.com/laityjet/mammoth/v0/internal/worker"
)

func NewMux(ctx context.Context, grpcConn *grpc.ClientConn) (http.Handler, error) {

	ctx = pctx.Child(ctx, "restgateway")

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
		if s != "Content-Length" {
			return s, true
		}
		return s, false
	}))

	var errs error
	if err := pps.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register PPS"))
	}
	if err := pfs.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register PFS"))
	}
	if err := worker.RegisterWorkerHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register worker"))
	}
	if err := proxy.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register proxy"))
	}
	if err := logs.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register logs"))
	}
	if err := admin.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register admin"))
	}
	if err := auth.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register auth"))
	}
	if err := license.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register license"))
	}
	if err := identity.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register identity"))
	}
	if err := debug.RegisterDebugHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register debug"))
	}
	if err := enterprise.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register enterprise"))
	}
	if err := transaction.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register transaction"))
	}
	if err := versionpb.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register version"))
	}
	if err := metadata.RegisterAPIHandler(ctx, mux, grpcConn); err != nil {
		errors.JoinInto(&errs, errors.Wrap(err, "register metadata"))
	}
	return mux, errs
}
