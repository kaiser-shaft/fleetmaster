include .env
export

run:
	@go run ./cmd/app/main.go

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

migrate-up:
	@docker compose up migrate

migrate-down:
	@docker compose run --rm --entrypoint /bin/sh migrate \
	    -c 'migrate -path=/migrations/ -database="${DB_DSN}" down'
