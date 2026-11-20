.PHONY: run-server run-client up down test lint migrate-up migrate-down

# --- Docker ---
# Сборка и запуск всего в Docker (сервер + базы)
up:
	docker compose -f deploy/docker-compose.yml up -d --build

# Запуск только баз данных (Postgres + Redis)
db-up:
	docker compose -f deploy/docker-compose.yml up -d postgres redis

# Остановка всех контейнеров
down:
	docker compose -f deploy/docker-compose.yml down

# --- Запуск приложения ---
# Запуск сервера локально (нужны запущенные базы)
run-server:
	go run ./cmd/server

# Запуск CLI-клиента локально
run-client:
	go run ./cmd/client

# --- Миграции ---
# Накатить все миграции (нужен установленный goose)
migrate-up:
	goose -dir migrations postgres "postgres://converter:converter@localhost:5432/converter?sslmode=disable" up

# Откатить последнюю миграцию
migrate-down:
	goose -dir migrations postgres "postgres://converter:converter@localhost:5432/converter?sslmode=disable" down

# --- Тесты и линтер ---
# Запуск всех тестов
test:
	go test -v ./...

# Запуск линтера
lint:
	golangci-lint run