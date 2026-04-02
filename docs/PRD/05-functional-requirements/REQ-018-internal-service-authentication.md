# REQ-018 内部服务认证

| 项目 | 内容 |
|------|------|
| **优先级** | P0 |
| **估时** | 6 人天 |
| **关联用户故事** | US-015 |

**背景：** IAM 不能只覆盖用户登录，还需要为平台自有的 OA、CRM、ERP、调度任务、后台服务等内部系统提供统一认证能力，避免各业务系统自行维护静态密钥和重复鉴权逻辑。

**目标：**

- 支持平台统一注册内部客户端
- 支持通过 `AK/SK` 换取短期 JWT Access Token
- 支持客户端级别的 `scope` 授权
- 支持凭证轮换、禁用、审计和限流
- 复用现有 API 网关 JWT 校验链路

**功能描述：**

### 1. 客户端模型

1. 内部系统以平台级 `Client` 身份接入 IAM
2. 每个客户端包含 `client_id`、名称、状态、允许的 `scopes`、Token TTL 等属性
3. 客户端不等同于租户级 `Application`，不参与租户数据隔离建模

### 2. 凭证管理

1. 平台管理员可为客户端创建一组 `AK/SK`
2. `SK` 明文仅在创建或轮换时展示一次
3. 服务端仅保存 `SK` 哈希值和必要的提示信息
4. 支持凭证禁用、过期和轮换

### 3. Token 申请

1. 内部客户端调用 `POST /api/v1/auth/token`
2. 请求携带 `grant_type=client_credentials`
3. 服务端校验 `AK/SK`、客户端状态和允许的 `scope`
4. 校验通过后签发短期 JWT Access Token
5. 客户端 Token 默认有效期 10 分钟，不提供 Refresh Token

### 4. 客户端授权

1. 客户端权限使用独立 `scope` 模型，不复用用户 RBAC 角色
2. `scope` 用于描述可访问的 API 或资源范围，例如：
   - `user:read`
   - `user:write`
   - `tenant:read`
   - `audit:read`
3. 网关根据 `subject_type=client` 和 `scopes` 执行授权校验

### 5. 预留扩展

1. 预留代表用户调用的代理字段，例如 `act`
2. 首版本不支持正式的 on-behalf-of 流程
3. 当客户端尝试申请代理用户能力时，返回明确错误

**客户端 Token Claim：**

| Claim | 说明 | 示例 |
|------|------|------|
| `sub` | 客户端主体 ID | `crm-service` |
| `subject_type` | 主体类型 | `client` |
| `client_id` | 客户端 ID | `crm-service` |
| `scopes` | 客户端权限范围 | `["user:read", "tenant:read"]` |
| `aud` | 目标服务 | `api-gateway` |
| `jti` | Token 唯一标识 | `ct_xxxxxxxxx` |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| `AK/SK` 错误 | 返回 401，错误码 `INVALID_CLIENT_CREDENTIALS` |
| 客户端已禁用 | 返回 403，错误码 `CLIENT_DISABLED` |
| 凭证已过期 | 返回 401，错误码 `CLIENT_CREDENTIAL_EXPIRED` |
| 申请未授权的 `scope` | 返回 403，错误码 `INSUFFICIENT_SCOPE` |
| 尝试代理用户调用 | 返回 400，错误码 `UNSUPPORTED_DELEGATION` |

**安全策略：**

| 策略 | 说明 |
|------|------|
| **短期 Token** | 客户端 Access Token 默认 10 分钟 |
| **单次展示** | `SK` 明文只展示一次 |
| **哈希存储** | 服务端不保存明文 `SK` |
| **轮换机制** | 支持主动轮换，旧凭证立即失效 |
| **最小权限** | 客户端默认不授予任何 `scope` |
| **统一审计** | 记录凭证创建、轮换、禁用、取 Token 和调用行为 |

**API 接口：**

```
POST /api/v1/auth/token                      # 客户端凭证换取 Token
POST /api/v1/clients                         # 创建客户端
GET  /api/v1/clients                         # 客户端列表
GET  /api/v1/clients/:id                     # 客户端详情
POST /api/v1/clients/:id/credentials         # 创建 AK/SK
POST /api/v1/clients/:id/credentials/rotate  # 轮换 AK/SK
POST /api/v1/clients/:id/scopes              # 配置客户端 scopes
POST /api/v1/clients/:id/disable             # 禁用客户端
```

**数据库设计：**

**客户端表（auth_clients）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 1001 |
| client_id | VARCHAR(64) | 是 | 客户端标识，全局唯一 | crm-service |
| client_name | VARCHAR(128) | 是 | 客户端名称 | CRM Internal Service |
| client_type | VARCHAR(20) | 是 | 客户端类型 | internal |
| status | VARCHAR(20) | 是 | 状态 | active/disabled |
| access_token_ttl_sec | INT | 是 | Access Token 有效期（秒） | 600 |
| created_at | DATETIME | - | 创建时间 | 2026-04-02 10:00:00 |
| updated_at | DATETIME | - | 更新时间 | 2026-04-02 10:00:00 |

**索引**：`uk_client_id` (client_id)

---

**客户端凭证表（auth_client_credentials）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 2001 |
| client_id | VARCHAR(64) | 是 | 客户端 ID | crm-service |
| access_key_id | VARCHAR(64) | 是 | AK 标识 | ak_xxxxx |
| secret_hash | VARCHAR(255) | 是 | SK 哈希值 | hash_xxxxx |
| secret_hint | VARCHAR(16) | 否 | SK 提示信息 | `...8F3A` |
| status | VARCHAR(20) | 是 | 状态 | active/disabled/expired |
| expires_at | DATETIME | 否 | 过期时间 | 2026-10-01 00:00:00 |
| last_used_at | DATETIME | 否 | 最近使用时间 | 2026-04-02 11:00:00 |
| rotated_at | DATETIME | 否 | 轮换时间 | 2026-04-02 10:30:00 |
| created_at | DATETIME | - | 创建时间 | 2026-04-02 10:00:00 |

**索引**：`uk_access_key_id` (access_key_id)、`idx_client_status` (client_id, status)

**验收标准：**

- [ ] 平台管理员可创建和查看内部客户端
- [ ] 客户端可通过 `AK/SK` 成功换取短期 JWT Access Token
- [ ] 错误、过期、禁用凭证会被正确拒绝
- [ ] 客户端 `scope` 授权正确生效
- [ ] 轮换后旧凭证立即失效
- [ ] 取 Token 和关键调用行为被正确审计
