# Freegames: grant request body

Grants a free-spin offer to a player (Redgenn). This runs two provider calls in order:
`create_offer` then `add_offer` (XML over `GET /admservice`). `login`, `password` and
`wlcode` are injected from `--config`.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `offer` | yes | string | Unique offer code (≤100 chars). | `welcome-100` |
| `game` | yes | string | Game code. | `bgaming_ufo_pyramids` |
| `spins` | yes | int | Number of free spins. | `10` |
| `totalbet` | yes | int | Total bet per spin, in cents. | `1000` |
| `currency` | yes | string | Currency code (ISO 4217). | `EUR` |
| `wlid` | yes | string | Player id (used by `add_offer`). | `183019832` |
| `startLive` | no | string | `create_offer` start; `DD/MM/YYYY HH:MM` UTC. | `30/09/2024 13:43` |
| `endDateAddOffer` | no | string | `create_offer` end; `DD/MM/YYYY HH:MM` UTC. | `30/09/2025 13:43` |
| `lifetime` | no | string | `add_offer` expiry; `DD/MM/YYYY HH:MM` UTC. | `31/12/2025 12:55` |
| `cmId` | no | string | `add_offer` idempotency key (uuid). | `6a862e9a-c1b1-44d4-a33d-eccedbef41f3` |

> The saved output is the `add_offer` response; the `create_offer` response is shown in the request diagnostics (stderr). If `create_offer` does not return `status="ok"`, `add_offer` is not sent.

## Example

```json
{
  "offer": "welcome-100",
  "game": "bgaming_ufo_pyramids",
  "spins": 10,
  "totalbet": 1000,
  "currency": "EUR",
  "wlid": "183019832"
}
```
