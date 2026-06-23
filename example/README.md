# Examples

Ready-to-use, valid JSON for every config and request body, so you can get a
working call without reverse-engineering the structs. Each file is import-ready:
pipe it straight into the matching command (the CLI reads `-` from STDIN).

```
example/
  env/example.json                 # a kube target (used by `apikit kube ...`)
  redgenn/
    config/example.json            # `apikit redgenn config import`
    requests/*.json                # `--body` payloads for redgenn calls
  vegangster/
    config/example.json
    requests/*.json                # `--body` payloads for vegangster calls
  softgamings/
    config/example.json
    requests/*.json                # `--query` params for softgamings calls
```

## Two different "env"s — don't confuse them

- **Build-time `.env`** (repo root, see `.env.example`) configures the *binary
  build* (pod image tag, Kube Job timeouts). Values are baked in at `task build`
  time. It has nothing to do with the files here.
- **Runtime `apikit env`** (`example/env/example.json`) describes a *Kubernetes
  target* (`namespace`, `context`) that `apikit kube ...` runs calls through.
  `apikit local ...` ignores it entirely.

## Credentials are injected from config

Request examples contain **only the fields you supply**. Fixed/credential fields
(`login`, `password`, `cm`, `mode`, `wl_code`, `token`, `operator_id`,
`brand_id`, softgamings `TID`/`Hash`, …) are injected by the CLI from your config
at call time, so they are intentionally absent from the request files. The same
content is shown by each endpoint's `--body-schema` / `--query-schema`.

## Quick start

1. (Kube mode only) Import an env — its name comes from the command argument,
   the JSON supplies `namespace` and `context`:

   ```sh
   apikit env import demo - < example/env/example.json
   ```

2. Import a config under an alias (here `demo`):

   ```sh
   apikit redgenn config import demo - < example/redgenn/config/example.json
   ```

   Replace the placeholder values with real credentials. For **vegangster**,
   swap in real PEM keys for `privateKey` / `publicKey` before making a call —
   they are parsed lazily when the config is used, not at import time.

3. Make a call, feeding a request file on STDIN. Use `local` to run from this
   machine (ignores env) or `kube` to run through the config's bound env:

   ```sh
   # redgenn / vegangster — payload goes to --body
   apikit local redgenn call game list --config demo --body - < example/redgenn/requests/game_list.json

   # softgamings — params go to --query
   apikit local softgamings call game launch --config demo --query - < example/softgamings/requests/game_launch.json
   ```

   `--body "$(cat <file>)"` works too if you prefer not to use STDIN.

## Which flag per module

| Module | Call payload flag | Schema flag |
|---|---|---|
| redgenn | `--body` | `--body-schema` |
| vegangster | `--body` | `--body-schema` |
| softgamings | `--query` | `--query-schema` |

## Endpoint → file map

| Module | Endpoint command | Example file |
|---|---|---|
| redgenn | `game list` | `redgenn/requests/game_list.json` |
| redgenn | `game launch` | `redgenn/requests/game_launch.json` |
| redgenn | `game demo-launch` | `redgenn/requests/game_demo_launch.json` |
| redgenn | `freegames bet-levels` | `redgenn/requests/freegames_bet_levels.json` |
| vegangster | `game launch` | `vegangster/requests/game_launch.json` |
| vegangster | `game demo-launch` | `vegangster/requests/game_demo_launch.json` |
| vegangster | `freegames grant` | `vegangster/requests/freegames_grant.json` |
| vegangster | `freegames cancel` | `vegangster/requests/freegames_cancel.json` |
| vegangster | `freegames bet-amounts` | `vegangster/requests/freegames_bet_amounts.json` |
| vegangster | `freegames bet-amount-list` | `vegangster/requests/freegames_bet_amount_list.json` |
| softgamings | `game launch` | `softgamings/requests/game_launch.json` |
| softgamings | `game demo-launch` | `softgamings/requests/game_demo_launch.json` |
| softgamings | `user get` | `softgamings/requests/user_get.json` |
| softgamings | `freegames grant` | `softgamings/requests/freegames_grant.json` |
| softgamings | `freegames cancel` | `softgamings/requests/freegames_cancel.json` |

> `redgenn`/`vegangster` `game list` take no user-supplied body (parameters are
> built from config), so they have no request file.
