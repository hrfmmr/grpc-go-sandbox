.PHONY: generate
generate:
	buf generate

.PHONY: build
build:
	@go build cmd/server/main.go

.PHONY: run-server
run-server:
	@go run cmd/server/main.go

.PHONY: test-grpcurl
test-grpcurl:
	grpcurl -plaintext -d '{"name": "john"}' \
		localhost:8080 \
		greeter.v1.GreetingService.Hello

.PHONY: test-breaking
test-breaking:
	cd .. && buf breaking \
		--against '.git#branch=feat/use-buf'
