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

The installer detects your OS/arch, downloads the matching archive from
[GitHub Releases](https://github.com/dkryvak/apikit/releases), **verifies its SHA-256 against the release's
`checksums.txt`**, and puts `apikit` on PATH (on macOS it also clears the Gatekeeper quarantine flag). Verify:

```sh
apikit --version
```

#### Options

Both scripts take the same settings via a flag or an environment variable:

| Setting | macOS / Linux | Windows | Default |
|---|---|---|---|
| Version | `-v v0.2.0` / `APIKIT_VERSION` | `-Version v0.2.0` / `$env:APIKIT_VERSION` | latest release |
| Install dir | `-d DIR` / `APIKIT_INSTALL_DIR` | `-InstallDir DIR` / `$env:APIKIT_INSTALL_DIR` | `/usr/local/bin` (Unix) · `%LOCALAPPDATA%\apikit\bin` (Windows) |
| Skip PATH edit | `--no-modify-path` / `APIKIT_NO_MODIFY_PATH` | `-NoModifyPath` / `$env:APIKIT_NO_MODIFY_PATH` | off (PATH is updated when needed) |

The default install dir differs by platform on purpose. On Unix it's
`/usr/local/bin` — already on `PATH`, so no shell config is touched (the script
uses `sudo` if the dir isn't writable). On Windows there's no equivalent
system dir on `PATH` by default, so the installer uses a per-user dir
(`%LOCALAPPDATA%\apikit\bin`) and registers it in the user `PATH`.

When you point `-d` at a directory that isn't on `PATH` (e.g. `~/.local/bin`),
the Unix installer appends an idempotent `export PATH=...` to the rc file for
your login shell (`~/.zshrc`, `~/.bash_profile` + `~/.bashrc`, or `~/.profile`)
and prints a `source` hint for the current session. Pass `--no-modify-path` to
skip this and set `PATH` yourself.

Pin a version and/or directory — download the script first when passing flags:

```sh
curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh -o install.sh
sh install.sh -v v0.2.0 -d ~/bin
```
```powershell
irm https://raw.githubusercontent.com/dkryvak/apikit/main/install.ps1 -OutFile install.ps1
.\install.ps1 -Version v0.2.0 -InstallDir C:\tools\apikit
```

…or pass the version through the pipe with an environment variable:

```sh
APIKIT_VERSION=v0.2.0 sh -c "$(curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh)"
```

### Manual install

Download the archive for your platform from the
[Releases page](https://github.com/dkryvak/apikit/releases) — assets are named
`apikit_<version>_<os>_<arch>.tar.gz` (`.zip` for Windows) — then verify the checksum and extract:

```sh
# verify against checksums.txt from the same release
shasum -a 256 apikit_0.2.0_darwin_arm64.tar.gz
grep apikit_0.2.0_darwin_arm64.tar.gz checksums.txt

tar -xzf apikit_0.2.0_darwin_arm64.tar.gz
sudo install -m 0755 apikit /usr/local/bin/
```

On Windows, unzip the archive and move `apikit.exe` to a directory on your `PATH`.

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

Ready-to-use config and request examples for every module live in
[`example/`](example/README.md) — pipe them straight in, e.g.
`apikit redgenn config import demo - < example/redgenn/config/example.json`.

## Commands

```
apikit env       create | import | list | show | delete
apikit <module>  config create | import | list | show | delete
apikit kube  <module> call <group> <endpoint> --config <alias> [--body - | --out f]
apikit local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

Configs and envs are stored under `~/.apikit/`.

## Documentation

- [`example/`](example/README.md) — import-ready config and request examples for every module.
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
