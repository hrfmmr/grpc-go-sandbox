CERTS_DEST ?= ssl
SERVER_CN ?= localhost

# Certificates Authority private key file (this shouldn't be shared)
CA_KEY := ca.key
CA_KEY_PW := caP@sswo0d
# Certificates Authority trust certificate (this should be shared with users)
CA_CERT := ca.crt
# Server private key, password protected (this shouldn't be shared)
SERVER_KEY := server.key
SERVER_KEY_PW ?= serverP@sswo0d
# Server certificate signing request (this should be shared with the CA owner)
SERVER_CSR := server.csr
# Server certificate signed by the CA (this would be sent back by the CA owner)
SERVER_CERT := server.crt
# Conversion of server.key into a format gRPC likes (this sholudn't be shared)
SERVER_PEM := server.pem

gen-greet-pb:
	@protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

run-greet-server:
	@go run greet/greet_server/server.go

run-greet-client:
	@go run greet/greet_client/client.go

gen-certs:
	@if [ ! -d "$(CERTS_DEST)" ]; then \
		mkdir $(CERTS_DEST); \
	fi
	"$(MAKE)" gen-ca-cert
	"$(MAKE)" gen-server-cert
	"$(MAKE)" gen-server-pem


# Generate CA certs
gen-ca-cert:
	openssl genrsa \
		-passout pass:$(CA_KEY_PW) \
		-aes256 \
		-out $(CERTS_DEST)/$(CA_KEY) \
		-rand /dev/urandom \
		4096
	openssl req \
		-passin pass:$(CA_KEY_PW) \
		-new \
		-x509 \
		-days 365 \
		-key $(CERTS_DEST)/$(CA_KEY) \
		-out $(CERTS_DEST)/$(CA_CERT) \
		-subj "/CN=${SERVER_CN}"

gen-server-cert:
	# Generate server private key
	openssl genrsa \
		-passout pass:$(SERVER_KEY_PW) \
		-aes256 \
		-out $(CERTS_DEST)/$(SERVER_KEY) \
		4096
	# Get a certificate signing request
	openssl req \
		-passin pass:$(SERVER_KEY_PW) \
		-new \
		-key $(CERTS_DEST)/$(SERVER_KEY) \
		-out $(CERTS_DEST)/$(SERVER_CSR) \
		-subj "/CN=${SERVER_CN}"
	# Sign the certificate with the CA (so called self-signed)
	openssl x509 \
		-req \
		-passin pass:$(CA_KEY_PW) \
		-days 365 \
		-in $(CERTS_DEST)/$(SERVER_CSR) \
		-CA $(CERTS_DEST)/$(CA_CERT) \
		-CAkey $(CERTS_DEST)/$(CA_KEY) \
		-set_serial 01 \
		-out $(CERTS_DEST)/$(SERVER_CERT)

gen-server-pem:
	# Convert the server private key to .pem format
	openssl pkcs8 \
		-topk8 \
		-nocrypt \
		-passin pass:$(SERVER_KEY_PW) \
		-in $(CERTS_DEST)/$(SERVER_KEY) \
		-out $(CERTS_DEST)/$(SERVER_PEM)
