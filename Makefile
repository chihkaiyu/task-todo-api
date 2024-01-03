.PHONY: build-api
build-api: gen-swagger
	docker build -t api -f $(shell pwd)/infra/app/api/Dockerfile .

.PHONY: build-swaggo
build-swaggo:
	docker build -t swaggo:1.20.12 -f $(shell pwd)/infra/docker/swaggo/Dockerfile .

.PHONY: build-migrate
build-migrate:
	docker build -t sql-migrate:latest $(shell pwd)/infra/docker/migration

.PHONY: gen-swagger
gen-swagger:
	docker run --rm \
		-v $(shell pwd):/go/src/github.com/chihkaiyu/task-todo-api \
		-v ${GOPATH}/pkg/mod:/go/pkg/mod \
		-e "GOPATH=/go" \
		swaggo:1.20.12 sh -c "cd /go/src/github.com/chihkaiyu/task-todo-api && swag init -g ./cmd/api/main.go -o ./cmd/api/docs"

.PHONY: migrate
migrate:
	docker run --rm --network task-todo-api_default -v $(shell pwd)/infra/databases:/app/databases sql-migrate:latest up api

