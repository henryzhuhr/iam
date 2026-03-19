# 模板目录

## 基础服务模板

路径：`assets/basic-service-template/`

当用户需要一个小型但结构完整的 Go 后端服务时，使用这个模板。它包含：

- `cmd/server/main.go`：启动入口与依赖装配。
- `internal/config/config.go`：基于环境变量的配置加载。
- `internal/platform/logging/logger.go`：`slog` 日志初始化。
- `internal/platform/postgres/postgres.go`：PostgreSQL 连接池初始化与连通性检查。
- `internal/platform/redis/redis.go`：Redis 客户端初始化与连通性检查。
- `internal/platform/httpserver/server.go`：带超时配置和优雅退出的 HTTP 服务。
- `internal/handler/health.go`：存活性与就绪性检查接口。

## 独立代码片段

路径：`assets/snippets/`

当只需要单个公共模块时，使用这些片段：

- `config_env.go`：最小化环境配置加载代码。
- `postgres_pgxpool.go`：PostgreSQL 连接池创建与关闭。
- `postgres_tx.go`：事务执行辅助封装。
- `redis_client.go`：Redis 客户端创建与连通性检查。
- `http_server.go`：HTTP 服务启动与优雅退出。
- `logger_slog.go`：结构化日志初始化。

## 选择规则

- 如果请求里同时提到“项目骨架”、“脚手架”、“service”、“API 服务”或多个基础设施模块，优先使用完整模板。
- 如果请求只涉及单个模块，例如“连接 PG”、“Redis 客户端”或“优雅退出”，优先使用独立片段。
- 只有在完整模板中尚未包含所需能力时，才将完整模板和独立片段组合使用。
