APP_NAME=ggt
TOOLS=.play/.bin

.PHONY: bin fmt lint

bin:
	@go build \
		-o .bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

fmt:
	@$(TOOLS)/golangci-lint run --fix -v ./... || true
	@cd .play && make fmt

lint:
	@$(TOOLS)/golangci-lint run -v ./...

share: bin
	@sudo cp .bin/$(APP_NAME) /usr/local/bin/$(APP_NAME)
