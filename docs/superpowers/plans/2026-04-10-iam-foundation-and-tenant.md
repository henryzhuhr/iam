# IAM 基础架构与租户管理实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 搭建 IAM 系统基础架构（数据库连接、实体层、错误处理、中间件框架、CI 流水线），并以租户 CRUD 作为首个可验证的功能闭环。

**Architecture:** 模块化单体，按层分目录（handler/service/repository/dto/entity 按领域划分子目录）。通过 ServiceContext 注入所有依赖。多租户逻辑隔离（tenant_id 强制过滤）。

**Tech Stack:** Go 1.25 + go-zero 1.10 + MySQL 8.0 + Redis 7.0 + Kafka (apache/kafka 4.1) + testify + Docker Compose

---

## 文件总览

**新建文件（按任务顺序）：**

| 文件 | 职责 | 任务 |
|------|------|------|
| `sql/001_init.sql` | 全部 16 张表的 DDL | Task 1 |
| `infra/database/mysql.go` | MySQL 连接封装 | Task 2 |
| `infra/cache/redis.go` | Redis 连接封装 | Task 2 |
| `infra/queue/kafka.go` | Kafka Producer 封装 | Task 2 |
| `internal/entity/tenant.go` | tenants 表实体 | Task 3 |
| `internal/entity/user.go` | users 相关实体 | Task 3 |
| `internal/entity/role.go` | RBAC 相关实体 | Task 3 |
| `internal/entity/permission.go` | permissions 实体 | Task 3 |
| `internal/entity/app.go` | applications 实体 | Task 3 |
| `internal/entity/client.go` | clients 实体 | Task 3 |
| `internal/entity/audit.go` | audit_logs/login_logs 实体 | Task 3 |
| `internal/entity/password_policy.go` | 密码策略实体 | Task 3 |
| `internal/entity/group.go` | user_groups 相关实体 | Task 3 |
| `internal/config/config.go` | 修改：增加 DB/Redis/Kafka 配置 | Task 4 |
| `internal/constant/errors.go` | 错误码常量定义 | Task 5 |
| `internal/constant/tenant.go` | 租户相关常量 | Task 5 |
| `internal/middleware/tenant.go` | 租户隔离中间件 | Task 6 |
| `internal/middleware/auth.go` | JWT 认证中间件（骨架） | Task 6 |
| `internal/middleware/audit.go` | 审计日志中间件（骨架） | Task 6 |
| `internal/middleware/error.go` | 全局错误处理中间件 | Task 6 |
| `internal/dto/tenant/tenant.go` | 租户 DTO | Task 7 |
| `internal/repository/tenant_repo.go` | 租户数据访问层 | Task 8 |
| `internal/service/tenant/tenant.go` | 租户业务逻辑层 | Task 9 |
| `internal/handler/tenant/tenant.go` | 租户 HTTP 处理器 | Task 10 |
| `internal/routes/tenant/tenant.go` | 租户路由注册 | Task 11 |
| `internal/svc/servicecontext.go` | 修改：注入所有依赖 | Task 12 |
| `internal/routes/routes.go` | 修改：注册租户路由 | Task 12 |
| `etc/dev.yaml` | 修改：增加 DB/Redis/Kafka 配置 | Task 4 |
| `docker-compose.yml` | 修改：增加健康检查验证 | Task 13 |
| `.github/workflows/01-ci.yaml` | CI 流水线 | Task 14 |
| `scripts/ci-local.sh` | 本地 CI 验证脚本 | Task 15 |
| `internal/tests/integration/tenant_test.go` | 租户集成测试 | Task 16 |

---

### Task 1: SQL Schema — 全部 16 张表 DDL

**Files:**
- Create: `sql/001_init.sql`

- [ ] **Step 1: 编写完整的 SQL 初始化脚本**

```sql
-- sql/001_init.sql
-- IAM 系统数据库初始化脚本
-- MySQL 8.0

CREATE DATABASE IF NOT EXISTS `iam` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `iam`;

-- ============================================
-- 1. tenants — 租户表
-- ============================================
CREATE TABLE `tenants` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(100) NOT NULL,
    `status`     TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用 3=过期',
    `max_users`  INT          NOT NULL DEFAULT 100 COMMENT '最大用户配额',
    `max_apps`   INT          NOT NULL DEFAULT 10  COMMENT '最大应用配额',
    `expire_at`  DATETIME              DEFAULT NULL COMMENT '租户过期时间',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- ============================================
-- 2. users — 用户表
-- ============================================
CREATE TABLE `users` (
    `id`                  BIGINT        NOT NULL AUTO_INCREMENT,
    `tenant_id`           BIGINT        NOT NULL COMMENT '租户 ID',
    `email`               VARCHAR(100)  NOT NULL COMMENT '邮箱（登录账号）',
    `phone`               VARCHAR(20)            DEFAULT NULL COMMENT '手机号',
    `password_hash`       VARCHAR(255)  NOT NULL COMMENT 'bcrypt 哈希',
    `status`              TINYINT       NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用 3=锁定',
    `mfa_enabled`         TINYINT       NOT NULL DEFAULT 0 COMMENT '是否开启 MFA',
    `mfa_secret`          VARCHAR(100)           DEFAULT NULL COMMENT 'TOTP 密钥',
    `last_login_at`       DATETIME               DEFAULT NULL COMMENT '最后登录时间',
    `password_changed_at` DATETIME               DEFAULT NULL COMMENT '密码最后修改时间',
    `created_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_email` (`tenant_id`, `email`),
    INDEX `idx_tenant_status` (`tenant_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 3. user_groups — 用户组表
-- ============================================
CREATE TABLE `user_groups` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `name`        VARCHAR(100) NOT NULL COMMENT '用户组名称',
    `description` TEXT                  DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组表';

