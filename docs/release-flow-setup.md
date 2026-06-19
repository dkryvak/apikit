# Release Flow: Tag Flow with main + hotfix

Summary of the configured git flow and CI/CD for the repository.

## 1. Overall model (Tag Flow)

**Operational branches** (CI jobs and rulesets key off these): `main`, `hotfix/*`

**Working branches** (merged into `main` via PR): `feature/*`, `bugfix/*`, `chore/*`, `refactor/*`, `docs/*`, `spike/*`

**Deploy triggers**:
- `push` to `main` → deploy to **stage**
- `push` of a tag (`v*`) → deploy to **prod** (regardless of which branch the tag technically points to)

```
main      ──●──●──●──●──●──●── (every merge → stage)
                        │
                      tag v1.4.0 ──→ deploy prod
```

## 2. Naming

| Entity | Format | Example |
|---|---|---|
| Working branch | `<type>/<description>` (free-form name) | `feature/add-endpoint` |
| Working branch types | `feature`, `bugfix`, `chore`, `refactor`, `docs`, `spike` | `bugfix/config-validation` |
| Hotfix branch | `hotfix/x.y.z` (version it prepares) | `hotfix/1.4.1` |
| Tag | SemVer `vX.Y.Z` | `v1.4.0`, `v1.4.1` |

## 3. Normal release flow

```bash
git checkout main
git pull
# the tag is created ONLY via the CI workflow (see §5), never by hand
```

## 4. Hotfix flow (detailed)

Principle: **the fix first goes through a normal PR into `main`** (review, CI), and only then is the
**finished squash commit cherry-picked** into the hotfix branch. This is safer than "hotfix first,
forward-port later" — if you forget the step, nothing gets deployed (an explicit failure) rather than the
bug silently coming back later (a silent failure).

```bash
# 1. Fix via a normal PR into main (squash merge)
git checkout main
git checkout -b bugfix/wallet-rollback
git commit -m "fix: wallet balance rollback on timeout"
# PR → squash merge into main

# 2. Cherry-pick the squash commit into the hotfix branch
git checkout -b hotfix/1.4.1 v1.4.0   # if it doesn't exist yet
git cherry-pick <sha-of-squash-merge-in-main>
git push origin hotfix/1.4.1
# then the tag is created via the CI workflow (§5), from the hotfix/1.4.1 branch
```

```
main         ──●──●──●──[fix]──●──●──
                          │
hotfix/1.4.1    ●─────────┴──●
                tag: v1.4.1   tag: v1.4.2
```

If another patch is needed on the same prod version — continue the same hotfix branch, repeating the cherry-pick.

## 5. CI/CD (GitHub Actions)

### 5.1 `ci.yml` — build + tests

Triggered on PRs to `main` and on `push` to `main` / `hotfix/**`. The job is named `build`
(required by the `main` ruleset). On a tag it also publishes the GHCR image via a `docker` job.

```yaml
name: CI

on:
  pull_request:
    branches: [main]
  push:
    branches:
      - main
      - "hotfix/**"

permissions:
  contents: read

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.26"

      - name: gofmt
        run: |
          out="$(gofmt -l . | grep -v '^vendor/' || true)"
          if [ -n "$out" ]; then echo "Unformatted files:"; echo "$out"; exit 1; fi

      - name: go vet
        run: go vet ./...

      - name: go test
        run: go test ./...

      - name: go build
        run: go build ./...
```

### 5.2 `create-release-tag.yml` — the only legitimate way to create a tag

Manual trigger (`workflow_dispatch`), verifies the run is from `main` or `hotfix/*`, and pushes the tag via a
**deploy key** (not `GITHUB_TOKEN` — that one is not supported in the ruleset bypass list).

```yaml
name: Create Release Tag

on:
  workflow_dispatch:
    inputs:
      version:
        description: 'Tag version (e.g. v1.4.0)'
        required: true
        type: string

permissions:
  contents: write

jobs:
  create-tag:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Verify branch is allowed
        run: |
          REF_NAME="${{ github.ref_name }}"
          if [[ "$REF_NAME" == "main" ]]; then
            echo "OK: main"
          elif [[ "$REF_NAME" =~ ^hotfix/ ]]; then
            echo "OK: hotfix branch ($REF_NAME)"
          else
            echo "::error::Tags may only be created from main or hotfix/*. Current branch: $REF_NAME"
            exit 1
          fi

      - name: Validate version format
        run: |
          VERSION="${{ inputs.version }}"
          if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "::error::Version must match vX.Y.Z, got: $VERSION"
            exit 1
          fi

      - name: Check tag doesn't already exist
        run: |
          VERSION="${{ inputs.version }}"
          if git rev-parse "$VERSION" >/dev/null 2>&1; then
            echo "::error::Tag $VERSION already exists"
            exit 1
          fi

      - name: Load deploy key
        uses: webfactory/ssh-agent@v0.9.0
        with:
          ssh-private-key: ${{ secrets.RELEASE_TAG_DEPLOY_KEY }}

      - name: Switch remote to SSH
        run: git remote set-url origin git@github.com:${{ github.repository }}.git

      - name: Create and push tag
        run: |
          git config user.name "release-bot"
          git config user.email "release-bot@users.noreply.github.com"
          git tag -a "${{ inputs.version }}" -m "Release ${{ inputs.version }} from ${{ github.ref_name }}"
          git push origin "${{ inputs.version }}"
```

