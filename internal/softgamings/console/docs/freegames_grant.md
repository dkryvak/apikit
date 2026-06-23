# Freegames: grant query params

Grants a freegames bonus to a player.
Endpoint: `GET /{KEY}/Freerounds/Add/`

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `Login` | yes | string | Player wallet ID (UUID). | `550e8400-e29b-41d4-a716-446655440000` |
| `Operator` | yes | string | Provider name. | `Belatra` |
| `Games` | yes | string | External game ID. | `lucky_streak` |
| `Count` | yes | integer | Number of free rounds to grant. | `10` |
| `Expire` | yes | string | Expiration, format `yyyy-MM-dd HH:mm:ss`. | `2026-12-31 23:59:59` |
| `TID` | yes | string | Wallet bonus balance ID (UUID); used for hash calculation. | `550e8400-e29b-41d4-a716-446655440001` |
| `Country` | no | string | Player country, ISO 3166-1 alpha-2. | `US` |

> `Hash` is required by the API but generated and injected automatically.

## Example

```json
{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Operator": "Belatra",
  "Games": "lucky_streak",
  "Count": 10,
  "Expire": "2026-12-31 23:59:59",
  "TID": "550e8400-e29b-41d4-a716-446655440001",
  "Country": "US"
}
```
