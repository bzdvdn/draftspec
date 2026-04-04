## Draftspec

Основной контекст проекта — `.draftspec/`. Языки: docs=[DOCS_LANGUAGE], agent=[AGENT_LANGUAGE], comments=[COMMENTS_LANGUAGE]

Цепочка workflow: `constitution → spec → inspect → plan → tasks → implement → verify → archive`
- `/draftspec.constitution`: создать или обновить `.draftspec/constitution.md`
- `/draftspec.spec`: создать или уточнить `.draftspec/specs/<slug>/spec.md`; `--amend` для точечных правок
- `/draftspec.inspect`: проверить одну фичу на согласованность и качество
- `/draftspec.plan`: создать или обновить `.draftspec/plans/<slug>/`; `--update` для точечных правок, `--research` для research-first
- `/draftspec.tasks`: создать или обновить `.draftspec/plans/<slug>/tasks.md`
- `/draftspec.implement`: выполнить незавершённые задачи из `tasks.md`
- `/draftspec.verify`: проверить один feature package; `--deep` для полной per-AC валидации по коду
- `/draftspec.archive`: архивировать в `.draftspec/archive/` (move-based); `--copy` оставляет оригиналы, `--restore` восстанавливает

Опциональные (в любой момент): `/draftspec.challenge` (адверсариальная проверка; `--spec`/`--plan`), `/draftspec.handoff` (передача сессии), `/draftspec.hotfix` (экстренное исправление ≤ 3 файлов), `/draftspec.scope` (проверка границ; `--plan`/`--tasks`), `/draftspec.recap` (обзор проекта)

Дисциплина чтения:
- Не пропускайте фазы; по умолчанию загружайте только текущий feature slug
- Предпочитайте readiness scripts перед чтением более глубоких артефактов; для CLI используйте `./.draftspec/scripts/run-draftspec.sh`
- Не загружайте: нерелевантные specs/plans, широкие сканы репозитория, исходники scripts, файлы уже прочитанные в сессии (если сами не редактировали)
- Используйте настроенный язык комментариев для нового/изменяемого кода; сохраняйте существующие соглашения файла

Перед значимыми изменениями: просмотрите `constitution.md`, релевантный `specs/<slug>/spec.md` и `plans/<slug>/` если есть. После изменений: поддерживайте согласованность specs, plans, tasks и реализации.
