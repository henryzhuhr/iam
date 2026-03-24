# MFA 与 TOTP 动态验证码

> 最后更新：2026-03-25
> 适用场景：IAM 多因素认证

## 1. MFA 是什么

**MFA (Multi-Factor Authentication)** 多因素认证，要求用户提供两种或以上的认证因素才能完成身份验证。

### 1.1 认证因素类型

| 类型 | 说明 | 示例 |
|------|------|------|
| **知识因素** (Something you know) | 用户知道的信息 | 密码、PIN 码 |
| **持有因素** (Something you have) | 用户拥有的物品 | 手机、硬件 Token、智能卡 |
| **生物因素** (Something you are) | 用户的生物特征 | 指纹、人脸、虹膜 |

### 1.2 为什么需要 MFA

| 单因素认证风险 | MFA 如何缓解 |
|----------------|--------------|
| 密码泄露 | 攻击者没有第二因素 |
| 撞库攻击 | 即使密码正确也无法登录 |
| 钓鱼攻击 | 动态验证码无法被钓鱼 |
| 暴力破解 | 增加攻击难度 |

---

## 2. TOTP 动态验证码

**TOTP (Time-based One-Time Password)** 基于时间的一次性密码，是 MFA 最常用的实现方式。

### 2.1 TOTP 原理

```
TOTP = HOTP(T)
T = (CurrentTime - T0) / X
```

- **HOTP**: 基于计数器的 HMAC 算法
- **T0**: 起始时间（通常为 0）
- **X**: 时间步长（通常 30 秒）

### 2.2 工作流程

```
┌─────────────┐         ┌──────────────┐
│  用户手机    │         │   IAM 服务器  │
│ (Authenticator)       │              │
└──────┬──────┘         └──────┬───────┘
       │                       │
       │  1. 展示二维码        │
       │  (包含 Secret)        │
       │<──────────────────────┤
       │                       │
       │  2. 扫描并存储 Secret │
       │                       │
       │  3. 每 30 秒生成 6 位码   │
       │     TOTP = HMAC(Secret, Time) │
       │                       │
       │  4. 输入验证码        │
       │──────────────────────>│
       │                       │
       │  5. 服务器用相同算法验证 │
       │     (允许±1 个时间步长)  │
       │                       │
       │  6. 验证结果          │
       │<──────────────────────┤
       │                       │
```

### 2.3 二维码格式

TOTP 二维码使用 `otpauth://` URL 格式：

```
otpauth://totp/IssuerName:AccountName?secret=SecretString&issuer=IssuerName&algorithm=SHA1&digits=6&period=30
```

示例：
```
otpauth://totp/IAM:zhangsan@example.com?secret=JBSWY3DPEHPK3PXP&issuer=IAM&algorithm=SHA1&digits=6&period=30
```

参数说明：

| 参数 | 说明 | 示例值 |
|------|------|--------|
| `totp` | 协议类型 | totp / hotp |
| `IssuerName` | 服务名称 | IAM |
| `AccountName` | 账号（通常是邮箱） | zhangsan@example.com |
| `secret` | 共享密钥（Base32 编码） | JBSWY3DPEHPK3PXP |
| `issuer` | 发行者 | IAM |
| `algorithm` | 哈希算法 | SHA1 / SHA256 / SHA512 |
| `digits` | 验证码位数 | 6 |
| `period` | 有效期（秒） | 30 |

---

## 3. 主流 Authenticator App

| 应用 | 平台 | 特点 |
|------|------|------|
| **Google Authenticator** | iOS/Android | 最常用，简单易用 |
| **Microsoft Authenticator** | iOS/Android | 支持云备份 |
| **Authy** | iOS/Android/Desktop | 支持多设备同步 |
| **1Password** | 全平台 | 密码管理器内置 |

---

## 4. TOTP 实现代码（Go）

### 4.1 生成 Secret 和二维码

```go
import (
    "bytes"
    "image/png"

    "github.com/pquerna/otp/totp"
    "github.com/pquerna/otp"
)

type MFASetup struct {
    Secret     string
    QRCodeURL  string
}

func GenerateMFASecret(accountName, issuer string) (*MFASetup, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      issuer,
        AccountName: accountName,
    })
    if err != nil {
        return nil, err
    }

    return &MFASetup{
        Secret:     key.Secret(),
        QRCodeURL:  key.URL(),
    }, nil
}

// 生成二维码图片
func GenerateQRCode(qrCodeURL string) ([]byte, error) {
    img, err := otp.NewQRCode(qrCodeURL)
    if err != nil {
        return nil, err
    }

    pngBytes, err := img.PNG(256)
    if err != nil {
        return nil, err
    }

    return pngBytes, nil
}
```

