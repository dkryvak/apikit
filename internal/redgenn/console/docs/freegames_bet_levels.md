# Freegames: bet-levels request body

Lists allowed bet amounts. Request body must be a JSON object.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `currency` | yes | string | Currency code. Multiple values comma-separated. | `EUR` |
| `game` | yes* | string | Game ID. Multiple values comma-separated. | `bgaming_ufo_pyramids` |
| `producer` | no* | string | Producer name (single value). | `bgaming` |

> \*One of `game` or `producer` must be passed.
>
> `login`, `password`, `cm` and `wlcode` are required by the API but injected from `--config`.

## Example

```json
{
  "currency": "EUR",
  "game": "bgaming_ufo_pyramids"
}
```
