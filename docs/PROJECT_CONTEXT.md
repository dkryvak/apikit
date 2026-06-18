# apikit — Project Context / Handoff

> Самодостатній брифінг для перенесення контексту в іншу сесію/проєкт.
> Прочитай цей файл — і ти повністю в курсі стану, рішень і наступних кроків.
> Дата зрізу: 2026-06-17.

---

## 1. Що це і навіщо

`apikit` — Go CLI для **підписаних HTTP-запитів до ігрових агрегаторів** (vegangster, softgamings,
redgenn; alea поки осторонь). Ключове обмеження: агрегатори приймають запити **лише з whitelisted IP**,
який має тільки Kubernetes-кластер. Тож CLI уміє виконувати виклик **усередині пода** (`kube`) або
**прямо з машини** (`local`).

Початково було два репозиторії: Go-код (`tools/apikit`) і купа shell-скриптів-обгорток
(`denyskryvak/apikit`) для запуску в k8s. Мета рефактора — прибрати весь shell-шар і складність
(pod-per-request, kubectl exec/cp/gzip, ручне редагування тіл у `.sh`).

## 2. Фінальна архітектура (узгоджено й реалізовано)

- **Один бінар `apikit`** з self-remoting (без окремого бінара, без прапора `--remote`).
- **Режим — явний неймспейс команди**: `apikit kube …` (через кластер) vs `apikit local …` (прямо).
- **Self-remoting = запуск `apikit local …` усередині пода** → рекурсія неможлива by design,
  жодних guard-ів не треба (детект режиму через `os.Args[1]`).
- **env і config — спільні** (одне сховище `~/.apikit/`), не per-mode. Той самий config (напр.
  `stage-winthrone`) працює в обох режимах.
  - **env** (`~/.apikit/env/<name>.json`): `{name, namespace, context}` — таргетування кластера,
    **без `type`**, **без `image`**. CRUD-only (`apikit env create|import|list|show|delete`).
    Образ — фіксована константа `kube.Image` (бо має містити apikit).
  - **config** (`~/.apikit/<module>/config/<alias>.json`): credentials + опційне поле `env`
    (потрібне для kube, ігнорується в local).
- **Що важливо в якому режимі**: `local` — лише config; `kube` — і env, і config.
- **Диспетч у `kit/middleware`**: `kube.IsRemote() && !schema → kube.SelfRemote(...)`, інакше хендлер
  локально. Ендпойнт-хендлери НЕ змінювались.
- **kube-флоу**: config→env→кластер → ensure Job-pod (reuse) → `apikit <module> config import <alias> -`
  (stdin=config JSON) → `apikit local <module> call … --config <alias> --body - --out -` → ловимо stdout.
- **Транспорт — client-go** (без зовнішнього `kubectl`). `aws` лишається пререквізитом kube-режиму
  (EKS: `aws eks get-token` + `aws sso login`); `local` aws не торкається.
- **Pod lifecycle — неявний**: kube-call авто-піднімає/перевикористовує Job; жорсткий ліміт 30 хв
  (`sleep 1800` + `ActiveDeadlineSeconds 1800` + `TTLSecondsAfterFinished 120`).
- **Composability**: `--out -` → raw body у stdout; діагностика запиту/відповіді → stderr;
  усі launch/demo-ендпойнти тепер мають `--out` і пишуть body.

