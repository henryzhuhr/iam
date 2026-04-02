# REQ-012 Token 管理

| 项目 | 内容 |
|------|------|
| **优先级** | P0 |
| **估时** | 4 人天 |
| **关联用户故事** | US-011、US-012 |

**背景：** Token 是 IAM 统一认证系统的核心凭证，需要同时服务于用户认证和内部服务认证，具备完整的生成、验证、刷新、撤销和审计能力，确保安全性、性能和可扩展性。

**目标：**

- 支持用户主体和客户端主体两类 Token
- 支持 JWT Access Token 和 Refresh Token 双令牌机制
- Access Token 有效期可配置（默认 30 分钟）
- Refresh Token 有效期可配置（默认 7 天）
- 支持 Token 黑名单机制，实现主动撤销
- 支持多设备登录控制策略
- 支持内部客户端短期 Token（默认 10 分钟）
- Token 验证性能 < 10ms

**功能描述：**

### 1. Token 生成

1. 用户登录成功后，系统生成用户 Access Token 和 Refresh Token
2. 内部服务客户端通过 `AK/SK` 申请客户端 Access Token，不颁发 Refresh Token
3. Access Token 采用 JWT 格式，包含以下通用 Claims：
   - `sub`: 主体 ID（用户 ID 或客户端 ID）
   - `subject_type`: 主体类型（`user` / `client`）
   - `iat`: 签发时间
   - `exp`: 过期时间
   - `jti`: Token 唯一标识
   - `aud`: 目标服务标识
4. 用户 Token 可附带以下 Claims：
   - `tid`: 租户 ID
   - `roles`: 角色列表
   - `apps`: 可访问应用列表
   - `dev`: 设备指纹（可选）
5. 客户端 Token 可附带以下 Claims：
   - `client_id`: 客户端 ID
   - `scopes`: 机器权限范围
   - `tenant_id`: 目标租户（预留，可选）
   - `act`: 代理用户上下文（预留，可选）
6. 用户 Refresh Token 采用不透明随机字符串，存储于 Redis，并与设备指纹绑定

### 2. Token 验证

1. 拦截器校验 `Authorization: Bearer <token>` 头
2. 验证 JWT 签名有效性
3. 验证 Token 是否过期
4. 查询 Token 黑名单，确认未被撤销
5. 根据 `subject_type` 分流校验用户主体或客户端主体
6. 用户 Token 验证 `tenant_id` 与当前请求租户匹配
7. 客户端 Token 验证 `scopes` 是否满足接口访问要求

### 3. Token 刷新

1. 客户端使用 Refresh Token 请求 `/api/v1/auth/refresh`
2. 服务端验证 Refresh Token 有效性（未过期、未撤销）
3. 可选：刷新后使旧 Refresh Token 失效（滚动刷新）
4. 返回新的 Access Token 和 Refresh Token

### 4. 客户端 Token 申请

1. 内部客户端使用 `AK/SK` 调用 `/api/v1/auth/token`
2. 服务端验证客户端状态、凭证有效性和允许的 `scope`
3. 返回短期 JWT Access Token（默认 10 分钟）
4. 不返回 Refresh Token，客户端在过期前重新申请

### 5. Token 撤销

1. 用户登出时，将 Access Token 加入黑名单
2. 删除 Redis 中的 Refresh Token
3. 批量撤销：管理员可撤销指定用户的所有 Token
4. 客户端凭证轮换或禁用时，使对应客户端后续申请 Token 立即失效
5. Token 黑名单 TTL 设置为 Access Token 剩余有效期

### 6. 多设备登录控制

支持三种策略，可配置：

| 策略 | 说明 | 适用场景 |
|------|------|----------|
| `ALLOW_MULTI` | 允许同一账号多设备同时在线 | 默认策略，适用于大多数场景 |
| `SINGLE_DEVICE` | 新登录会使旧设备 Token 失效 | 高安全要求场景 |
| `SINGLE_DEVICE_PER_TYPE` | 同类型设备只允许一个在线（Web/App/移动端） | 平衡安全与体验 |

