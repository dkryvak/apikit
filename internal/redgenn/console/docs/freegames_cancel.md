# Freegames: cancel request body

Cancels (removes) a free-spin offer from a player (Redgenn `remove_offer`, XML over `GET /admservice`). Only possible while the series has not started. `login` and `password` are injected from `--config`.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `giftId` | yes | string (UUID) | The gift id of the offer assigned to the player. | `dc921b76-f5ec-4d33-85e6-e73ca6dc754d` |

## Example

```json
{
  "giftId": "dc921b76-f5ec-4d33-85e6-e73ca6dc754d"
}
```
