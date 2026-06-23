# User: get user data query params

Returns user data by wallet ID.
Endpoint: `GET /{KEY}/User/GetUserData/`

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `Login` | yes | string | Player wallet ID (UUID). | `550e8400-e29b-41d4-a716-446655440000` |

> `TID` and `Hash` are required by the API but generated and injected automatically.

## Example

```json
{
  "Login": "550e8400-e29b-41d4-a716-446655440000"
}
```
