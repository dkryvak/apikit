# Freegames: bet-amounts request body

Lists free-spin cost levels per currency for a game (Alea `freeSpinsCurrencies`).
The request body is the GraphQL query variables as a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `gameId` | yes | int | Alea game id. | `15836` |
| `currencies` | no | string[] | Currency codes (ISO 4217) to include. | `["EUR","USD"]` |

> `jurisdictionCode` is required by the API but injected from `--config`.

## Example

```json
{
  "gameId": 15836,
  "currencies": ["EUR", "USD", "CAD", "NZD", "AUD"]
}
```
