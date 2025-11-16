# PR Reviewer Assignment Service

Сервис для управления pull request и назначения ревьюверов.

## Архитектура

- **handlers** — HTTP уровень (Gin)
- **services** — бизнес-логика  
- **repository** — работа с БД (GORM)
- **models** — сущности данных

## Запуск

```bash
docker-compose up --build
```

Приложение будет доступно по адресу: http://localhost:8080

## API Endpoints

### Teams
- **POST /team/add** — Создать команду с участниками
- **GET /team/get** — Получить команду с участниками

### Users
- **POST /users/setIsActive** — Установить флаг активности пользователя
- **GET /users/getReview** — Получить PR'ы пользователя для ревью

### Pull Requests
- **POST /pullRequest/create** — Создать PR и назначить ревьюверов
- **POST /pullRequest/merge** — Пометить PR как MERGED
- **POST /pullRequest/reassign** — Переназначить ревьювера

### Stats
- **GET /stats/reviewers** - Статистика ревью по юзерам

## Примеры запросов

```bash
# Создание команды
curl -X POST http://localhost:8080/team/add -H "Content-Type: application/json" -d "{"team_name":"backend","members":[{"user_id":"u1","username":"Alice","is_active":true},{"user_id":"u2","username":"Bob","is_active":true},{"user_id":"u3","username":"Charlie","is_active":true},{"user_id":"u4","username":"David","is_active":true},{"user_id":"u5","username":"Eve","is_active":true}]}"

# Получение команды
curl -X GET "http://localhost:8080/team/get?team_name=backend"

# Создание PR
curl -X POST http://localhost:8080/pullRequest/create -H "Content-Type: application/json" -d "{"pull_request_id":"pr-1001","pull_request_name":"Add search feature","author_id":"u1"}"

# Получение PR для ревью
curl -X GET "http://localhost:8080/users/getReview?user_id=u2"

# Переназначение ревьювера
curl -X POST http://localhost:8080/pullRequest/reassign -H "Content-Type: application/json" -d "{"pull_request_id":"pr-1001","old_user_id":"u2"}"

# Изменение активности пользователя
curl -X POST http://localhost:8080/users/setIsActive -H "Content-Type: application/json" -d "{"user_id":"u4","is_active":false}"

# Мердж PR
curl -X POST http://localhost:8080/pullRequest/merge -H "Content-Type: application/json" -d "{"pull_request_id":"pr-1001"}"

# Статистика
curl http://localhost:8080/stats/reviewers
```

## Дополнительные задания
### E2E тесты
Тестовое окружение будет доступно по адресу: http://localhost:8081
Запустить тестовое окружение:
```bash
docker-compose -f docker-compose.test.yml up -d --build
```
Запустить тесты:
```bash
docker-compose -f docker-compose.test.yml exec tests go test -v ./e2e/api_test.go
```
## Линтер
Запустить линтер
```bash
golangci-lint run
```
### Нагрузочное тестирование
Запустить тестовое окружение:
```bash
docker-compose -f docker-compose.test.yml up -d --build
```
Запустить нагрузочное тестирование:
```bash
docker-compose -f docker-compose.test.yml exec loadtest k6 run /scripts/loadtest.js
```
Отчет: loadtest/loadtest_report.md

### Эндпоинт статистики
Возвращает количество открытых ревью и общее количество ревью за всё время по каждому пользователю.
**GET /stats/reviewers**
Пример ответа:
```bash
{
  "reviewer_stats": [
    {
      "open_reviews": 1,
      "total_reviews_ever": 2,
      "user_id": "u1"
    },
    {
      "open_reviews": 1,
      "total_reviews_ever": 4,
      "user_id": "u2"
    },
    ...
    }
  ]
}
```

## Комментарий к ТЗ

1. При создании PR автоматически назначаются **до двух** активных ревьюверов из **команды автора**, исключая самого автора.
2. Переназначение заменяет одного ревьювера на случайного **активного** участника **из команды заменяемого** ревьювера.

Получается, ревьюеры могут быть в любом случае только из команды автора: при переназначении PR юзер берется из команды заменяемого юзера, следовательно, из команды автора.