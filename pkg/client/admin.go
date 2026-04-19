package client

import (
	"github.com/laityjet/mammoth/v0/internal/admin"
	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/version/versionpb"
)

// InspectCluster retrieves cluster state
func (c APIClient) InspectCluster() (*admin.ClusterInfo, error) {
	clusterInfo, err := c.AdminAPIClient.InspectCluster(c.Ctx(), &admin.InspectClusterRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to inspect cluster")
	}
	return clusterInfo, nil
}

// InspectClusterWithVersion retrieves cluster state, and sends the server its
// version for the server to validate.
func (c APIClient) InspectClusterWithVersion(v *versionpb.Version) (*admin.ClusterInfo, error) {
	clusterInfo, err := c.AdminAPIClient.InspectCluster(c.Ctx(), &admin.InspectClusterRequest{
		ClientVersion: v,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to inspect cluster")
	}
	return clusterInfo, nil
}
