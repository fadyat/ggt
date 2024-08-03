.PHONY: tools

tools:
	@cd ./tools && go mod tidy && go mod verify && go generate --tags=tools -x

generate:
	@go generate ./...