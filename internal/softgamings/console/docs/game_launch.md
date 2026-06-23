# Game: launch query params

Creates a real-money game session URL for a player.
Endpoint: `GET /{KEY}/User/AuthHTML/`

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `Login` | yes | string | Player wallet ID (UUID). | `550e8400-e29b-41d4-a716-446655440000` |
| `Password` | yes | string | Player ID (UUID). | `550e8400-e29b-41d4-a716-446655440001` |
| `Currency` | yes | string | Player currency, ISO 4217. | `USD` |
| `ExtParam` | yes | string | Game session ID (UUID). | `550e8400-e29b-41d4-a716-446655440002` |
| `System` | yes | string | External provider ID. | `789` |
| `Page` | yes | string | External game ID. | `123` |
| `UserAutoCreate` | no | integer | Auto-create the player if missing: `1` yes, `0` no. | `1` |
| `IsMobile` | no | integer | Device type: `1` mobile, `0` desktop. | `0` |
| `UserIP` | no | string | Player real IP (for licensing); server IP used as fallback. | `192.168.1.1` |

> `TID` and `Hash` are required by the API but generated and injected automatically.

## Example

```json
{
  "Login": "550e8400-e29b-41d4-a716-446655440000",
  "Password": "550e8400-e29b-41d4-a716-446655440001",
  "Currency": "USD",
  "ExtParam": "550e8400-e29b-41d4-a716-446655440002",
  "System": "789",
  "Page": "123",
  "UserAutoCreate": 1,
  "IsMobile": 0,
  "UserIP": "192.168.1.1"
}
```
