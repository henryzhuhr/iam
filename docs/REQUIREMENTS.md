# 身份认证与访问管理系统 需求分析文档

> 身份认证与访问管理系统 (Identity and Access Management)
> 版本：v0.1.2-draft

---

## 修订历史

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v0.1.0-draft | 2026-03-17 | - | 初始版本 |
| v0.1.1-draft | 2026-03-17 | - | 新增核心概念说明、完善数据库表结构设计 |
| v0.1.2-draft | 2026-03-18 | - | 新增 OAuth2 第三方登录需求 |

## 1. 项目概述

### 1.1 产品定位

为 SaaS 多租户应用提供完整的身份认证与访问管理能力，支持多租户隔离、用户管理、认证授权、权限控制等核心功能。

### 1.2 目标用户

- **平台运营管理员**：管理租户的开通、配置、监控
- **租户管理员**：管理企业内的用户、角色、权限
- **终端用户**：使用 SaaS 系统的普通员工

### 1.3 核心概念

| 概念 | 说明 |
|------|------|
| **租户 (Tenant)** | SaaS 平台中独立的企业客户，数据相互隔离。租户是数据隔离的基本单位 |
| **用户 (User)** | 属于某个租户的具体个人，在租户内进行角色分配和权限管理 |
| **JWT Token** | JSON Web Token，包含 Header、Payload、Signature，用于身份认证 |
| **Access Token** | 短期访问令牌（15-30 分钟），用于 API 请求认证 |
| **Refresh Token** | 长期刷新令牌（7-30 天），用于获取新的 Access Token |

### 1.4 技术栈

- **后端**: Golang + go-zero 框架
- **数据库**: MySQL
- **缓存**: Redis
- **消息队列**: Kafka
- **容器化**: Docker + Docker Compose

### 1.5 Token 方案选型

IAM 系统采用 **JWT + 双 Token** 方案：

| 特性 | Access Token | Refresh Token |
|------|--------------|---------------|
| **类型** | JWT | 不透明 Token |
| **有效期** | 15-30 分钟 | 7-30 天 |
| **存储** | 客户端（内存/LocalStorage） | 服务端（Redis）+ HttpOnly Cookie |
| **用途** | API 请求认证 | 刷新 Access Token |
| **撤销** | 加入黑名单 | 从 Redis 删除 |

**选型理由：**

1. **无状态认证** - JWT 自包含用户信息，减少数据库查询
2. **高性能** - 签名验证快，适合高并发场景
3. **安全性** - Access Token 短期有效，泄露风险低
4. **用户体验** - Refresh Token 长期有效，无需频繁登录

---

## 2. 核心功能模块

### 2.1 用户管理 (User Management)

#### 2.1.1 用户基础管理

- [ ] 用户 CRUD（创建、查询、更新、删除）
- [ ] 用户状态管理（启用/禁用）
- [ ] 用户批量导入/导出
- [ ] 用户头像上传

#### 2.1.2 用户组管理

- [ ] 用户组 CRUD
- [ ] 用户组层级结构（树形）
- [ ] 用户组成员管理
- [ ] 用户组权限继承

#### 2.1.3 租户管理

- [ ] 租户 CRUD
- [ ] 租户配额管理（用户数、角色数等）
- [ ] 租户状态管理（激活/冻结）
- [ ] 租户管理员分配

### 2.2 认证管理 (Authentication)

#### 2.2.1 基础认证

- [ ] 用户名密码登录
- [ ] 邮箱验证码登录
- [ ] 手机号验证码登录
- [ ] 登出功能

#### 2.2.2 第三方登录 (OAuth2)

- [ ] GitHub OAuth2 登录
- [ ] Google OAuth2 登录（预留）
- [ ] 钉钉 OAuth2 登录（预留）
- [ ] 第三方账号绑定/解绑
- [ ] OAuth2 回调处理

#### 2.2.3 Token 管理

- [ ] JWT Token 生成/验证/刷新
- [ ] Token 黑名单
- [ ] 多设备登录控制
- [ ] Token 过期策略

#### 2.2.4 多因素认证 (MFA)

- [ ] TOTP 动态验证码（Google Authenticator）
- [ ] 短信验证码
- [ ] 邮箱验证码
- [ ] MFA 启用/禁用

#### 2.2.5 密码策略

- [ ] 密码强度校验
- [ ] 密码过期策略
- [ ] 密码历史记录
- [ ] 密码重置（忘记密码）

### 2.3 权限管理 (Authorization)

#### 2.3.1 RBAC 模型

- [ ] 角色 CRUD
- [ ] 权限点 CRUD
- [ ] 角色 - 权限关联
- [ ] 用户 - 角色关联

#### 2.3.2 资源授权

