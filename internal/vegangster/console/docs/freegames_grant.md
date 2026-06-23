# Freegames: grant request body

Grants a freegames bonus to a player. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `player_id` | yes | string | Unique player ID in the operator's system. | `64651509b8c355917ec34421` |
| `country` | yes | string | Player country, ISO 3166-1 alpha-2. | `US` |
| `ip` | yes | string | Player IP address. | `145.22.35.62` |
| `reference` | yes | string | Operator-side reference for the grant. | `Ref_1_fsefgrh` |
| `game_code` | yes | string | Unique game identifier (from `game list`). | `playson.100_joker_staxx` |
| `rounds` | yes | integer | Number of rounds in the freegame session. | `10` |
| `rounds_bet` | yes | integer | Bet per round in subunits (100 = 1.00 of currency). | `100` |
| `currency` | yes | string | Currency code, ISO 4217. | `USD` |
| `end_date` | yes | date | ISO 8601 datetime `YYYY-MM-DDThh:mm:ss`. | `2023-05-11T12:22:35` |
| `offer_end_date` | yes | date | ISO 8601 datetime `YYYY-MM-DDThh:mm:ss`. | `2023-05-11T12:22:35` |
| `start_date` | no | date | ISO 8601 datetime `YYYY-MM-DDThh:mm:ss`. | `2023-05-11T12:22:35` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "player_id": "64651509b8c355917ec34421",
  "country": "US",
  "ip": "145.22.35.62",
  "reference": "Ref_1_fsefgrh",
  "game_code": "playson.100_joker_staxx",
  "rounds": 10,
  "rounds_bet": 100,
  "currency": "USD",
  "end_date": "2023-05-11T12:22:35",
  "offer_end_date": "2023-05-11T12:22:35",
  "start_date": "2023-05-11T12:22:35"
}
```
