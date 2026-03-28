# REQ-010 登录日志记录

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 2 人天 |
| **关联用户故事** | US-014、US-021 |

**背景：** 需要记录用户登录行为，用于安全分析、异常检测和责任追溯。

**目标：**

- 记录所有登录尝试（成功/失败）
- 记录登录 IP、设备、地理位置
- 支持异常登录检测
- 支持租户管理员查看本企业登录日志
- 登录日志保留 180 天

**功能描述：**

### 1. 日志记录范围

记录以下登录场景：

| 场景 | 说明 |
|------|------|
| 密码登录成功 | 用户输入正确密码登录成功 |
| 密码登录失败 | 密码错误、账号不存在、账号被禁用 |
| 验证码登录成功 | 邮箱/手机验证码登录成功 |
| 验证码登录失败 | 验证码错误、验证码过期 |
| MFA 登录完成 | MFA 验证通过，登录完成 |
| MFA 验证失败 | MFA 验证失败 |
| 第三方登录成功 | GitHub 等第三方登录成功 |
| 第三方登录失败 | 第三方授权失败 |
| Token 刷新 | 使用 Refresh Token 刷新 Access Token |
| 登出 | 用户主动登出 |

### 2. 日志内容

每条登录日志包含以下字段：

| 字段 | 说明 | 示例 |
|------|------|------|
| login_id | 登录记录唯一标识 | `log_xxxxxxxxx` |
| tenant_id | 租户 ID | `1001` |
| user_id | 用户 ID（失败时为空） | `2001` |
| username | 尝试登录的用户名 | `zhangsan@example.com` |
| result | 登录结果 | `SUCCESS` / `FAILURE` |
| failure_reason | 失败原因 | `INVALID_PASSWORD` / `ACCOUNT_LOCKED` |
| ip_address | 登录 IP | `192.168.1.100` |
| location | 地理位置 | `北京市海淀区` |
| country | 国家 | `中国` |
| province | 省份 | `北京市` |
| city | 城市 | `北京市` |
| isp | 运营商 | `中国电信` |
| device_type | 设备类型 | `DESKTOP` / `MOBILE` / `TABLET` |
| os | 操作系统 | `macOS 14.0` |
| browser | 浏览器 | `Chrome 120.0` |
| user_agent | 原始 User-Agent | `Mozilla/5.0...` |
| is_new_device | 是否新设备 | `true` / `false` |
| is_new_location | 是否新地点 | `true` / `false` |
| risk_score | 风险评分 | `0-100` |
| session_id | 会话 ID | `sess_xxxxxxxxx` |
| created_at | 登录时间 | `2026-03-25 10:30:00` |

### 3. 地理位置解析

1. 使用 IP 地址解析地理位置
2. 支持离线 IP 库（如 GeoIP2）
3. 解析失败时，地理位置显示"未知"
4. 缓存解析结果，避免重复查询

### 4. 异常登录检测

检测以下异常场景：

| 异常类型 | 检测规则 | 处置 |
|----------|----------|------|
| 新设备登录 | 设备指纹不在可信列表 | 标记、可选 MFA |
| 新地点登录 | 地理位置与历史差异大 | 标记、可选 MFA |
| 频繁登录失败 | 5 分钟内失败≥5 次 | 告警、临时锁定 |
| 异地并发登录 | 同一账号短时间内两地登录 | 告警、可选阻断 |
| 非常用时间登录 | 凌晨等非常用时段 | 标记 |
| 暴力破解嫌疑 | 同一 IP 尝试多个账号 | 封禁 IP |
| 账号共享嫌疑 | 同一账号多设备频繁切换 | 告警 |

风险评分模型：
- 基础分 0 分
- 新设备 +20 分
- 新地点 +30 分
- 频繁失败 +25 分
- 异地并发 +40 分
- 非常用时间 +10 分

风险等级：
- 0-30 分：低风险
- 31-60 分：中风险
- 61-100 分：高风险

### 5. 登录日志查询

支持以下查询条件：

| 条件 | 说明 |
|------|------|
| 时间范围 | 支持自定义起止时间（最大 31 天） |
| 用户 | 按用户 ID 或用户名搜索 |
| 登录结果 | 成功/失败 |
| IP 地址 | 按 IP 或 IP 段搜索 |
| 设备类型 | 按设备类型筛选 |
| 地理位置 | 按省份、城市筛选 |
| 风险等级 | 按风险等级筛选 |

查询结果支持：
- 分页加载
- 按时间倒序/正序排序
- 导出 CSV/Excel

### 6. 登录告警通知

支持以下告警配置：

| 告警类型 | 通知方式 | 接收人 |
|----------|----------|--------|
| 高风险登录 | 邮件 + 短信 | 用户本人、管理员 |
| 频繁失败 | 邮件 | 用户本人、管理员 |
| 异地并发 | 邮件 + 短信 | 用户本人、管理员 |
| 账号锁定 | 邮件 | 用户本人、管理员 |

### 7. 用户登录历史页面

用户可查看自己的登录历史：
- 最近登录时间、IP、设备
- 活跃会话列表
- 可远程登出其他设备

