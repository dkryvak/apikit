# Freegames: cancel request body

Cancels a previously granted free-spin bonus (Alea `DELETE /api/v1/free-spins/{casinoBonusId}`, signed). A successful response is HTTP 204 with an empty body.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `casinoBonusId` | yes | string | The id of the bonus grant to cancel. | `856543-603` |

> `environmentToken`, `secretKey`, `casinoId` and `environmentId` come from `--config`. The `Digest` header signs `casinoBonusId`.

## Example

```json
{
  "casinoBonusId": "856543-603"
}
```
