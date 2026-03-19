# 定制检查清单

复制模板后，按下面的顺序调整：

1. 替换 `go.mod` 中的模块路径。
2. 按目标仓库结构重命名 `internal/` 下的包路径。
3. 如果项目已有约定，统一替换环境变量命名。
4. 根据实际需求调整 HTTP 超时、数据库连接池大小，以及 Redis 地址和 DB 配置。
5. 如果需要更严格的探活逻辑，将示例 `/healthz` 接口替换为真实的存活性或就绪性检查。
6. 不要把 SQL 或 Redis 访问逻辑直接写进 handler，而是补充 migration、repository 和 service 层。
7. 集成完成后执行 `gofmt` 和 `go test ./...`。
