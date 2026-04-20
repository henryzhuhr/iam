# 需求缺口收敛清单

## 问题描述

当前 IAM 项目的 PRD、TDD 与代码实现之间已经出现明显漂移。问题并不只是“文档还没写完”，而是有若干核心需求仍未收口，但下游设计和 SQL 已经按某个假设继续推进，继续实现会放大返工成本。

## 背景

本轮项目阅读后，已确认以下现象同时存在：

- PRD 内部仍存在若干待澄清项，但部分内容已经被 TDD 和 SQL 当作既定结论落地
- 部分需求文档的用户故事、接口路径、数据库表设计之间不一致
- 权限模型的复杂度已经超出当前 v1 的清晰边界
- 代码仓库当前只有少量基础路由落地，设计范围远大于实现范围，状态表达不准确

## 解决方案

### 一、P0 收口项

- [ ] 明确注册模型与租户归属规则，见 [002-clarify-registration-and-tenant-binding.md](./002-clarify-registration-and-tenant-binding.md)
- [ ] 确认 REQ-017 的最终方案与需求状态，见 [003-finalize-application-isolation-scope.md](./003-finalize-application-isolation-scope.md)
- [ ] 对齐 REQ-018 的用户故事、接口与表结构设计，见 [004-align-internal-service-auth-design.md](./004-align-internal-service-auth-design.md)
- [ ] 冻结 v1 权限模型边界，见 [005-freeze-permission-model-v1.md](./005-freeze-permission-model-v1.md)

### 二、P1 补充项

- [ ] 补充验证码、邮件、短信、邀请注册等外部依赖需求，明确接入边界、失败处理、频控与模板要求
- [ ] 补充审计日志与登录日志的保留期、导出、脱敏、查询权限、合规要求
- [ ] 为 PRD、TDD、代码建立统一状态口径，避免“待开发 / 已设计 / 已实现”混用

### 三、按层面拆解的待办

#### PRD 层

- [ ] 在产品摘要和需求总览中删除或收敛自相矛盾描述
- [ ] 为所有 P0 需求补齐主体、边界、依赖、异常处理和验收标准
- [ ] 补齐缺失或错误的用户故事映射
- [ ] 明确应用管理是否作为独立模块或继续并入 REQ-017 / REQ-007

#### TDD 层

- [ ] 根据最终 PRD 结论回写 API 路径、表结构、模块边界
- [ ] 补齐用户组权限、内部客户端凭证、应用授权等高复杂度模块的实现约束
- [ ] 为登录、注册、应用授权、客户端认证等核心流程补完整时序与错误处理

#### 代码与数据层

- [ ] 在需求冻结前，避免继续扩散新的接口命名和表结构命名
- [ ] 将“已实现”“占位实现”“仅设计未实现”三种状态区分开
- [ ] 后续开发按冻结后的文档回填测试和路由，不再让 SQL/TDD 先行漂移

### 四、建议执行顺序

1. 先完成 002、003、004、005 四个 issue 的决策与文档回写
2. 再统一更新 PRD 总览、TDD 和 SQL 命名
3. 最后基于冻结后的需求拆实现任务与测试计划

### 五、完成标准

- [ ] PRD 不再包含已知自相矛盾项
- [ ] TDD 与 PRD 在接口、表结构、主体模型上保持一致
- [ ] SQL 命名与需求/设计一致，不再出现并行命名
- [ ] 需求状态可以明确区分“待决策 / 已定义 / 已设计 / 已实现”

## 相关链接

- [docs/PRD/00-executive-summary.md](../docs/PRD/00-executive-summary.md)
- [docs/PRD/04-requirements-overview.md](../docs/PRD/04-requirements-overview.md)
- [docs/PRD/05-functional-requirements/REQ-002-user-registration.md](../docs/PRD/05-functional-requirements/REQ-002-user-registration.md)
- [docs/PRD/05-functional-requirements/REQ-017-application-isolation.md](../docs/PRD/05-functional-requirements/REQ-017-application-isolation.md)
- [docs/PRD/05-functional-requirements/REQ-018-internal-service-authentication.md](../docs/PRD/05-functional-requirements/REQ-018-internal-service-authentication.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
- [internal/routes/routes.go](../internal/routes/routes.go)
- [sql/001_init.sql](../sql/001_init.sql)
