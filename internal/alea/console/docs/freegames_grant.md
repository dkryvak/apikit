# Freegames: grant request body

Grants a free-spin bonus to a player (Alea `POST /api/v1/free-spins`, signed with the `Digest` header).

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `casinoPlayerId` | yes | string | Operator player id (same as Games Launch). | `183019832` |
| `casinoBonusId` | yes | string | Unique bonus-grant id (idempotency key). | `856543-603` |
| `country` | yes | string | Player country (ISO 3166). | `MT` |
| `currency` | yes | string | Player currency (ISO 4217). | `EUR` |
| `softwareId` | yes | int | Software (provider) id. | `1` |
| `level` | yes | int | Value per spin, 1–5. | `2` |
| `amount` | yes | int | Number of spins. | `10` |
| `games` | yes | array | Exactly one `{ "id": <int> }`. | `[{"id":1}]` |
| `casinoCampaignId` | no | string | Operator label for the offer. | `WelcomeBonus` |
| `games[].device` | no | string | `MOBILE` or `DESKTOP`. | `DESKTOP` |
| `startAt` | no | datetime | UTC `yyyy-mm-ddThh:mm:ssZ`; future, ≤30 days. | `2026-11-01T00:00:00Z` |
| `expireAt` | no | datetime | UTC `yyyy-mm-ddThh:mm:ssZ`; ≤350 days. | `2026-12-01T00:00:00Z` |

> `environmentToken`, `secretKey`, `casinoId` and `environmentId` come from `--config` (headers + Digest signature).

## Example

```json
{
  "casinoPlayerId": "183019832",
  "casinoBonusId": "856543-603",
  "casinoCampaignId": "WelcomeBonus",
  "country": "MT",
  "currency": "EUR",
  "softwareId": 1,
  "level": 2,
  "amount": 10,
  "games": [{ "id": 1 }]
}
```
