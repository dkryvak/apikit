# Release Flow: Tag Flow з main + hotfix

Підсумок налаштованого гіт-флоу та CI/CD для репозиторію.

## 1. Загальна модель (Tag Flow)

**Гілки**: `main`, `feature/*`, `hotfix/*`
**Без**: `develop`, `release/*`, `bugfix/*` (свідомо не використовуються)

**Тригери деплою**:
- `push` у `main` → деплой на **stage**
- `push` тегу (`v*`) → деплой на **prod** (незалежно, з якої гілки технічно стоїть тег)

```
main      ──●──●──●──●──●──●── (кожен мердж → stage)
                        │
                      tag v1.4.0 ──→ deploy prod
```

## 2. Нейминг

| Сутність | Формат | Приклад |
|---|---|---|
| Feature-гілка | `feature/TICKET-короткий-опис` | `feature/CASINO-456-bulk-update` |
| Hotfix-гілка | `hotfix/x.y.z` (версія, яку готує) | `hotfix/1.4.1` |
| Тег | SemVer `vX.Y.Z` | `v1.4.0`, `v1.4.1` |

## 3. Звичайний release-потік

```bash
git checkout main
git pull
# тег створюється ТІЛЬКИ через CI workflow (див. п.5), не руками
```

## 4. Hotfix-потік (детально)

Принцип: **фікс спочатку йде через звичайний PR у `main`** (review, CI), і лише потім **cherry-pick готового squash-коміту** в hotfix-гілку. Це безпечніше, ніж "hotfix спочатку, forward-port пізніше" — якщо забути крок, нічого не задеплоїться (явний фейл), а не тихо повернеться баг пізніше (тихий фейл).

```bash
# 1. Фікс через звичайний PR у main (squash merge)
git checkout main
git checkout -b fix/CASINO-789-wallet-rollback
git commit -m "fix: wallet balance rollback on timeout"
# PR → squash merge у main

# 2. Cherry-pick squash-коміту в hotfix-гілку
git checkout -b hotfix/1.4.1 v1.4.0   # якщо ще не існує
git cherry-pick <sha-зі-squash-merge-в-main>
git push origin hotfix/1.4.1
# далі тег створюється через CI workflow (п.5), з гілки hotfix/1.4.1
```

```
main         ──●──●──●──[fix]──●──●──
                          │
hotfix/1.4.1    ●─────────┴──●
                tag: v1.4.1   tag: v1.4.2
```

Якщо потрібен ще один патч на ту саму прод-версію — продовжуєш ту саму hotfix-гілку, повторюючи cherry-pick.

## 5. CI/CD (GitHub Actions)

### 5.1 `ci.yml` — build + тести на кожен коміт

Тригериться на `push` у `main` та `hotfix/**`.

```yaml
name: CI

on:
  push:
    branches:
      - main
      - 'hotfix/**'

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up JDK 21
        uses: actions/setup-java@v4
        with:
          java-version: '21'
          distribution: 'temurin'
          cache: 'gradle'

      - name: Build and run tests (Gradle)
        run: ./gradlew build

      - name: Publish test report
        if: always()
        uses: dorny/test-reporter@v1
        with:
          name: Test Results
          path: '**/build/test-results/test/*.xml'
          reporter: java-junit
```

### 5.2 `create-release-tag.yml` — єдиний легітимний спосіб створити тег

Manual trigger (`workflow_dispatch`), перевіряє, що запущено з `main` або `hotfix/*`, пушить тег через **deploy key** (а не `GITHUB_TOKEN` — той не підтримується ruleset bypass list).

```yaml
name: Create Release Tag

on:
  workflow_dispatch:
    inputs:
      version:
        description: 'Версія тегу (наприклад v1.4.0)'
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
            echo "::error::Тег можна створювати тільки з main або hotfix/*. Поточна гілка: $REF_NAME"
            exit 1
          fi

      - name: Validate version format
        run: |
          VERSION="${{ inputs.version }}"
          if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "::error::Версія має бути у форматі vX.Y.Z, отримано: $VERSION"
            exit 1
          fi

      - name: Check tag doesn't already exist
        run: |
          VERSION="${{ inputs.version }}"
          if git rev-parse "$VERSION" >/dev/null 2>&1; then
            echo "::error::Тег $VERSION вже існує"
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

**Налаштування deploy key** (одноразово):
1. `ssh-keygen -t ed25519 -C "release-tag-bot" -f release_tag_key -N ""`
2. **Settings → Deploy keys → Add deploy key** — публічний ключ, з галочкою **"Allow write access"**
3. **Settings → Secrets and variables → Actions → New repository secret** — приватний ключ під назвою `RELEASE_TAG_DEPLOY_KEY`

⚠️ Bypass для deploy keys діє на **категорію** "будь-який deploy key з write-доступом", не на конкретний ключ. Тримати в репо лише один write deploy key — той, що для тегів.

## 6. GitHub Rulesets

### 6.1 Tag ruleset — `release-tags-protection`

| Налаштування | Значення |
|---|---|
| Enforcement | Active |
| Target tags | `v*` |
| Bypass list | **Deploy keys** (Always allow) |
| Restrict creations | ✅ |
| Restrict updates | ✅ |
| Restrict deletions | ✅ |
| Block force pushes | ✅ |

Результат: створити тег `v*` можна тільки через `create-release-tag.yml` (єдиний шлях з write deploy key), пряме `git push --tags` людиною — відхиляється.

### 6.2 Branch ruleset — `main`

| Налаштування | Значення |
|---|---|
| Target branches | Include default branch |
| Bypass list | порожній |
| Require a pull request before merging | ✅ (Required approvals: **0** — соло-проєкт, review не потрібен) |
| Allowed merge methods | **Squash** тільки |
| Require status checks to pass | ✅ → `build-and-test` |
| Require branches to be up to date before merging | ✅ |
| Restrict deletions | ✅ |
| Block force pushes | ✅ |

### 6.3 Branch ruleset — `hotfix/**`

Мінімалістичний навмисно — основний CI-гейт тут дає `ci.yml` (push-тригер), а не ruleset.

| Налаштування | Значення |
|---|---|
| Target branches | Include by pattern: `hotfix/**` |
| Bypass list | порожній |
| Block force pushes | ✅ |
| Все інше | ❌ вимкнено (немає PR-флоу в hotfix, фікс заходить прямим push/cherry-pick) |

## 7. Repo-level налаштування (Settings → General → Pull Requests)

- ✅ Allow squash merging
- ❌ Allow merge commits
- ❌ Allow rebase merging
- ✅ Automatically delete head branches

## 8. Підсумкова схема захисту тегів

```
Людина з Write/Maintain/Admin
        │
        ✗  git push origin v1.4.0  ──> ВІДХИЛЕНО ruleset'ом
        │
        ✓  Запускає workflow_dispatch (create-release-tag.yml)
              │
        CI job: перевіряє гілку (main / hotfix/*)
              │ ✓
        CI job: пушить тег через deploy key
              │
        Ruleset: deploy key є в bypass list → push дозволено
```
