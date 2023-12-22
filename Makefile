.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	go test ./...

.PHONY: codegen
codegen:
	go generate ./...
