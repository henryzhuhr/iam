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

### 已落定结论

- [x] v1 坚持纯 RBAC，不纳入用户组权限与继承
- [x] 不允许用户直接权限，只允许角色授权
- [x] v1 不支持拒绝策略，并已从 REQ-014 中移除
- [x] 角色支持在用户分配时附加应用维度生效范围
- [x] 用户最终权限计算公式已确定

### 需要回写的文档

- [x] 更新产品摘要中的权限模型描述
- [x] 更新 `REQ-005`、`REQ-006`、`REQ-014` 的边界和验收标准
- [x] 更新系统架构中的授权原则与权限计算说明
- [x] 更新 TDD 中权限、角色、用户组模块的依赖与表设计

### 推荐补充内容

- [x] 给出 v1 权限模型总览图，明确 User / Role / Group / Permission / Application 的关系
- [x] 给出权限计算示例，覆盖“多角色”“多组”“应用范围”“冲突角色”场景
- [x] 给出哪些能力明确延期到 roadmap，例如 deny 规则、动态组、复杂 SoD

### 完成标准

- [x] 用一句公式即可描述 v1 最终权限计算方式
- [x] `REQ-006` 与 `REQ-014` 不再互相扩张同一能力边界
- [x] TDD 可以据此稳定定义鉴权中间件和查询逻辑

## 处理结果

v1 已冻结为“用户最终权限 = 所有有效角色权限的并集”，角色可在用户分配时附加应用范围；用户组仅做组织管理，不参与权限计算。

## 相关链接

- [docs/PRD/00-executive-summary.md](../docs/PRD/00-executive-summary.md)
- [docs/PRD/05-functional-requirements/REQ-005-role-management.md](../docs/PRD/05-functional-requirements/REQ-005-role-management.md)
- [docs/PRD/05-functional-requirements/REQ-006-permission-assignment.md](../docs/PRD/05-functional-requirements/REQ-006-permission-assignment.md)
- [docs/PRD/05-functional-requirements/REQ-014-user-group-management.md](../docs/PRD/05-functional-requirements/REQ-014-user-group-management.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
