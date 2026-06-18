package command

import (
	"apikit/internal/kit/file"
	"apikit/internal/kit/flag"
	"apikit/internal/kit/validator"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func PrintSchemaOrParseBody[T any](command *cli.Command, printSchema func(), inject func(*T)) (*T, error) {
	if command.Bool(flag.BodySchema.Name) {
		if printSchema != nil {
			printSchema()
		}
		return nil, nil
	}

	rawBody := strings.TrimSpace(command.String(flag.Body.Name))
	if rawBody == "" {
		return nil, fmt.Errorf(
			"required flag %q not set. Use --%s to see the expected request body format",
			flag.Body.Name,
			flag.BodySchema.Name,
		)
	}

	var body T
	if err := ParseJson(rawBody, &body); err != nil {
		return nil, err
	}

	if inject != nil {
		inject(&body)
	}

	if err := validator.Struct(body); err != nil {
		return nil, err
	}

	return &body, nil
}

func PrintSchemaOrParseQuery[T any](command *cli.Command, printSchema func(), inject func(*T)) (*T, error) {
	if command.Bool(flag.QuerySchema.Name) {
		if printSchema != nil {
			printSchema()
		}
		return nil, nil
	}

	rawBody := strings.TrimSpace(command.String(flag.Query.Name))
	if rawBody == "" {
		return nil, fmt.Errorf(
			"required flag %q not set. Use --%s to see the expected request body format",
			flag.Query.Name,
			flag.QuerySchema.Name,
		)
	}

	var body T
	if err := ParseJson(rawBody, &body); err != nil {
		return nil, err
	}

	if inject != nil {
		inject(&body)
	}

	if err := validator.Struct(body); err != nil {
		return nil, err
	}

	return &body, nil
}

func ParseJson[T any](raw string, t *T) error {
	if raw == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		raw = strings.TrimSpace(string(b))
		if raw == "" {
			return fmt.Errorf("invalid JSON: empty stdin")
		}
	}

	if err := json.Unmarshal([]byte(raw), t); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

func WriteOutIfNeeded(content []byte, suffix string, command *cli.Command, metadata Metadata, fileType FileType) error {
	out := strings.TrimSpace(command.String(flag.Out.Name))
	if out == "" {
		return nil
	}

	filename := fmt.Sprintf("%s_%s_%s.%s", metadata.Alias, metadata.Module, suffix, fileType)
	return file.WriteBytes(content, out, filename)
}
