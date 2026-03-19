# Template Catalog

## Basic Service Template

Path: `assets/basic-service-template/`

Use this template when the user wants a small but structured Go backend service. It includes:

- `cmd/server/main.go`: Startup wiring.
- `internal/config/config.go`: Environment-based config loading.
- `internal/platform/logging/logger.go`: `slog` logger construction.
- `internal/platform/postgres/postgres.go`: PostgreSQL pool setup and ping.
- `internal/platform/redis/redis.go`: Redis client setup and ping.
- `internal/platform/httpserver/server.go`: HTTP server with timeouts and graceful shutdown.
- `internal/handler/health.go`: Liveness/readiness endpoint.

## Standalone Snippets

Path: `assets/snippets/`

Use these when only one module is needed:

- `config_env.go`: Minimal environment config loader.
- `postgres_pgxpool.go`: PostgreSQL pool creation and shutdown.
- `postgres_tx.go`: Transaction helper wrapper.
- `redis_client.go`: Redis client creation and ping.
- `http_server.go`: Server bootstrap and graceful shutdown.
- `logger_slog.go`: Structured logger initialization.

## Selection Rules

- Prefer the full template if the request mentions "项目骨架", "脚手架", "service", "API 服务", or multiple infrastructure pieces at once.
- Prefer snippets if the request asks for a single module such as "连接 PG", "redis 客户端", or "优雅退出".
- Combine the full template with snippets only when adding a capability not already present.
