# Game: demo launch request body

Creates a demo (free-play) game session URL. No player authentication required.

## Fields

| Field | Required | Type | Description | Example |
|---|---|---|---|---|
| `game_code` | yes | string | Unique game identifier code. | `bgaming_lucky_streak` |
| `country` | no | string | Player country, ISO 3166-1 alpha-2. | `US` |
| `platform` | no | string | `desktop` or `mobile` browser. | `desktop` |
| `exit_url` | no | string | Redirect URL for the exit/lobby button. | `https://example.com/lobby` |
| `lang` | no | string | Interface language, ISO 639-1. | `en` |

> `mode`, `wl_code` and `token` are required by the API but injected automatically.

## Example

```json
{
  "game_code": "bgaming_lucky_streak",
  "country": "US",
  "platform": "desktop",
  "exit_url": "https://example.com/lobby",
  "lang": "en"
}
```
