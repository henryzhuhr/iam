# 参考资料目录

本目录收录 IAM 项目相关的技术概念、协议标准和设计方案参考资料，供研发团队查阅。

## 资料清单

| 主题 | 说明 | 文档 |
|------|------|------|
|  OAuth 2.0  | OAuth 2.0 协议核心概念、授权模式、流程说明 | [oauth-2.0-basics.md](./oauth-2.0-basics.md) |
|  权限模型  | RBAC、ABAC、ReBAC 等权限模型对比与选型 | [permission-models.md](./permission-models.md) |
|  JWT  | JWT Token 结构、签名算法、最佳实践 | [jwt-basics.md](./jwt-basics.md) |
|  MFA/TOTP  | 多因素认证与 TOTP 动态验证码原理 | [mfa-totp.md](./mfa-totp.md) |
|  多租户架构  | SaaS 多租户数据隔离方案对比 | [multi-tenancy-architecture.md](./multi-tenancy-architecture.md) |
|  双 Token 方案  | Access Token + Refresh Token 设计详解 | [jwt-basics.md](./jwt-basics.md) |

## 推荐阅读顺序

1. 先阅读 [权限模型](./permission-models.md)，理解 IAM 的核心设计基础
2. 再阅读 [OAuth 2.0](./oauth-2.0-basics.md)，了解第三方登录协议
3. 接着阅读 [JWT](./jwt-basics.md)，理解 Token 认证机制
4. 根据需求选读 [MFA/TOTP](./mfa-totp.md) 和 [多租户架构](./multi-tenancy-architecture.md)

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
