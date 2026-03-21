# IAM 产品需求文档 (PRD)

> 身份认证与访问管理系统 (Identity and Access Management)
> 当前版本：v0.1.9-draft

本文档采用目录结构维护，便于按主题阅读、按需求迭代和按模块维护。

## 文档定位

- 面向产品经理：管理需求范围、优先级和迭代计划
- 面向研发团队：按模块和 `REQ` 粒度定位实现范围
- 面向测试团队：根据验收标准和业务流程设计测试用例
- 面向项目管理：跟踪版本计划、风险依赖和成功指标

## 快速入口

如需在几分钟内理解需求范围、优先级和主要风险，先阅读 [00-executive-summary.md](./00-executive-summary.md)。

## 产品目标

| 目标类型 | 具体目标 |
|----------|----------|
| 业务目标 | 为 SaaS 产品提供开箱即用的 IAM 能力，降低开发成本 70% |
| 技术目标 | 支持千万级用户、99.9% 可用性、API 响应 < 100ms |
| 安全目标 | 通过 SOC2 合规要求，支持 MFA、审计日志、数据加密 |
| 体验目标 | 开发者友好，API 文档完善，SDK 支持主流语言 |

## 需求概览

| 优先级 | 数量 | 估时 |
|--------|------|------|
| P0 | 7 | 21 人天 |
| P1 | 3 | 10 人天 |
| P2 | 1 | 5 人天 |

当前共 11 个需求，涵盖认证管理、用户管理、权限管理、租户管理和审计日志五个模块。最新新增需求为 [REQ-011 OAuth2 第三方登录](./05-functional-requirements/REQ-011-oauth2-third-party-login.md)。

## 阅读导航

| 主题 | 说明 | 文档 |
|------|------|------|
| 一页摘要 | 快速理解产品定位、范围、优先级和待澄清项 | [00-executive-summary.md](./00-executive-summary.md) |
| 文档概述 | 目的、适用范围、名词解释、修订历史 | [01-document-overview.md](./01-document-overview.md) |
| 产品概述 | 背景、定位、目标、技术栈、Token 方案 | [02-product-overview.md](./02-product-overview.md) |
| 用户分析 | 用户角色、场景、租户与用户关系 | [03-user-analysis.md](./03-user-analysis.md) |
| 需求总览 | 需求总表、状态、优先级、迭代计划 | [04-requirements-overview.md](./04-requirements-overview.md) |
| 功能需求详情 | 按 `REQ` 粒度拆分的详细需求文档 | [05-functional-requirements/README.md](./05-functional-requirements/README.md) |
| 非功能需求 | 性能、安全、可用性、多租户隔离 | [06-non-functional-requirements.md](./06-non-functional-requirements.md) |
| 业务流程 | 登录、权限分配、租户开通流程 | [07-business-flows.md](./07-business-flows.md) |
| 风险与依赖 | 技术风险、关键基础设施依赖 | [08-risks-and-dependencies.md](./08-risks-and-dependencies.md) |
| 成功指标 | 业务、技术、用户体验指标 | [09-success-metrics.md](./09-success-metrics.md) |
| 附录 | 用户故事清单 | [10-appendix-user-stories.md](./10-appendix-user-stories.md) |

## 功能需求入口

| 模块 | 需求 |
|------|------|
| 认证管理 | `REQ-001`、`REQ-002`、`REQ-003`、`REQ-008`、`REQ-011` |
| 用户管理 | `REQ-004` |
| 权限管理 | `REQ-005`、`REQ-006` |
| 租户管理 | `REQ-007` |
| 审计日志 | `REQ-009`、`REQ-010` |

详细文档请从 [05-functional-requirements/README.md](./05-functional-requirements/README.md) 进入。
