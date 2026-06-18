package ip

import (
	httpclient "apikit/internal/kit/http/client"
	httptypes "apikit/internal/kit/http/types"
	"context"
	"strings"
)

func GetCurrentIp(ctx context.Context) (string, error) {
	response, err := httpclient.NewHttpClient("https://checkip.amazonaws.com").Do(ctx, &httptypes.Request{
		Method: httptypes.GET,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(response.Body))
	return ip, nil
}
