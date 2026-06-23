# Game: demo launch request body

Creates a demo (free-play) game session URL. No player authentication required.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `game_code` | yes | string | Unique game identifier code. | `bgaming_lucky_streak` |
| `platform` | yes | string | `desktop` or `mobile` browser. | `desktop` |
| `currency` | yes | string | Currency used in demo mode, ISO 4217. | `USD` |
| `lang` | no | string | Interface language, ISO 639-1. | `en` |
| `country` | no | string | Player country, ISO 3166-1 alpha-2. | `US` |
| `ip` | no | string | Player IP address. | `192.168.1.1` |
| `lobby_url` | no | string | Redirect URL for the lobby/exit button. | `https://example.com/lobby` |

> `operator_id` and `brand_id` are required by the API but injected from `--config`.

## Example

```json
{
  "game_code": "bgaming_lucky_streak",
  "platform": "desktop",
  "currency": "USD",
  "lang": "en",
  "country": "US",
  "ip": "192.168.1.1",
  "lobby_url": "https://example.com/lobby"
}
```
