---
name: go-code-library
description: 面向 Go 后端场景的可复用代码模板库，覆盖 PostgreSQL 连接、Redis 连接、配置加载、HTTP 服务启动、日志、优雅退出，以及 repository/service 分层等公共模块。用于 Codex 需要搭建 Go 服务骨架、提供参考代码片段，或快速组装新 Go 项目的基础设施模块时。
---

# Go 代码模板库

## 概览

使用技能内置的模板快速组装面向生产的 Go 后端代码。优先从 `assets/` 复制模板并按需调整名称、配置键和业务逻辑，而不是每次从零编写基础设施代码。

## 工作流程

1. 判断目标形态：
   - 当用户需要一个小型服务骨架时，使用 `assets/basic-service-template/`。
   - 当用户只需要单个公共模块或一段聚焦代码时，使用 `assets/snippets/`。
2. 阅读 [references/catalog.md](references/catalog.md) 找到最接近的模板。
3. 将相关文件或目录复制到当前工作区。
4. 重命名包路径，调整 `go.mod` 模块路径，并替换占位配置值。
5. 保持 handler 层轻量，把集成逻辑放在 `internal/platform` 或 `internal/repository` 中。

## 约定

- PostgreSQL 连接池优先使用 `pgx/v5` 和 `pgxpool`。
- Redis 客户端优先使用 `go-redis/v9`。
- 配置放到独立 package 中，并在进程启动时统一加载。
- 仅在测试隔离确实需要时，才为基础设施客户端额外抽接口。
- 所有 I/O 路径都传递 `context.Context`。
- 在添加业务接口前，先补齐健康检查或就绪检查接口。

## 资源

- [references/catalog.md](references/catalog.md)：模板索引与适用场景说明。
- [references/customization-checklist.md](references/customization-checklist.md)：复制模板后需要修改的关键项。
- `assets/basic-service-template/`：最小服务骨架，包含配置、日志、PG、Redis、HTTP 服务和优雅退出。
- `assets/snippets/`：常见公共模块的独立代码片段。
