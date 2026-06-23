# Game: list request body

Actual game list can be obtained through the SpinByte API using the `cm=games` command.

## Fields

All fields are optional.

| Field | Type | Description | Example |
|---|---|---|---|
| `disabled` | integer | Filter by state: `0` only enabled, `1` only disabled. All games if omitted. | `0` |
| `producer` | string | Producer name. | `bgaming` |
| `provider` | string | Provider name. | |
| `promo_freespins` | integer | `1` games with freebets, `0` without. | `1` |
| `game_id` | string | Search by exact game id (game code). | |
| `title` | string | Full-text search by game title; any match is shown. | `lucky` |
| `limit` | integer | Max entries returned (first 1000 if not set). | `10` |
| `offset` | integer | Number of entries to skip. | `0` |

> `login`, `password` and `cm` are required by the API but injected from `--config`.

## Example

```json
{
  "disabled": 0,
  "producer": "bgaming",
  "promo_freespins": 1,
  "limit": 10,
  "offset": 0
}
```
