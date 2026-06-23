# Freegames: cancel request body

Cancels a freegames bonus. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `id` | yes | string | ID of the freegame offer generated on the Vegangster side. | `64e4c3900e812e194e6e3767` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "id": "64e4c3900e812e194e6e3767"
}
```
