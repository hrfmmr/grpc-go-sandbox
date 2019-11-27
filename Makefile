gen-greet-pb:
	@protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

run-greet-server:
	@go run greet/greet_server/server.go

run-greet-client:
	@go run greet/greet_client/client.go