### 8. 日志保留策略

1. 默认保留 180 天
2. 租户可配置保留期限：30 天 / 90 天 / 180 天 / 365 天
3. 过期日志自动清理（每日凌晨执行）

**日志配置项：**

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `login_log_retention_days` | 180 | 日志保留天数 |
| `login_log_async_enabled` | true | 是否异步写入 |
| `login_risk_detection_enabled` | true | 是否启用风险检测 |
| `login_alert_high_risk` | true | 高风险登录告警 |
| `new_device_mfa_required` | false | 新设备是否要求 MFA |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| IP 解析失败 | 地理位置显示"未知" |
| 日志写入失败 | 记录到本地文件，后续补写 |
| User-Agent 解析失败 | 设备信息保存原始字符串 |
| 告警发送失败 | 记录告警日志，重试 3 次 |

**安全策略：**

| 策略 | 说明 |
|------|------|
| **租户隔离** | 租户管理员只能查看本企业日志 |
| **隐私保护** | 用户只能查看自己的登录历史 |
| **数据脱敏** | 导出日志敏感信息脱敏 |
| **防篡改** | 日志写入后不可修改 |

**API 接口：**

```
# 登录日志查询（管理员）
GET    /api/v1/login-logs                 # 登录日志列表
GET    /api/v1/login-logs/:id             # 登录日志详情
GET    /api/v1/login-logs/export          # 导出登录日志
GET    /api/v1/login-logs/statistics      # 统计数据

# 用户个人登录历史
GET    /api/v1/my/login-history          # 我的登录历史
GET    /api/v1/my/sessions               # 我的活跃会话
DELETE /api/v1/my/sessions/:id           # 登出指定会话
DELETE /api/v1/my/sessions/all           # 登出所有会话
```

**数据库设计：**

**登录日志表（login_logs）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 1001 |
| tenant_id | BIGINT | 是 | 租户 ID | 100 |
| user_id | BIGINT | 否 | 用户 ID（失败时为空） | 2001 |
| username | VARCHAR(100) | 否 | 尝试登录的用户名 | zhangsan@example.com |
| result | VARCHAR(20) | 是 | 登录结果 | SUCCESS/FAILURE |
| failure_reason | VARCHAR(50) | 否 | 失败原因 | INVALID_PASSWORD/ACCOUNT_LOCKED |
| ip_address | VARCHAR(45) | 否 | 登录 IP | 192.168.1.100 |
| location | VARCHAR(100) | 否 | 地理位置 | 北京市海淀区 |
| country | VARCHAR(50) | 否 | 国家 | 中国 |
| province | VARCHAR(50) | 否 | 省份 | 北京市 |
| city | VARCHAR(50) | 否 | 城市 | 北京市 |
| isp | VARCHAR(50) | 否 | 运营商 | 中国电信 |
| device_type | VARCHAR(20) | 否 | 设备类型 | DESKTOP/MOBILE/TABLET |
| os | VARCHAR(50) | 否 | 操作系统 | macOS 14.0 |
| browser | VARCHAR(50) | 否 | 浏览器 | Chrome 120.0 |
| user_agent | VARCHAR(500) | 否 | 原始 User-Agent | Mozilla/5.0... |
| is_new_device | BOOLEAN | - | 是否新设备 | true/false |
| is_new_location | BOOLEAN | - | 是否新地点 | true/false |
| risk_score | INT | - | 风险评分 | 0-100 |
| session_id | VARCHAR(64) | 否 | 会话 ID | sess_xxxxxxxxx |
| created_at | DATETIME | - | 登录时间 | 2026-03-25 10:30:00 |

**索引**：
- `idx_tenant_time` (tenant_id, created_at)
- `idx_user_time` (tenant_id, user_id, created_at)
- `idx_result` (tenant_id, result, created_at)
- `idx_ip` (tenant_id, ip_address, created_at)

---

**登录告警记录表（login_alerts）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 2001 |
| tenant_id | BIGINT | 是 | 租户 ID | 100 |
| user_id | BIGINT | 是 | 用户 ID | 2001 |
| alert_type | VARCHAR(50) | 是 | 告警类型 | HIGH_RISK_LOGIN/FREQUENT_FAILURE |
| risk_level | VARCHAR(20) | 是 | 风险等级 | LOW/MEDIUM/HIGH |
| login_log_id | BIGINT | 否 | 关联登录日志 ID | 1001 |
| is_sent | BOOLEAN | - | 是否已发送 | true/false |
| sent_at | DATETIME | 否 | 发送时间 | 2026-03-25 10:35:00 |
| created_at | DATETIME | - | 创建时间 | 2026-03-25 10:30:00 |

**索引**：`idx_tenant` (tenant_id, created_at)、`idx_user` (tenant_id, user_id, is_sent)

**验收标准：**

- [ ] 登录日志完整记录
- [ ] 地理位置解析正确
- [ ] 租户管理员只能查看本企业日志
- [ ] 异常登录检测告警正常
- [ ] 风险评分计算准确
- [ ] 用户可查看和管理自己的会话
- [ ] 日志导出功能正常
- [ ] 过期日志自动清理

