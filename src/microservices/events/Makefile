include .env

gofmt:
	gofmt -s -w ./
.PHONY: gofmt


server: server-run
.PHONY: server


server-run: server-build
	./cmd/server/server \
		-a="${SERVER_ADDRESS}:${SERVER_PORT}" \
		-l="${LOG_LEVEL}"
.PHONY: server-run


server-build:
	go build -o ./cmd/server/server ./cmd/server/
.PHONY: server-build
