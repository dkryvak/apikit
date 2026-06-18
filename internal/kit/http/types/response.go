package types

import "net/http"

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}
