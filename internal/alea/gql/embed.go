package gql

import _ "embed"

// FreeSpinsCurrencies returns free-spin cost levels per currency for a game.
// Exposed under the CLI as `freegames bet-amounts`.
//
//go:embed freeSpinsCurrencies.graphql
var FreeSpinsCurrencies string
