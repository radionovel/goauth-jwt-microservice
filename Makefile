run: build
	@./bin/auth

build:
	@mkdir -p ./bin
	@go build -o ./bin/auth ./cmd/auth/