package graphql

import (
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"context"
)

// Path is the single GraphQL endpoint path appended to the API host.
const Path = "/api/graphql"

// Request is the GraphQL HTTP envelope: a query string plus its variables.
type Request struct {
	Query     string `json:"query"`
	Variables any    `json:"variables,omitempty"`
}

// Do posts a GraphQL {query, variables} envelope to host+Path with the given
// headers via the shared HTTP client and returns the raw response. JSON content
// headers are set automatically; caller headers (auth, etc.) are merged on top.
func Do(ctx context.Context, host string, headers map[string]string, query string, variables any) (*httptypes.Response, error) {
	merged := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	for k, v := range headers {
		merged[k] = v
	}

	return httpclient.NewHttpClient(host).Do(ctx, &httptypes.Request{
		Method:  httptypes.POST,
		Path:    Path,
		Headers: merged,
		Body:    Request{Query: query, Variables: variables},
	})
}
