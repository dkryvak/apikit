package types

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
)

type Request struct {
	Method  RequestMethod
	Path    string
	Headers map[string]string
	Query   map[string]string
	Body    any
}