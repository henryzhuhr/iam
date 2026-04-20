# 确认 REQ-017 应用隔离最终方案

## 问题描述

REQ-017 当前仍标记为“待决策”，但架构、TDD、SQL 已经按照“应用是 IAM 一等实体”的方向继续扩展。需求状态与设计状态不一致，已经形成文档漂移。

## 背景

- PRD 总览将 REQ-017 标为 P0 且纳入版本计划
- REQ-017 文档中已经给出推荐方案 A
- TDD 已定义 `/apps` 模块和相关 API
- SQL 已创建 `applications` 与 `user_app_authorizations` 表

如果不先冻结方案，后续 token claim、应用授权、审计维度、租户边界都会反复调整。

## 解决方案

### 已落定结论

- [x] 应用作为 IAM 内的一等资源统一管理
- [x] 应用管理能力继续挂在 REQ-017，并明确为“应用管理与数据隔离功能”
- [x] 登录态 token 固化 `apps` claim
- [x] v1 应用授权粒度为用户-应用，角色通过用户分配时附加应用范围
- [x] 审计日志和登录日志增加可选 `app_id`

### 需要回写的文档

- [x] 更新 `REQ-017` 的状态，从“待决策”改为明确结论
- [x] 更新产品摘要、需求总览、系统架构中的模块命名和边界
- [x] 更新 TDD 中应用模块 API、表结构和登录时序
- [x] 更新 SQL 命名和字段说明，避免继续出现并行叫法

### 推荐补充内容

- [x] 给出应用与内部 `Client` 的对照表，避免业务应用与机器主体混淆
- [x] 给出应用授权何时生效、何时失效、禁用应用如何影响已有 token
- [x] 补充应用删除、冻结、迁移等异常场景

### 完成标准

- [x] PRD、TDD、SQL 对“Application”的定义完全一致
- [x] 能明确回答应用授权是登录时固化还是请求时实时计算
- [x] REQ-017 不再保留“待决策”状态

## 处理结果

已将 REQ-017 重写为最终需求文档，明确 IAM 统一管理应用、登录固化 `apps` claim、日志支持 `app_id`，并统一 PRD/TDD/SQL 命名。

## 相关链接

- [docs/PRD/04-requirements-overview.md](../docs/PRD/04-requirements-overview.md)
- [docs/PRD/05-functional-requirements/REQ-017-application-isolation.md](../docs/PRD/05-functional-requirements/REQ-017-application-isolation.md)
- [docs/PRD/03-system-architecture.md](../docs/PRD/03-system-architecture.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
- [sql/001_init.sql](../sql/001_init.sql)
