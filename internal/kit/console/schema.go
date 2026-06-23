package console

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

// RenderMarkdown writes the given Markdown to stdout. When stdout is a terminal
// it is styled with glamour (auto dark/light theme); otherwise the raw Markdown
// is printed unchanged, so redirecting (e.g. `--body-schema > schema.md`) yields
// clean, valid Markdown.
func RenderMarkdown(md []byte) {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if out, err := glamour.Render(string(md), "auto"); err == nil {
			fmt.Print(out)
			return
		}
	}
	fmt.Print(string(md))
}
