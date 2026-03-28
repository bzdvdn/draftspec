# Roadmap

Этот roadmap сфокусирован на ближайших практических итерациях Draftspec, а не на длинном спекулятивном backlog.

## Направление

Draftspec должен и дальше занимать позицию между тяжелыми spec-driven системами и более свободными change-driven системами:

- строже, чем OpenSpec
- легче, чем spec-kit
- лучше приспособлен для agent-first workflow на реальных кодовых базах

## Итерация 1

### Главная цель

Усилить `inspect` как центральный слой качества.

### Планируемая работа

- определить и уточнить явную структуру inspection report и семантику verdict
- усилить проверки `constitution <-> spec`
- усилить проверки `spec <-> plan`
- усилить проверки `plan <-> tasks`
- усилить acceptance-to-task traceability

### Anti-Bloat Notes

Безопасное направление:

- более сильные structural checks
- более четкая семантика verdict
- лучшая traceability через стабильные acceptance IDs

Требует осторожности:

- делать persisted inspect reports обязательными для каждой фичи
- читать implementation code по умолчанию во время inspect
- превращать inspect в широкий review engine

### Почему это важно

Если `inspect` сильный, все downstream-фазы становятся качественнее при меньшем количестве пустой реализации и переделок.

## Итерация 2

### Главная цель

Добавить легкий post-implementation verification layer.

Статус: легкий contract, prompt, readiness script и report template уже есть. Дальше нужно усиливать проверки, не расширяя default context.

### Планируемая работа

- ввести небольшой `verify` или review-oriented workflow после `implement`
- проверять, что завершенные tasks соответствуют реальному состоянию implementation
- проверять, что implementation по-прежнему соответствует intent из spec и plan
- следить, чтобы memory и task state оставались синхронизированными

### Anti-Bloat Notes

Безопасное направление:

- helper-скрипты для проверки состояния задач
- проверки синхронизации memory/tasks
- optional persisted verify reports

Требует осторожности:

- читать код по умолчанию во время verify
- превращать verify в тяжелый review или QA engine
- требовать verify-артефакты перед каждым следующим действием

### Почему это важно

Это закрывает разрыв между "tasks выполнены" и "фича реально соответствует задуманному дизайну".

## Итерация 3

### Главная цель

Усилить brownfield ergonomics и machine-readable outputs.

### Планируемая работа

- улучшить archive summaries и связи архива с `memory.md`
- удерживать проверки completed-архива дешевыми за счет переиспользования task-state verification
- добавить machine-readable outputs вроде `doctor --json`
- улучшить config-aware поведение scripts и будущих утилит
- продолжить выравнивать многоязычную консистентность docs и prompts

### Anti-Bloat Notes

Безопасное направление:

- machine-readable outputs для уже существующих проверок
- более удобная индексация и summaries архива
- config-aware helpers, уменьшающие повторный reasoning

Требует осторожности:

- archive-flow, который требует чтения широкой истории репозитория
- новые automation outputs, создающие обязательные артефакты
- brownfield helpers, которые незаметно расширяют default context

### Почему это важно

Это сделает Draftspec удобнее для автоматизации, удобнее для эксплуатации в больших проектах и сильнее для долгоживущих brownfield-кодовых баз.

## Постоянная Работа Над Качеством

Параллельно с feature work Draftspec должен продолжать улучшать:

- консистентность документации
- unit test coverage
- ergonomics CLI
- ясность prompts и token discipline
- качество brownfield workflow

## Что Пока Не Планируется

Draftspec не стоит тащить в эти стороны без очень сильной причины:

- тяжелый orchestration engine
- обязательные checkpoint systems
- approval-gate бюрократию
- большие default prompt contexts
- обязательное разрастание артефактов для каждой фичи
