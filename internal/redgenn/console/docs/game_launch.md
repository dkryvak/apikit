# Game: launch request body

Creates a real-money game session URL for a player.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `token` | yes | string | Player authentication token issued by the operator. | `eyJhbGciOiJSUzI1NiJ9...` |
| `game_code` | yes | string | Unique game identifier code. | `bgaming_lucky_streak` |
| `country` | no | string | Player country, ISO 3166-1 alpha-2. | `US` |
| `platform` | no | string | `desktop` or `mobile` browser. | `desktop` |
| `exit_url` | no | string | Redirect URL for the exit/lobby button. | `https://example.com/lobby` |
| `cashier_url` | no | string | Redirect URL for the cashier/deposit button. | `https://example.com/deposit` |
| `lang` | no | string | Interface language, ISO 639-1. | `en` |

> `mode` and `wl_code` are required by the API but injected from `--config`.

## Example

```json
{
  "token": "eyJhbGciOiJSUzI1NiJ9...",
  "game_code": "bgaming_lucky_streak",
  "country": "US",
  "platform": "desktop",
  "exit_url": "https://example.com/lobby",
  "cashier_url": "https://example.com/deposit",
  "lang": "en"
}
```