### 4.2 验证 TOTP 验证码

```go
func VerifyTOTP(secret, passcode string) bool {
    valid := totp.Validate(passcode, secret)
    return valid
}

// 允许时间偏差（前后各一个周期）
func VerifyTOTPWithDrift(secret, passcode string) bool {
    valid := totp.ValidateCustom(passcode, secret, time.Now().UTC(),
        totp.ValidateOpts{
            Period:    30,
            Skew:      1,  // 允许±1 个周期
            Digits:    otp.DigitsSix,
            Algorithm: otp.AlgorithmSHA1,
        })
    return valid
}
```

### 4.3 生成备用码

```go
func GenerateBackupCodes(count int) ([]string, error) {
    codes := make([]string, count)
    for i := 0; i < count; i++ {
        // 生成 10 位随机码，格式 XXXXX-XXXXX
        code, err := generateRandomCode()
        if err != nil {
            return nil, err
        }
        codes[i] = code
    }
    return codes, nil
}

func generateRandomCode() (string, error) {
    // 使用 crypto/rand 生成安全随机数
    // 格式：ABCDE-FGHIJ
    // ...
}
```

---

## 5. 备用码（Backup Codes）

### 5.1 什么是备用码

备用码是 TOTP 失效时的应急方案，通常是一组一次性使用的随机码。

### 5.2 备用码特点

| 特性 | 说明 |
|------|------|
| 数量 | 通常 10 个 |
| 格式 | `XXXXX-XXXXX`（10 位，便于输入） |
| 使用规则 | 每个码只能用一次 |
| 有效期 | 永不过期，直到用完或重置 |
| 生成时机 | 绑定 TOTP 时展示（仅此一次） |

### 5.3 备用码使用场景

- 手机丢失/损坏
- Authenticator App 数据丢失
- 更换设备

---

## 6. MFA 安全最佳实践

### 6.1 Secret 保护

| 措施 | 说明 |
|------|------|
| 加密存储 | Secret 加密后存入数据库 |
| 传输加密 | 二维码通过 HTTPS 传输 |
| 不重复使用 | 每个用户的 Secret 唯一 |

### 6.2 验证策略

| 策略 | 说明 |
|------|------|
| 允许时间偏差 | ±1 个周期（90 秒窗口） |
| 速率限制 | 5 次/分钟，防暴力破解 |
| 失败锁定 | 连续 5 次失败锁定 15 分钟 |

### 6.3 可信设备

| 功能 | 说明 |
|------|------|
| 信任此设备 | 30 天内无需 MFA |
| 设备指纹 | 记录设备特征 |
| 随时撤销 | 用户可清除可信设备列表 |

---

## 7. MFA 策略配置

### 7.1 触发条件

| 策略 | 说明 | 适用场景 |
|------|------|----------|
| `ALWAYS` | 每次登录都需要 MFA | 高安全场景 |
| `NEW_DEVICE_ONLY` | 仅新设备需要 MFA | 平衡安全与体验 |
| `ROLE_BASED` | 特定角色强制 MFA | 管理员必须，普通用户可选 |
| `RISK_BASED` | 根据风险评分动态决定 | 智能风控 |

### 7.2 推荐配置

| 用户类型 | MFA 策略 |
|----------|----------|
| 平台管理员 | 强制，Always |
| 租户管理员 | 强制，Always |
| 普通用户 | 可选，New Device Only |
| 敏感操作 | 强制验证（删除、导出、改密码） |

---

## 8. 常见问题

### Q1: TOTP 验证码一直错误怎么办？

可能原因：
1. 设备时间不准确（检查手机时间设置）
2. 网络延迟导致服务器时间偏差
3. Secret 传输错误

解决方案：
- 校准设备时间
- 使用备用码
- 重新绑定 TOTP

### Q2: 用户手机丢了怎么办？

1. 使用备用码登录
2. 联系管理员重置 MFA
3. 通过邮箱验证临时禁用 MFA

### Q3: TOTP 和短信验证码哪个更好？

| 维度 | TOTP | 短信 |
|------|------|------|
| 安全性 | 高（本地生成） | 中（SIM 卡劫持风险） |
| 成本 | 免费 | 按条收费 |
| 离线可用 | 是 | 否 |
| 用户体验 | 中（需要 App） | 高（直接收短信） |

推荐：TOTP 为主，短信为辅。

---

## 9. 参考链接

- RFC 6238 (TOTP): https://tools.ietf.org/html/rfc6238
- RFC 4226 (HOTP): https://tools.ietf.org/html/rfc4226
- Google Authenticator 文档：https://github.com/google/google-authenticator

---

## 10. 相关需求文档

- [REQ-008 MFA 多因素认证](../05-functional-requirements/REQ-008-mfa.md)
