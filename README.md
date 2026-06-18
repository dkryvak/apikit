# apikit

CLI for signed HTTP requests to gaming aggregators (vegangster, softgamings, redgenn).
Aggregators accept requests only from the cluster's whitelisted IP, so `apikit` can run a call
**inside a Kubernetes pod** (`kube`) or **directly from the machine** (`local`).

## Installation

### macOS / Linux
```sh
curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh | sh
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/dkryvak/apikit/main/install.ps1 | iex
```

The script detects the OS/arch, downloads the binary from [GitHub Releases](https://github.com/dkryvak/apikit/releases),
and puts it on PATH (on macOS it also clears the Gatekeeper quarantine). For manual installation, download the
archive from the Releases page.

```sh
apikit --version
```

## Prerequisites

- **`local`** — nothing extra.
- **`kube`** — the `aws` CLI with a configured SSO profile and cluster access (kubeconfig).
  On the first call apikit checks the session and runs `aws sso login` if needed.
  Manage credentials differently? `APIKIT_SKIP_SSO=1` skips the check.

## Quick start

```sh
apikit env create                                          # (for kube) describe the cluster
apikit vegangster config create                            # config + binding to an env
apikit local vegangster call game list --config <alias> --out games.json   # direct
apikit kube  vegangster call game list --config <alias> --out games.json   # via the cluster
```

Body/query — as a string or from stdin (`--body - < body.json`). Output: stdout (`--out -` or no `--out`) or
a file (`--out f`); diagnostics go to stderr.

## Commands

```
apikit env       create | import | list | show | delete
apikit <module>  config create | import | list | show | delete
apikit kube  <module> call <group> <endpoint> --config <alias> [--body - | --out f]
apikit local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

Configs and envs are stored under `~/.apikit/`.

## Documentation

- [`docs/release-flow-setup.md`](docs/release-flow-setup.md) — git flow, releases, CI/CD.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — how to contribute.

## Development

```sh
task build          # build the binary into ./bin
task docker:build   # pod image dkryvak/apikit
```

Git flow, releases and CI/CD — see [`docs/release-flow-setup.md`](docs/release-flow-setup.md) and
[`CONTRIBUTING.md`](CONTRIBUTING.md). A release is a `vX.Y.Z` tag created via the **Create Release Tag**
workflow (not by hand).
