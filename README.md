# Сервис назначения ревьюверов для Pull Requestов

Сервис для назначения ревьюверов на Pull Requestы внутри команды.

## Документация

- Условия тестового задания: [docs/task.md](docs/task.md)
- Описание схемы БД: [docs/db.md](docs/db.md)
- Описание переменных окружения: [docs/env.md](docs/env.md)
- Принятые допущения: [docs/assumptions.md](docs/assumptions.md)
- HTTP API (OpenAPI/Swagger): [docs/openapi.yml](docs/openapi.yml)

## Подготовка

Скопируйте файл шаблон `.env.example` в `.env` и настройте под свою среду:

```bash
cp .env.example .env
```

## Запуск

### Через Docker

Поднятие контейнеров:

```bash
make docker-up
```

Остановка контейнеров:

```bash
make docker-down
```

### Локальный запуск

Сборка:

```bash
make build
```

Запуск:

```bash
make run
```

Сервис по умолчанию доступен на порту `8080`.
