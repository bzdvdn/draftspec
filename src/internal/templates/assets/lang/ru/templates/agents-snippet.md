## Draftspec

Основной контекст проекта хранится в `.draftspec/`.

Предпочтительные языки:
- Документация: [DOCS_LANGUAGE]
- Агент: [AGENT_LANGUAGE]
- Комментарии в коде: [COMMENTS_LANGUAGE]

Workflow-команды:
- `/draftspec.constitution`: patch-обновить `.draftspec/constitution.md`
- `/draftspec.spec`: создать или уточнить один файл `.draftspec/specs/<slug>.md` и работать из `feature/<slug>`
- `/draftspec.inspect`: проанализировать одну фичу на согласованность и качество до или после планирования
- `/draftspec.plan`: создать или обновить `.draftspec/plans/<slug>/plan.md`, `data-model.md` и `contracts/`
- `/draftspec.tasks`: создать или обновить `.draftspec/plans/<slug>/tasks.md`
- `/draftspec.implement`: выполнить незавершенные задачи
- `/draftspec.archive`: архивировать один feature package в `.draftspec/archive/`

Дисциплина чтения:
- Следуйте цепочке `constitution -> spec -> inspect -> plan -> tasks -> implement -> archive`
- Не пропускайте prerequisites
- По умолчанию загружайте только текущий feature slug
- Предпочитайте readiness scripts каждой фазы перед чтением более глубоких артефактов
- `/draftspec.spec` поддерживает `--name`, optional `--slug` и optional `--branch`; для chat-based ввода описание фичи может прийти следующим сообщением
- Для file-based входа в `/draftspec.spec` предпочитайте `name:` и опциональный `slug:` в начале файла, а не fallback на filename
- Разрешайте явный `--branch <name>` override для repository-specific branch naming conventions, например Jira keys
- В `tasks` начинайте с `plan.md` и грузите более глубокие артефакты только при необходимости
- В `implement` начинайте с `tasks.md` и грузите более глубокие артефакты только при необходимости

Никогда не загружайте по умолчанию:
- нерелевантные спецификации или plan packages
- широкие сканы репозитория
- исходники scripts (используйте readiness scripts)

Дисциплина языка реализации:
- Считайте настроенный язык комментариев в коде основным для новых или изменяемых комментариев
- Сохраняйте устойчивое локальное соглашение файла, если редактируете комментарии в уже существующем коде
- Не смешивайте языки комментариев в одной локальной области кода без сильной причины

Перед значимыми изменениями:
- Просмотреть `.draftspec/constitution.md`
- Изучить релевантный `.draftspec/specs/<slug>.md`
- Изучить релевантный feature package в `.draftspec/plans/<slug>/`, если он есть

После значимых решений или изменений:
- Поддерживать согласованность спецификаций, планов, задач, archive state и реализации
