MAKEFLAGS += --no-print-directory

include .env.development
export

run:
	@CONFIG_PATH=.env.development go run ./cmd/app/main.go

env-up:
	@docker compose up -d db redis;

env-down:
	@docker compose down

migrate-create:
	@if [ -z "$(name)" ]; then\
		echo "Ошибка: нужно указать имя миграции. Пример: make migrate-create name=init"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir ./migrations -seq $(name);

migrate-action:
	@if [ -z "$(action)" ]; then\
		echo "Ошибка: нужно указать action. Пример: make migrate-action action=\"up 2\""; \
		exit 1; \
	fi; \
	migrate -database ${DB_DSN} -path=./migrations $(action)

migrate-up:
	@make migrate-action action="up"

migrate-down:
	@make migrate-action action="down"
