# REQ-008 MFA 多因素认证

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 5 人天 |
| **关联用户故事** | US-012 |

**背景：** 高安全场景下需要多因素认证，防止账号被盗用。

**目标：**

- 支持 TOTP 动态验证码（Google Authenticator、Microsoft Authenticator）
- 支持邮箱验证码作为 MFA 方式
- 支持短信验证码作为 MFA 方式
- 支持按角色配置 MFA 强制启用
- 支持可信设备豁免 MFA
- 目标 MFA 覆盖率 > 60%（管理员 100%）

**功能描述：**

### 1. MFA 方式

支持以下 MFA 认证方式：

| 方式 | 说明 | 优先级 |
|------|------|--------|
| TOTP | 基于时间的一次性验证码（Google/Microsoft Authenticator） | 推荐 |
| 邮箱验证码 | 发送 6 位数字验证码到绑定邮箱 | 备选 |
| 短信验证码 | 发送 6 位数字验证码到绑定手机 | 备选 |

用户可同时绑定多种 MFA 方式，登录时选择任一方式验证。

### 2. MFA 绑定流程（TOTP）

1. 用户进入 MFA 设置页面
2. 系统生成 TOTP Secret（Base32 编码）
3. 展示二维码（包含 issuer、账号、Secret）
4. 用户使用 Authenticator App 扫描二维码
5. 用户输入 6 位验证码完成绑定验证
6. 展示并下载备用码（10 个，一次性使用）
7. 加密存储 TOTP Secret 于数据库

### 3. MFA 触发条件

支持以下触发策略配置：

| 策略 | 说明 | 默认 |
|------|------|------|
| `ALWAYS` | 每次登录都需要 MFA | 默认 |
| `NEW_DEVICE_ONLY` | 仅新设备/新地点登录时需要 | 可选 |
| `ROLE_BASED` | 根据角色配置决定是否强制 | 可选 |
| `RISK_BASED` | 根据风险评分动态决定 | 预留 |

**可信设备豁免：**
- 用户可选择「信任此设备 30 天」
- 可信设备记录设备指纹
- 可信设备有效期内登录无需 MFA
- 用户可随时在设置中清除可信设备

### 4. 强制 MFA 策略

租户管理员可配置：

1. 指定角色必须启用 MFA（如管理员角色）
2. 指定操作必须 MFA 验证（如删除、导出）
3. 新创建用户强制启用 MFA
4. MFA 启用宽限期（如 7 天内必须绑定）

### 5. MFA 验证流程

1. 用户通过密码认证后，检查 MFA 状态
2. 如需要 MFA，返回 `MFA_REQUIRED` 状态和可用的 MFA 方式
3. 用户选择 MFA 方式并提交验证码
4. 校验验证码正确性
5. 验证通过，颁发 Token
6. 如设备被信任，记录可信设备指纹

### 6. 备用码

1. 绑定 TOTP 时生成 10 个备用码
2. 每个备用码只能使用一次
3. 备用码格式：`XXXXX-XXXXX`（10 位，含分隔符）
4. 用户可重新生成备用码（原备用码全部失效）
5. 备用码用于丢失 Authenticator 时应急

**MFA 配置项：**

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `mfa_totp_enabled` | true | 是否启用 TOTP |
| `mfa_email_enabled` | true | 是否启用邮箱验证码 |
| `mfa_sms_enabled` | false | 是否启用短信验证码 |
| `mfa_trust_device_days` | 30 | 可信设备有效期（天） |
| `mfa_backup_code_count` | 10 | 备用码数量 |
| `mfa_required_roles` | [] | 强制 MFA 的角色列表 |
| `mfa_grace_period_days` | 0 | MFA 启用宽限期（天） |

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| TOTP 验证码错误 | 提示「验证码错误，剩余重试次数 N」 |
| TOTP 验证码超时 | 提示「验证码已过期，请刷新后重试」 |
| 备用码已使用 | 提示「备用码已失效，请使用未使用的备用码」 |
| 备用码全部用完 | 提示「备用码已用完，请联系管理员」 |
| 连续 5 次 MFA 失败 | 账号锁定 15 分钟 |
| 无可用 MFA 方式 | 提示「请先绑定 MFA 设备」 |