- [ ] API 资源定义
- [ ] 菜单资源定义
- [ ] 数据范围授权
- [ ] 资源权限校验中间件

#### 2.3.3 权限策略

- [ ] 允许/拒绝策略
- [ ] 策略优先级
- [ ] 条件策略（时间、IP 等）
- [ ] 权限继承与覆盖

### 2.4 审计日志 (Audit Log)

#### 2.4.1 操作日志

- [ ] 用户操作记录
- [ ] 管理员操作记录
- [ ] 操作详情存储
- [ ] 操作日志查询

#### 2.4.2 登录日志

- [ ] 登录成功/失败记录
- [ ] 登录 IP/设备/地理位置
- [ ] 异常登录检测
- [ ] 登录历史查询

---

## 3. 非功能性需求

### 3.1 性能要求

- [ ] API 响应时间 < 100ms (P95)
- [ ] 支持 1000+ QPS
- [ ] 支持 100 万 + 用户

### 3.2 安全要求

- [ ] 密码加密存储（bcrypt/argon2）
- [ ] 通信加密（HTTPS/TLS）
- [ ] SQL 注入防护
- [ ] XSS 防护
- [ ] CSRF 防护
- [ ] 敏感操作二次验证

### 3.3 可用性要求

- [ ] 99.9% 可用性
- [ ] 支持水平扩展
- [ ] 数据库主从复制
- [ ] Redis 集群
- [ ] 服务优雅关闭

### 3.4 多租户隔离

- [ ] 数据隔离（tenant_id）
- [ ] 资源隔离
- [ ] 配额隔离

---

## 4. API 接口设计

### 4.1 接口规范

- RESTful API 风格
- 统一响应格式
- 统一错误码
- JWT 认证
- API 版本管理（`/api/v1/`）

### 4.2 接口列表

#### 用户管理

```
POST   /api/v1/users           # 创建用户
GET    /api/v1/users           # 用户列表
GET    /api/v1/users/:id       # 用户详情
PUT    /api/v1/users/:id       # 更新用户
DELETE /api/v1/users/:id       # 删除用户
POST   /api/v1/users/:id/enable  # 启用用户
POST   /api/v1/users/:id/disable # 禁用用户
```

#### 认证管理

```
POST   /api/v1/auth/login      # 登录
POST   /api/v1/auth/logout     # 登出
POST   /api/v1/auth/refresh    # 刷新 Token
POST   /api/v1/auth/password/reset  # 重置密码
POST   /api/v1/mfa/bind        # 绑定 MFA
POST   /api/v1/mfa/verify      # 验证 MFA
```

#### 第三方登录 (OAuth2)

```
GET    /api/v1/auth/github     # 发起 GitHub 登录
GET    /api/v1/auth/github/callback  # GitHub 回调处理
POST   /api/v1/auth/bind       # 绑定第三方账号
POST   /api/v1/auth/unbind     # 解绑第三方账号
GET    /api/v1/auth/providers  # 获取可用的第三方登录方式
```

#### 权限管理

```
POST   /api/v1/roles           # 创建角色
GET    /api/v1/roles           # 角色列表
PUT    /api/v1/roles/:id       # 更新角色
DELETE /api/v1/roles/:id       # 删除角色
POST   /api/v1/roles/:id/permissions  # 分配权限
GET    /api/v1/permissions     # 权限列表
```

#### 租户管理

```
POST   /api/v1/tenants         # 创建租户
GET    /api/v1/tenants         # 租户列表
PUT    /api/v1/tenants/:id     # 更新租户
DELETE /api/v1/tenants/:id     # 删除租户
```

---

## 5. 数据库设计

### 5.1 核心表

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `tenants` | 租户表 | id, name, status, quota, created_at |
| `users` | 用户表 | id, tenant_id, email, password_hash, status, created_at |
| `user_groups` | 用户组表 | id, tenant_id, name, parent_id, level |
| `user_group_members` | 用户组成员关系表 | id, tenant_id, group_id, user_id |
| `roles` | 角色表 | id, tenant_id, name, description, is_builtin |
| `permissions` | 权限表 | id, name, resource_type, resource_value, action |
| `role_permissions` | 角色权限关系表 | id, role_id, permission_id |
| `user_roles` | 用户角色关系表 | id, user_id, role_id |
| `auth_tokens` | Token 表 | id, user_id, token_hash, type, expires_at, revoked |
| `user_third_party_accounts` | 第三方账号绑定表 | id, user_id, provider, open_id, access_token, created_at |
| `operation_logs` | 操作日志表 | id, tenant_id, user_id, action, resource, result, ip, created_at |
| `login_logs` | 登录日志表 | id, tenant_id, user_id, result, ip, device, user_agent, created_at |

