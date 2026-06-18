package console

const (
	sepLine              = "────────────────────────────────────────────────"
	indent2              = "  "
	maxResponseBodyBytes = 2 * 1024
)

var sensitiveHeaderNames = []string{
	//    "authorization", "cookie", "x-api-signature", "x-api-key",
}
