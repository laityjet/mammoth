package client

import (
	"github.com/laityjet/mammoth/v0/internal/grpcutil"
	"github.com/laityjet/mammoth/v0/internal/version"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Version returns the version of pachd as a string.
func (c APIClient) Version() (string, error) {
	v, err := c.VersionAPIClient.GetVersion(c.Ctx(), &emptypb.Empty{})
	if err != nil {
		return "", grpcutil.ScrubGRPC(err)
	}
	return version.PrettyPrintVersion(v), nil
}
