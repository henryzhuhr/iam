# 参考资料目录

本目录收录 IAM 项目相关的技术概念、协议标准和设计方案参考资料，供研发团队查阅。

## 快速导航

| 分类 | 文档数量 | 核心文档 |
|------|----------|----------|
| [**认证管理**](#认证管理) | 7 篇 | JWT, OAuth 2.0, OIDC, MFA |
| [**权限管理**](#权限管理) | 3 篇 | RBAC, 数据权限 |
| [**安全与架构**](#安全与架构) | 3 篇 | OWASP Top 10, API 网关 |

---

## 资料清单

## 认证管理

| 主题 | 说明 | 文档 |
|------|------|------|
| OAuth 2.0 | OAuth 2.0 协议核心概念、授权模式、流程说明 | [oauth-2.0-basics.md](./oauth-2.0-basics.md) |
| OpenID Connect (OIDC) | 基于 OAuth 2.0 的身份层协议，用户身份认证标准 | [oidc-basics.md](./oidc-basics.md) |
| JWT | JWT Token 结构、签名算法、最佳实践 | [jwt-basics.md](./jwt-basics.md) |
| Token 方案选型 | JWT vs Opaque Token 对比与选型依据 | [token-strategy.md](./token-strategy.md) |
| MFA/TOTP | 多因素认证与 TOTP 动态验证码原理 | [mfa-totp.md](./mfa-totp.md) |
| 密码安全 | 密码加密存储 (bcrypt/argon2)、加盐、彩虹表防护 | [password-security.md](./password-security.md) |
| 会话管理 | 会话生命周期、并发控制、会话固定攻击防护 | [session-management.md](./session-management.md) |

## 权限管理

| 主题 | 说明 | 文档 |
|------|------|------|
| 权限模型 | RBAC、ABAC、ReBAC 等权限模型对比与选型 | [permission-models.md](./permission-models.md) |
| RBAC 详细设计 | 角色层级、权限继承、约束 RBAC | [rbac-design.md](./rbac-design.md) |
| 数据权限 | 行级/列级数据访问控制 | [data-permission.md](./data-permission.md) |

## 安全与架构

| 主题 | 说明 | 文档 |
|------|------|------|
| 多租户架构 | SaaS 多租户数据隔离方案对比 | [multi-tenancy-architecture.md](./multi-tenancy-architecture.md) |
| OWASP Top 10 | Web 应用安全风险与防护 | [owasp-top10.md](./owasp-top10.md) |
| API 网关安全 | API 认证、限流、WAF | [api-gateway-security.md](./api-gateway-security.md) |

## 推荐阅读顺序

### 入门路径（必读）

1. [权限模型](./permission-models.md) - 理解 IAM 的核心设计基础
2. [OAuth 2.0](./oauth-2.0-basics.md) - 了解第三方登录协议
3. [JWT](./jwt-basics.md) - 理解 Token 认证机制
4. [Token 方案选型](./token-strategy.md) - 了解 JWT vs Opaque Token 的选型依据

### 认证管理进阶

5. [OpenID Connect (OIDC)](./oidc-basics.md) - 深入理解身份认证标准
6. [MFA/TOTP](./mfa-totp.md) - 多因素认证原理
7. [密码安全](./password-security.md) - 密码加密存储方案
8. [会话管理](./session-management.md) - 会话生命周期管理

### 权限管理进阶

9. [RBAC 详细设计](./rbac-design.md) - 权限模型具体实现
10. [数据权限](./data-permission.md) - 细粒度数据访问控制

### 安全与架构

11. [多租户架构](./multi-tenancy-architecture.md) - 数据隔离架构
12. [OWASP Top 10](./owasp-top10.md) - Web 安全风险与防护
13. [API 网关安全](./api-gateway-security.md) - API 层安全防护

## 外部参考资源

### 协议标准

| 资源 | 链接 |
|------|------|
| RFC 6749 (OAuth 2.0) | https://tools.ietf.org/html/rfc6749 |
| RFC 7519 (JWT) | https://tools.ietf.org/html/rfc7519 |
| OpenID Connect Core | https://openid.net/specs/openid-connect-core-1_0.html |
| TOTP 算法 (RFC 6238) | https://tools.ietf.org/html/rfc6238 |

### 技术文章

| 资源 | 说明 |
|------|------|
| OAuth.net | OAuth 协议官方资源站 |
| jwt.io | JWT 在线调试工具与文档 |
| Auth0 Blog | https://auth0.com/blog - 身份认证领域高质量技术博客 |

### 开源项目参考

| 项目 | 说明 |
|------|------|
| Keycloak | https://github.com/keycloak/keycloak - RedHat 开源 IAM 解决方案 |
| Casbin | https://github.com/casbin/casbin - Go 语言权限模型库 |
| Ory Hydra | https://github.com/ory/hydra - OAuth 2.0 服务器实现 |
