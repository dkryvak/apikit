package flag

import "github.com/urfave/cli/v3"

var Config = &cli.StringFlag{
	Name:  "config",
	Usage: "Configuration alias (must exist). Example: --config stage",
}

var Out = &cli.StringFlag{
	Name:  "out",
	Usage: "Save response body to a file or directory (e.g. output/, output/result.json)",
}

var Body = &cli.StringFlag{
	Name:  "body",
	Usage: "Request body as JSON string. Use --body-schema to print an example.",
}

var BodySchema = &cli.BoolFlag{
	Name:  "body-schema",
	Usage: "Print JSON body schema and exit",
}

var Query = &cli.StringFlag{
	Name:  "query",
	Usage: "Request query as JSON string. Use --query-schema to print an example.",
}

var QuerySchema = &cli.BoolFlag{
	Name:  "query-schema",
	Usage: "Print JSON query schema and exit",
}
