# Freegames: bet-amount list request body

Lists allowed bet amounts for multiple games. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `game_codes` | yes | string | JSON-encoded array of game identifiers. | `["playson.100_joker_staxx"]` |
| `currencies` | no | string | JSON-encoded array of ISO 4217 currency codes. | `["EUR"]` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.
>
> `game_codes` and `currencies` are sent as strings containing a JSON array (escape the quotes in the file).

## Example

```json
{
  "game_codes": "[\"playson.100_joker_staxx\"]",
  "currencies": "[\"EUR\"]"
}
```
