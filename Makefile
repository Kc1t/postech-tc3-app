.PHONY: help run build up down logs tidy swagger test

APP_NAME=workshop-api
MAIN=./cmd/api
DOCS_DIR=docs

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "  run       - roda a API localmente (sem Docker)"
	@echo "  build     - compila o binario"
	@echo "  up        - sobe todos os containers (build incluso)"
	@echo "  up-d      - sobe todos os containers em background"
	@echo "  down      - para e remove os containers"
	@echo "  logs      - exibe logs dos containers"
	@echo "  swagger   - gera a documentacao Swagger"
	@echo "  tidy      - baixa e limpa dependencias"
	@echo "  test      - roda os testes"

run:
	go run $(MAIN)

build:
	go build -o $(APP_NAME) $(MAIN)

up:
	docker compose up --build

up-d:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

swagger:
	swag init -g $(MAIN)/main.go -o $(DOCS_DIR)

tidy:
	go mod tidy

test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
