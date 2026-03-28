# REQ-009 操作审计日志

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 3 人天 |
| **关联用户故事** | US-020 |

**背景：** 需要记录用户的敏感操作，满足安全审计和合规要求，支持事后追溯和责任认定。

**目标：**

- 记录所有敏感操作（用户管理、权限变更、配置修改）
- 记录操作人、操作时间、操作内容、IP 地址
- 支持按条件查询审计日志
- 支持日志导出和统计分析
- 日志可配置保留期限（默认 180 天）

**功能描述：**

### 1. 日志记录范围

记录以下操作类型：

| 模块 | 操作类型 | 示例 |
|------|----------|------|
| 用户管理 | USER_CREATE, USER_UPDATE, USER_DELETE, USER_ENABLE, USER_DISABLE | 创建用户、禁用用户 |
| 角色管理 | ROLE_CREATE, ROLE_UPDATE, ROLE_DELETE, ROLE_ASSIGN | 创建角色、分配角色 |
| 权限管理 | PERMISSION_ASSIGN, PERMISSION_REVOKE | 分配权限、撤销权限 |
| 租户管理 | TENANT_CREATE, TENANT_UPDATE, TENANT_FREEZE, TENANT_QUOTA | 创建租户、冻结租户 |
| 密码管理 | PASSWORD_CHANGE, PASSWORD_RESET, PASSWORD_EXPIRE | 修改密码、重置密码 |
| MFA 管理 | MFA_BIND, MFA_UNBIND, MFA_VERIFY | 绑定 MFA、解绑 MFA |
| 系统配置 | CONFIG_UPDATE, POLICY_UPDATE | 修改配置、更新策略 |

### 2. 日志内容

每条日志包含以下字段：

| 字段 | 说明 | 示例 |
|------|------|------|
| operation_id | 操作唯一标识 | `op_xxxxxxxxx` |
| tenant_id | 租户 ID | `1001` |
| user_id | 操作人 ID | `2001` |
| user_name | 操作人姓名 | `张三` |
| operation_type | 操作类型 | `USER_CREATE` |
| resource_type | 资源类型 | `USER` |
| resource_id | 资源 ID | `2002` |
| resource_name | 资源名称 | `李四` |
| action | 具体动作 | `CREATE` |
| request_body | 请求参数（脱敏） | `{email: "li***@example.com"}` |
| response_status | 响应状态 | `SUCCESS` / `FAILURE` |
| error_message | 错误信息 | 失败时的错误消息 |
| ip_address | 操作 IP | `192.168.1.100` |
| user_agent | 用户代理 | `Mozilla/5.0...` |
| device_info | 设备信息 | `Chrome 120 / macOS` |
| location | 地理位置 | `北京市海淀区` |
| created_at | 操作时间 | `2026-03-25 10:30:00` |

### 3. 日志采集方式

1. **注解方式**：通过 `@AuditLog` 注解标记需要记录的操作
2. **AOP 切面**：自动拦截注解标记的方法
3. **异步写入**：通过 Kafka 异步写入，不阻塞主流程
4. **失败降级**：写入失败时降级到本地文件

### 4. 日志查询功能

支持以下查询条件：

| 条件 | 说明 |
|------|------|
| 时间范围 | 支持自定义起止时间（最大 31 天） |
| 操作人 | 按用户 ID 或用户名搜索 |
| 操作类型 | 单选或多选操作类型 |
| 资源类型 | 按资源类型筛选 |
| 操作结果 | 成功/失败 |
| IP 地址 | 按 IP 或 IP 段搜索 |
| 关键词 | 全文检索资源名称、操作详情 |

查询结果支持：
- 分页加载
- 按时间倒序/正序排序
- 导出 CSV/Excel

### 5. 日志导出

1. 支持导出当前查询结果
2. 导出格式：CSV、Excel
3. 大文件异步生成，完成后通知下载
4. 导出记录审计（记录谁在什么时候导出了日志）

### 6. 日志保留策略

