build:
	go build ./...

run: build
	go run ./cmd/engine

test:
	go test ./...

migration:
	migrate create -ext sql -dir migrations $(name)

migrate-up:
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "Get-Content .env | ForEach-Object { if ($$_.Trim() -and -not $$_.Trim().StartsWith('#')) { $$parts = ($$_.Split('=',2)); [Environment]::SetEnvironmentVariable($$parts[0].Trim(), $$parts[1].Trim(), 'Process') } }; migrate -path migrations -database \"postgres://$${env:POSTGRES_USER}:$${env:POSTGRES_PASSWORD}@localhost:$${env:DB_PORT}/$${env:POSTGRES_DB}?sslmode=disable\" up"
else
	@set -a; . .env; migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DB_PORT}/${POSTGRES_DB}?sslmode=disable" up; set +a
endif

migrate-down:
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "Get-Content .env | ForEach-Object { if ($$_.Trim() -and -not $$_.Trim().StartsWith('#')) { $$parts = ($$_.Split('=',2)); [Environment]::SetEnvironmentVariable($$parts[0].Trim(), $$parts[1].Trim(), 'Process') } }; migrate -path migrations -database \"postgres://$${env:POSTGRES_USER}:$${env:POSTGRES_PASSWORD}@localhost:$${env:DB_PORT}/$${env:POSTGRES_DB}?sslmode=disable\" down 1"
else
	@set -a; . .env; migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${DB_PORT}/${POSTGRES_DB}?sslmode=disable" down 1; set +a
endif

swagger:
	swag init -g ./cmd/engine/main.go -o ./docs --parseInternal --parseDependency  --outputTypes json
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '..\ui\packages\wf-engine-api\docs' ; Copy-Item -Path './docs/swagger.json' -Destination '..\ui\packages\wf-engine-api\docs\swagger.json' -Force"
else
	mkdir -p ../ui/packages/wf-engine-api/docs
	cp ./docs/swagger.json ../ui/packages/wf-engine-api/docs/swagger.json
endif