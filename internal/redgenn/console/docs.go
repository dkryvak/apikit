package console

import (
	"embed"
	"fmt"
	"os"

	kconsole "apikit/internal/kit/console"
)

//go:embed docs/*.md
var docsFS embed.FS

// renderDoc renders an embedded Markdown doc (by file name under docs/) to
// stdout via the shared TTY-aware renderer.
func renderDoc(name string) {
	b, err := docsFS.ReadFile("docs/" + name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "schema doc %q not found\n", name)
		return
	}
	kconsole.RenderMarkdown(b)
}