1. 默认保留 180 天
2. 租户可配置保留期限：30 天 / 90 天 / 180 天 / 365 天 / 永久
3. 过期日志自动清理（每日凌晨执行）
4. 清理前可归档到冷存储（可选）

### 7. 日志统计分析

支持以下统计维度：

| 维度 | 说明 |
|------|------|
| 操作趋势 | 按日/周/月统计操作数量趋势 |
| 操作人排行 | 统计活跃操作人 TOP10 |
| 操作类型分布 | 各类操作占比 |
| 失败操作分析 | 失败操作类型、原因分析 |
| 时间段分布 | 24 小时内操作分布热力图 |

**日志配置项：**

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `audit_log_retention_days` | 180 | 日志保留天数 |
| `audit_log_async_enabled` | true | 是否异步写入 |
| `audit_log_sensitive_mask` | true | 敏感数据脱敏 |
| `audit_log_max_query_days` | 31 | 最大查询时间跨度 |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| 日志写入失败 | 记录到本地文件，后续补写 |
| 日志存储超限 | 自动清理最早 10% 日志 |
| 查询时间跨度过大 | 限制最大查询范围为 31 天 |
| 导出量大 | 异步生成，完成后邮件通知 |

**安全策略：**

| 策略 | 说明 |
|------|------|
| **日志防篡改** | 日志写入后不可修改 |
| **敏感数据脱敏** | 密码、手机号、邮箱脱敏显示 |
| **访问控制** | 只有管理员可查询审计日志 |
| **导出审计** | 导出操作本身被记录 |

**API 接口：**

```
GET    /api/v1/audit-logs                 # 审计日志列表
GET    /api/v1/audit-logs/:id             # 审计日志详情
GET    /api/v1/audit-logs/export          # 导出审计日志
GET    /api/v1/audit-logs/statistics      # 统计数据
GET    /api/v1/audit-logs/operation-types # 获取操作类型列表
```

**数据库设计：**

**操作审计日志表（operation_audit_logs）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 1001 |
| tenant_id | BIGINT | 是 | 租户 ID | 100 |
| user_id | BIGINT | 否 | 操作人 ID | 2001 |
| user_name | VARCHAR(100) | 否 | 操作人姓名 | 张三 |
| operation_type | VARCHAR(50) | 是 | 操作类型 | USER_CREATE/ROLE_ASSIGN |
| resource_type | VARCHAR(50) | 否 | 资源类型 | USER/ROLE/TENANT |
| resource_id | BIGINT | 否 | 资源 ID | 2002 |
| resource_name | VARCHAR(255) | 否 | 资源名称 | 李四 |
| action | VARCHAR(50) | 否 | 具体动作 | CREATE/UPDATE/DELETE |
| request_body | JSON | 否 | 请求参数（脱敏） | {"email": "li***@example.com"} |
| response_status | VARCHAR(20) | 是 | 响应状态 | SUCCESS/FAILURE |
| error_message | TEXT | 否 | 错误信息 | 失败时的错误消息 |
| ip_address | VARCHAR(45) | 否 | 操作 IP | 192.168.1.100 |
| user_agent | VARCHAR(500) | 否 | 用户代理 | Mozilla/5.0... |
| device_info | VARCHAR(100) | 否 | 设备信息 | Chrome 120 / macOS |
| location | VARCHAR(100) | 否 | 地理位置 | 北京市海淀区 |
| created_at | DATETIME | - | 操作时间 | 2026-03-25 10:30:00 |

**索引**：
- `idx_tenant_time` (tenant_id, created_at) — 租户 + 时间查询
- `idx_user` (tenant_id, user_id, created_at) — 按用户查询
- `idx_operation` (tenant_id, operation_type, created_at) — 按操作类型查询
- `idx_resource` (tenant_id, resource_type, resource_id) — 按资源查询

**验收标准：**

- [ ] 敏感操作完整记录
- [ ] 日志内容准确
- [ ] 查询功能正常（支持所有条件）
- [ ] 日志导出正常
- [ ] 过期日志自动清理
- [ ] 统计数据准确
- [ ] 敏感数据脱敏显示

