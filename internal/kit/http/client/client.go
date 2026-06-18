package client

import (
	"apikit/internal/kit/http/console"
	"apikit/internal/kit/http/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewHttpClient(baseUrl string) *HttpClient {
	return &HttpClient{
		baseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *HttpClient) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	start := time.Now()
	console.PrintRequest(c.baseUrl, request)

	fullUrl, err := prepareUrl(c.baseUrl, request)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare fullUrl: %w", err)
	}

	var bodyReader io.Reader
	if request.Body != nil {
		bodyBytes, err := prepareBody(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, string(request.Method), fullUrl, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	response := &types.Response{
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Body:       resBody,
	}

	console.PrintResponse(response, time.Since(start))

	return response, nil
}

func prepareUrl(baseUrl string, req *types.Request) (string, error) {
	baseURL, err := url.Parse(strings.TrimSpace(baseUrl))
	if err != nil {
		return "", fmt.Errorf("invalid baseURL: %w", err)
	}

	basePath := strings.TrimRight(baseURL.Path, "/")
	reqPath := strings.TrimSpace(req.Path)
	if reqPath == "" {
		reqPath = "/"
	}
	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}
	baseURL.Path = basePath + reqPath

	if len(req.Query) > 0 {
		q := baseURL.Query()
		for k, v := range req.Query {
			q.Set(k, v)
		}
		baseURL.RawQuery = q.Encode()
	}

	return baseURL.String(), nil
}

func prepareBody(body any) ([]byte, error) {
	switch v := body.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(v)
	}
}