**安全策略：**

| 策略 | 说明 |
|------|------|
| **Secret 加密** | TOTP Secret 加密存储 |
| **备用码哈希** | 备用码哈希存储，不存明文 |
| **速率限制** | MFA 验证频率限制（5 次/分钟） |
| **失败锁定** | 连续失败后账号锁定 |
| **审计日志** | MFA 绑定、验证、禁用操作记录日志 |

**API 接口：**

```
# MFA 绑定与管理
POST   /api/v1/mfa/bind              # 发起绑定（返回 Secret 和二维码）
POST   /api/v1/mfa/bind/verify       # 验证并完成绑定
POST   /api/v1/mfa/unbind            # 解绑 MFA（需要验证）
GET    /api/v1/mfa/status            # 获取 MFA 状态
POST   /api/v1/mfa/backup-codes/generate  # 重新生成备用码
GET    /api/v1/mfa/backup-codes      # 查看备用码（需要二次验证）

# MFA 验证
POST   /api/v1/mfa/verify            # MFA 验证码验证
POST   /api/v1/mfa/trust-device      # 标记为可信设备
DELETE /api/v1/mfa/trust-device/:id  # 移除可信设备
GET    /api/v1/mfa/trusted-devices   # 可信设备列表
```

**数据库设计：**

**MFA 绑定（扩展 users 表）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| mfa_enabled | BOOLEAN | - | 是否启用 MFA | true/false |
| mfa_type | VARCHAR(20) | 否 | MFA 类型 | totp/email/sms |
| mfa_secret | VARCHAR(255) | 否 | TOTP Secret（加密存储） | 加密后的 Base32 字符串 |
| mfa_backup_codes_hash | VARCHAR(500) | 否 | 备用码哈希 | 哈希后的备用码列表 |

---

**可信设备表（mfa_trusted_devices）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 1001 |
| user_id | BIGINT | 是 | 用户 ID | 2001 |
| device_fingerprint | VARCHAR(64) | 是 | 设备指纹 | fp_xxxxxxxxxxxxx |
| device_type | VARCHAR(20) | 否 | 设备类型 | web/ios/android |
| ip_address | VARCHAR(45) | 否 | IP 地址 | 192.168.1.100 |
| expires_at | DATETIME | 是 | 过期时间 | 2026-04-28 10:00:00 |
| created_at | DATETIME | - | 创建时间 | 2026-03-28 10:00:00 |

**索引**：`uk_user_device` (user_id, device_fingerprint) — 唯一索引、`idx_expires` (expires_at)

---

**MFA 操作日志表（mfa_operation_logs）**

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| id | BIGINT | 是 | 主键 | 2001 |
| user_id | BIGINT | 是 | 用户 ID | 2001 |
| action | VARCHAR(50) | 是 | 操作类型 | bind/unbind/verify/backup_codes |
| result | VARCHAR(20) | 是 | 操作结果 | success/failure |
| ip_address | VARCHAR(45) | 否 | IP 地址 | 192.168.1.100 |
| user_agent | VARCHAR(255) | 否 | 用户代理 | Mozilla/5.0... |
| created_at | DATETIME | - | 创建时间 | 2026-03-28 10:00:00 |

**索引**：`idx_user_time` (user_id, created_at)

**验收标准：**

- [ ] 可成功绑定 TOTP（Google/Microsoft Authenticator）
- [ ] 登录时 MFA 校验生效
- [ ] 强制 MFA 策略生效（按角色）
- [ ] 备用码可正常使用
- [ ] MFA 失败锁定机制生效
- [ ] 可信设备豁免功能正常
- [ ] 邮箱/短信验证码作为 MFA 方式可用
- [ ] MFA 操作日志完整记录

