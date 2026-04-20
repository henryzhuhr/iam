# 冻结 v1 权限模型边界

## 问题描述

当前权限相关需求已经从基础 RBAC 扩展到用户组层级、权限继承、允许/拒绝策略、组管理员、应用维度角色等复杂能力，但 v1 并没有一个被明确冻结的权限模型说明。

## 背景

- 摘要已经将权限模型粒度列为待澄清项
- REQ-006 以用户-角色-权限并集为基础
- REQ-014 又引入了树形用户组、继承、拒绝策略、组管理员
- TDD 对用户组的 API 与数据结构承接不足，尚不足以支撑这些高级特性

如果继续在未冻结模型的情况下推进实现，权限计算和鉴权中间件会最先失控。

## 解决方案

### 待决策项

- [ ] v1 是否坚持纯 RBAC，还是纳入用户组权限与继承
- [ ] 是否允许“用户直接权限”，还是只允许角色授权
- [ ] 是否支持拒绝策略；如果不支持，应从 REQ-014 中移除
- [ ] 角色是否支持应用维度生效范围
- [ ] 用户最终权限计算公式在 v1 中到底是什么

### 需要回写的文档

- [ ] 更新产品摘要中的权限模型描述
- [ ] 更新 `REQ-005`、`REQ-006`、`REQ-014` 的边界和验收标准
- [ ] 更新系统架构中的授权原则与权限计算说明
- [ ] 更新 TDD 中权限、角色、用户组模块的依赖与表设计

### 推荐补充内容

- [ ] 给出 v1 权限模型总览图，明确 User / Role / Group / Permission / Application 的关系
- [ ] 给出权限计算示例，覆盖“多角色”“多组”“应用范围”“冲突角色”场景
- [ ] 给出哪些能力明确延期到 roadmap，例如 deny 规则、动态组、复杂 SoD

### 完成标准

- [ ] 用一句公式即可描述 v1 最终权限计算方式
- [ ] `REQ-006` 与 `REQ-014` 不再互相扩张同一能力边界
- [ ] TDD 可以据此稳定定义鉴权中间件和查询逻辑

## 相关链接

- [docs/PRD/00-executive-summary.md](../docs/PRD/00-executive-summary.md)
- [docs/PRD/05-functional-requirements/REQ-005-role-management.md](../docs/PRD/05-functional-requirements/REQ-005-role-management.md)
- [docs/PRD/05-functional-requirements/REQ-006-permission-assignment.md](../docs/PRD/05-functional-requirements/REQ-006-permission-assignment.md)
- [docs/PRD/05-functional-requirements/REQ-014-user-group-management.md](../docs/PRD/05-functional-requirements/REQ-014-user-group-management.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