-- ============================================
-- 4. user_group_members — 用户组成员表
-- ============================================
CREATE TABLE `user_group_members` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `group_id`   BIGINT   NOT NULL COMMENT '用户组 ID',
    `user_id`    BIGINT   NOT NULL COMMENT '用户 ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_user` (`group_id`, `user_id`),
    INDEX `idx_group` (`group_id`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组成员表';

-- ============================================
-- 5. permissions — 权限定义表
-- ============================================
CREATE TABLE `permissions` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `code`        VARCHAR(100) NOT NULL COMMENT '权限编码，如 user:read',
    `name`        VARCHAR(100) NOT NULL COMMENT '权限名称',
    `resource`    VARCHAR(100) NOT NULL COMMENT '资源类型',
    `action`      VARCHAR(50)  NOT NULL COMMENT '操作：read/write/delete',
    `app_code`    VARCHAR(50)           DEFAULT NULL COMMENT '归属应用（NULL=平台级）',
    `description` TEXT                  DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    INDEX `idx_app_code` (`app_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限定义表';

-- ============================================
-- 6. roles — 角色表
-- ============================================
CREATE TABLE `roles` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `name`        VARCHAR(100) NOT NULL COMMENT '角色名称',
    `code`        VARCHAR(100) NOT NULL COMMENT '角色编码',
    `type`        TINYINT      NOT NULL DEFAULT 2 COMMENT '1=系统内置 2=自定义',
    `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `description` TEXT                  DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================
-- 7. role_permissions — 角色权限关联表
-- ============================================
CREATE TABLE `role_permissions` (
    `id`            BIGINT     NOT NULL AUTO_INCREMENT,
    `role_id`       BIGINT     NOT NULL COMMENT '角色 ID',
    `permission_id` BIGINT     NOT NULL COMMENT '权限 ID',
    `data_scope`    VARCHAR(50) NOT NULL DEFAULT 'all' COMMENT 'all/dept/dept_and_sub/personal/custom',
    `created_at`    DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_perm` (`role_id`, `permission_id`),
    INDEX `idx_role` (`role_id`),
    INDEX `idx_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ============================================
-- 8. user_roles — 用户角色关联表
-- ============================================
CREATE TABLE `user_roles` (
    `id`        BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id` BIGINT       NOT NULL COMMENT '租户 ID',
    `user_id`   BIGINT       NOT NULL COMMENT '用户 ID',
    `role_id`   BIGINT       NOT NULL COMMENT '角色 ID',
    `app_code`  VARCHAR(50)           DEFAULT NULL COMMENT '应用编码（角色生效范围）',
    `created_at` DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user_role_app` (`tenant_id`, `user_id`, `role_id`, `app_code`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 9. role_constraints — 角色约束表（SoD 约束）
-- ============================================
CREATE TABLE `role_constraints` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT   NOT NULL COMMENT '租户 ID',
    `type`       TINYINT  NOT NULL COMMENT '1=静态SoD 2=动态SoD',
    `role_a`     BIGINT   NOT NULL COMMENT '冲突角色 A',
    `role_b`     BIGINT   NOT NULL COMMENT '冲突角色 B',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色约束表（SoD）';

-- ============================================
-- 10. applications — 应用表
-- ============================================
CREATE TABLE `applications` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT       NOT NULL COMMENT '租户 ID',
    `code`        VARCHAR(50)  NOT NULL COMMENT '应用编码（租户内唯一）',
    `name`        VARCHAR(100) NOT NULL COMMENT '应用名称',
    `description` TEXT                  DEFAULT NULL,
    `status`      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    INDEX `idx_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用表';

-- ============================================
-- 11. user_app_authorizations — 用户应用授权表
-- ============================================
CREATE TABLE `user_app_authorizations` (
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `tenant_id`  BIGINT   NOT NULL COMMENT '租户 ID',
    `user_id`    BIGINT   NOT NULL COMMENT '用户 ID',
    `app_id`     BIGINT   NOT NULL COMMENT '应用 ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_user_app` (`tenant_id`, `user_id`, `app_id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_app` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户应用授权表';

-- ============================================
-- 12. clients — 内部客户端表
-- ============================================
CREATE TABLE `clients` (
    `id`              BIGINT       NOT NULL AUTO_INCREMENT,
    `client_id`       VARCHAR(64)  NOT NULL COMMENT '客户端标识',
    `access_key`      VARCHAR(64)  NOT NULL COMMENT 'AK',
    `secret_key_hash` VARCHAR(255) NOT NULL COMMENT 'SK 哈希',
    `name`            VARCHAR(100) NOT NULL COMMENT '客户端名称',
    `allowed_scopes`  JSON         NOT NULL COMMENT '允许的 scopes',
    `status`          TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 2=禁用',
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='内部客户端表';

-- ============================================
-- 13. audit_logs — 操作审计日志表
-- ============================================
CREATE TABLE `audit_logs` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `tenant_id`     BIGINT       NOT NULL COMMENT '租户 ID',
    `user_id`       BIGINT       NOT NULL COMMENT '操作人',
    `action`        VARCHAR(100) NOT NULL COMMENT '操作类型',
    `resource_type` VARCHAR(50)  NOT NULL COMMENT '资源类型',
    `resource_id`   BIGINT                DEFAULT NULL COMMENT '资源 ID',
    `detail`        JSON                  DEFAULT NULL COMMENT '操作详情',
    `ip`            VARCHAR(45)           DEFAULT NULL COMMENT '操作 IP',
    `created_at`    DATETIME     NOT NULL COMMENT '操作时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_user` (`user_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作审计日志表';

-- ============================================
-- 14. login_logs — 登录日志表
-- ============================================
CREATE TABLE `login_logs` (
    `id`          BIGINT        NOT NULL AUTO_INCREMENT,
    `tenant_id`   BIGINT        NOT NULL COMMENT '租户 ID',
    `user_id`     BIGINT                 DEFAULT NULL COMMENT 'NULL=登录失败',
    `email`       VARCHAR(100)  NOT NULL COMMENT '登录账号',
    `status`      TINYINT       NOT NULL COMMENT '1=成功 2=失败 3=MFA待验证',
    `fail_reason` VARCHAR(200)           DEFAULT NULL COMMENT '失败原因',
    `login_type`  VARCHAR(30)   NOT NULL COMMENT 'password/code/oauth/mfa',
    `ip`          VARCHAR(45)            DEFAULT NULL COMMENT '登录 IP',
    `user_agent`  VARCHAR(500)           DEFAULT NULL COMMENT '用户代理',
    `created_at`  DATETIME      NOT NULL COMMENT '登录时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant` (`tenant_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ============================================
-- 15. password_policies — 密码策略表
-- ============================================
CREATE TABLE `password_policies` (
    `id`                  BIGINT  NOT NULL AUTO_INCREMENT,
    `tenant_id`           BIGINT  NOT NULL COMMENT '每租户一条',
    `min_length`          INT     NOT NULL DEFAULT 8 COMMENT '最小密码长度',
    `require_uppercase`   TINYINT NOT NULL DEFAULT 1 COMMENT '需要大写字母',
    `require_lowercase`   TINYINT NOT NULL DEFAULT 1 COMMENT '需要小写字母',
    `require_digit`       TINYINT NOT NULL DEFAULT 1 COMMENT '需要数字',
    `require_special`     TINYINT NOT NULL DEFAULT 1 COMMENT '需要特殊字符',
    `history_count`       INT     NOT NULL DEFAULT 3 COMMENT '历史密码检查次数',
    `expire_days`         INT     NOT NULL DEFAULT 0 COMMENT '密码过期天数（0=永不过期）',
    `max_login_attempts`  INT     NOT NULL DEFAULT 5 COMMENT '最大登录失败次数',
    `lockout_minutes`     INT     NOT NULL DEFAULT 30 COMMENT '锁定时长（分钟）',
    `updated_at`          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码策略表';

-- ============================================
-- 16. password_history — 密码历史表
-- ============================================
CREATE TABLE `password_history` (
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`       BIGINT       NOT NULL COMMENT '用户 ID',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '历史密码哈希',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='密码历史表';
```

- [ ] **Step 2: 验证 SQL 语法**

启动 MySQL 容器并执行脚本验证：

```bash
# 启动 MySQL
docker compose up -d mysql
# 等待 MySQL 就绪
docker compose exec mysql mysqladmin ping -h localhost --wait

# 执行初始化脚本
docker compose exec -T mysql mysql -uroot -prootpassword < sql/001_init.sql

# 验证表创建成功
docker compose exec mysql mysql -uroot -prootpassword -e "USE iam; SHOW TABLES;"
```

期望输出：16 张表全部列出。

- [ ] **Step 3: Commit**

```bash
git add sql/001_init.sql
git commit -m "feat: add complete SQL schema for all 16 tables"
```

---

### Task 2: 基础设施层 — DB / Redis / Kafka 连接

**Files:**
- Create: `infra/database/mysql.go`
- Create: `infra/cache/redis.go`
- Create: `infra/queue/kafka.go`

- [ ] **Step 1: MySQL 连接封装**

```go
// infra/database/mysql.go
package database

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
    Host     string `json:"Host"`
    Port     int    `json:"Port"`
    User     string `json:"User"`
    Password string `json:"Password"`
    Database string `json:"Database"`
}

// NewMySQL 创建 MySQL 连接池
func NewMySQL(cfg MySQLConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("open mysql: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ping mysql: %w", err)
    }

    return db, nil
}
```

- [ ] **Step 2: Redis 连接封装**

```go
// infra/cache/redis.go
package cache

import (
    "context"
    "fmt"
    "time"

    "github.com/zeromicro/go-zero/core/stores/redis"
)

type RedisConfig struct {
    Host     string `json:"Host"`
    Port     int    `json:"Port"`
    Password string `json:"Password,optional"`
    DB       int    `json:"DB,optional"`
}

// NewRedis 创建 Redis 客户端
func NewRedis(cfg RedisConfig) (*redis.Redis, error) {
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
    r := redis.New(addr, redis.WithPass(cfg.Password), redis.WithDB(cfg.DB))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := r.PingCtx(ctx); err != nil {
        return nil, fmt.Errorf("ping redis: %w", err)
    }

    return r, nil
}
```

- [ ] **Step 3: Kafka Producer 封装**

```go
// infra/queue/kafka.go
package queue

import (
    "context"
    "fmt"

    "github.com/zeromicro/go-zero/core/logx"
)

type KafkaConfig struct {
    Brokers []string `json:"Brokers"`
    Topic   string   `json:"Topic,optional"`
}

// KafkaProducer Kafka 生产者封装
type KafkaProducer struct {
    config  KafkaConfig
    stopped bool
}

// NewKafkaProducer 创建 Kafka Producer（当前为骨架实现）
func NewKafkaProducer(cfg KafkaConfig) (*KafkaProducer, error) {
    // TODO: 集成 Sarama 或 go-zero Kafka 后完善
    logx.Infof("kafka producer initialized (stub): brokers=%v", cfg.Brokers)
    return &KafkaProducer{config: cfg}, nil
}

// SendMessage 发送消息到 Kafka（当前为骨架实现）
func (p *KafkaProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
    if p.stopped {
        return fmt.Errorf("producer is stopped")
    }
    // TODO: 实际发送逻辑
    logx.Infof("kafka send message (stub): topic=%s key=%s", topic, string(key))
    return nil
}

// Close 关闭生产者
func (p *KafkaProducer) Close() {
    p.stopped = true
    logx.Info("kafka producer closed")
}
```

- [ ] **Step 4: Commit**

```bash
git add infra/database/mysql.go infra/cache/redis.go infra/queue/kafka.go
git commit -m "feat: add infrastructure layer for MySQL, Redis, Kafka"
```

---

### Task 3: Entity 层 — 全部 16 张表对应 Go 结构体

**Files:**
- Create: `internal/entity/tenant.go`
- Create: `internal/entity/user.go`
- Create: `internal/entity/role.go`
- Create: `internal/entity/permission.go`
- Create: `internal/entity/app.go`
- Create: `internal/entity/client.go`
- Create: `internal/entity/audit.go`
- Create: `internal/entity/password_policy.go`
- Create: `internal/entity/group.go`

- [ ] **Step 1: 租户实体**

```go
// internal/entity/tenant.go
package entity

import "time"

// Tenant 租户实体
type Tenant struct {
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    Status    int8      `db:"status"`    // 1=启用 2=禁用 3=过期
    MaxUsers  int       `db:"max_users"`
    MaxApps   int       `db:"max_apps"`
    ExpireAt  time.Time `db:"expire_at"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// TenantStatus 租户状态常量
const (
    TenantStatusActive  int8 = 1
    TenantStatusDisabled int8 = 2
    TenantStatusExpired  int8 = 3
)
```

- [ ] **Step 2: 用户实体**

```go
// internal/entity/user.go
package entity

import "time"

// User 用户实体
type User struct {
    ID                int64      `db:"id"`
    TenantID          int64      `db:"tenant_id"`
    Email             string     `db:"email"`
    Phone             string     `db:"phone"`
    PasswordHash      string     `db:"password_hash"`
    Status            int8       `db:"status"`       // 1=启用 2=禁用 3=锁定
    MFAEnabled        int8       `db:"mfa_enabled"`
    MFASecret         string     `db:"mfa_secret"`
    LastLoginAt       *time.Time `db:"last_login_at"`
    PasswordChangedAt *time.Time `db:"password_changed_at"`
    CreatedAt         time.Time  `db:"created_at"`
    UpdatedAt         time.Time  `db:"updated_at"`
}

// UserStatus 用户状态常量
const (
    UserStatusActive   int8 = 1
    UserStatusDisabled int8 = 2
    UserStatusLocked   int8 = 3
)
```

- [ ] **Step 3: 角色实体**

```go
// internal/entity/role.go
package entity

import "time"

// Role 角色实体
type Role struct {
    ID          int64     `db:"id"`
    TenantID    int64     `db:"tenant_id"`
    Name        string    `db:"name"`
    Code        string    `db:"code"`
    Type        int8      `db:"type"`        // 1=系统内置 2=自定义
    Status      int8      `db:"status"`      // 1=启用 2=禁用
    Description string    `db:"description"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}

// RoleType 角色类型常量
const (
    RoleTypeSystem  int8 = 1
    RoleTypeCustom  int8 = 2
)

// RoleStatus 角色状态常量
const (
    RoleStatusActive   int8 = 1
    RoleStatusDisabled int8 = 2
)

// RolePermission 角色权限关联实体
type RolePermission struct {
    ID           int64     `db:"id"`
    RoleID       int64     `db:"role_id"`
    PermissionID int64     `db:"permission_id"`
    DataScope    string    `db:"data_scope"`
    CreatedAt    time.Time `db:"created_at"`
}

// UserRole 用户角色关联实体
type UserRole struct {
    ID        int64     `db:"id"`
    TenantID  int64     `db:"tenant_id"`
    UserID    int64     `db:"user_id"`
    RoleID    int64     `db:"role_id"`
    AppCode   string    `db:"app_code"`
    CreatedAt time.Time `db:"created_at"`
}

// RoleConstraint 角色约束实体（SoD）
type RoleConstraint struct {
    ID        int64     `db:"id"`
    TenantID  int64     `db:"tenant_id"`
    Type      int8      `db:"type"`
    RoleA     int64     `db:"role_a"`
    RoleB     int64     `db:"role_b"`
    CreatedAt time.Time `db:"created_at"`
}
```

- [ ] **Step 4: 权限实体**

```go
// internal/entity/permission.go
package entity

import "time"

// Permission 权限实体
type Permission struct {
    ID          int64     `db:"id"`
    Code        string    `db:"code"`
    Name        string    `db:"name"`
    Resource    string    `db:"resource"`
    Action      string    `db:"action"`
    AppCode     string    `db:"app_code"`
    Description string    `db:"description"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}
```

- [ ] **Step 5: 应用实体**

```go
// internal/entity/app.go
package entity

import "time"

// Application 应用实体
type Application struct {
    ID          int64     `db:"id"`
    TenantID    int64     `db:"tenant_id"`
    Code        string    `db:"code"`
    Name        string    `db:"name"`
    Description string    `db:"description"`
    Status      int8      `db:"status"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}

// UserAppAuthorization 用户应用授权实体
type UserAppAuthorization struct {
    ID        int64     `db:"id"`
    TenantID  int64     `db:"tenant_id"`
    UserID    int64     `db:"user_id"`
    AppID     int64     `db:"app_id"`
    CreatedAt time.Time `db:"created_at"`
}
```

- [ ] **Step 6: 客户端实体**

```go
// internal/entity/client.go
package entity

import "time"

// Client 内部客户端实体
type Client struct {
    ID            int64     `db:"id"`
    ClientID      string    `db:"client_id"`
    AccessKey     string    `db:"access_key"`
    SecretKeyHash string    `db:"secret_key_hash"`
    Name          string    `db:"name"`
    AllowedScopes string    `db:"allowed_scopes"` // JSON 字符串
    Status        int8      `db:"status"`
    CreatedAt     time.Time `db:"created_at"`
    UpdatedAt     time.Time `db:"updated_at"`
}
```

- [ ] **Step 7: 审计日志实体**

```go
// internal/entity/audit.go
package entity

import "time"

// AuditLog 操作审计日志实体
type AuditLog struct {
    ID           int64     `db:"id"`
    TenantID     int64     `db:"tenant_id"`
    UserID       int64     `db:"user_id"`
    Action       string    `db:"action"`
    ResourceType string    `db:"resource_type"`
    ResourceID   *int64    `db:"resource_id"`
    Detail       string    `db:"detail"` // JSON 字符串
    IP           string    `db:"ip"`
    CreatedAt    time.Time `db:"created_at"`
}

// LoginLog 登录日志实体
type LoginLog struct {
    ID         int64     `db:"id"`
    TenantID   int64     `db:"tenant_id"`
    UserID     *int64    `db:"user_id"`
    Email      string    `db:"email"`
    Status     int8      `db:"status"` // 1=成功 2=失败 3=MFA待验证
    FailReason string    `db:"fail_reason"`
    LoginType  string    `db:"login_type"`
    IP         string    `db:"ip"`
    UserAgent  string    `db:"user_agent"`
    CreatedAt  time.Time `db:"created_at"`
}
```

- [ ] **Step 8: 密码策略实体**

```go
// internal/entity/password_policy.go
package entity

import "time"

// PasswordPolicy 密码策略实体
type PasswordPolicy struct {
    ID               int64     `db:"id"`
    TenantID         int64     `db:"tenant_id"`
    MinLength        int       `db:"min_length"`
    RequireUppercase int8      `db:"require_uppercase"`
    RequireLowercase int8      `db:"require_lowercase"`
    RequireDigit     int8      `db:"require_digit"`
    RequireSpecial   int8      `db:"require_special"`
    HistoryCount     int       `db:"history_count"`
    ExpireDays       int       `db:"expire_days"`
    MaxLoginAttempts int       `db:"max_login_attempts"`
    LockoutMinutes   int       `db:"lockout_minutes"`
    UpdatedAt        time.Time `db:"updated_at"`
}

// PasswordHistory 密码历史实体
type PasswordHistory struct {
    ID           int64     `db:"id"`
    UserID       int64     `db:"user_id"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
}
```

- [ ] **Step 9: 用户组实体**

```go
// internal/entity/group.go
package entity

import "time"

// UserGroup 用户组实体
type UserGroup struct {
    ID          int64     `db:"id"`
    TenantID    int64     `db:"tenant_id"`
    Name        string    `db:"name"`
    Description string    `db:"description"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}

// UserGroupMember 用户组成员实体
type UserGroupMember struct {
    ID        int64     `db:"id"`
    GroupID   int64     `db:"group_id"`
    UserID    int64     `db:"user_id"`
    CreatedAt time.Time `db:"created_at"`
}
```

- [ ] **Step 10: Commit**

```bash
git add internal/entity/
git commit -m "feat: add entity definitions for all 16 tables"
```

---

### Task 4: 更新配置文件 — 增加 DB/Redis/Kafka 配置

**Files:**
- Modify: `internal/config/config.go`
- Modify: `etc/dev.yaml`

- [ ] **Step 1: 更新 Config 结构体**

```go
// internal/config/config.go
package config

import (
    "iam/infra/cache"
    "iam/infra/database"
    "iam/infra/queue"

    "github.com/zeromicro/go-zero/rest"
)

type Config struct {
    rest.RestConf
    LocaleDir string                   `json:"LocaleDir,optional"`
    DB        database.MySQLConfig       `json:"DB"`
    Redis     cache.RedisConfig          `json:"Redis"`
    Kafka     queue.KafkaConfig          `json:"Kafka,optional"`
    Pprof     PprofConfig                `json:"Pprof,optional"`
}

type PprofConfig struct {
    Enabled bool `json:"Enabled,default=false"`
    Port    int  `json:"Port,default=6060"`
}
```

- [ ] **Step 2: 更新 dev.yaml**

```yaml
Name: iam
Host: 0.0.0.0
Port: 8888
Timeout: 10000
Log:
  Level: debug
  Stat: false
  Encoding: plain

LocaleDir: ./locales

DB:
  Host: localhost
  Port: 35069
  User: root
  Password: rootpassword
  Database: iam

Redis:
  Host: localhost
  Port: 33308
  Password: ""
  DB: 0

Kafka:
  Brokers:
    - localhost:39092
  Topic: iam-audit-logs

Pprof:
  Enabled: true
  Port: 6060
```

- [ ] **Step 3: 验证编译通过**

```bash
go build ./...
```

期望：编译通过，无报错。

- [ ] **Step 4: Commit**

```bash
git add internal/config/config.go etc/dev.yaml
git commit -m "feat: add DB, Redis, Kafka config to Config struct and dev.yaml"
```

---

### Task 5: 错误码与常量定义

**Files:**
- Create: `internal/constant/errors.go`

- [ ] **Step 1: 定义错误码常量**

```go
// internal/constant/errors.go
package constant

import "errors"

// 业务错误码（模块2位 + 业务3位）
const (
    CodeOK          = 0
    CodeAuthFailed  = 10001 // 用户名或密码错误
    CodeAuthLocked  = 10002 // 账号被锁定
    CodeAuthMFAFail = 10003 // MFA 验证失败
    CodeAuthCodeErr   = 10004 // 验证码错误
    CodeAuthCodeExp   = 10005 // 验证码已过期

    CodeTokenExpired  = 11001 // Token 已过期
    CodeTokenInvalid = 11002 // Token 无效
    CodeTokenRevoked  = 11003 // Token 已撤销
    CodeTokenRefresh  = 11004 // Refresh Token 无效

    CodeUserNotFound    = 20001 // 用户不存在
    CodeUserEmailExists = 20002 // 邮箱已存在
    CodeUserDisabled      = 20003 // 用户已禁用
    CodeUserPasswordPolicy = 20004 // 密码不满足策略

    CodeRoleNotFound   = 30001 // 角色不存在
    CodeRoleCodeExists = 30002 // 角色编码重复
    CodeRoleSoDConflict = 30003 // 权限冲突（SoD）
    CodeRolePermDenied = 30004 // 权限不足

    CodeTenantNotFound   = 40001 // 租户不存在
    CodeTenantExpired    = 40002 // 租户已过期
    CodeTenantDisabled   = 40003 // 租户已禁用
    CodeTenantQuotaExceed = 40004 // 超出配额

    CodeAppNotFound    = 50001 // 应用不存在
    CodeAppDisabled    = 50002 // 应用已禁用
    CodeAppNotAuthorized = 50003 // 用户未授权此应用

    CodeClientNotFound = 60001 // 客户端不存在
    CodeClientAKSKInvalid = 60002 // AK/SK 无效
    CodeClientDisabled = 60003 // 客户端已禁用

    CodeInternalError   = 99001 // 内部错误
    CodeDBError         = 99002 // 数据库错误
    CodeRedisError      = 99003 // Redis 错误
    CodeKafkaError      = 99004 // Kafka 错误
)

// 常见错误 sentinel
var (
    ErrRecordNotFound = errors.New("record not found")
    ErrDuplicateEntry = errors.New("duplicate entry")
    ErrTenantExpired  = errors.New("tenant expired")
    ErrTenantDisabled = errors.New("tenant disabled")
)
```

- [ ] **Step 2: Commit**

```bash
git add internal/constant/errors.go
git commit -m "feat: add error code constants for all modules"
```

---

### Task 6: 中间件层 — 租户隔离 / 认证 / 审计 / 错误处理

**Files:**
- Create: `internal/middleware/tenant.go`
- Create: `internal/middleware/auth.go`
- Create: `internal/middleware/audit.go`
- Create: `internal/middleware/error.go`

- [ ] **Step 1: 租户隔离中间件**

```go
// internal/middleware/tenant.go
package middleware

import (
    "context"
    "net/http"
    "strconv"

    "iam/internal/constant"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// contextKey 用于存储上下文值的键
type contextKey string

const tenantIDKey contextKey = "tenant_id"

// TenantMiddleware 租户隔离中间件
// 从 JWT Token 中提取 tenant_id 并注入到 context
func TenantMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // 尝试从 header 获取 tenant_id（后续由 AuthMiddleware 从 JWT 提取）
            tenantIDStr := r.Header.Get("X-Tenant-ID")
            if tenantIDStr == "" {
                // 暂不强制校验（Auth 中间件未实现），仅注入默认值
                // TODO: Auth 中间件实现后，从 JWT claims 中提取 tenant_id
                ctx := context.WithValue(r.Context(), tenantIDKey, int64(0))
                next(w, r.WithContext(ctx))
                return
            }

            tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
            if err != nil {
                logx.Errorf("invalid tenant_id: %s", tenantIDStr)
                httpx.WriteJson(w, http.StatusBadRequest, map[string]any{
                    "code":    constant.CodeTenantNotFound,
                    "message": "invalid tenant_id",
                    "data":    nil,
                })
                return
            }

            ctx := context.WithValue(r.Context(), tenantIDKey, tenantID)
            next(w, r.WithContext(ctx))
        }
    }
}

// GetTenantID 从 context 中提取 tenant_id
func GetTenantID(ctx context.Context) int64 {
    if v := ctx.Value(tenantIDKey); v != nil {
        return v.(int64)
    }
    return 0
}
```

- [ ] **Step 2: JWT 认证中间件（骨架）**

```go
// internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"

    "iam/internal/constant"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// AuthMiddleware JWT 认证中间件
// 验证 Authorization header 中的 Bearer Token
func AuthMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                httpx.WriteJson(w, http.StatusUnauthorized, map[string]any{
                    "code":    constant.CodeTokenInvalid,
                    "message": "missing Authorization header",
                    "data":    nil,
                })
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            if token == authHeader {
                httpx.WriteJson(w, http.StatusUnauthorized, map[string]any{
                    "code":    constant.CodeTokenInvalid,
                    "message": "invalid Authorization format, expected 'Bearer <token>'",
                    "data":    nil,
                })
                return
            }

            // TODO: 实际 JWT 验证逻辑（解析 token、验证签名、提取 claims）
            // 当前仅检查格式，后续 Task 中完善

            logx.Debugf("auth middleware: token format ok (stub validation)")
            next(w, r)
        }
    }
}

// SkipAuthPaths 不需要认证的路径
var SkipAuthPaths = map[string]bool{
    "/api/v1/auth/login":            true,
    "/api/v1/auth/register":         true,
    "/api/v1/auth/password/reset":   true,
    "/api/v1/auth/code/send":        true,
    "/api/v1/auth/code/login":       true,
    "/health":                       true,
    "/api/v1/clients/token":         true,
}
```

- [ ] **Step 3: 审计日志中间件（骨架）**

```go
// internal/middleware/audit.go
package middleware

import (
    "net/http"
    "time"

    "github.com/zeromicro/go-zero/core/logx"
)

// AuditMiddleware 审计日志中间件（骨架）
// 记录请求信息，后续通过 Kafka 异步写入审计日志
func AuditMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // TODO: 通过 responseWriter wrapper 捕获响应状态码
            next(w, r)

            duration := time.Since(start)
            tenantID := GetTenantID(r.Context())

            logx.WithDuration(duration).Infof(
                "[audit] method=%s path=%s tenant_id=%d ip=%s duration=%s",
                r.Method, r.URL.Path, tenantID, r.RemoteAddr, duration,
            )

            // TODO: 通过 Kafka Producer 异步发送审计日志
        }
    }
}
```

- [ ] **Step 4: 全局错误处理中间件**

```go
// internal/middleware/error.go
package middleware

import (
    "net/http"

    "iam/internal/constant"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// ErrorMiddleware 全局错误处理中间件
func ErrorMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logx.Errorf("panic recovered: %v", err)
                    httpx.WriteJson(w, http.StatusInternalServerError, map[string]any{
                        "code":    constant.CodeInternalError,
                        "message": "internal server error",
                        "data":    nil,
                    })
                }
            }()

            next(w, r)
        }
    }
}

// WriteSuccess 统一成功响应
func WriteSuccess(w http.ResponseWriter, data any) {
    httpx.WriteJson(w, http.StatusOK, map[string]any{
        "code":    constant.CodeOK,
        "message": "success",
        "data":    data,
    })
}

// WriteError 统一错误响应
func WriteError(w http.ResponseWriter, httpStatus int, code int, message string) {
    httpx.WriteJson(w, httpStatus, map[string]any{
        "code":    code,
        "message": message,
        "data":    nil,
    })
}
```

- [ ] **Step 5: Commit**

```bash
git add internal/middleware/
git commit -m "feat: add middleware layer (tenant, auth, audit, error handling)"
```

---

### Task 7: Tenant DTO — 请求/响应结构体

**Files:**
- Create: `internal/dto/tenant/tenant.go`

- [ ] **Step 1: 定义租户 DTO**

```go
// internal/dto/tenant/tenant.go
package tenant

import "time"

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
    Name      string `json:"name" validate:"required,min=1,max=100"`
    MaxUsers  int    `json:"max_users" validate:"min=1"`
    MaxApps   int    `json:"max_apps" validate:"min=1"`
    ExpireAt  string `json:"expire_at" validate:"omitempty,datetime=2006-01-02 15:04:05"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
    Name     string `json:"name" validate:"omitempty,min=1,max=100"`
    MaxUsers *int   `json:"max_users" validate:"omitempty,min=1"`
    MaxApps  *int   `json:"max_apps" validate:"omitempty,min=1"`
    ExpireAt string `json:"expire_at" validate:"omitempty,datetime=2006-01-02 15:04:05"`
}

// TenantResponse 租户响应
type TenantResponse struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Status    int8      `json:"status"`
    MaxUsers  int       `json:"max_users"`
    MaxApps   int       `json:"max_apps"`
    ExpireAt  time.Time `json:"expire_at"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// TenantListResponse 租户列表响应
type TenantListResponse struct {
    Items    []TenantResponse `json:"items"`
    Total    int64            `json:"total"`
    Page     int              `json:"page"`
    PageSize int              `json:"page_size"`
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/dto/tenant/tenant.go
git commit -m "feat: add tenant DTO definitions"
```

---

### Task 8: Tenant Repository — 数据访问层

**Files:**
- Create: `internal/repository/tenant_repo.go`

- [ ] **Step 1: 实现 TenantRepository**

```go
// internal/repository/tenant_repo.go
package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "iam/internal/entity"
    "iam/internal/middleware"
)

// TenantRepository 租户数据访问层
type TenantRepository struct {
    db *sql.DB
}

// NewTenantRepository 创建租户 Repository
func NewTenantRepository(db *sql.DB) *TenantRepository {
    return &TenantRepository{db: db}
}

// Create 创建租户
func (r *TenantRepository) Create(ctx context.Context, t *entity.Tenant) error {
    query := `INSERT INTO tenants (name, status, max_users, max_apps, expire_at) VALUES (?, ?, ?, ?, ?)`
    result, err := r.db.ExecContext(ctx, query,
        t.Name, t.Status, t.MaxUsers, t.MaxApps, t.ExpireAt)
    if err != nil {
        return fmt.Errorf("create tenant: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("get last insert id: %w", err)
    }
    t.ID = id
    return nil
}

// GetByID 根据 ID 获取租户
func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*entity.Tenant, error) {
    query := `SELECT id, name, status, max_users, max_apps, expire_at, created_at, updated_at FROM tenants WHERE id = ?`
    var t entity.Tenant
    var expireAt sql.NullTime
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &t.ID, &t.Name, &t.Status, &t.MaxUsers, &t.MaxApps,
        &expireAt, &t.CreatedAt, &t.UpdatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("get tenant by id: %w", err)
    }
    if expireAt.Valid {
        t.ExpireAt = expireAt.Time
    }
    return &t, nil
}

// List 获取租户列表（分页）
func (r *TenantRepository) List(ctx context.Context, page, pageSize int) ([]entity.Tenant, int64, error) {
    offset := (page - 1) * pageSize
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    // 查询总数
    var total int64
    countQuery := `SELECT COUNT(*) FROM tenants`
    if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("count tenants: %w", err)
    }

    // 查询列表
    query := `SELECT id, name, status, max_users, max_apps, expire_at, created_at, updated_at FROM tenants ORDER BY id DESC LIMIT ? OFFSET ?`
    rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("list tenants: %w", err)
    }
    defer rows.Close()

    var tenants []entity.Tenant
    for rows.Next() {
        var t entity.Tenant
        var expireAt sql.NullTime
        if err := rows.Scan(&t.ID, &t.Name, &t.Status, &t.MaxUsers, &t.MaxApps, &expireAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
            return nil, 0, fmt.Errorf("scan tenant: %w", err)
        }
        if expireAt.Valid {
            t.ExpireAt = expireAt.Time
        }
        tenants = append(tenants, t)
    }

    return tenants, total, nil
}

// Update 更新租户
func (r *TenantRepository) Update(ctx context.Context, t *entity.Tenant) error {
    query := `UPDATE tenants SET name=?, status=?, max_users=?, max_apps=?, expire_at=? WHERE id=?`
    _, err := r.db.ExecContext(ctx, query,
        t.Name, t.Status, t.MaxUsers, t.MaxApps, t.ExpireAt, t.ID)
    if err != nil {
        return fmt.Errorf("update tenant: %w", err)
    }
    return nil
}

// Delete 删除租户
func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM tenants WHERE id = ?`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("delete tenant: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("tenant not found: id=%d", id)
    }
    return nil
}

// UpdateStatus 更新租户状态
func (r *TenantRepository) UpdateStatus(ctx context.Context, id int64, status int8) error {
    query := `UPDATE tenants SET status=? WHERE id=?`
    result, err := r.db.ExecContext(ctx, query, status, id)
    if err != nil {
        return fmt.Errorf("update tenant status: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("tenant not found: id=%d", id)
    }
    return nil
}

// CheckNameExists 检查租户名称是否已存在
func (r *TenantRepository) CheckNameExists(ctx context.Context, name string, excludeID int64) (bool, error) {
    var count int64
    query := `SELECT COUNT(*) FROM tenants WHERE name = ?`
    if excludeID > 0 {
        query += ` AND id != ?`
        if err := r.db.QueryRowContext(ctx, query, name, excludeID).Scan(&count); err != nil {
            return false, fmt.Errorf("check name exists: %w", err)
        }
    } else {
        if err := r.db.QueryRowContext(ctx, query, name).Scan(&count); err != nil {
            return false, fmt.Errorf("check name exists: %w", err)
        }
    }
    return count > 0, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/repository/tenant_repo.go
git commit -m "feat: add TenantRepository with CRUD operations"
```

---

### Task 9: Tenant Service — 业务逻辑层

**Files:**
- Create: `internal/service/tenant/tenant.go`

- [ ] **Step 1: 实现 TenantService**

```go
// internal/service/tenant/tenant.go
package tenant

import (
    "context"
    "fmt"
    "time"

    "iam/internal/constant"
    "iam/internal/dto/tenant"
    "iam/internal/entity"
    "iam/internal/middleware"
    "iam/internal/repository"
)

// TenantService 租户业务逻辑
type TenantService struct {
    repo *repository.TenantRepository
}

// NewTenantService 创建租户 Service
func NewTenantService(repo *repository.TenantRepository) *TenantService {
    return &TenantService{repo: repo}
}

// Create 创建租户
func (s *TenantService) Create(ctx context.Context, req tenant.CreateTenantRequest) (*entity.Tenant, int, string, error) {
    // 检查名称唯一性
    exists, err := s.repo.CheckNameExists(ctx, req.Name, 0)
    if err != nil {
        return nil, 0, "", err
    }
    if exists {
        return nil, 409, "tenant name already exists", fmt.Errorf("tenant name '%s' already exists", req.Name)
    }

    var expireAt time.Time
    if req.ExpireAt != "" {
        expireAt, err = time.Parse("2006-01-02 15:04:05", req.ExpireAt)
        if err != nil {
            return nil, 400, "invalid expire_at format, expected 'YYYY-MM-DD HH:MM:SS'", err
        }
    }

    t := &entity.Tenant{
        Name:     req.Name,
        Status:   entity.TenantStatusActive,
        MaxUsers: req.MaxUsers,
        MaxApps:  req.MaxApps,
        ExpireAt: expireAt,
    }

    if err := s.repo.Create(ctx, t); err != nil {
        return nil, 500, "failed to create tenant", err
    }

    return t, 0, "", nil
}

// GetByID 获取租户详情
func (s *TenantService) GetByID(ctx context.Context, id int64) (*entity.Tenant, int, string, error) {
    t, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, 500, "failed to get tenant", err
    }
    if t == nil {
        return nil, 404, "tenant not found", fmt.Errorf("tenant id=%d not found", id)
    }
    return t, 0, "", nil
}

// List 获取租户列表
func (s *TenantService) List(ctx context.Context, page, pageSize int) ([]entity.Tenant, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    return s.repo.List(ctx, page, pageSize)
}

// Update 更新租户
func (s *TenantService) Update(ctx context.Context, id int64, req tenant.UpdateTenantRequest) (*entity.Tenant, int, string, error) {
    t, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, 500, "failed to get tenant", err
    }
    if t == nil {
        return nil, 404, "tenant not found", fmt.Errorf("tenant id=%d not found", id)
    }

    if req.Name != "" {
        exists, err := s.repo.CheckNameExists(ctx, req.Name, id)
        if err != nil {
            return nil, 500, "failed to check tenant name", err
        }
        if exists {
            return nil, 409, "tenant name already exists", fmt.Errorf("tenant name '%s' already exists", req.Name)
        }
        t.Name = req.Name
    }
    if req.MaxUsers != nil {
        t.MaxUsers = *req.MaxUsers
    }
    if req.MaxApps != nil {
        t.MaxApps = *req.MaxApps
    }
    if req.ExpireAt != "" {
        t.ExpireAt, err = time.Parse("2006-01-02 15:04:05", req.ExpireAt)
        if err != nil {
            return nil, 400, "invalid expire_at format", err
        }
    }

    if err := s.repo.Update(ctx, t); err != nil {
        return nil, 500, "failed to update tenant", err
    }

    return t, 0, "", nil
}

// Delete 删除租户
func (s *TenantService) Delete(ctx context.Context, id int64) (int, string, error) {
    err := s.repo.Delete(ctx, id)
    if err != nil {
        return 500, "failed to delete tenant", err
    }
    return 0, "", nil
}

// UpdateStatus 更新租户状态
func (s *TenantService) UpdateStatus(ctx context.Context, id int64, status int8) (int, string, error) {
    err := s.repo.UpdateStatus(ctx, id, status)
    if err != nil {
        return 500, "failed to update tenant status", err
    }
    return 0, "", nil
}

// toResponse 实体转响应 DTO
func toResponse(t *entity.Tenant) tenant.TenantResponse {
    return tenant.TenantResponse{
        ID:        t.ID,
        Name:      t.Name,
        Status:    t.Status,
        MaxUsers:  t.MaxUsers,
        MaxApps:   t.MaxApps,
        ExpireAt:  t.ExpireAt,
        CreatedAt: t.CreatedAt,
        UpdatedAt: t.UpdatedAt,
    }
}

// ToResponseList 实体列表转响应 DTO
func ToResponseList(items []entity.Tenant, total int64, page, pageSize int) tenant.TenantListResponse {
    respItems := make([]tenant.TenantResponse, len(items))
    for i, t := range items {
        respItems[i] = toResponse(&t)
    }
    return tenant.TenantListResponse{
        Items:    respItems,
        Total:    total,
        Page:     page,
        PageSize: pageSize,
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/service/tenant/
git commit -m "feat: add TenantService with CRUD business logic"
```

---

### Task 10: Tenant Handler — HTTP 处理器

**Files:**
- Create: `internal/handler/tenant/tenant.go`

- [ ] **Step 1: 实现租户 HTTP 处理器**

```go
// internal/handler/tenant/tenant.go
package tenant

import (
    "encoding/json"
    "net/http"
    "strconv"

    "iam/internal/constant"
    "iam/internal/dto/tenant"
    "iam/internal/middleware"
    "iam/internal/service/tenant"

    "github.com/zeromicro/go-zero/core/logx"
)

// TenantHandler 租户 HTTP 处理器
type TenantHandler struct {
    svc *tenant.TenantService
}

// NewTenantHandler 创建租户 Handler
func NewTenantHandler(svc *tenant.TenantService) *TenantHandler {
    return &TenantHandler{svc: svc}
}

// CreateTenant 创建租户 POST /api/v1/tenants
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
    var req tenant.CreateTenantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
        return
    }

    t, httpStatus, msg, err := h.svc.Create(r.Context(), req)
    if err != nil {
        logx.Errorf("create tenant failed: %v", err)
        middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
        return
    }

    middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// GetTenant 获取租户详情 GET /api/v1/tenants/:id
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
        return
    }

    t, httpStatus, msg, err := h.svc.GetByID(r.Context(), id)
    if err != nil {
        logx.Errorf("get tenant failed: %v", err)
        middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
        return
    }

    middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// ListTenants 获取租户列表 GET /api/v1/tenants
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

    items, total, err := h.svc.List(r.Context(), page, pageSize)
    if err != nil {
        logx.Errorf("list tenants failed: %v", err)
        middleware.WriteError(w, http.StatusInternalServerError, constant.CodeInternalError, "failed to list tenants")
        return
    }

    middleware.WriteSuccess(w, tenant.ToResponseList(items, total, page, pageSize))
}

// UpdateTenant 更新租户 PUT /api/v1/tenants/:id
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
        return
    }

    var req tenant.UpdateTenantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
        return
    }

    t, httpStatus, msg, err := h.svc.Update(r.Context(), id, req)
    if err != nil {
        logx.Errorf("update tenant failed: %v", err)
        middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
        return
    }

    middleware.WriteSuccess(w, middleware.ToResponse(t))
}

// DeleteTenant 删除租户 DELETE /api/v1/tenants/:id
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
        return
    }

    httpStatus, msg, err := h.svc.Delete(r.Context(), id)
    if err != nil {
        logx.Errorf("delete tenant failed: %v", err)
        middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
        return
    }

    middleware.WriteSuccess(w, nil)
}

// UpdateTenantStatus 更新租户状态 PUT /api/v1/tenants/:id/status
func (h *TenantHandler) UpdateTenantStatus(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeTenantNotFound, "invalid tenant id")
        return
    }

    var req struct {
        Status int8 `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        middleware.WriteError(w, http.StatusBadRequest, constant.CodeAuthFailed, "invalid request body")
        return
    }

    httpStatus, msg, err := h.svc.UpdateStatus(r.Context(), id, req.Status)
    if err != nil {
        logx.Errorf("update tenant status failed: %v", err)
        middleware.WriteError(w, httpStatus, constant.CodeInternalError, msg)
        return
    }

    middleware.WriteSuccess(w, nil)
}
```

**Note:** `middleware.ToResponse` 需要在 `middleware` 包中添加一个通用的 entity→map 转换函数：

```go
// internal/middleware/response.go
package middleware

import (
    "encoding/json"
    "net/http"
)

// WriteSuccess 统一成功响应
func WriteSuccess(w http.ResponseWriter, data any) {
    httpx.WriteJson(w, http.StatusOK, map[string]any{
        "code":    constant.CodeOK,
        "message": "success",
        "data":    data,
    })
}

// ToResponse 将实体转换为响应 map（通用转换）
func ToResponse(entity any) map[string]any {
    if entity == nil {
        return nil
    }
    data, _ := json.Marshal(entity)
    var result map[string]any
    json.Unmarshal(data, &result)
    return result
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/handler/tenant/ internal/middleware/response.go
git commit -m "feat: add tenant HTTP handlers for CRUD operations"
```

---

### Task 11: Tenant Routes — 路由注册

**Files:**
- Create: `internal/routes/tenant/tenant.go`

- [ ] **Step 1: 实现租户路由注册**

```go
// internal/routes/tenant/tenant.go
package tenant

import (
    "net/http"

    "iam/internal/handler/tenant"
    "iam/internal/svc"

    "github.com/zeromicro/go-zero/rest"
)

// TenantRouter 租户路由注册器
type TenantRouter struct {
    server  *rest.Server
    svcCtx  *svc.ServiceContext
    handler *tenant.TenantHandler
}

// NewTenantRouter 创建租户路由注册器
func NewTenantRouter(server *rest.Server, svcCtx *svc.ServiceContext) *TenantRouter {
    return &TenantRouter{
        server:  server,
        svcCtx:  svcCtx,
        handler: tenant.NewTenantHandler(svcCtx.TenantService),
    }
}

// Register 注册租户路由
func (r *TenantRouter) Register() {
    // 基础路径: /api/v1/tenants
    r.server.AddRoutes(
        []rest.Route{
            {Method: http.MethodGet, Path: "/tenants", Handler: r.handler.ListTenants},
            {Method: http.MethodGet, Path: "/tenants/:id", Handler: r.handler.GetTenant},
            {Method: http.MethodPost, Path: "/tenants", Handler: r.handler.CreateTenant},
            {Method: http.MethodPut, Path: "/tenants/:id", Handler: r.handler.UpdateTenant},
            {Method: http.MethodDelete, Path: "/tenants/:id", Handler: r.handler.DeleteTenant},
            {Method: http.MethodPut, Path: "/tenants/:id/status", Handler: r.handler.UpdateTenantStatus},
        },
        rest.WithPrefix("/api/v1"),
    )
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/routes/tenant/
git commit -m "feat: add tenant route registration"
```

---

### Task 12: 更新 ServiceContext 和主路由 — 注入所有依赖

**Files:**
- Modify: `internal/svc/servicecontext.go`
- Modify: `internal/routes/routes.go`

- [ ] **Step 1: 更新 ServiceContext**

```go
// internal/svc/servicecontext.go
package svc

import (
    "context"
    "crypto/rsa"
    "database/sql"
    "fmt"
    "iam/infra/cache"
    "iam/infra/database"
    "iam/infra/queue"
    "iam/internal/config"
    "iam/internal/repository"
    tenantsvc "iam/internal/service/tenant"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
    Config config.Config
    Logger logx.Logger

    // Infra
    DB           *sql.DB
    Redis        *redis.Redis
    KafkaProducer *queue.KafkaProducer

    // Repositories
    TenantRepo *repository.TenantRepository

    // Services
    TenantService *tenantsvc.TenantService

    // JWT（骨架，后续完善）
    JWTKey       *rsa.PrivateKey
    JWTPubKey    *rsa.PublicKey
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
    ctx := context.Background()
    logger := logx.WithContext(ctx)

    // 初始化 MySQL
    db, err := database.NewMySQL(c.DB)
    if err != nil {
        return nil, fmt.Errorf("init mysql: %w", err)
    }
    logger.Info("mysql connected")

    // 初始化 Redis
    redisClient, err := cache.NewRedis(c.Redis)
    if err != nil {
        return nil, fmt.Errorf("init redis: %w", err)
    }
    logger.Info("redis connected")

    // 初始化 Kafka（骨架）
    kafkaProducer, err := queue.NewKafkaProducer(c.Kafka)
    if err != nil {
        return nil, fmt.Errorf("init kafka: %w", err)
    }
    logger.Info("kafka producer initialized (stub)")

    // 初始化 Repository
    tenantRepo := repository.NewTenantRepository(db)

    // 初始化 Service
    tenantSvc := tenantsvc.NewTenantService(tenantRepo)

    return &ServiceContext{
        Config:        c,
        Logger:        logger,
        DB:            db,
        Redis:         redisClient,
        KafkaProducer: kafkaProducer,
        TenantRepo:    tenantRepo,
        TenantService: tenantSvc,
    }, nil
}

func (sc *ServiceContext) Close() error {
    if sc.DB != nil {
        sc.DB.Close()
    }
    if sc.Redis != nil {
        // go-zero redis 没有 Close 方法
    }
    if sc.KafkaProducer != nil {
        sc.KafkaProducer.Close()
    }
    return nil
}
```

- [ ] **Step 2: 更新主路由注册**

```go
// internal/routes/routes.go
package routes

import (
    "iam/internal/routes/health"
    "iam/internal/routes/tenant"
    "iam/internal/svc"

    "github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 注册健康检查相关路由
    healthRouter := health.NewHealthRouter(server, serverCtx)
    healthRouter.Register()

    // 注册租户管理路由
    tenantRouter := tenant.NewTenantRouter(server, serverCtx)
    tenantRouter.Register()
}
```

- [ ] **Step 3: 验证编译通过并测试**

```bash
go build ./...
```

期望：编译通过。

- [ ] **Step 4: Commit**

```bash
git add internal/svc/servicecontext.go internal/routes/routes.go
git commit -m "feat: wire up ServiceContext with DB/Redis/Kafka and tenant service, register tenant routes"
```

---

### Task 13: Docker Compose — 确保开发环境可运行

**Files:**
- Modify: `docker-compose.yml`

- [ ] **Step 1: 确保 MySQL 端口映射正确**

当前 `docker-compose.yml` 已配置 MySQL 端口 `35069:3306`，Redis 端口 `33308`（已注释，需要取消注释以便本地连接），Kafka `PLAINTEXT://kafka-broker:9092`（需要添加主机端口映射）。

修改 `docker-compose.yml` 中 redis 和 kafka 部分：

```yaml
  redis:
    image: redis:8.6
    restart: always
    networks:
      - default
    ports:
      - "33308:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  kafka-broker:
    # ... 保持现有配置不变 ...
    ports:
      - "39092:9092"
    # ...
```

- [ ] **Step 2: Commit**

```bash
git add docker-compose.yml
git commit -m "fix: expose Redis and Kafka ports for local development"
```

---

### Task 14: CI 流水线 — GitHub Actions

**Files:**
- Create: `.github/workflows/01-ci.yaml`

- [ ] **Step 1: 创建 CI workflow**

```yaml
# .github/workflows/01-ci.yaml
name: ci-pipeline

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      - name: Run lint
        run: golangci-lint run ./...

  build:
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Build
        run: go build ./...

  unit-test:
    runs-on: ubuntu-latest
    needs: build
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Run unit tests
        run: go test ./... -race -cover -short

  integration-test:
    runs-on: ubuntu-latest
    needs: build
    services:
      mysql:
        image: mysql:9.6
        env:
          MYSQL_ROOT_PASSWORD: rootpassword
          MYSQL_DATABASE: iam
        ports:
          - "35069:3306"
        options: >-
          --health-cmd "mysqladmin ping -h localhost"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:8.6
        ports:
          - "33308:6379"
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Initialize database
        run: mysql -h 127.0.0.1 -P 35069 -uroot -prootpassword iam < sql/001_init.sql
      - name: Run integration tests
        run: go test -tags=integration ./... -race -v
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/01-ci.yaml
git commit -m "ci: add CI pipeline with lint, build, unit and integration tests"
```

---

### Task 15: 本地 CI 验证脚本

**Files:**
- Create: `scripts/ci-local.sh`

- [ ] **Step 1: 创建脚本**

```bash
#!/bin/bash
# scripts/ci-local.sh
# 本地 CI 验证脚本，在推送前运行
set -e

echo "=== Step 1: Build ==="
go build ./...
echo "✅ Build passed"

echo "=== Step 2: Unit tests ==="
go test ./... -race -cover -short
echo "✅ Unit tests passed"

echo ""
echo "=== All checks passed ==="
echo "Run 'go test -tags=integration ./...' with Docker services for full integration tests."
```

```bash
chmod +x scripts/ci-local.sh
```

- [ ] **Step 2: 运行验证**

```bash
bash scripts/ci-local.sh
```

期望：build 和 unit test 通过。

- [ ] **Step 3: Commit**

```bash
git add scripts/ci-local.sh
git commit -m "feat: add local CI validation script"
```

---

### Task 16: 租户集成测试

**Files:**
- Create: `internal/tests/integration/tenant_test.go`

- [ ] **Step 1: 编写集成测试**

```go
// internal/tests/integration/tenant_test.go
//go:build integration

package integration

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "iam/internal/dto/tenant"
    "iam/internal/entity"
    tenanthandler "iam/internal/handler/tenant"
    "iam/internal/middleware"
    tenantsvc "iam/internal/service/tenant"
    "iam/internal/repository"

    _ "github.com/go-sql-driver/mysql"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/zeromicro/go-zero/rest"
)

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
        "root", "rootpassword", "127.0.0.1", 35069, "iam")

    db, err := sql.Open("mysql", dsn)
    require.NoError(t, err)

    // 清理租户表
    _, err = db.Exec("DELETE FROM tenants")
    require.NoError(t, err)

    t.Cleanup(func() {
        _, _ = db.Exec("DELETE FROM tenants")
        db.Close()
    })

    return db
}

func setupTenantHandler(t *testing.T) (*tenanthandler.TenantHandler, *sql.DB) {
    t.Helper()
    db := setupTestDB(t)
    repo := repository.NewTenantRepository(db)
    svc := tenantsvc.NewTenantService(repo)
    return tenanthandler.NewTenantHandler(svc), db
}

func TestTenant_CreateAndGet(t *testing.T) {
    handler, _ := setupTestHandler(t)
    server := rest.MustNewServer(rest.RestConf{Port: 8888})
    defer server.Stop()

    // 注册路由
    server.AddRoutes([]rest.Route{
        {Method: http.MethodPost, Path: "/tenants", Handler: handler.CreateTenant},
        {Method: http.MethodGet, Path: "/tenants/:id", Handler: handler.GetTenant},
    }, rest.WithPrefix("/api/v1"))

    // 1. 创建租户
    createReq := tenant.CreateTenantRequest{
        Name:     "Test Tenant",
        MaxUsers: 100,
        MaxApps:  5,
    }
    body, _ := json.Marshal(createReq)
    req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    server.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)
    var createResp map[string]any
    json.Unmarshal(recorder.Body.Bytes(), &createResp)
    assert.Equal(t, float64(0), createResp["code"])
    assert.Equal(t, "Test Tenant", createResp["data"].(map[string]any)["name"])

    // 2. 获取租户详情
    data := createResp["data"].(map[string]any)
    tenantID := int64(data["id"].(float64))
    getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tenants/%d", tenantID), nil)
    getRecorder := httptest.NewRecorder()
    server.ServeHTTP(getRecorder, getReq)

    assert.Equal(t, http.StatusOK, getRecorder.Code)
    var getResp map[string]any
    json.Unmarshal(getRecorder.Body.Bytes(), &getResp)
    assert.Equal(t, float64(0), getResp["code"])
    assert.Equal(t, "Test Tenant", getResp["data"].(map[string]any)["name"])
}

func TestTenant_CreateDuplicateName(t *testing.T) {
    handler, _ := setupTestHandler(t)
    server := rest.MustNewServer(rest.RestConf{Port: 8888})
    defer server.Stop()

    server.AddRoutes([]rest.Route{
        {Method: http.MethodPost, Path: "/tenants", Handler: handler.CreateTenant},
    }, rest.WithPrefix("/api/v1"))

    createReq := tenant.CreateTenantRequest{Name: "Dup Tenant", MaxUsers: 10, MaxApps: 1}
    body, _ := json.Marshal(createReq)

    // 第一次创建
    req1 := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
    req1.Header.Set("Content-Type", "application/json")
    rec1 := httptest.NewRecorder()
    server.ServeHTTP(rec1, req1)
    assert.Equal(t, http.StatusOK, rec1.Code)

    // 第二次创建（重复名称）
    req2 := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader(body))
    req2.Header.Set("Content-Type", "application/json")
    rec2 := httptest.NewRecorder()
    server.ServeHTTP(rec2, req2)
    assert.Equal(t, http.StatusConflict, rec2.Code)
}

func TestTenant_GetNotFound(t *testing.T) {
    handler, _ := setupTestHandler(t)
    server := rest.MustNewServer(rest.RestConf{Port: 8888})
    defer server.Stop()

    server.AddRoutes([]rest.Route{
        {Method: http.MethodGet, Path: "/tenants/:id", Handler: handler.GetTenant},
    }, rest.WithPrefix("/api/v1"))

    req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/999999", nil)
    recorder := httptest.NewRecorder()
    server.ServeHTTP(recorder, req)

    assert.Equal(t, http.StatusNotFound, recorder.Code)
    var resp map[string]any
    json.Unmarshal(recorder.Body.Bytes(), &resp)
    assert.Equal(t, float64(40001), resp["code"]) // CodeTenantNotFound
}
```

- [ ] **Step 2: 添加 testify 依赖**

```bash
go get github.com/stretchr/testify@latest
```

- [ ] **Step 3: 运行集成测试**

```bash
# 确保 Docker 服务已启动
docker compose up -d mysql redis
docker compose exec -T mysql mysql -uroot -prootpassword < sql/001_init.sql

# 运行集成测试
go test -tags=integration ./internal/tests/integration/... -v
```

期望：3 个测试用例全部通过。

- [ ] **Step 4: Commit**

```bash
git add internal/tests/integration/tenant_test.go go.mod go.sum
git commit -m "test: add tenant integration tests (create, duplicate name, not found)"
```

---

## 计划完成标志

以上 16 个 Task 全部完成后，系统具备：

1. 完整的数据库 Schema（16 张表）
2. 全部实体定义（Go 结构体）
3. 基础设施连接（MySQL/Redis/Kafka）
4. 错误码体系与中间件框架
5. 租户 CRUD 完整闭环（DTO → Handler → Service → Repository → DB）
6. CI 流水线
7. 集成测试覆盖

后续 REQ 的实现将复用此模式。
