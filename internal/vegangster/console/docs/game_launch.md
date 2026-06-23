# Game: launch request body

Creates a real-money game session URL for a player. Requires a player authentication token.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `player_id` | yes | string | Unique player identifier in the operator's system. | `player123` |
| `token` | yes | string | Player authentication token issued by the operator. | `eyJhbGciOiJSUzI1NiJ9...` |
| `game_code` | yes | string | Unique game identifier code. | `bgaming_lucky_streak` |
| `platform` | yes | string | `desktop` or `mobile` browser. | `desktop` |
| `currency` | yes | string | Player currency, ISO 4217. | `USD` |
| `lang` | yes | string | Interface language, ISO 639-1. | `en` |
| `country` | yes | string | Player country, ISO 3166-1 alpha-2. | `US` |
| `ip` | yes | string | Player IP address. | `192.168.1.1` |
| `lobby_url` | no | string | Redirect URL for the lobby/exit button. | `https://example.com/lobby` |
| `deposit_url` | no | string | Redirect URL for the deposit button. | `https://example.com/deposit` |
| `player_nick` | no | string | Player display nickname shown in the game. | `Lucky Player` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "player_id": "player123",
  "token": "eyJhbGciOiJSUzI1NiJ9...",
  "game_code": "bgaming_lucky_streak",
  "platform": "desktop",
  "currency": "USD",
  "lang": "en",
  "country": "US",
  "ip": "192.168.1.1",
  "lobby_url": "https://example.com/lobby",
  "deposit_url": "https://example.com/deposit",
  "player_nick": "Lucky Player"
}
```
