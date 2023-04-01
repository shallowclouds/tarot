.PHONY: lint

lint:
	go mod tidy
	gofmt -w .
	goimports -w .

.PHONY: divinetest
# Generate some random divine results.
divinetest:
	go run cmd/divinetest/main.go
