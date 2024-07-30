APP_NAME=ggt

.PHONY: bin fmt lint

bin:
	@go build \
		-o .bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

fmt:
	@golangci-lint run --fix -v ./...

lint:
	@golangci-lint run -v ./...