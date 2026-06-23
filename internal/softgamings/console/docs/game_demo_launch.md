# Game: demo launch query params

Creates a demo (free-play) game session URL.
Endpoint: `GET /{KEY}/User/AuthHTML/`

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `System` | yes | string | External provider ID. | `634` |
| `Page` | yes | string | External game ID. | `235` |
| `UserAutoCreate` | no | integer | Auto-create the player if missing: `1` yes, `0` no. | `1` |
| `IsMobile` | no | integer | Device type: `1` mobile, `0` desktop. | `0` |
| `UserIP` | no | string | Player real IP (for licensing); server IP used as fallback. | `192.168.1.1` |

> `Login`, `Password` and `Demo` are required by the API but injected automatically.
>
> `TID` and `Hash` are required by the API but generated and injected automatically.

## Example

```json
{
  "System": "634",
  "Page": "235",
  "UserAutoCreate": 1,
  "IsMobile": 0,
  "UserIP": "192.168.1.1"
}
```