### 7. Token 配置项

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `access_token_ttl` | 30m | Access Token 有效期 |
| `refresh_token_ttl` | 7d | Refresh Token 有效期 |
| `client_access_token_ttl` | 10m | 内部客户端 Access Token 有效期 |
| `multi_device_policy` | ALLOW_MULTI | 多设备登录策略 |
| `refresh_token_rolling` | false | 是否启用 Refresh Token 滚动刷新 |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| Token 格式错误 | 返回 401，错误码 `INVALID_TOKEN` |
| Token 签名无效 | 返回 401，错误码 `INVALID_SIGNATURE` |
| Token 已过期 | 返回 401，错误码 `TOKEN_EXPIRED` |
| Token 在黑名单中 | 返回 401，错误码 `TOKEN_REVOKED` |
| Refresh Token 无效 | 返回 401，错误码 `INVALID_REFRESH_TOKEN` |
| 多设备策略冲突 | 撤销旧 Token，返回新 Token，记录日志 |
| 客户端凭证无效 | 返回 401，错误码 `INVALID_CLIENT_CREDENTIALS` |
| 客户端 Scope 不足 | 返回 403，错误码 `INSUFFICIENT_SCOPE` |

**安全策略：**

| 策略 | 说明 |
|------|------|
| **签名算法** | 使用 HS256 或 RS256 |
| **密钥轮换** | 支持密钥版本，平滑轮换 |
| **最短有效期** | Access Token 最短 5 分钟 |
| **最长有效期** | Refresh Token 最长 30 天 |
| **客户端 Token** | 默认 10 分钟，无 Refresh Token |
| **黑名单清理** | 过期 Token 自动从黑名单移除 |

**API 接口：**

```
POST   /api/v1/auth/refresh      # 刷新 Token
POST   /api/v1/auth/logout       # 登出（撤销 Token）
POST   /api/v1/auth/logout/all   # 撤销所有设备 Token
POST   /api/v1/auth/token        # 客户端凭证换取 Token
GET    /api/v1/auth/sessions     # 获取活跃会话列表
DELETE /api/v1/auth/sessions/:id # 撤销指定会话
```

**数据库设计：**

**Token 黑名单表（token_blacklist）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 1001 |
| jti | VARCHAR(64) | 是 | Token 唯一标识 | at_xxxxxxxxxxxxx |
| subject_type | VARCHAR(20) | 是 | 主体类型 | user/client |
| user_id | BIGINT | 否 | 用户 ID | 2001 |
| client_id | VARCHAR(64) | 否 | 客户端 ID | crm-service |
| expire_at | DATETIME | 是 | 过期时间 | 2026-03-28 11:00:00 |
| created_at | DATETIME | - | 创建时间 | 2026-03-28 10:30:00 |

**索引**：`idx_jti` (jti)、`idx_expire` (expire_at) — 便于定时清理

---

**活跃会话表（active_sessions）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 3001 |
| user_id | BIGINT | 是 | 用户 ID | 2001 |
| refresh_token_hash | VARCHAR(64) | 是 | Refresh Token 哈希 | hash_xxxxxxxxxxxxx |
| device_type | VARCHAR(20) | 否 | 设备类型 | web/ios/android |
| device_fingerprint | VARCHAR(64) | 否 | 设备指纹 | fp_xxxxxxxxxxxxx |
| ip_address | VARCHAR(45) | 否 | IP 地址 | 192.168.1.100 |
| user_agent | VARCHAR(255) | 否 | 用户代理 | Mozilla/5.0... |
| last_active_at | DATETIME | 否 | 最后活跃时间 | 2026-03-28 10:30:00 |
| expires_at | DATETIME | 否 | 过期时间 | 2026-04-04 10:30:00 |
| created_at | DATETIME | - | 创建时间 | 2026-03-28 10:30:00 |

**索引**：`idx_user` (user_id)、`idx_refresh` (refresh_token_hash)

**验收标准：**

- [ ] Access Token 和 Refresh Token 正确生成
- [ ] 过期 Token 被正确拒绝
- [ ] 黑名单 Token 被正确拒绝
- [ ] Refresh Token 可成功刷新 Access Token
- [ ] 多设备策略正确生效
- [ ] 登出后 Token 被撤销
- [ ] 客户端可通过 `AK/SK` 获取短期 Access Token
- [ ] 客户端 Token 的 `subject_type` 和 `scopes` 正确下发
- [ ] 活跃会话列表可查询和管理
- [ ] Token 验证性能 < 10ms
