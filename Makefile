APP_NAME=ggt

.PHONY: bin fmt lint

bin:
	@go build \
		-o .bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

fmt:
	@golangci-lint run --fix -v ./...

lint:
	@golangci-lint run -v ./...

share: bin
	@sudo cp .bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)
