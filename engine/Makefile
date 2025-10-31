proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

build:
	go build -o ./out/engine ./cmd/engine

run: build
	./out/engine