# dev:
# 	go run ./cmd/market-data-service
# sql:
# 	sqlc generate

# # example run: make new_migration name=00001_init_schema
# new_migration:
# 	migrate create -ext sql -dir sql/schema -seq $(name)

# migrate_up:
# 	migrate -path sql/schema -database "postgresql://postgres:Kiloma123@@localhost:5432/Crypto?sslmode=disable" -verbose up

# migrate_down:
# 	migrate -path sql/schema -database "postgresql://postgres:Kiloma123@@localhost:5432/Crypto?sslmode=disable" -verbose down


WORKER_IMAGE=1.24.0-alpine3.21

DIRS := $(shell find service-internal -mindepth 1 -maxdepth 1 -type d ! -name 'no_get')

DONT_STOP := db redis

.PHONY: tsl-generate build-services stop-services start-services reset-services doc-generate swagger-2-to-3 new-gateway


go-lint:
	@for dir in $(DIRS); do \
		if [ "$$(basename $$dir)" != "proto" ]; then \
			echo "Running golangci-lint in $$dir"; \
			cd $$dir && golangci-lint run && cd ..; \
		else \
			echo "Skipping $$dir"; \
		fi \
	done

run-all:
	@docker compose -f docker-compose.dev.yml up -d