.PHONY: lint

lint:
	go mod tidy
	gofmt -w .
	goimports -w .
