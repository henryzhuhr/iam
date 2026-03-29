# 数据库设计规范

> 最后更新：2026-03-29
> 适用范围：IAM 项目所有 MySQL 数据库设计

---

## 1. 命名规范

### 1.1 表名

- 小写，复数形式
- 下划线分隔
- 语义清晰

```sql
-- ✅ 正确
users
user_roles
tenants
applications

-- ❌ 错误
User        -- 大写
tbl_user    -- 前缀冗余
t_user      -- 前缀冗余
```

### 1.2 字段名

- 小写，下划线分隔
- 主键统一用 `id`
- 外键用 `表名_id` 格式

```sql
-- ✅ 正确
id              -- 主键
tenant_id       -- 外键
user_id         -- 外键
created_at      -- 创建时间
updated_at      -- 更新时间

-- ❌ 错误
userId          -- 驼峰式
tenantId        -- 驼峰式
createTime      -- 非标准命名
```

### 1.3 索引名

| 索引类型 | 命名规范 | 示例 |
|----------|----------|------|
| 主键索引 | `PRIMARY` | `PRIMARY` |
| 唯一索引 | `uk_字段名` | `uk_tenant_email` |
| 普通索引 | `idx_字段名` | `idx_status` |
| 联合索引 | `idx_字段 1_字段 2` | `idx_tenant_status` |

---

## 2. 表设计规范

### 2.2 必填字段

每张表必须包含以下字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT | 主键，自增 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

### 2.3 多租户字段

租户相关表必须包含：

| 字段 | 类型 | 说明 |
|------|------|------|
| `tenant_id` | BIGINT | 租户 ID，数据隔离 |

### 2.4 示例

```sql
CREATE TABLE users (
    id              BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
    tenant_id       BIGINT NOT NULL COMMENT '租户 ID',
    email           VARCHAR(100) NOT NULL COMMENT '邮箱',
    password_hash   VARCHAR(255) NOT NULL COMMENT '密码哈希',
    status          TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_tenant_email (tenant_id, email),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

---

## 3. 字段类型规范

### 3.1 整数类型

| 类型 | 字节 | 范围 | 使用场景 |
|------|------|------|----------|
| TINYINT | 1 | -128 ~ 127 | 状态、布尔值 |
| SMALLINT | 2 | -32768 ~ 32767 | 小数值计数 |
| INT | 4 | ±21 亿 | 一般计数 |
| BIGINT | 8 | ±922 亿亿 | 主键、用户 ID |

### 3.2 字符串类型

| 类型 | 说明 | 使用场景 |
|------|------|----------|
| VARCHAR(n) | 变长字符串 | 邮箱、名称等 |
| TEXT | 长文本 | 描述、内容 |
| JSON | JSON 数据 | 配置、扩展字段 |

### 3.3 时间类型

| 类型 | 说明 | 使用场景 |
|------|------|----------|
| DATETIME | 日期时间 | 创建时间、更新时间 |
| TIMESTAMP | 时间戳 | 过期时间 |

---

## 4. 索引规范

### 4.1 索引设计原则

- 主键优先使用自增 BIGINT
- 外键建立索引
- 查询频繁的字段建立索引
- 避免过多索引（单表不超过 5 个）

### 4.2 联合索引

- 最左前缀原则
- 高频查询字段在前

```sql
-- ✅ 正确：联合索引
KEY idx_tenant_status (tenant_id, status)

-- 查询可以使用索引
SELECT * FROM users WHERE tenant_id = 100 AND status = 1;
SELECT * FROM users WHERE tenant_id = 100;

-- 查询不能使用索引
SELECT * FROM users WHERE status = 1;
```

---

## 5. SQL 编写规范

### 5.1 查询语句

```sql
-- ✅ 正确：明确列名
SELECT id, email, name FROM users WHERE tenant_id = 100;

-- ❌ 错误：使用 SELECT *
SELECT * FROM users;
```

### 5.2 插入语句

```sql
-- ✅ 正确：明确列名和值
INSERT INTO users (tenant_id, email, password_hash)
VALUES (100, 'user@example.com', 'hashed_password');
```

### 5.3 更新语句

```sql
-- ✅ 正确：带 WHERE 条件，包含 updated_at
UPDATE users 
SET name = '新名字', updated_at = NOW()
WHERE id = 12345;
```

### 5.4 删除语句

```sql
-- ✅ 正确：软删除（推荐）
UPDATE users SET deleted_at = NOW() WHERE id = 12345;

-- ❌ 错误：硬删除（除非必要）
DELETE FROM users WHERE id = 12345;
```

---

## 6. 数据变更流程

### 6.1 DDL 变更

1. 编写 SQL 脚本，放在 `sql/` 目录
2. 脚本命名：`V 版本号_描述.sql`
3. 评审 SQL 脚本
4. 在测试环境执行验证
5. 生产环境执行

### 6.2 脚本示例

```sql
-- V1_001_add_user_table.sql
-- 创建用户表

CREATE TABLE IF NOT EXISTS users (
    id              BIGINT NOT NULL AUTO_INCREMENT,
    tenant_id       BIGINT NOT NULL,
    email           VARCHAR(100) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    status          TINYINT NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_tenant_email (tenant_id, email),
    KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

---

## 7. 参考链接

- MySQL 官方文档：https://dev.mysql.com/doc/
- 阿里巴巴 Java 开发手册 - 数据库部分
- 《高性能 MySQL》
