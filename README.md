# OCPI Service

OCPI is a Go-based implementation of the **Open Charge Point Interface (OCPI)** hub/backend.  
It provides:

- OCPI 2.2.1-compatible HTTP API endpoints
- gRPC API for internal integration
- Postgres-backed storage for locations, tariffs, sessions, tokens, CDRs, etc.
- SDK for interacting with OCPI endpoints
- Tools for running migrations, tests, and building Docker images

---

## Project Layout

This repository contains multiple services. The main ones are:

- `ocpi/` – OCPI hub/backend service (this README)
- `ocpiem/` – OCPI emulator service
- `ocpp/`, `ocpp2/` – OCPP-related services

Each service is a separate Go module with its own `go.mod`, `config.yml`, and `Makefile`.

---

## Requirements

- Go 1.21+
- Make
- Docker (optional, for container builds)
- PostgreSQL 13+ (for production / local DB)
- `goose` migration tool (for DB migrations)
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` (if regenerating protobufs)
- `mockery` (if regenerating mocks)
- `swag` (if regenerating Swagger)

---

## Getting Started (OCPI Service)

All commands in this section assume you are in the `ocpi/` directory.

### 1. Clone

```
bash git clone [https://github.com/](https://github.com/)<your-org>/ocpi.git cd ocpi
```

### 2. Configure `ROOT` and config

`ROOT` is used to locate configuration files. For local development:


Ensure you have a configuration file:

- `ocpi/config.yml` – main service configuration
- Optional `.env` file in `ocpi/.env` to override settings (use `.env-default` as a template)

You can start by copying the default:

### 3. Dependencies

```
bash make dep make vendor
```

---

## Database Setup

The service uses PostgreSQL and migrations managed by `goose`.

**Environment variables** (can be set in `.env`):

- `DB_HOST`, `DB_PORT`, `DB_NAME`
- `DB_ADMIN_USER`, `DB_ADMIN_PASSWORD`
- `DB_OCPI_USER`, `DB_OCPI_PASSWORD`

### Initialize schema (create DB objects)

````
make db-init-schema
````

This produces a binary in `ocpi/bin/main`.

### Run locally

Make sure `ROOT` is set and DB is reachable, then:


---

## API & SDK

The OCPI service exposes:

- **HTTP API** – OCPI 2.2.1 endpoints (locations, tariffs, tokens, sessions, CDRs, commands, credentials, etc.)
- **gRPC API** – for internal communication and SDK usage
- **SDK package** – Go client to interact with OCPI endpoints

To regenerate protobufs (if you modify `proto/*.proto`):


---

## Configuration Overview

The main configuration (`config.yml`) covers:

- gRPC server (host, port, tracing)
- HTTP server (port, timeouts, tracing)
- Storages (Postgres master/slave, migration path)
- Logging (level, format, contextual fields)
- Monitoring (metrics endpoint, Go runtime metrics)
- Profiling (pprof settings)
- OCPI:
    - Local platform (id, name, roles, token A, versions)
    - Local party (ids, roles)
    - Webhook behaviour (mock/real, timeout)
    - Remote platforms (timeout, mock mode)
    - Emulator integration
- Test config (e.g. webhook URL for tests)

Most values can be overridden via environment variables (see `config.yml` for the full list).

---

## Makefile Targets (OCPI)

Common targets:

- `dep` – tidy Go modules
- `vendor` – vendor dependencies
- `lint` – `go vet` + `go fmt`
- `build` – build binary to `./bin/main`
- `proto` – generate protobuf Go and gRPC stubs
- `swagger` – generate Swagger documentation
- `test`, `test-with-coverage`, `test-integration` – run tests
- `db-init-schema`, `db-up`, `db-down`, `db-status`, `db-create` – DB schema & migrations
- `docker-build`, `docker-build-test`, `docker-push`, `docker-push-test`, `docker-run` – Docker workflow


---

## Development Notes

- Set `ROOT` to the absolute path of the repo root so the service can find `config.yml` and `.env`.
- Use `.env` to store local, non-committed overrides.
- For working on multiple services (`ocpi`, `ocpiem`, `ocpp`, `ocpp2`), treat each subdirectory as a separate Go module and follow its own `README`/`Makefile` conventions.
- For mocks & CI checks, see `mock` and `ci-*` targets in the Makefiles.

---

## License

Apache 2.0


