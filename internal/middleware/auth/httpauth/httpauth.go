// Package httpauth extracts auth information from an HTTP request.
package httpauth

import (
	"context"
	"net/http"

	"github.com/laityjet/mammoth/v0/internal/constants"
	"github.com/laityjet/mammoth/v0/internal/client"
	"github.com/laityjet/mammoth/v0/internal/log"
	"go.uber.org/zap"
)

// ClientWithToken extracts an auth token from the HTTP request (special header or query parameter),
// and returns a Pach client that will use that token for future requests.
func ClientWithToken(ctx context.Context, c *client.APIClient, req *http.Request) *client.APIClient {
	if token := req.URL.Query().Get(constants.ContextTokenKey); token != "" {
		log.Debug(ctx, "using authn-token from URL query", zap.Int("len", len(token)))
		c.SetAuthToken(token)
		return c
	}
	if token := req.Header.Get(constants.ContextTokenKey); token != "" {
		log.Debug(ctx, "using authn-token from HTTP header", zap.Int("len", len(token)))
		c.SetAuthToken(token)
		return c
	}
	return c
}
