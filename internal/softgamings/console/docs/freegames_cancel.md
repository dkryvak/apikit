# Freegames: cancel query params

Cancels a freegames bonus for a player.
Endpoint: `GET /{KEY}/Freerounds/Remove/`

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `Login` | yes | string | Player wallet ID (UUID). | `550e8400-e29b-41d4-a716-446655440000` |
| `Operator` | yes | string | Provider name. | `Belatra` |
| `ExtID` | yes | string | Wallet bonus balance ID (UUID). | `550e8400-e29b-41d4-a716-446655440001` |

> `TID` and `Hash` are required by the API but generated and injected automatically.

## Example

```json
{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Operator": "Belatra",
  "ExtID": "550e8400-e29b-41d4-a716-446655440001"
}
```
