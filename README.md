## Go service


This repo comes with a set of `make` commands to streamline development, code generation, testing, linting, and Docker image building. Here's a quick guide to what each command does.

---

## 🔧 Setup & Code Generation

### `make di-generate`
Generate dependency injection using [wire](https://github.com/google/wire).  
> 🔄 Typically used after updating service constructors.

### `make code-gen`
Generate database schema and code from Ent.  
> 🎉 Required whenever the schema is updated.

---

## 🧪 Testing & Linting

### `make test`
Run all tests with verbose output.

### `make test-silent`
Run all tests without verbose output.

### `make test-cov`
Run tests and display coverage information.

### `make go-lint`
Run static analysis using `golangci-lint`.

### `make go-lint-install`
Install `golangci-lint` and set up Git pre-commit hook.  
> ⚠️ Run this once after cloning the repo.

---

## 🧬 Migrations

### `make migrate-gen`
Generate a new blank migration **after** running `code-gen`.  
You'll be prompted to enter a migration name (e.g., `create-users-table`).

### `make migrate-new`
Generate a new blank migration without triggering code-gen first.

### `make migrate-up`
Apply all pending migrations.

### `make migrate-down`
Rollback to a specific migration version.  
You'll be prompted to input the version (e.g., `20250304025256`).

### `make migrate-hash`
Recalculate integrity hash for the migration directory.

### `make migrate-status`
Check the current migration status.  
> 🧐 Useful for debugging what has/hasn’t been applied.

---

## 🐳 Docker Builds

### `make build-image-local`
Build Docker image for the service using the dev Dockerfile.

### `make build-linux`
Compile the binary for Linux (amd64).  
> ⛓ Required before creating the Docker image.

### `make download-go-mod`
Download Go module dependencies.  
> 📦 Ensures all dependencies are fetched before build.

### `make run-all`
Run docker compose for all service
---

## 💡 Tips

- Run `make go-lint-install` once after cloning to set up hooks.
- Always run `make code-gen` before generating migrations.
- Keep `.env` files and config up to date before running Docker builds.
"""
# go-service