### 5.2 表结构设计

#### tenants - 租户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| name | VARCHAR(100) | 租户名称 |
| status | TINYINT | 状态（1-激活，0-冻结） |
| max_users | INT | 最大用户数配额 |
| max_roles | INT | 最大角色数配额 |
| expires_at | DATETIME | 租户过期时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### users - 用户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| tenant_id | BIGINT | 租户 ID（外键） |
| email | VARCHAR(100) | 邮箱（租户内唯一） |
| password_hash | VARCHAR(255) | 密码哈希（bcrypt） |
| status | TINYINT | 状态（1-启用，0-禁用） |
| mfa_enabled | BOOLEAN | 是否启用 MFA |
| mfa_secret | VARCHAR(255) | MFA 密钥（加密存储） |
| last_login_at | DATETIME | 最后登录时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

**索引：**
- `idx_tenant_email`: (tenant_id, email) - 租户内邮箱唯一
- `idx_status`: (status) - 状态查询

#### user_third_party_accounts - 第三方账号绑定表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| user_id | BIGINT | 用户 ID（外键） |
| provider | VARCHAR(20) | 第三方提供商（github/google/dingtalk） |
| open_id | VARCHAR(100) | 第三方用户唯一标识 |
| union_id | VARCHAR(100) | 第三方统一标识（可选，用于跨应用识别） |
| access_token | VARCHAR(500) | 第三方 Access Token（加密存储） |
| refresh_token | VARCHAR(500) | 第三方 Refresh Token（加密存储） |
| token_expires_at | DATETIME | Token 过期时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

**索引：**
- `idx_user_provider`: (user_id, provider) - 用户 + 提供商唯一
- `idx_open_id`: (provider, open_id) - 第三方用户唯一

#### roles - 角色表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| tenant_id | BIGINT | 租户 ID（外键） |
| name | VARCHAR(50) | 角色名称 |
| description | VARCHAR(255) | 角色描述 |
| is_builtin | BOOLEAN | 是否预置角色（预置角色不可删除） |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### permissions - 权限表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| name | VARCHAR(50) | 权限名称 |
| resource_type | VARCHAR(20) | 资源类型（API/MENU/DATA） |
| resource_value | VARCHAR(255) | 资源标识（如 GET /api/users） |
| action | VARCHAR(20) | 操作类型（create/read/update/delete） |
| description | VARCHAR(255) | 权限描述 |

### 5.3 ER 关系图

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────┐
│  tenants    │ 1   N │     users       │ N   1 │  roles      │
│─────────────│───────│─────────────────│───────│─────────────│
│  id         │       │  id             │       │  id         │
│  name       │       │  tenant_id (FK) │       │  tenant_id  │
│  status     │       │  email          │       │  name       │
└─────────────┘       │  password_hash  │       └─────────────┘
                      │  status         │              │
                      └─────────────────┘              │ N
                              │                        │
                              │ N                      │
                              ▼                        ▼
                      ┌─────────────────┐       ┌─────────────────┐
                      │  user_roles     │       │  role_permissions│
                      │─────────────────│       │─────────────────│
                      │  user_id (FK)   │       │  role_id (FK)   │
                      │  role_id (FK)   │       │  permission_id  │
                      └─────────────────┘       └─────────────────┘

## 6. 项目里程碑

### Phase 1 - 基础框架 (Week 1-2)

- [ ] 项目脚手架搭建
- [ ] 数据库设计
- [ ] 用户管理 CRUD
- [ ] 基础登录认证

### Phase 2 - 核心功能 (Week 3-4)

- [ ] RBAC 权限模型
- [ ] JWT Token 管理
- [ ] 租户管理
- [ ] API 权限校验中间件

### Phase 3 - 高级功能 (Week 5-6)

- [ ] 多因素认证 MFA
- [ ] 审计日志
- [ ] 密码策略
- [ ] 用户组管理

### Phase 4 - 完善优化 (Week 7-8)

- [ ] OAuth2 第三方登录（GitHub）
- [ ] 性能优化
- [ ] 安全加固
- [ ] 文档完善
- [ ] 单元测试

---

## 7. 待确认事项

- [x] 是否需要支持第三方登录（GitHub、Google、微信等）？ → **已确认：需要，优先支持 GitHub，P2 需求**
- [ ] 是否需要支持 LDAP/AD 集成？
- [ ] 是否需要支持 OIDC/SAML 协议？
- [ ] 是否需要提供管理后台 UI？
- [ ] 多租户数据隔离级别（数据库隔离 / Schema 隔离 / 数据行隔离）？

---

## 修订历史

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v0.1.0-draft | 2026-03-17 | - | 初始草稿 |
