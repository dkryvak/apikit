package console

import "apikit/internal/kit/console"

func PrintGetUserQuerySchema() {
	console.PrintSchema(console.Schema{
		Title: "User: get user data query params",
		Intro: "Returns user data by wallet ID. Endpoint: GET /{KEY}/User/GetUserData/",
		Fields: []console.SchemaField{
			{
				Name:        "Login",
				Required:    true,
				Type:        "string",
				Description: "Player wallet ID (UUID)",
				Example:     "550e8400-e29b-41d4-a716-446655440000",
			},
		},
		Example: `{
  "Login": "550e8400-e29b-41d4-a716-446655440000"
}`,
		Notes: []string{
			`Fields "TID" and "Hash" are required by API, but CLI generates and injects them automatically.`,
		},
	})
}
