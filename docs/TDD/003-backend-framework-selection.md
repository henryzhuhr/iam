# 003 后端框架选型设计

> 日期：2026-04-24
> 作者：IAM Team
> 状态：已决策

---

## 1. 决策结论

IAM 后端主选型调整为 **Gin + grpc-go**：

- Gin 承载 HTTP REST API，面向前端控制台和开放接口
- grpc-go 承载 gRPC API，面向内部微服务调用和被调用
- HTTP 和 gRPC 契约分开维护，但必须共享同一套业务 Service
- 单进程双协议部署，HTTP 默认端口 `8888`，gRPC 默认端口 `9090`

**Kratos** 保留为唯一备选方案。后续如果 IAM 需要框架级 transport 统一、注册发现、配置中心、链路追踪治理深度整合，再评估迁移到 Kratos。

go-zero 作为历史实现背景保留在代码迁移记录中，不再作为目标架构推荐方案。当前选型不继续评估 Gin + grpc-go 和 Kratos 之外的其他 Go 框架。

## 2. 选型依据

IAM 当前阶段是模块化单体，业务优先级是认证安全、租户隔离、权限一致性和可测试性。HTTP 与 gRPC 的使用场景不同：

| 协议 | 使用方 | 设计重点 |
|------|--------|----------|
| HTTP | 前端控制台、开放 REST API | 易调试、统一响应格式、分页和表单错误表达 |
| gRPC | 内部微服务 | 强类型契约、低延迟调用、metadata 传递、服务间错误码 |

Gin + grpc-go 的优势是边界清晰、依赖可控、贴近 Go 标准生态。项目可以在不重写业务层的前提下分别演进 HTTP 和 gRPC 契约。

Kratos 的优势是更完整的双协议微服务框架能力，但当前阶段引入会提高框架约束和迁移成本。因此 Kratos 保留为后续治理增强场景的备选。

## 3. 目标架构

```text
Frontend / Open API
  -> HTTP :8888
     -> Gin router / middleware / handler / HTTP DTO
        -> shared service layer

Internal services
  -> gRPC :9090
     -> grpc-go server / interceptor / protobuf
        -> shared service layer

shared service layer
  -> repository
  -> MySQL / Redis / Kafka
```

分层约束：

- HTTP Handler 只处理 REST 入参、响应、状态码和前端友好的错误表达
- gRPC Server 只处理 proto 消息、metadata、gRPC status 和服务间调用语义
- 认证、租户隔离、权限校验、审计、Token 等业务规则统一沉到 Service 层
- Repository 和 infra 层不感知 HTTP 或 gRPC 协议

## 4. 治理与迁移边界

Gin + grpc-go 不内置完整微服务治理，项目通过明确组件补齐：

| 能力 | 当前方案 |
|------|----------|
| 日志 | HTTP middleware + gRPC interceptor 注入 request_id |
| 鉴权 | HTTP middleware + gRPC interceptor 统一解析身份上下文 |
| 超时 | HTTP Server timeout + gRPC deadline |
| 限流 | HTTP middleware + gRPC interceptor |
| 健康检查 | `/health` + gRPC health service |
| 指标与追踪 | 预留 Prometheus / OpenTelemetry |
| 服务发现 | 后续按部署环境选择 Kubernetes DNS、Consul、etcd 或配置中心 |

迁移到 Kratos 的触发条件：

- 需要框架统一管理 HTTP/gRPC transport 生命周期
- 需要注册发现、配置中心、链路追踪等治理能力深度集成
- Gin + grpc-go 的手工组合成本超过业务收益

## 5. 后续实施要求

- 先迁移 HTTP 服务适配层，再新增 gRPC 服务入口，业务层保持稳定
- proto 契约放在 `api/proto/`，生成代码与手写业务代码分离
- HTTP 和 gRPC 契约允许字段形态不同，但同一业务动作必须调用同一套 Service 方法
- 框架迁移期间必须保持现有 `/health` 和租户管理 API 测试通过
- 新增 gRPC health check，并补充至少一个 gRPC smoke 测试
