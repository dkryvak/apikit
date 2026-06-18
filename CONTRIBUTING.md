# Contributing

The full git flow, CI/CD, rulesets and settings are described in
[`docs/release-flow-setup.md`](docs/release-flow-setup.md). A short cheat sheet below.

## Branching model (Tag Flow)

**Operational branches** — CI jobs and rulesets key off these:
- `main` — protected, always release-ready. No direct pushes (PR + squash only).
- `hotfix/x.y.z` — branch for a patch release on an older prod version (version-named, e.g. `hotfix/1.4.1`).

**Working branches** — `<type>/<description>` (the name is free-form; only the type prefix matters),
merged into `main` via PR (squash):
- `feature/` — new functionality (e.g. `feature/add-endpoint`)
- `bugfix/` — bug fix (e.g. `bugfix/config-validation`)
- `chore/` — maintenance (e.g. `chore/update-deps`)
- `refactor/` — refactoring (e.g. `refactor/http-client`)
- `docs/` — documentation (e.g. `docs/update-readme`)
- `spike/` — research / investigation (e.g. `spike/new-transport`)

Deploy triggers: push to `main` → stage; tag `v*` → prod.

## Commits — Conventional Commits

```
<type>[scope]: <imperative summary, ≤72 chars>
```

Types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `chore`, `ci`.
SemVer: `fix` → patch, `feat` → minor, `feat!:`/`BREAKING CHANGE:` → major.
The squash commit title must follow this format (it ends up in main's history).

## Pull Request

- PR title follows Conventional Commits; the description links the issue via `Closes #<issue>`.
- Merge is **squash** only; the branch must be up to date with `main`; `build` must be green.

## Release and hotfix

- **A tag is created ONLY via the `Create Release Tag` workflow** (`workflow_dispatch`) from `main` or
  `hotfix/*`. A manual `git push --tags` is blocked by the tag ruleset.
- **Hotfix**: the fix first goes through a normal PR into `main` (`bugfix/*`), then the squash commit is
  cherry-picked into `hotfix/x.y.z` → tag from that branch. Details — `docs/release-flow-setup.md` §4.

## Local development

```sh
task build          # binary in ./bin
go test ./...
go vet ./...
gofmt -l .          # empty output = OK
task docker:build   # pod image
```

Dependencies are vendored (`go mod vendor`). Go 1.26+.
