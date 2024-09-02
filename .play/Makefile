.PHONY: tools

tools:
	@cd ./tools && go mod tidy && go mod verify && go generate --tags=tools -x

generate:
	@go generate ./...

fmt:
	@.bin/golangci-lint run --fix -v ./...