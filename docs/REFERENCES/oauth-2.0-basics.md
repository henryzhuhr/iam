# OAuth 2.0 协议基础

> 最后更新：2026-03-25
> 适用场景：IAM 第三方登录集成

## 1. OAuth 2.0 是什么

OAuth 2.0 是一个**授权协议**，允许第三方应用有限度地访问用户在资源服务器上的资源，而无需暴露用户的凭据。

### 1.1 核心概念

| 角色 | 说明 | 示例 |
|------|------|------|
| **Resource Owner** (资源所有者) | 用户，拥有资源的所有权 | 使用 GitHub 登录的用户 |
| **Client** (客户端) | 请求访问用户资源的应用 | 我们的 IAM 系统 |
| **Resource Server** (资源服务器) | 存储用户资源的服务器 | GitHub API 服务器 |
| **Authorization Server** (授权服务器) | 颁发 Token 的服务器 | GitHub OAuth 服务器 |

### 1.2 为什么需要 OAuth

**传统方式的问题：**
- 第三方应用需要存储用户用户名密码
- 密码泄露风险高
- 用户无法限制第三方权限范围

**OAuth 的优势：**
- 第三方不接触用户密码
- 用户可以授权特定范围（Scopes）
- 可以随时撤销授权

---

## 2. OAuth 2.0 授权模式

OAuth 2.0 定义了 4 种授权模式，IAM 主要使用 **Authorization Code** 模式。

### 2.1 Authorization Code（授权码模式） ⭐

**适用场景：** 服务端 Web 应用（有后端服务器）

**流程：**

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│  用户   │         │  IAM(Client) │         │GitHub(AuthZ)│
└────┬────┘         └──────┬───────┘         └──────┬──────┘
     │                     │                         │
     │  1. 点击 GitHub 登录  │                         │
     ├────────────────────>│                         │
     │                     │                         │
     │                     │  2. 重定向到 GitHub     │
     │<────────────────────┼─────────────────────────┤
     │                     │                         │
     │  3. 访问 GitHub 授权页│                         │
     ├──────────────────────────────────────────────>│
     │                     │                         │
     │  4. 用户同意授权     │                         │
     ├──────────────────────────────────────────────>│
     │                     │                         │
     │  5. 回调 + code     │                         │
     │<──────────────────────────────────────────────┤
     │                     │                         │
     │  6. 传递 code 给 IAM  │                         │
     ├────────────────────>│                         │
     │                     │                         │
     │                     │  7. code 换 access_token│
     │                     ├────────────────────────>│
     │                     │                         │
     │                     │  8. 返回 access_token   │
     │                     ├────────────────────────<│
     │                     │                         │
     │                     │  9. 用 token 获取用户信息│
     │                     ├────────────────────────>│
     │                     │                         │
     │                     │  10. 返回用户信息       │
     │                     ├────────────────────────<│
     │                     │                         │
     │  11. 登录成功       │                         │
     │<────────────────────┤                         │
     │                     │                         │
```

**步骤说明：**

1. 用户在 IAM 点击「使用 GitHub 登录」
2. IAM 重定向用户到 GitHub 授权页，带上 `client_id` 和 `redirect_uri`
3. 用户在 GitHub 页面输入用户名密码（IAM 不接触）
4. 用户同意授权
5. GitHub 重定向回 IAM 的 `redirect_uri`，带上 `code`（授权码）
6. 前端将 `code` 传给 IAM 后端
7. IAM 后端用 `code` + `client_secret` 换取 `access_token`
8. GitHub 返回 `access_token`
9. IAM 用 `access_token` 调用 GitHub API 获取用户邮箱等信息
10. GitHub 返回用户信息
11. IAM 匹配/创建本地用户，颁发自己的 JWT Token，登录完成

**为什么有 code 和 token 两次交换？**
- `code` 是一次性的，且只能通过 HTTPS 传输
- `token` 不直接暴露在浏览器 URL 中，更安全
- 后端存储 `client_secret`，前端无法获取

---

### 2.2 Implicit Grant（隐式模式）

**适用场景：** 纯前端应用（已不推荐，被 PKCE 替代）

**特点：**
- 没有 `code` 交换环节，直接返回 `access_token`
- Token 暴露在 URL 中，有安全风险
- 现代 OAuth 2.1 已弃用

---

### 2.3 Resource Owner Password Credentials（密码模式）

**适用场景：** 高度信任的第一方应用

**特点：**
- 用户直接输入用户名密码给客户端
- 客户端用密码换取 token
- **不适用于第三方**，仅用于自己的官方应用

---

### 2.4 Client Credentials（客户端模式）

**适用场景：** 机器对机器（M2M）通信

**特点：**
- 没有用户参与
- 应用用自己的身份获取 token
- 适用于后台服务、定时任务等

---

## 3. OAuth 2.0 关键参数

### 3.1 注册应用时获取

| 参数 | 说明 | 保密性 |
|------|------|--------|
| `client_id` | 应用的唯一标识 | 公开 |
| `client_secret` | 应用密钥 | **严格保密** |
| `redirect_uri` | 授权回调地址 | 公开，但必须精确匹配 |

### 3.2 授权请求参数

```
GET https://github.com/login/oauth/authorize?
  client_id=YOUR_CLIENT_ID&
  redirect_uri=https://yourdomain.com/auth/github/callback&
  scope=read:user,user:email&
  state=random_state_string&
  response_type=code
