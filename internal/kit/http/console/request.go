package console

import (
	"apikit/internal/kit/http/types"
	"fmt"
)

func PrintRequest(baseUrl string, req *types.Request) {
	cprintln(sepLine)
	cprintf("REQUEST %s %s%s\n", fmt.Sprint(req.Method), baseUrl, req.Path)
	cprintln(sepLine)

	if len(req.Query) > 0 {
		cprintln("Query:")
		for _, k := range sortedKeys(req.Query) {
			cprintf("%s%s=%s\n", indent2, k, req.Query[k])
		}
	}

	if len(req.Headers) > 0 {
		cprintln("Headers:")
		for _, k := range sortedKeys(req.Headers) {
			v := req.Headers[k]
			if isSensitiveHeader(k) {
				v = "<redacted>"
			}
			cprintf("%s%s: %s\n", indent2, k, v)
		}
	}

	if req.Body != nil {
		cprintln("Body:")
		printIndentedBlock(indent2, bodyToPrettyString(req.Body))
	}
}
