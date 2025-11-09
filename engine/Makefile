proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

build:
	go build -o ./out/engine ./cmd/engine

run: build
	./out/engine

swagger:
	swag init -g ./cmd/engine/main.go -o ./docs --parseInternal --parseDependency

swagger-gen-ts: swagger
	docker run --rm \
	  -v $(CURDIR)/../ui/packages/wf-engine-api/src:/output \
      -v $(CURDIR):/local openapitools/openapi-generator-cli generate \
      -i /local/docs/swagger.json \
      -g typescript-fetch \
      -o /output \
      --additional-properties=stringEnums=true
	clean-swagger-output.bat

