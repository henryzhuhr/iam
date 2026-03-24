# REQ-015 验证码登录

| 项目 | 内容 |
|------|------|
| **优先级** | P1 |
| **估时** | 3 人天 |
| **关联用户故事** | US-018、US-019 |

**背景：** 部分用户场景下（如忘记密码、无密码登录），需要提供基于验证码的登录方式，支持邮箱和手机号两种渠道。

**目标：**

- 支持邮箱验证码登录
- 支持手机号验证码登录
- 验证码有效期可配置
- 支持发送频率限制
- 目标验证码登录成功率 > 90%

**功能描述：**

### 1. 邮箱验证码登录

1. 用户输入邮箱地址
2. 系统生成 6 位数字验证码（或 6 位字母数字混合）
3. 发送验证码到用户邮箱
4. 用户输入验证码
5. 校验验证码正确性和有效性
6. 登录成功，颁发 Token

### 2. 手机号验证码登录

1. 用户输入手机号码
2. 系统生成 6 位数字验证码
3. 通过短信服务商发送验证码
4. 用户输入验证码
5. 校验验证码正确性和有效性
6. 登录成功，颁发 Token

### 3. 验证码生成规则

| 规则 | 说明 | 默认值 |
|------|------|--------|
| 验证码长度 | 验证码位数 | 6 位 |
| 验证码类型 | 数字/字母混合 | 纯数字 |
| 有效期 | 验证码有效时长 | 10 分钟 |
| 最大验证次数 | 单验证码可验证次数 | 5 次 |
| 发送间隔 | 同一目标最小发送间隔 | 60 秒 |
| 发送上限 | 单 IP/单账号每日发送上限 | 10 次/天 |

### 4. 发送频率限制

防刷策略：
- 同一邮箱/手机号：60 秒内只能发送 1 次
- 同一 IP：1 小时内最多发送 5 次
- 同一账号：24 小时内最多发送 10 次
- 超过限制后，提示「发送过于频繁，请稍后再试」

### 5. 验证码安全

1. 验证码加密存储（不存明文）
2. 验证后立即失效（一次性使用）
3. 连续验证错误 3 次，验证码失效
4. 支持图形验证码前置（触发风控时）

### 6. 注册/登录自动判断

1. 验证码登录时，自动判断用户是否存在
2. 用户不存在：自动创建账号并登录
3. 用户已存在：直接登录
4. 需配置是否允许自动注册

**验证码模板：**

邮箱验证码模板：
```
主题：【IAM】您的登录验证码

您好！

您的登录验证码为：123456

该验证码 10 分钟内有效，请勿泄露给他人。

如非本人操作，请忽略此邮件。

IAM 团队
```

短信验证码模板：
```
【IAM】您的登录验证码为 123456，10 分钟内有效。
如非本人操作，请忽略。
```

**异常情况：**

| 异常场景 | 系统处理 |
|----------|----------|
| 邮箱/手机号格式错误 | 提示格式错误 |
| 验证码发送失败 | 提示「发送失败，请重试」 |
| 验证码已过期 | 提示「验证码已过期，请重新获取」 |
| 验证码错误 | 提示「验证码错误，剩余 N 次机会」 |
| 发送频率超限 | 提示「发送过于频繁，X 分钟后再试」 |
| 图形验证码触发 | 要求先完成图形验证 |

**配置项：**

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `code_length` | 6 | 验证码长度 |
| `code_ttl_minutes` | 10 | 验证码有效期 |
| `code_max_verify` | 5 | 最大验证次数 |
| `send_interval_seconds` | 60 | 发送间隔 |
| `send_limit_per_hour` | 5 | 每小时发送上限 |
| `send_limit_per_day` | 10 | 每日发送上限 |
| `auto_register_enabled` | true | 是否允许自动注册 |

**API 接口：**

```
# 邮箱验证码
POST   /api/v1/auth/code/email/send    # 发送邮箱验证码
POST   /api/v1/auth/code/email/login   # 邮箱验证码登录

# 手机验证码
POST   /api/v1/auth/code/phone/send    # 发送手机验证码
POST   /api/v1/auth/code/phone/login   # 手机验证码登录

# 验证码校验（可选，用于分步验证）
POST   /api/v1/auth/code/verify        # 校验验证码
```

**数据库设计：**

```sql
-- 验证码发送记录表
CREATE TABLE verification_codes (
    id BIGINT PRIMARY KEY,
    target VARCHAR(100) NOT NULL,      -- 邮箱/手机号
    code_hash VARCHAR(64) NOT NULL,    -- 验证码哈希
    type VARCHAR(20) NOT NULL,         -- email_login/phone_login/password_reset
    expires_at DATETIME NOT NULL,
    verify_count INT DEFAULT 0,        -- 已验证次数
    max_verify_count INT DEFAULT 5,    -- 最大验证次数
    is_used BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_target_type (target, type),
    INDEX idx_expires (expires_at)
);

-- 验证码发送记录（用于频率限制）
CREATE TABLE code_send_logs (
    id BIGINT PRIMARY KEY,
    target VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    ip_address VARCHAR(45),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_target_time (target, type, created_at),
    INDEX idx_ip_time (ip_address, created_at)
);

-- 登录日志（复用已有设计）
-- 见 REQ-010 登录日志
```

**安全策略：**

| 策略 | 说明 |
|------|------|
| **防枚举** | 不提示目标是否存在 |
| **防重放** | 验证码一次性使用 |
| **防暴力** | 连续错误后验证码失效 |
| **防刷** | 频率限制 + IP 限制 |
| **加密存储** | 验证码哈希存储 |
| **风控触发** | 异常情况要求图形验证码 |

**验收标准：**

- [ ] 邮箱验证码可正常发送和验证
- [ ] 手机验证码可正常发送和验证
- [ ] 验证码登录成功颁发 Token
- [ ] 频率限制正确生效
- [ ] 验证码过期后正确拒绝
- [ ] 自动注册功能正常工作
- [ ] 发送记录完整可查
