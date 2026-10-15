# Freegames: bet-amount list request body

Lists allowed bet amounts for multiple games. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `game_codes` | yes | array of strings | Game identifiers. | `["playson.100_joker_staxx"]` |
| `currencies` | no | array of strings | ISO 4217 currency codes. | `["EUR"]` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "game_codes": ["playson.100_joker_staxx"],
  "currencies": ["EUR"]
}
```
