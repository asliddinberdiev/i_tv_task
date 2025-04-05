.PHONY:
.SILENT:
.DEFAULT_GOAL := run

CURRENT_DIR := $(shell pwd)
APP := $(shell basename ${CURRENT_DIR})
APP_CMD_DIR := ${CURRENT_DIR}/cmd/app

-include ./config/dev.env
export APP_ENVIRONMENT=dev

POSTGRESQL_URL = 'postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSLMODE}'

docker-build:
	@docker build -t $(APP) .

docker-run:
	@docker run -d --name $(APP) -p 8000:8000 $(APP)

docker-stop:
	@docker stop $(APP)
	@docker rm $(APP)

docker-clean:
	@docker rmi $(APP)

migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path migrations/postgres up

migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path migrations/postgres down

migrate-new: # make migrate-new name=file_name
	@migrate create -ext sql -dir migrations/postgres -seq $(name)

swag:
	@swag init -g internal/delivery/http/v1/routes.go

run: 
	@go run ${APP_CMD_DIR}/main.go