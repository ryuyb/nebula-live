.PHONY: init
# init env
init:
	go install github.com/air-verse/air@latest

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: run
# run
run:
	go run ./...

.PHONY: air
# air
air:
	air -c .air.toml

.PHONY: build
# build
build:
	mkdir -p bin/ & go build -o ./bin/ ./...

.PHONY: swagger
swagger:
	swag fmt
	swag init -g cmd/app/main.go --output docs

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [targe]'
	@echo ''
	@echo 'Targets'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
