package grpcutil

import (
	"google.golang.org/grpc/status"

	"github.com/laityjet/mammoth/v0/internal/errors"
)

// ScrubGRPC removes GRPC error code information from 'err' if it came from
// GRPC (and returns it unchanged otherwise)
func ScrubGRPC(err error) error {
	if err == nil {
		return nil
	}
	if s, ok := status.FromError(err); ok {
		return errors.New(s.Message())
	}
	return err
}
