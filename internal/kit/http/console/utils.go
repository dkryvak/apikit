package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
)

// diagOut is where request/response diagnostics go. It is stderr so that stdout
// stays clean for the raw response body (pipeable, capturable by self-remote).
var diagOut io.Writer = os.Stderr

func cprintf(format string, a ...any) { _, _ = fmt.Fprintf(diagOut, format, a...) }
func cprintln(a ...any)               { _, _ = fmt.Fprintln(diagOut, a...) }

func printIndentedBlock(indent, text string) {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		cprintf("%s<empty>\n", indent)
		return
	}
	for _, line := range strings.Split(text, "\n") {
		cprintf("%s%s\n", indent, line)
	}
}

func bodyToPrettyString(v any) string {
    switch x := v.(type) {
    case string:
        return prettyJSONString(x)
    case []byte:
        return prettyJSONBytes(x)
    default:
        if b, err := json.MarshalIndent(x, "", "  "); err == nil {
            return string(b)
        }
        return fmt.Sprintf("%v", x)
    }
}

func prettyJSONString(s string) string {
    return prettyJSONBytes([]byte(strings.TrimSpace(s)))
}

func prettyJSONBytes(b []byte) string {
    trimmed := bytes.TrimSpace(b)
    if len(trimmed) == 0 {
        return ""
    }
    if trimmed[0] != '{' && trimmed[0] != '[' {
        return string(b)
    }
    var out bytes.Buffer
    if err := json.Indent(&out, trimmed, "", "  "); err != nil {
        return string(b)
    }
    return out.String()
}

func truncateBytes(b []byte, limit int) (preview []byte, truncated bool) {
    if limit <= 0 || len(b) <= limit {
        return b, false
    }
    return b[:limit], true
}

func extractErrorName(body []byte) string {
    trimmed := bytes.TrimSpace(body)
    if len(trimmed) == 0 || trimmed[0] != '{' {
        return ""
    }

    var m map[string]any
    if err := json.Unmarshal(trimmed, &m); err != nil {
        return ""
    }

    for _, k := range []string{"error", "code", "message"} {
        if v, ok := m[k]; ok {
            if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
                return s
            }
        }
    }
    return ""
}

func sortedKeys(m map[string]string) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys
}

func sortedHeaderKeys(h http.Header) []string {
    keys := make([]string, 0, len(h))
    for k := range h {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys
}

func isSensitiveHeader(name string) bool {
    name = strings.ToLower(strings.TrimSpace(name))
    return slices.Contains(sensitiveHeaderNames, name)
}
