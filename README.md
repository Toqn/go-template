# go-template

A minimal Go HTTP service skeleton using a ports-and-adapters layout.

## Layout

| Path                   | Purpose                                                                 |
| ---------------------- | ----------------------------------------------------------------------- |
| `cmd/api`              | Entrypoint and composition root: wires config, logging, adapters, HTTP. |
| `api/v1`               | Public wire types (request/response bodies).                            |
| `internal/core`        | Business rules. May import only selected stdlib, `api/...`, `internal/ports`. |
| `internal/ports`       | Interfaces consumed by core and implemented by adapters.                |
| `internal/adapters`    | Implementations of ports (DBs, queues, external APIs). Must not import core. |
| `internal/transport`   | Inbound transports (HTTP handlers, gRPC, ...). Must not import adapters. |
| `internal/platform`    | Cross-cutting plumbing: config, logging, tracing.                       |
| `internal/testkit`     | Shared test helpers and fakes.                                          |

The import rules are enforced by `depguard` in `.golangci.yml`.

## Configuration

| Variable           | Default | Values                         |
| ------------------ | ------- | ------------------------------ |
| `APP_ENV`          | `dev`   | `dev`, `test`, `prod`          |
| `HTTP_ADDR`        | `:8080` | any `host:port`                |
| `LOG_LEVEL`        | `info`  | `debug`, `info`, `warn`, `error` (case-insensitive) |
| `SHUTDOWN_TIMEOUT` | `5s`    | Go duration, e.g. `10s`        |

## Make targets

```bash
make run    # build and start the server
make test   # go test -race -shuffle=on ./...
make lint   # golangci-lint (pinned version, fetched via go run)
make all    # fmt, lint, test, build
```

## Endpoints

- `GET /healthz` returns `{"status":"ok"}`.

Every response carries an `X-Request-ID` header. A well-formed incoming
`X-Request-ID` is reused; otherwise a random one is generated.