```

| 参数 | 必填 | 说明 |
|------|------|------|
| `client_id` | 是 | 应用 ID |
| `redirect_uri` | 推荐 | 回调地址，必须与注册一致 |
| `scope` | 否 | 请求的权限范围 |
| `state` | 推荐 | 防 CSRF 攻击的随机字符串 |
| `response_type` | 是 | 固定为 `code` |

### 3.3 Token 交换参数

```
POST https://github.com/login/oauth/access_token
Content-Type: application/json

{
  "client_id": "YOUR_CLIENT_ID",
  "client_secret": "YOUR_CLIENT_SECRET",
  "code": "AUTHORIZATION_CODE",
  "redirect_uri": "https://yourdomain.com/auth/github/callback"
}
```

---

## 4. GitHub OAuth 实战

### 4.1 注册 GitHub OAuth App

1. 访问 https://github.com/settings/developers
2. 点击 "OAuth Apps" → "New OAuth App"
3. 填写应用信息：
   - Application name: IAM System
   - Homepage URL: https://iam.yourdomain.com
   - Authorization callback URL: https://iam.yourdomain.com/auth/github/callback

### 4.2 获取的凭据

```
Client ID: Iv1.xxxxxxxxxxxxxxxx
Client Secret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 4.3 Scope 说明（GitHub）

| Scope | 说明 |
|-------|------|
| `read:user` | 读取用户基本信息 |
| `user:email` | 读取用户邮箱 |
| `repo` | 访问仓库（不需要） |
| `admin:org` | 访问组织（不需要） |

IAM 推荐：`scope=read:user,user:email`

---

## 5. 安全注意事项

### 5.1 必须做的

| 措施 | 原因 |
|------|------|
| 使用 HTTPS | 防止中间人攻击 |
| 验证 `state` 参数 | 防止 CSRF 攻击 |
| 后端存储 `client_secret` | 不能暴露给前端 |
| `redirect_uri` 精确匹配 | 防止回调劫持 |
| `code` 一次性使用 | 防止重放攻击 |

### 5.2 建议做的

| 措施 | 原因 |
|------|------|
| `state` 使用加密随机数 | 不可预测 |
| 设置合理的 `scope` | 最小权限原则 |
| 记录 OAuth 登录日志 | 审计追溯 |
| Token 加密存储 | 防止泄露 |

---

## 6. 常见问题

### Q1: OAuth 和 OIDC 有什么区别？

- **OAuth 2.0** 是**授权**协议，用于访问资源
- **OIDC (OpenID Connect)** 是**身份认证**协议，基于 OAuth 2.0
- OIDC 增加了 `id_token`，包含用户身份信息
- GitHub OAuth 使用标准 OAuth 2.0，不支持 OIDC

### Q2: 用户没有 GitHub 账号怎么办？

OAuth 流程会引导用户去 GitHub 注册，注册完成后再回调。

### Q3: 如何解绑第三方账号？

在用户设置中提供解绑功能，删除 `user_third_party_accounts` 表中的绑定记录。

### Q4: 一个本地账号可以绑定多个第三方账号吗？

可以，设计上支持一对多绑定（一个用户绑定 GitHub、Google 等多个账号）。

---

## 7. 相关需求文档

- [REQ-011 OAuth2 第三方登录](../05-functional-requirements/REQ-011-oauth2-third-party-login.md)
- [REQ-012 Token 管理](../05-functional-requirements/REQ-012-token-management.md)

---

## 8. 参考链接

- RFC 6749: https://tools.ietf.org/html/rfc6749
- GitHub OAuth 文档：https://docs.github.com/en/developers/apps/building-oauth-apps
- OAuth.net: https://oauth.net/2/
