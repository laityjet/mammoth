package client

import (
	"io"

	"github.com/laityjet/mammoth/v0/internal/debug"
	"github.com/laityjet/mammoth/v0/internal/grpcutil"
	"github.com/laityjet/mammoth/v0/internal/pctx"
)

// Profile collects a set of pprof profiles.
func (c APIClient) Profile(profile *debug.Profile, filter *debug.Filter, w io.Writer) (retErr error) {
	defer func() {
		retErr = grpcutil.ScrubGRPC(retErr)
	}()
	profileC, err := c.DebugClient.Profile(c.Ctx(), &debug.ProfileRequest{
		Profile: profile,
		Filter:  filter,
	})
	if err != nil {
		return err
	}
	return grpcutil.WriteFromStreamingBytesClient(profileC, w)
}

// Binary collects a set of binaries.
func (c APIClient) Binary(filter *debug.Filter, w io.Writer) (retErr error) {
	defer func() {
		retErr = grpcutil.ScrubGRPC(retErr)
	}()
	binaryC, err := c.DebugClient.Binary(c.Ctx(), &debug.BinaryRequest{Filter: filter})
	if err != nil {
		return err
	}
	return grpcutil.WriteFromStreamingBytesClient(binaryC, w)
}

// Dump collects a standard set of debugging information.
func (c APIClient) Dump(filter *debug.Filter, limit int64, w io.Writer) (retErr error) {
	defer func() {
		retErr = grpcutil.ScrubGRPC(retErr)
	}()
	ctx, cf := pctx.WithCancel(c.Ctx())
	defer cf()
	dumpC, err := c.DebugClient.Dump(ctx, &debug.DumpRequest{
		Filter: filter,
		Limit:  limit,
	})
	if err != nil {
		return err
	}
	return grpcutil.WriteFromStreamingBytesClient(dumpC, w)
}
