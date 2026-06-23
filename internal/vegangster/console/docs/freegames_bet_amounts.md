# Freegames: bet-amounts request body

Lists allowed bet amounts for a game. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `game_code` | yes | string | Unique game identifier (from `game list`). | `playson.100_joker_staxx` |
| `currency` | no | string (ISO 4217) | Currency code. | `EUR` |
| `country` | no | string (ISO 3166-1 alpha-2) | Player country code. | `US` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "game_code": "playson.100_joker_staxx",
  "currency": "EUR",
  "country": "US"
}
```
