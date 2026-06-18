# apikit — план налаштування репозиторію та процесів

> ⚠️ **ОНОВЛЕНО:** фактичний гіт-флоу, реліз і CI/CD — **Tag Flow** (`main`+`feature/*`+`hotfix/*`, теги
> лише через workflow з deploy-key, rulesets) — описані в **[release-flow-setup.md](./release-flow-setup.md)**
> і **превалюють** над секцією «GitHub Flow» нижче. GoReleaser/Homebrew/Scoop/install-скрипти — **відкладено**.
> Цей документ лишається для решти (Issues/Projects, README/докси).

> Статус: узгоджені рішення, реалізація поетапна · Дата: 2026-06-17
> Рішення: репо `apikit`; GitHub Flow + теги; GitHub Issues + Projects; Docker-образ авто на тег.

---

## 0. Рішення (зафіксовані)

| Тема | Рішення |
|------|---------|
| Назва репо | `apikit` (публічний GitHub, `dkryvak/apikit`) |
| Гілки | **GitHub Flow** + SemVer-теги |
| Версії/коміти | SemVer + Conventional Commits → авто-changelog |
| Задачі | GitHub Issues + шаблони + **Projects** (Kanban) |
| Дистрибуція | GitHub Releases + install-скрипти (без brew/scoop) |
| Docker | авто build+push на релізний тег |

> Розробка нині в Gitea (`gitea.homelab/dkryvak/apikit`). Для дистрибуції потрібен **публічний GitHub**
> `dkryvak/apikit` — homelab колеги ззовні не дістануть. Gitea можна лишити як дзеркало або перейти на GitHub.

---

## 1. Гілки (GitHub Flow)

```
main ─────●────────●────────●───────●──────►   (завжди реліз-придатний, захищений)
           \        \                /
   feat/env-import   fix/sso-check  (короткоживучі → PR → squash-merge)
```

- `main` — захищений, завжди зелений і реліз-придатний.
- Робочі гілки короткоживучі, з префіксом за типом:
  `feat/*`, `fix/*`, `chore/*`, `docs/*`, `refactor/*`, `ci/*`.
- Кожна зміна → **PR у main** → зелений CI + 1 рев'ю → **squash-merge**.
- **Реліз** = SemVer-тег `vX.Y.Z` на `main` (тригерить release + docker workflow).
- `release/x.y` — заводимо **лише за потреби** стабілізації/бекпорту хотфіксів у стару мінорну лінію.
- Хотфікс: `fix/*` → PR у `main` → патч-тег `vX.Y.(Z+1)`.

**Branch protection (main):** заборонити прямий push; вимагати PR, зелений `ci`, ≥1 approve,
лінійну історію (squash), актуальність гілки.

---

## 2. Коміти та версії

- **Conventional Commits**: `feat:`, `fix:`, `perf:`, `refactor:`, `docs:`, `test:`, `chore:`, `ci:`.
  Заголовок ≤ 72 симв.; `feat!:`/`BREAKING CHANGE:` → major.
- **SemVer**: `fix:`→patch, `feat:`→minor, breaking→major.
- Changelog генерує GoReleaser із заголовків PR/комітів (секція `changelog` у `.goreleaser.yaml`).
- Squash-merge: заголовок squash-коміта має відповідати Conventional Commits (його й бачить changelog).

---

## 3. Задачі та PR

- **Issues**: шаблони `bug_report.md`, `feature_request.md` (у `.github/ISSUE_TEMPLATE/`).
- **PR**: `.github/pull_request_template.md` (опис, тип зміни, чек-лист, лінк на issue `Closes #N`).
- **Лейбли**: `type: bug|feature|chore|docs`, `priority: p0|p1|p2`, `area: kube|config|ci|module/*`,
  `status: blocked|in-review`.
- **CODEOWNERS**: `* @dkryvak` (авто-реквест рев'ю).
- **Projects (Kanban)**: колонки `Backlog → Todo → In Progress → In Review → Done`;
  авто-додавання issue/PR, авто-перехід In Review при відкритті PR, Done при merge.

---

## 4. CI/CD (GitHub Actions)

| Workflow | Тригер | Кроки |
|----------|--------|-------|
| `ci.yml` | PR + push до `main` | `golangci-lint`, `go vet`, `go test ./...`, `go build ./...` |
| `release.yml` | тег `v*` | setup-go 1.26 → `goreleaser release` (бінарі + GitHub Release) ✅ є |
| `docker.yml` | тег `v*` | build+push `dkryvak/apikit:{vX.Y.Z, latest}` (образ пода) |

- Секрети: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN` (для docker push). `GITHUB_TOKEN` — дефолтний.
- `.golangci.yml` — базовий набір лінтерів (govet, staticcheck, errcheck, ineffassign, gofmt).
- Docker-теги синхронні з релізом → у кластері завжди образ, що відповідає випущеній версії CLI.

---

## 5. README та докси

- **README.md** — що це / навіщо / встановлення (mac/win/manual) / пререквізити / quickstart / команди. ✅ чернетка є.
- **CONTRIBUTING.md** — модель гілок, Conventional Commits, як відкрити PR, локальний `task build/test`.
- **docs/** — design-доки (`remote-design.md`, `refactor-plan.md`, цей план).
- **LICENSE** — ⚠️ потребує рішення: MIT (відкрито) чи пропрієтарна/`All rights reserved` (внутрішній інструмент).
- **CHANGELOG** — автоматичний від GoReleaser (у Release notes); окремий файл за бажанням.

---

## 6. Поетапний rollout

**Фаза A — Фундамент репо**
- [ ] Створити публічний `github.com/dkryvak/apikit`; запушити код; (опц.) дзеркало з Gitea.
- [ ] Обрати й додати LICENSE; фіналізувати README; перевірити `.gitignore` (vendor/bin/output вже є).
- [ ] Закомітити вже готові `.goreleaser.yaml`, `.github/workflows/release.yml`, `install.sh`, `install.ps1`.

**Фаза B — Гілки й захист**
- [ ] `main` як default; branch protection (PR + CI + 1 review + squash).
- [ ] `CONTRIBUTING.md` (гілки, Conventional Commits).

**Фаза C — Issues / PR / Projects**
- [ ] Issue/PR-шаблони, лейбли, `CODEOWNERS`.
- [ ] Projects-дошка з автоматизаціями.

**Фаза D — CI**
- [ ] `.golangci.yml` + `ci.yml` (lint/vet/test/build).

**Фаза E — Реліз CD**
- [ ] Перший тег `v0.1.0`; перевірити, що GoReleaser зібрав асети; протестувати install-скрипти на mac/win.

**Фаза F — Docker CD**
- [ ] `docker.yml` (build+push на тег) + секрети DockerHub; звірити образ пода з релізом.

**Фаза G — Поліш**
- [ ] Бейджі в README, troubleshooting, перевірка end-to-end сценарію встановлення колегою.

---

## 7. Рішення дозакриті (2026-06-17)

1. **LICENSE** — поки **без ліцензії** (додамо пізніше за потреби).
2. **Gitea** — **повний перехід на GitHub** (Gitea більше не використовується).
3. **Owner** — `dkryvak` (підтверджено).
