package env

import (
	"apikit/internal/kit/console"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptEnv runs an interactive wizard to build an Env.
func PromptEnv() (*Env, error) {
	reader := bufio.NewReader(os.Stdin)
	var e Env

	name, err := console.Prompt(reader, "Env name (e.g., 'prod', 'stage')")
	if err != nil {
		return nil, err
	}
	e.Name = name

	if Exists(name) {
		overwrite, err := console.PromptYesNo(reader, fmt.Sprintf("Env '%s' already exists. Overwrite?", name))
		if err != nil {
			return nil, err
		}
		if !overwrite {
			return nil, fmt.Errorf("env creation cancelled")
		}
	}

	namespace, err := console.PromptOptional(reader, fmt.Sprintf("Namespace (default '%s-casino')", name))
	if err != nil {
		return nil, err
	}
	e.Namespace = strings.TrimSpace(namespace)

	context, err := console.PromptOptional(reader, "Kube context (optional, resolved by name if empty)")
	if err != nil {
		return nil, err
	}
	e.Context = strings.TrimSpace(context)

	e.ApplyDefaults()
	return &e, nil
}

func PrintEnv(e *Env) {
	fmt.Printf("========================================\n")
	fmt.Printf("Env: %s\n", e.Name)
	fmt.Printf("========================================\n")
	fmt.Printf("Namespace: %s\n", e.Namespace)
	if e.Context != "" {
		fmt.Printf("Context: %s\n", e.Context)
	} else {
		fmt.Printf("Context: (resolved by name)\n")
	}
	fmt.Printf("========================================\n")
}
