# 5. 功能需求详情

本目录按 `REQ` 粒度拆分功能需求，便于单需求设计、开发、测试和评审。模块分组仅用于导航，具体优先级、估时和验收标准以各需求文档为准。

## 认证管理

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-001 | 用户登录功能 | P0 | [REQ-001-user-login.md](./REQ-001-user-login.md) |
| REQ-002 | 用户注册功能 | P0 | [REQ-002-user-registration.md](./REQ-002-user-registration.md) |
| REQ-003 | 密码重置功能 | P0 | [REQ-003-password-reset.md](./REQ-003-password-reset.md) |
| REQ-008 | MFA 多因素认证 | P1 | [REQ-008-mfa.md](./REQ-008-mfa.md) |
| REQ-011 | OAuth2 第三方登录 | P2 | [REQ-011-oauth2-third-party-login.md](./REQ-011-oauth2-third-party-login.md) |
| REQ-012 | Token 管理 | P0 | [REQ-012-token-management.md](./REQ-012-token-management.md) |
| REQ-013 | 密码策略管理 | P1 | [REQ-013-password-policy.md](./REQ-013-password-policy.md) |
| REQ-015 | 验证码登录 | P1 | [REQ-015-code-login.md](./REQ-015-code-login.md) |

## 用户管理

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-004 | 用户管理功能 | P0 | [REQ-004-user-management.md](./REQ-004-user-management.md) |
| REQ-014 | 用户组管理功能 | P1 | [REQ-014-user-group-management.md](./REQ-014-user-group-management.md) |

## 权限管理

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-005 | 角色管理功能 | P0 | [REQ-005-role-management.md](./REQ-005-role-management.md) |
| REQ-006 | 权限分配功能 | P0 | [REQ-006-permission-assignment.md](./REQ-006-permission-assignment.md) |

## 租户管理

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-007 | 租户管理功能 | P0 | [REQ-007-tenant-management.md](./REQ-007-tenant-management.md) |

## 审计日志

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-009 | 操作审计日志 | P1 | [REQ-009-operation-audit-log.md](./REQ-009-operation-audit-log.md) |
| REQ-010 | 登录日志记录 | P1 | [REQ-010-login-log.md](./REQ-010-login-log.md) |

## 系统基础

| REQ | 需求名称 | 优先级 | 文档 |
|-----|----------|--------|------|
| REQ-016 | API 限流和配额管理 | P1 | [REQ-016-rate-limit-quota.md](./REQ-016-rate-limit-quota.md) |
| REQ-017 | 应用级数据隔离 | P0 | [REQ-017-application-isolation.md](./REQ-017-application-isolation.md) |
