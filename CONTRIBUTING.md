# Contributing

Повний опис гіт-флоу, CI/CD, rulesets і налаштувань — у [`docs/release-flow-setup.md`](docs/release-flow-setup.md).
Нижче — стисла шпаргалка.

## Модель гілок (Tag Flow)

- `main` — захищений, завжди реліз-придатний. Прямий push заборонено (тільки через PR, squash).
- `feature/TICKET-опис` — нова функціональність (напр. `feature/CASINO-456-bulk-update`).
- `fix/TICKET-опис` — виправлення (теж через PR у `main`).
- `hotfix/x.y.z` — гілка для патч-релізу на стару прод-версію (cherry-pick готового squash-коміту).
- Без `develop`, `release/*`, `bugfix/*`.

Деплой-тригери: push у `main` → stage; тег `v*` → prod.

## Коміти — Conventional Commits

```
<type>[scope]: <опис у наказовій формі, ≤72 симв.>
```

Типи: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `chore`, `ci`.
SemVer: `fix` → patch, `feat` → minor, `feat!:`/`BREAKING CHANGE:` → major.
Заголовок squash-коміта має відповідати формату (він іде в історію main).

## Pull Request

- Назва PR — за Conventional Commits; в описі: що/навіщо, `Closes #<issue>`.
- Merge — лише **squash**; гілка має бути актуальною щодо `main`; `build-and-test` зелений.

## Реліз і хотфікс

- **Тег створюється ТІЛЬКИ через workflow `Create Release Tag`** (`workflow_dispatch`) з `main` або `hotfix/*`.
  Ручний `git push --tags` блокується tag-ruleset'ом.
- **Хотфікс**: фікс спочатку йде звичайним PR у `main` (`fix/*`), потім cherry-pick squash-коміту в
  `hotfix/x.y.z` → тег із цієї гілки. Деталі — `docs/release-flow-setup.md` §4.

## Локальна розробка

```sh
task build          # бінар у ./bin
go test ./...
go vet ./...
gofmt -l .          # порожній вивід = ок
task docker:build   # образ пода
```

Залежності завендорені (`go mod vendor`). Go 1.26+.
