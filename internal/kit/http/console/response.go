package console

import (
	"apikit/internal/kit/http/types"
	"net/http"
	"strings"
	"time"
)

func PrintResponse(res *types.Response, duration time.Duration) {
	d := duration.Round(time.Millisecond)

	cprintln(sepLine)
	if res == nil {
		cprintf("RESPONSE <nil> (%s)\n", d)
		cprintln(sepLine)
		return
	}

	cprintf("RESPONSE %d %s (%s)\n", res.StatusCode, http.StatusText(res.StatusCode), d)
	cprintln(sepLine)

	if len(res.Headers) > 0 {
		cprintln("Headers:")
		for _, k := range sortedHeaderKeys(res.Headers) {
			v := strings.Join(res.Headers.Values(k), ", ")
			if isSensitiveHeader(k) {
				v = "<redacted>"
			}
			cprintf("%s%s: %s\n", indent2, k, v)
		}
	}

	if errName := extractErrorName(res.Body); errName != "" {
		cprintf("Error: %s\n", errName)
	}

	if len(res.Body) > 0 {
		cprintln("Body:")

		preview, truncated := truncateBytes(res.Body, maxResponseBodyBytes)
		text := prettyJSONBytes(preview)

		printIndentedBlock(indent2, text)

		if truncated {
			cprintf("%s… (truncated, %d/%d bytes shown)\n", indent2, len(preview), len(res.Body))
		}
	}
}
