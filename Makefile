.PHONY: lint

lint:
	gofmt -w .
	goimports -w .
