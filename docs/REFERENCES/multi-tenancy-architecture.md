# 多租户架构设计方案

> 最后更新：2026-03-25
> 适用场景：IAM SaaS 多租户设计

## 1. 多租户是什么

**多租户 (Multi-Tenancy)** 是一种软件架构模式，允许多个租户（客户）共享同一套应用实例和基础设施，同时保持数据隔离。

### 1.1 核心概念

| 概念 | 说明 | 示例 |
|------|------|------|
| **租户 (Tenant)** | 使用系统的独立组织或客户 | 公司 A、公司 B |
| **租户隔离** | 不同租户的数据互不可见 | 公司 A 看不到公司 B 的用户 |
| **配额 (Quota)** | 租户可使用的资源限制 | 最多 1000 用户、100 角色 |

### 1.2 为什么需要多租户

| 单租户架构 | 多租户架构 |
|------------|------------|
| 每个客户部署一套 | 所有客户共享一套 |
| 运维成本高 | 运维成本低 |
| 资源利用率低 | 资源利用率高 |
| 升级麻烦 | 统一升级 |

---

## 2. 多租户隔离级别

### 2.1 三种隔离方案对比

| 方案 | 说明 | 隔离性 | 成本 | 扩展性 | 适用场景 |
|------|------|--------|------|--------|----------|
| **数据库隔离** | 每个租户独立数据库 | ★★★★★ | 高 | 中 | 高安全/合规要求 |
| **Schema 隔离** | 共享 DB，独立 Schema | ★★★★☆ | 中 | 中 | 中等规模 SaaS |
| **数据行隔离** | 共享 DB+Schema，`tenant_id` 区分 | ★★★☆☆ | 低 | 高 | 互联网 SaaS |

---

### 2.2 方案一：数据库隔离（Database per Tenant）

```
┌─────────────────────────────────────────┐
│           应用服务器                     │
└─────────────────────────────────────────┘
         │         │         │
         ▼         ▼         ▼
    ┌────────┐ ┌────────┐ ┌────────┐
    │ DB-A   │ │ DB-B   │ │ DB-C   │
    │ TenantA│ │ TenantB│ │ TenantC│
    └────────┘ └────────┘ └────────┘
```

**优点：**
- 数据隔离最彻底
- 备份恢复简单
- 支持定制数据库结构

**缺点：**
- 数据库数量随租户增长
- 连接池开销大
- 跨租户分析困难

**IAM 选择：** 不适用于当前场景（成本高，维护复杂）

---

### 2.3 方案二：Schema 隔离（Schema per Tenant）

```
┌─────────────────────────────────────────┐
│           应用服务器                     │
└─────────────────────────────────────────┘
                    │
                    ▼
         ┌──────────────────────┐
         │    MySQL Database    │
         ├──────────────────────┤
         │  Schema: tenant_a    │
         │  Schema: tenant_b    │
         │  Schema: tenant_c    │
         └──────────────────────┘
```

**优点：**
- 逻辑隔离清晰
- 备份恢复较简单
- 支持租户级数据导出

**缺点：**
- Schema 数量有限制
- 表结构变更需要同步
- 跨 Schema 查询复杂

**IAM 选择：** 不适用于当前场景（MySQL Schema 隔离收益不明显）

---

### 2.4 方案三：数据行隔离（Row-level Isolation）⭐

```
┌─────────────────────────────────────────┐
│           应用服务器                     │
└─────────────────────────────────────────┘
                    │
                    ▼
         ┌──────────────────────┐
         │    MySQL Database    │
         ├──────────────────────┤
         │  tenants 表          │
         │  users 表 (tenant_id)│
         │  roles 表 (tenant_id)│
         │  ...                 │
         └──────────────────────┘
```

**优点：**
- 成本最低
- 扩展性最好
- 运维简单
- 跨租户分析容易

**缺点：**
- 需要代码保证隔离
- 误操作风险（需严格控制）
- 备份恢复复杂

**IAM 选择：** **当前方案**

---

## 3. IAM 数据行隔离实现

### 3.1 表设计要求

所有业务表必须包含 `tenant_id` 字段：

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,  -- 租户 ID（必填）
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- 租户内邮箱唯一
    UNIQUE KEY uk_tenant_email (tenant_id, email),
    -- 索引加速租户查询
    INDEX idx_tenant (tenant_id)
);
```

### 3.2 查询时必须带 tenant_id

```go
// ✅ 正确：始终带 tenant_id
func GetUserByID(tenantID, userID int64) (*User, error) {
    var user User
    err := db.QueryRow(
        "SELECT id, tenant_id, email FROM users WHERE id = ? AND tenant_id = ?",
        userID, tenantID,
    ).Scan(&user.ID, &user.TenantID, &user.Email)
    return &user, err
}

