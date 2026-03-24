# REQ-016 API 限流和配额管理

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 3 人天 |
| **关联用户故事** | US-022、US-023 |

**背景：** 多租户场景下，需要防止单个租户或用户过度占用系统资源，保障整体系统稳定性和公平性，同时满足租户配额管理需求。

**目标：**

- 支持按租户/API/用户维度的限流
- 支持租户配额管理（用户数、角色数等）
- 限流策略可配置
- 配额超限后优雅降级
- 支持配额预警通知

**功能描述：**

### 1. API 限流策略

支持以下限流维度：

| 维度 | 说明 | 示例 |
|------|------|------|
| 全局限流 | 整个系统的总 QPS 限制 | 10000 QPS |
| 租户限流 | 单个租户的 QPS 限制 | 1000 QPS |
| 用户限流 | 单个用户的 QPS 限制 | 100 QPS |
| IP 限流 | 单个 IP 的 QPS 限制 | 50 QPS |
| API 限流 | 单个 API 接口的 QPS 限制 | 500 QPS |

限流算法：
- **令牌桶算法**：适用于平滑限流
- **滑动窗口**：适用于精确限流
- **固定窗口**：适用于简单场景

### 2. 限流配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `rate_limit_global_qps` | 10000 | 全局 QPS 上限 |
| `rate_limit_tenant_qps` | 1000 | 单租户 QPS 上限 |
| `rate_limit_user_qps` | 100 | 单用户 QPS 上限 |
| `rate_limit_ip_qps` | 50 | 单 IP QPS 上限 |
| `rate_limit_enabled` | true | 是否启用限流 |
| `rate_limit_mode` | STRICT | STRICT/SLACK/BURST |

限流模式：
- `STRICT`: 严格模式，超限立即拒绝
- `SLACK`: 宽松模式，允许短暂超限
- `BURST`: 突发模式，允许配置突发流量

### 3. 租户配额管理

支持以下配额类型：

| 配额类型 | 说明 | 默认值 |
|----------|------|--------|
| `max_users` | 最大用户数 | 1000 |
| `max_roles` | 最大角色数 | 50 |
| `max_permissions` | 最大权限点数 | 200 |
| `max_user_groups` | 最大用户组数 | 20 |
| `max_admin_users` | 最大管理员数 | 10 |
| `api_call_monthly` | 月度 API 调用次数 | 1000 万 |
| `storage_mb` | 存储空间（MB） | 1024 |

### 4. 配额检查

1. 创建资源前检查配额
2. 配额不足时拒绝创建
3. 返回明确的配额错误码
4. 提示当前使用量和配额上限

### 5. 配额预警

| 预警阈值 | 通知方式 | 接收人 |
|----------|----------|--------|
| 使用量达到 80% | 邮件 | 租户管理员 |
| 使用量达到 90% | 邮件 + 短信 | 租户管理员 |
| 使用量达到 100% | 邮件 + 短信 | 租户管理员、平台运营 |
| 连续 7 天超过 90% | 邮件 | 平台运营（推荐升级） |

### 6. 限流响应

超限时的响应格式：

```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "请求频率超限，请稍后重试",
  "data": {
    "limit": 1000,
    "current": 1050,
    "reset_at": "2026-03-25T10:31:00Z",
    "retry_after": 60
  }
}
```

HTTP 状态码：`429 Too Many Requests`

响应头：
- `X-RateLimit-Limit`: 限制值
- `X-RateLimit-Remaining`: 剩余量
- `X-RateLimit-Reset`: 重置时间戳
- `Retry-After`: 建议重试时间（秒）

### 7. 配额豁免

以下场景可豁免限流：

| 场景 | 说明 |
|------|------|
| 健康检查接口 | `/health`, `/ready` |
| 平台管理员操作 | 平台运营紧急操作 |
| 白名单租户 | VIP 租户或测试租户 |
| 内部服务调用 | 服务间调用（带内部 Token） |

### 8. 配额查询 API

租户可查询自己的配额使用情况：

```
GET /api/v1/quotas       # 配额列表和使用情况
GET /api/v1/quotas/:type # 指定配额详情
```

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| 限流配置错误 | 使用默认配置，记录告警 |
| Redis 不可用 | 降级为本地限流或放行 |
| 配额检查超时 | 允许操作，记录日志 |
| 配额数据不一致 | 异步校准，告警通知 |

**API 接口：**

```
# 限流配置（平台运营）
GET    /api/v1/rate-limits/config        # 获取限流配置
PUT    /api/v1/rate-limits/config        # 更新限流配置

# 配额管理
GET    /api/v1/quotas                    # 配额列表
GET    /api/v1/quotas/:type              # 配额详情
PUT    /api/v1/quotas                    # 更新配额（平台运营）
POST   /api/v1/quotas/alert-config       # 配置配额告警

# 使用情况统计
GET    /api/v1/usage/summary             # 使用汇总
GET    /api/v1/usage/api-calls           # API 调用统计
GET    /api/v1/usage/resource-usage      # 资源使用情况
```

**数据库设计：**

```sql
-- 租户配额表
CREATE TABLE tenant_quotas (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL UNIQUE,
    max_users INT DEFAULT 1000,
    max_roles INT DEFAULT 50,
    max_permissions INT DEFAULT 200,
    max_user_groups INT DEFAULT 20,
    max_admin_users INT DEFAULT 10,
    api_call_monthly BIGINT DEFAULT 10000000,
    storage_mb INT DEFAULT 1024,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 配额使用记录表
CREATE TABLE quota_usage (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    quota_type VARCHAR(50) NOT NULL,
    used_count INT DEFAULT 0,
    limit_count INT NOT NULL,
    last_checked_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_type (tenant_id, quota_type)
);

-- 配额告警记录表
CREATE TABLE quota_alerts (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    quota_type VARCHAR(50) NOT NULL,
    threshold_percent INT NOT NULL,
    current_usage INT NOT NULL,
    is_sent BOOLEAN DEFAULT FALSE,
    sent_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant (tenant_id, created_at)
);

-- API 调用统计（用于限流和计费）
CREATE TABLE api_usage_stats (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT,
    api_path VARCHAR(255) NOT NULL,
    call_count BIGINT DEFAULT 0,
    stat_date DATE NOT NULL,
    UNIQUE KEY uk_tenant_user_api_date (tenant_id, user_id, api_path, stat_date),
    INDEX idx_tenant_date (tenant_id, stat_date)
);
```

**验收标准：**

- [ ] API 限流正确生效
- [ ] 租户配额检查正确
- [ ] 配额超限后拒绝创建
- [ ] 配额预警通知正常发送
- [ ] 限流响应头和状态码正确
- [ ] 配额使用情况可查询
- [ ] 支持配额豁免配置