**Deploy key setup** (one-time):
1. `ssh-keygen -t ed25519 -C "release-tag-bot" -f release_tag_key -N ""`
2. **Settings → Deploy keys → Add deploy key** — the public key, with **"Allow write access"** checked.
3. **Settings → Secrets and variables → Actions → New repository secret** — the private key, named `RELEASE_TAG_DEPLOY_KEY`.

⚠️ The bypass for deploy keys applies to the **category** "any deploy key with write access", not a specific
key. Keep only one write deploy key in the repo — the one used for tags.

### 5.3 GoReleaser — cross-platform binaries + GitHub Release

On a `v*` tag, `ci.yml` runs a `release` job (GoReleaser) **alongside** `docker`. Division of labor:

| Job | Produces |
|---|---|
| `docker` | The in-cluster pod image (`ghcr.io/dkryvak/apikit:<tag>` + `:latest`), from the `build` artifact |
| `release` | Downloadable binaries, `checksums.txt`, and the GitHub Release with an auto-generated changelog |

`release` is independent of `build` — GoReleaser does its own cross-compilation (it doesn't reuse the
linux/amd64 artifact). Config lives in [`.goreleaser.yaml`](../.goreleaser.yaml):

- **Platforms:** darwin (arm64+amd64), linux (arm64+amd64), windows (amd64).
- **ldflags:** the same build-time config as the `build` job — `config.version` = the tag, `config.image` =
  `ghcr.io/dkryvak/apikit:<tag>`, and `POD_*` from repo Variables (unset → the binary's built-in defaults).
- **Changelog:** generated from Conventional Commits since the previous tag (grouped into Features / Bug
  fixes / Performance; `docs`/`test`/`chore`/`ci`/`build`/`style`/`refactor` excluded).
- **Output:** per-platform archives (`tar.gz`, `zip` for Windows) + `checksums.txt`, attached to the Release.

```yaml
  release:
    if: startsWith(github.ref, 'refs/tags/v')
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: "1.26"
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          POD_LIFETIME_SECONDS: ${{ vars.POD_LIFETIME_SECONDS }}
          JOB_TTL_SECONDS: ${{ vars.JOB_TTL_SECONDS }}
          POD_READY_TIMEOUT_SECONDS: ${{ vars.POD_READY_TIMEOUT_SECONDS }}
```

The default `GITHUB_TOKEN` (with `contents: write`) is enough — no extra secret. Validate the config
locally before tagging with `task release-check`, and dry-run the full build with `task release-snapshot`
(output in `./dist`, nothing published).

## 6. GitHub Rulesets

### 6.1 Tag ruleset — `release-tags-protection`

| Setting | Value |
|---|---|
| Enforcement | Active |
| Target tags | `v*` |
| Bypass list | **Deploy keys** (Always allow) |
| Restrict creations | ✅ |
| Restrict updates | ✅ |
| Restrict deletions | ✅ |
| Block force pushes | ✅ |

Result: a `v*` tag can only be created via `create-release-tag.yml` (the only path with a write deploy key);
a direct `git push --tags` by a human is rejected.

### 6.2 Branch ruleset — `main`

| Setting | Value |
|---|---|
| Target branches | Include default branch |
| Bypass list | empty |
| Require a pull request before merging | ✅ (Required approvals: **0** — solo project, no review needed) |
| Allowed merge methods | **Squash** only |
| Require status checks to pass | ✅ → `build` |
| Require branches to be up to date before merging | ✅ |
| Restrict deletions | ✅ |
| Block force pushes | ✅ |

### 6.3 Branch ruleset — `hotfix/**`

Intentionally minimal — the main CI gate here comes from `ci.yml` (push trigger), not the ruleset.

| Setting | Value |
|---|---|
| Target branches | Include by pattern: `hotfix/**` |
| Bypass list | empty |
| Block force pushes | ✅ |
| Everything else | ❌ disabled (no PR flow in hotfix; the fix arrives via direct push/cherry-pick) |

## 7. Repo-level settings (Settings → General → Pull Requests)

- ✅ Allow squash merging
- ❌ Allow merge commits
- ❌ Allow rebase merging
- ✅ Automatically delete head branches

## 8. Tag protection summary

```
Person with Write/Maintain/Admin
        │
        ✗  git push origin v1.4.0  ──> REJECTED by the ruleset
        │
        ✓  Runs workflow_dispatch (create-release-tag.yml)
              │
        CI job: verifies the branch (main / hotfix/*)
              │ ✓
        CI job: pushes the tag via the deploy key
              │
        Ruleset: the deploy key is in the bypass list → push allowed
```