// ❌ 错误：缺少 tenant_id 条件
func GetUserByID(userID int64) (*User, error) {
    var user User
    err := db.QueryRow(
        "SELECT id, tenant_id, email FROM users WHERE id = ?",
        userID,
    ).Scan(&user.ID, &user.TenantID, &user.Email)
    return &user, err
}
```

### 3.3 中间件自动注入 tenant_id

```go
// 从 JWT Token 中提取 tenant_id，注入到上下文
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := GetClaimsFromContext(r.Context())
        if claims == nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        // 将 tenant_id 注入上下文
        ctx := context.WithValue(r.Context(), "tenant_id", claims.TenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3.4 Repository 层封装

```go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) List(ctx context.Context, req ListUsersRequest) ([]User, error) {
    // 从上下文获取 tenant_id，不能信任用户传入的参数
    tenantID := ctx.Value("tenant_id").(int64)

    query := "SELECT id, email, status FROM users WHERE tenant_id = ?"
    rows, err := r.db.QueryContext(ctx, query, tenantID)
    // ...
}
```

---

## 4. 多租户关键设计

### 4.1 租户数据结构

```sql
CREATE TABLE tenants (
    id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,      -- 租户名称
    status TINYINT DEFAULT 1,        -- 1-激活 0-冻结
    max_users INT DEFAULT 1000,      -- 最大用户数配额
    max_roles INT DEFAULT 50,        -- 最大角色数配额
    expires_at DATETIME,             -- 租户过期时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 4.2 租户状态管理

| 状态 | 说明 | 可登录 | 可创建资源 |
|------|------|--------|------------|
| 激活 (Active) | 正常状态 | ✅ | ✅ |
| 冻结 (Frozen) | 欠费或违规 | ❌ | ❌ |
| 过期 (Expired) | 超过有效期 | ❌ | ❌ |

### 4.3 配额检查

```go
func (s *TenantService) CheckQuota(ctx context.Context, resourceType string) error {
    tenantID := GetTenantID(ctx)
    quota, err := s.GetTenantQuota(tenantID)
    if err != nil {
        return err
    }

    usage, err := s.GetUsage(tenantID, resourceType)
    if err != nil {
        return err
    }

    if usage >= quota.Limit {
        return &QuotaExceededError{
            ResourceType: resourceType,
            Limit:        quota.Limit,
            Usage:        usage,
        }
    }

    return nil
}
```

---

## 5. 多租户隔离检查清单

### 5.1 数据隔离

- [ ] 所有业务表包含 `tenant_id` 字段
- [ ] 所有查询包含 `tenant_id` 条件
- [ ] 所有 `tenant_id` 相关的唯一索引正确创建
- [ ] 外键关联包含 `tenant_id`

### 5.2 API 隔离

- [ ] 从 Token 中提取 `tenant_id`，不信任用户传入参数
- [ ] 批量操作（如导入）校验所有记录的 `tenant_id`
- [ ] 跨租户操作（如平台运营）需要特殊权限

### 5.3 日志隔离

- [ ] 审计日志包含 `tenant_id`
- [ ] 租户管理员只能查看本租户日志
- [ ] 日志导出自动过滤 `tenant_id`

### 5.4 缓存隔离

- [ ] Cache Key 包含 `tenant_id` 前缀
- [ ] Redis 按租户设置 TTL
- [ ] 租户删除时清理缓存

---

## 6. 常见问题

### Q1: 全局唯一 ID 如何生成？

使用雪花算法（Snowflake）生成全局唯一 ID：

```go
// 示例：使用 sonyflake
import "github.com/sony/sonyflake"

fl := sonyflake.NewSonyflake(sonyflake.Settings{})
id, err := fl.NextID()  // 64 位唯一 ID
```

### Q2: 租户数据如何导出？

```sql
-- 导出租户 A 的所有数据
SELECT * FROM users WHERE tenant_id = 1;
SELECT * FROM roles WHERE tenant_id = 1;
SELECT * FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE u.tenant_id = 1;
```

### Q3: 租户删除后数据如何处理？

| 策略 | 说明 |
|------|------|
| 软删除 | 标记 `deleted_at`，保留数据 |
| 硬删除 | 级联删除所有相关数据 |
| 归档 | 迁移到冷存储 |

推荐：**软删除 + 定期归档**

### Q4: 如何支持跨租户协作？

当前设计不支持跨租户协作（一个用户只能属于一个租户）。如需支持：

1. 设计 `tenant_user_mapping` 表
2. 一个用户可映射到多个租户
3. 登录时选择目标租户

---

## 7. 参考链接

- AWS 多租户架构：https://docs.aws.amazon.com/wellarchitected/latest/saas-lens/saas-architecture.html
- SaaS Mag：https://saasmag.com/ （SaaS 架构博客）

---

## 8. 相关需求文档

- [REQ-007 租户管理功能](../05-functional-requirements/REQ-007-tenant-management.md)
- [REQ-016 API 限流和配额管理](../05-functional-requirements/REQ-016-rate-limit-quota.md)
