# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Balance Service

## What this service does
Сервис управления балансами пользователей.
Append-only хранилище событий с периодическими снапшотами.
Текущий баланс кэшируется в Redis и пересчитывается как
сумма событий от последнего снапшота.
Спроектирован для работы в нескольких инстансах одновременно.

## API — Users
POST   /users              — создать пользователя
GET    /users/:id          — получить пользователя
PUT    /users/:id          — обновить пользователя
DELETE /users/:id          — удалить пользователя

Поля: id, name, overdraft_limit, created_at, updated_at, deleted_at.

## API — Balance
GET    /users/:id/balance            — текущий баланс
POST   /users/:id/balance/credit     — пополнение
POST   /users/:id/balance/debit      — списание
POST   /users/:id/balance/freeze     — заморозка части суммы
POST   /users/:id/balance/unfreeze   — разморозка

## API — Transactions
GET    /users/:id/transactions           — список транзакций
GET    /users/:id/transactions/:tx_id    — транзакция по ID
Query params: limit, offset, from, to

## Business rules
- Баланс может уйти в минус только в пределах overdraft_limit
- Заморозка блокирует конкретную сумму — остаток доступен
- Заморозка снимается по запросу или по timeout (передаётся при создании)
- Все операции идемпотентны — клиент передаёт transaction_id
- Повторный запрос с тем же transaction_id возвращает результат первого вызова
- События не изменяются и не удаляются никогда

## Data model
- events      — все изменения баланса, append-only
- snapshots   — снапшоты баланса на момент создания
- users       — пользователи с overdraft_limit
- Текущий баланс = последний снапшот + сумма событий после него

## Snapshot strategy
Настраивается через конфиг:
- По количеству транзакций: SNAPSHOT_EVERY_N
- По расписанию: SNAPSHOT_CRON

## Concurrency
Конкурентные операции на один баланс через optimistic locking
(version field на снапшоте).
При конфликте — retry на уровне use case, максимум 3 попытки.

## Configuration
Все параметры через environment variables.
Обязательные: DATABASE_URL, REDIS_URL, PORT.
Опциональные: SNAPSHOT_EVERY_N, SNAPSHOT_CRON.

## Tech stack
- Go 1.25+
- PostgreSQL 15
- Redis 7
- gRPC + HTTP через clay
- Docker

## Validation gates — перед каждым коммитом
Три проверки в строгом порядке. Если любая падает — остановись и исправь:
1. go build ./...
2. make lint
3. go test ./...

## Git workflow
Коммиты после каждого завершённого этапа разработки.
Формат: [FEAT/FIX/REFACTOR/TEST/DOCS/CHORE](scope): Description
Незавершённый код и код с ошибками линтера не коммитить.

## What NOT to do
- Float для денежных значений — только shopspring/decimal
- UPDATE или DELETE в таблице events
- Бизнес-логика в delivery layer
- Прямые SQL запросы минуя repository layer
- Глобальное состояние
- Хардкод любых параметров конфигурации
- Комментарии в коде кроме godoc для публичных API
- Circular dependencies между слоями
