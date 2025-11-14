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
```

## Допущения

- Поле `needMoreReviewers` PullRequest отсутствует в openapi.yml, поэтому не реализовано

## Переменные окружения
 - DATABASE_URL = postgres://pr_user:pr_pass@db:5432/pr_db?sslmode=disable

## TODO

* Добавить простой эндпоинт статистики (например, количество назначений по пользователям и/или по PR).
* Добавить метод массовой деактивации пользователей команды и безопасную переназначаемость открытых PR (стремиться уложиться в 100 мс для средних объёмов данных).
- ТестоваяБД
- TeamMember or User
- Поменять ошибки на ErrResponse moldels.go

## Дополнительные задания
### E2E тесты
```bash
go test ./e2e/ -v
```
## Линтер
```bash
golangci-lint run
```
### Нагрузочное тестирование
```bash
k6 run loadtest.js
```
Результаты: loadtest_report.md