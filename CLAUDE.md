# CLAUDE.md — apikit

Контекст проєкту для Claude/Cowork. Стислий вступ; повний стан — у `docs/PROJECT_CONTEXT.md`.

## Що це
`apikit` — Go CLI для підписаних HTTP-запитів до ігрових агрегаторів (vegangster, softgamings, redgenn).
Агрегатори приймають запити лише з whitelisted IP кластера, тож CLI виконує виклик або **в Kubernetes-поді**
(`apikit kube …`), або **прямо з машини** (`apikit local …`).

## Ключова архітектура
- Один бінар; режим — явний неймспейс `kube`/`local`. Self-remoting = запуск `apikit local …` усередині пода.
- `env` і `config` — **спільні** (`~/.apikit/`). `local` потребує лише config; `kube` — і env, і config.
- Диспетч у `internal/kit/middleware`; транспорт — **client-go**; `aws` — пререквізит kube (EKS), `local` без нього.
- Job-pod: reuse + жорсткий ліміт 30 хв.

## Команди
```
apikit env       create|import|list|show|delete
apikit <module>  config create|import|list|show|delete
apikit kube|local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

## Документи (джерела істини)
- `docs/PROJECT_CONTEXT.md` — повний брифінг: стан, рішення, що далі.
- `docs/repo-setup-plan.md` — план репо/CI-CD/процесів.

## Стиль
Відповіді стисло й по суті, українською. Спершу узгодити дизайн, потім реалізовувати поетапно.
Стек користувача: Java, Spring Boot, Groovy, AutoTesting, DevOps.
Білд: `task build` (локально), `task docker:build` (образ пода). Go 1.26, vendored deps.
Гілки: **Tag Flow** (`main` + `feature/*` + `hotfix/*`); теги `v*` лише через workflow `Create Release Tag`;
коміти — Conventional Commits. Деталі — `docs/release-flow-setup.md`.
