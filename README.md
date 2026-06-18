# apikit

CLI для підписаних HTTP-запитів до ігрових агрегаторів (vegangster, softgamings, redgenn).
Агрегатори приймають запити лише з whitelisted IP кластера, тож `apikit` уміє виконувати виклик
**усередині Kubernetes-пода** (`kube`) або **прямо з машини** (`local`).

## Встановлення

### macOS / Linux
```sh
curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh | sh
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/dkryvak/apikit/main/install.ps1 | iex
```

Скрипт визначає ОС/архітектуру, тягне бінар з [GitHub Releases](https://github.com/dkryvak/apikit/releases),
кладе в PATH (на macOS знімає Gatekeeper-карантин). Ручне встановлення — завантаж архів зі сторінки Releases.

```sh
apikit --version
```

## Пререквізити

- **`local`** — нічого додаткового.
- **`kube`** — `aws` CLI з налаштованим SSO-профілем та доступ до кластера (kubeconfig).
  При першому виклику apikit перевірить сесію і за потреби запустить `aws sso login`.
  Керуєш кредами інакше? `APIKIT_SKIP_SSO=1` пропускає перевірку.

## Швидкий старт

```sh
apikit env create                                          # (для kube) описати кластер
apikit vegangster config create                            # конфіг + прив'язка до env
apikit local vegangster call game list --config <alias> --out games.json   # прямо
apikit kube  vegangster call game list --config <alias> --out games.json   # через кластер
```

Тіло/квері — рядком або зі stdin (`--body - < body.json`). Вивід: stdout (`--out -` чи без `--out`) або
файл (`--out f`); діагностика — у stderr.

## Команди

```
apikit env       create | import | list | show | delete
apikit <module>  config create | import | list | show | delete
apikit kube  <module> call <group> <endpoint> --config <alias> [--body - | --out f]
apikit local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

Конфіги й env зберігаються в `~/.apikit/`.

## Документація

- [`docs/PROJECT_CONTEXT.md`](docs/PROJECT_CONTEXT.md) — повний контекст проєкту.
- [`docs/repo-setup-plan.md`](docs/repo-setup-plan.md) — процеси, гілки, CI/CD.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — як контриб'ютити.

## Розробка

```sh
task build          # зібрати бінар у ./bin
task docker:build   # образ пода dkryvak/apikit
```

Гіт-флоу, реліз і CI/CD — [`docs/release-flow-setup.md`](docs/release-flow-setup.md) та
[`CONTRIBUTING.md`](CONTRIBUTING.md). Реліз — тег `vX.Y.Z` через workflow **Create Release Tag** (не вручну).
