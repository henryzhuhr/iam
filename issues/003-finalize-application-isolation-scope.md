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

### 待决策项

- [ ] 应用是否作为 IAM 内的一等资源统一管理
- [ ] 应用管理能力是否继续挂在 REQ-017，还是拆成独立“应用管理”需求
- [ ] 登录态 token 中是否固化 `apps` claim，还是改为实时查询
- [ ] 应用授权的最小粒度是什么：仅用户-应用，还是角色/权限也需要应用维度
- [ ] 审计日志和登录日志是否需要记录 `app_id`

### 需要回写的文档

- [ ] 更新 `REQ-017` 的状态，从“待决策”改为明确结论
- [ ] 更新产品摘要、需求总览、系统架构中的模块命名和边界
- [ ] 更新 TDD 中应用模块 API、表结构和登录时序
- [ ] 更新 SQL 命名和字段说明，避免继续出现并行叫法

### 推荐补充内容

- [ ] 给出应用与内部 `Client` 的对照表，避免业务应用与机器主体混淆
- [ ] 给出应用授权何时生效、何时失效、禁用应用如何影响已有 token
- [ ] 补充应用删除、冻结、迁移等异常场景

### 完成标准

- [ ] PRD、TDD、SQL 对“Application”的定义完全一致
- [ ] 能明确回答应用授权是登录时固化还是请求时实时计算
- [ ] REQ-017 不再保留“待决策”状态

## 相关链接

- [docs/PRD/04-requirements-overview.md](../docs/PRD/04-requirements-overview.md)
- [docs/PRD/05-functional-requirements/REQ-017-application-isolation.md](../docs/PRD/05-functional-requirements/REQ-017-application-isolation.md)
- [docs/PRD/03-system-architecture.md](../docs/PRD/03-system-architecture.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
- [sql/001_init.sql](../sql/001_init.sql)
