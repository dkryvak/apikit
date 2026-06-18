package kube

import (
	"apikit/internal/env"
	cfgkit "apikit/internal/kit/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// SelfRemote executes a call inside the cluster by re-invoking `apikit local ...`
// in an apikit pod. The config is imported into the pod first, the body/query is
// piped via stdin, and the pod's stdout (raw response) is captured back locally.
func SelfRemote(ctx context.Context, command *cli.Command, module, alias string, cfg cfgkit.Config) error {
	envName := strings.TrimSpace(cfg.BoundEnv())
	if envName == "" {
		return fmt.Errorf(
			"config %q has no bound env; set 'env' in the config to run via 'apikit kube' (or use 'apikit local')",
			alias,
		)
	}

	e, err := env.Load(envName)
	if err != nil {
		return err
	}

	if err := EnsureAWSSSO(ctx); err != nil {
		return err
	}

	configJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	payload, kind, err := readPayload(command)
	if err != nil {
		return err
	}

	executor, err := NewExecutor(e.Context, e.Namespace)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "→ kube: env=%s namespace=%s\n", e.Name, e.Namespace)

	pod, err := executor.EnsurePod(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "→ kube: using pod %s\n", pod)

	// Import the config into the pod (under the same module).
	importArgs := []string{"apikit", module, "config", "import", alias, "-"}
	if _, eErr, err := executor.Exec(ctx, pod, importArgs, configJSON); err != nil {
		return fmt.Errorf("config import in pod failed: %w\n%s", err, string(eErr))
	}

	// Run the call in local mode inside the pod.
	innerArgs := buildInnerArgs(alias, kind)
	stdout, stderr, err := executor.Exec(ctx, pod, innerArgs, payload)
	if len(stderr) > 0 {
		os.Stderr.Write(stderr)
	}
	if err != nil {
		return err
	}

	return writeResult(command, stdout)
}

// readPayload returns the request body or query payload (resolving "-" from
// stdin) and which kind it is ("body", "query", or "").
func readPayload(command *cli.Command) ([]byte, string, error) {
	if raw := strings.TrimSpace(command.String("body")); raw != "" {
		data, err := resolveRaw(raw)
		return data, "body", err
	}
	if raw := strings.TrimSpace(command.String("query")); raw != "" {
		data, err := resolveRaw(raw)
		return data, "query", err
	}
	return nil, "", nil
}

func resolveRaw(raw string) ([]byte, error) {
	if raw == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read stdin: %w", err)
		}
		return data, nil
	}
	return []byte(raw), nil
}

// buildInnerArgs reconstructs the command path from os.Args, swaps the `kube`
// namespace for `local`, and appends normalized flags so the pod emits the raw
// body to stdout.
func buildInnerArgs(alias, payloadKind string) []string {
	path := commandPath(os.Args[1:])
	if len(path) > 0 && path[0] == "kube" {
		path[0] = "local"
	}

	args := append([]string{"apikit"}, path...)
	args = append(args, "--config", alias)
	switch payloadKind {
	case "body":
		args = append(args, "--body", "-")
	case "query":
		args = append(args, "--query", "-")
	}
	args = append(args, "--out", "-")
	return args
}

// commandPath returns the leading non-flag tokens (the command chain).
func commandPath(raw []string) []string {
	var path []string
	for _, a := range raw {
		if strings.HasPrefix(a, "-") {
			break
		}
		path = append(path, a)
	}
	return path
}

// writeResult sends the captured raw body to the local --out target.
func writeResult(command *cli.Command, body []byte) error {
	out := strings.TrimSpace(command.String("out"))
	if out == "" || out == "-" {
		_, err := os.Stdout.Write(body)
		return err
	}
	if err := os.WriteFile(out, body, 0o644); err != nil {
		return fmt.Errorf("failed to write %q: %w", out, err)
	}
	fmt.Fprintf(os.Stderr, "✓ saved response to %s\n", out)
	return nil
}
