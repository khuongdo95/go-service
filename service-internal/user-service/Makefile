.PHONY: migrate-add migrate-run migrate-down di-generate test-cov build build-init go-lint generate-proto
di-generate:
	wire ./app

go-lint:
	golangci-lint run

# Test
test:
	go test -v ./...
test-silent:
	go test ./...

test-cov:
	go test -cover ./...

go-lint-install:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@cp hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit

# generate code for database schema
code-gen:
	@go generate cmd/ent/generate.go
	@echo "🎉 Database code generated successfully!"

# generate a blank migration
migrate-gen: code-gen
	@read -p "Enter migration name (e.g., create-todo-table): " name; \
			go run cmd/migrate/main.go ent_add $$name

# generate a blank migration
migrate-new:
	@read -p "Enter migration name (e.g., create-todo-table): " name; \
		go run cmd/migrate/main.go add $$name

# deploy migration
migrate-up:
	go run cmd/migrate/main.go up

# rollback migration
migrate-down:
	@read -p "Enter migration version (e.g., 20250304025256): " version; \
		go run cmd/migrate/main.go down $$version

# re-generate migration directory integrity hash 
migrate-hash:
	go run cmd/migrate/main.go hash

# check migration status
migrate-status:
	go run cmd/migrate/main.go statusa

SHELL := /bin/bash
CI_DOCKER_IMAGE=github
service=user

# build images from all services
.PHONY: build-image-local
build-image-local: build-linux
	docker build --platform linux/arm64 -f build/Dockerfile.dev -t user_service:latest . --build-arg service=${service}

.PHONY: build-linux
build-linux: download-go-mod
	GOOS=linux GOARCH=amd64 go build -mod=readonly -o out/cmd/main cmd/app/main.go

.PHONY: download-go-mod
download-go-mod: go.mod
	go mod download all
