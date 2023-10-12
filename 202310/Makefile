TARGET_PROTO_FILE := proto/greeter/v1/greeting.proto

.PHONY: protoc
protoc:
	@protoc \
		--go_out=. \
		--go-grpc_out=. \
		${TARGET_PROTO_FILE}
