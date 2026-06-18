# apikit

CLI для підписаних HTTP-запитів до ігрових агрегаторів (vegangster, softgamings, redgenn).
Запити дозволені лише з whitelisted IP кластера, тож `apikit` уміє виконувати виклик **усередині
Kubernetes-пода** (режим `kube`) або **прямо з машини** (режим `local`).

## Встановлення

### macOS / Linux
```sh
curl -fsSL https://raw.githubusercontent.com/dkryvak/apikit/main/install.sh | sh
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/dkryvak/apikit/main/install.ps1 | iex
```

Скрипт визначає ОС/архітектуру, тягне відповідний бінар з [GitHub Releases](https://github.com/dkryvak/apikit/releases),
кладе його в PATH (на macOS ще знімає Gatekeeper-карантин).

**Ручне встановлення:** завантаж архів під свою платформу зі сторінки Releases, розпакуй і поклади
`apikit` у будь-яку теку з PATH.

Перевірка:
```sh
apikit --version
```

## Пререквізити

- **`local`-режим** — нічого додаткового не треба.
- **`kube`-режим** — потрібні:
  - `aws` CLI з налаштованим SSO-профілем (для входу в кластер EKS),
  - доступ до кластера (kubeconfig із потрібними контекстами).

  При першому `kube`-виклику apikit сам перевірить сесію і за потреби запустить `aws sso login`.
  Якщо `aws` не встановлено — буде зрозуміла підказка. Керуєш кредами інакше (не SSO)?
  Постав `APIKIT_SKIP_SSO=1`, щоб пропустити перевірку.

## Швидкий старт

```sh
# 1. (лише для kube) описати кластер
apikit env create                     # напр. name=stage → namespace stage-casino

# 2. створити конфіг і прив'язати до env
apikit vegangster config create       # у візарді вкажи env (напр. stage)

# 3a. виклик прямо з машини
apikit local vegangster call game list --config stage-winthrone --out games.json

# 3b. виклик через кластер (whitelisted IP)
apikit kube  vegangster call game list --config stage-winthrone --out games.json
```

Тіло/квері можна передати рядком або зі stdin:
```sh
apikit kube vegangster call game launch --config stage-winthrone --body - < body.json
```

Вивід відповіді йде в stdout (`--out -` або без `--out`) чи у файл (`--out res.json`);
діагностика запиту/відповіді — у stderr.

## Команди

```
apikit env       create | import | list | show | delete      # реєстр кластерів (для kube)
apikit <module>  config create | import | list | show | delete
apikit kube  <module> call <group> <endpoint> --config <alias> [--body - | --out f]
apikit local <module> call <group> <endpoint> --config <alias> [--body - | --out f]
```

Конфіги й env зберігаються в `~/.apikit/`.

## Реліз (для мейнтейнера)

Релізи збирає [GoReleaser](https://goreleaser.com) через GitHub Actions на push тега:
```sh
git tag v1.0.0
git push origin v1.0.0
```
Workflow збере бінарі під macOS/Linux/Windows (amd64 + arm64), створить GitHub Release з архівами
та `checksums.txt`. Інсталяційні скрипти беруть `releases/latest/download/apikit_<os>_<arch>.*`.

> Розробка ведеться в Gitea (`gitea.homelab`), а для дистрибуції потрібен **публічний GitHub-репозиторій**
> `dkryvak/apikit` (homelab-домен колеги ззовні не дістануть). Налаштуй push/дзеркалення коду й тегів у GitHub.
