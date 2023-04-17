proto:
	protoc --proto_path=cmd/addressbook/transport/grpc/proto --go_out=cmd/addressbook/transport/grpc/pb --go_opt=paths=source_relative \
    --go-grpc_out=cmd/addressbook/transport/grpc/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=cmd/addressbook/transport/grpc/pb --grpc-gateway_opt=paths=source_relative \
    cmd/addressbook/transport/grpc/proto/*.proto