Команди:
```
apikit env       create|import|list|show|delete
apikit <module>  config create|import|list|show|delete
apikit kube  <module> call <group> <endpoint> --config <alias> [--body - | --out f]
apikit local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

## 3. Поточний стан реалізації

- **Фази 1–3 реалізовані й зібрані** (локально + Docker-образ). Self-remote **перевірено на реальному
  кластері — працює.**
- Додано client-go у залежності (`go get k8s.io/client-go/api/apimachinery`), go.mod пішов на `go 1.26.0`;
  Dockerfile піднято до `golang:1.26` + офлайн-збірка `-mod=vendor` (vendor закомічено в контекст, але
  у `.gitignore`).
- **AWS**: `kube.EnsureAWSSSO` перевіряє `aws` у PATH (зрозуміла помилка, якщо немає), підтримує
  `APIKIT_SKIP_SSO=1`; інтерактивний `aws sso login` при протермінованій сесії.
- **Дистрибуція (чернетки готові)**: GitHub Releases + install-скрипти (без brew/scoop).
  Файли: `.goreleaser.yaml`, `.github/workflows/release.yml`, `install.sh`, `install.ps1`, README,
  version-wiring (`main.version` через ldflags).

## 4. Репо та процеси (узгоджений план)

Деталі — `docs/repo-setup-plan.md`. Зафіксовані рішення:
- Репо: **`apikit`**, публічний **GitHub `dkryvak/apikit`**; **повний перехід з Gitea** (`gitea.homelab`).
- Гілки: **GitHub Flow** + SemVer-теги; `main` захищений; `feat/*`,`fix/*` → PR → squash.
- Коміти: **Conventional Commits** → авто-changelog (GoReleaser).
- Задачі: **GitHub Issues + Projects** (Kanban).
- CI/CD: `ci.yml` (lint/vet/test/build) + `release.yml` (тег→GoReleaser) + `docker.yml`
  (build+push образу на тег). Docker — **авто на тег**.
- LICENSE — поки немає (рішення відкладено).

## 5. Карта ключових файлів

```
cmd/apikit/main.go                 # entrypoint, var version
internal/app/app.go                # реєстрація env + <module> config + kube/local
internal/env/                      # сутність env (types/paths/store/console + command CRUD)
internal/kube/                     # mode.go, sso.go, executor.go (client-go), runner.go (SelfRemote)
internal/kit/middleware/           # диспетч local vs self-remote
internal/kit/config/               # Config interface (+ BoundEnv), store
internal/<module>/command/         # NewConfigCommand / NewCallCommand (call під kube/local)
internal/kit/file/writer.go        # --out - → stdout
internal/kit/http/console/         # діагностика → stderr
docs/remote-design.md              # АКТУАЛЬНА специфікація архітектури
docs/repo-setup-plan.md            # план репо/процесів
docs/refactor-plan.md              # історія рішень (проксі відхилено тощо)
.goreleaser.yaml / install.* / .github/workflows/release.yml
```

## 6. Що далі (не зроблено)

- **Repo rollout** (фази A–G з `repo-setup-plan.md`): створити GitHub-репо, push, branch protection,
  CONTRIBUTING, issue/PR-шаблони, CODEOWNERS, лейбли, Projects, `.golangci.yml`, `ci.yml`, `docker.yml`,
  перший тег `v0.1.0`, тест install-скриптів.
- **Фаза 5 рефактора** (з `remote-design.md`): міграція старих `denyskryvak/apikit/*/config/*.json` у
  `~/.apikit`, тіла запитів у `*/payloads/*.json`, депрекейт `k8s.sh` і per-endpoint скриптів,
  оновити `docs/overview.md`.
- **Рантайм-тести**: різні модулі/ендпойнти, body/query через stdin, reuse пода, великий вивід,
  помилки з пода, config без env.
- **Поза скоупом**: гігієна секретів (захардкоджений Alea JWT у `denyskryvak/apikit/alea/provider_list.sh`,
  приватні RSA-ключі в старих config JSON), модуль `alea`.

## 7. Стиль/преференції

Користувач (Denys): Java, Spring Boot, Groovy, AutoTesting, DevOps. Відповіді — **стисло й по суті**,
без зайвої води. Українською. Любить, щоб спершу узгодити дизайн, потім реалізовувати поетапно;
кожне рішення оформляти в design-док.

---

### Як використати цей контекст в іншому проєкті/сесії
1. Поклади цей файл у новий проєкт (або відкрий папку apikit у новій сесії — він уже в `docs/`).
2. На старті нової сесії скажи: «прочитай docs/PROJECT_CONTEXT.md» — і продовжуй з будь-якого пункту розділу 6.
