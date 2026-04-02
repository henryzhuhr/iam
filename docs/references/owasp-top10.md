# OWASP Top 10 安全风险与防护

> 最后更新：2026-03-29
> 适用场景：Web 应用安全设计、代码审计、渗透测试

---

## 1. 概述

**OWASP（Open Web Application Security Project）** 是一个开放的中立组织，致力于提高软件安全性。

**OWASP Top 10** 是最关键的 Web 应用安全风险列表，每 3-4 年更新一次。

| 版本 | 年份 |
|------|------|
| OWASP Top 10 2017 | 2017 |
| OWASP Top 10 2021 | 2021（当前） |

---

## 2. OWASP Top 10 2021

### 2.1 风险列表总览

| 编号 | 风险 | 说明 |
|------|------|------|
| **A01:2021** | 失效的访问控制 | 用户越权访问资源 |
| **A02:2021** | 加密机制失效 | 敏感数据未加密或加密不当 |
| **A03:2021** | 注入攻击 | SQL 注入、命令注入等 |
| **A04:2021** | 不安全设计 | 设计缺陷导致的安全漏洞 |
| **A05:2021** | 安全配置错误 | 默认配置、错误配置 |
| **A06:2021** | 易受攻击的组件 | 使用有漏洞的第三方库 |
| **A07:2021** | 认证和会话失效 | 身份认证和会话管理缺陷 |
| **A08:2021** | 软件和数据完整性故障 | 依赖链攻击、数据篡改 |
| **A09:2021** | 安全日志和监控缺失 | 无法检测和响应攻击 |
| **A10:2021** | 服务端请求伪造（SSRF） | 伪造服务端请求 |

### 2.2 与 2017 版本对比

| 变化 | 说明 |
|------|------|
| **新增** | A04 不安全设计、A08 软件和数据完整性故障、A10 SSRF |
| **合并** | A03 注入 + A05 安全配置错误 + A07 XSS → 重新分类 |
| **改名** | A05 失效的访问控制（原 A05 是越权） |

---

## 3. A01:2021 失效的访问控制

### 3.1 风险说明

攻击者通过修改请求参数，访问未授权的资源或功能。

```
常见场景：
- 修改 URL 中的 ID 访问他人数据
- 直接调用未授权的 API
- 越权删除/修改数据
```

### 3.2 攻击示例

```
# 用户 A 访问自己的订单
GET /api/orders/123

# 攻击者修改 ID，访问用户 B 的订单
GET /api/orders/124

# 如果服务端未校验所有权，攻击成功
```

### 3.3 防护措施

| 措施 | 说明 |
|------|------|
| **最小权限原则** | 默认拒绝，显式授权 |
| **所有权校验** | 每次请求校验资源归属 |
| **集中式授权** | 统一授权逻辑，避免分散 |
| **审计日志** | 记录所有访问控制失败 |

### 3.4 Go 代码示例

```go
// ❌ 错误：未校验资源所有权
func GetOrder(c *gin.Context) {
    orderID := c.Param("id")
    var order Order
    db.First(&order, orderID)
    c.JSON(200, order) // 任何人都可以访问
}

// ✅ 正确：校验资源所有权
func GetOrder(c *gin.Context) {
    orderID := c.Param("id")
    userID := GetUserIDFromToken(c) // 从 Token 获取当前用户

    var order Order
    result := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order)

    if result.Error != nil {
        c.JSON(404, gin.H{"error": "订单不存在"})
        return
    }

    c.JSON(200, order)
}
```

---

## 4. A02:2021 加密机制失效

### 4.1 风险说明

敏感数据未加密传输或存储，导致数据泄露。

```
常见场景：
- 使用 HTTP 而非 HTTPS
- 密码明文存储
- 敏感数据（身份证号、银行卡）明文存储
- 使用弱加密算法（MD5、DES）
```

### 4.2 防护措施

| 场景 | 防护措施 |
|------|----------|
| **传输中** | 强制 HTTPS，HSTS |
| **密码存储** | bcrypt/argon2 哈希加盐 |
| **敏感数据** | AES-256 加密存储 |
| **密钥管理** | 使用 KMS，不硬编码密钥 |

### 4.3 Go 代码示例

```go
// ❌ 错误：使用 MD5 存储密码
func HashPassword(password string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(password))) // 不安全！
}

// ✅ 正确：使用 bcrypt
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func CheckPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## 5. A03:2021 注入攻击

### 5.1 风险说明

攻击者通过输入恶意数据，改变后端查询或命令的逻辑。

```
常见类型：
- SQL 注入
- 命令注入
- LDAP 注入
- NoSQL 注入
```

### 5.2 SQL 注入示例

```
# 正常输入
用户名：zhangsan
密码：123456

SQL: SELECT * FROM users WHERE username='zhangsan' AND password='123456'

# 恶意输入
用户名：admin' --
密码：任意

SQL: SELECT * FROM users WHERE username='admin' -- ' AND password='xxx'
     ↑ 注释掉密码检查，直接以 admin 身份登录
```

### 5.3 防护措施

| 措施 | 说明 |
|------|------|
| **参数化查询** | 使用预编译语句，不拼接 SQL |
| **输入验证** | 白名单验证输入格式 |
| **ORM 框架** | 使用 GORM 等 ORM 自动防注入 |
| **最小权限** | 数据库账号只授予必要权限 |

### 5.4 Go 代码示例

```go
// ❌ 错误：SQL 拼接导致注入
func GetUserByUsername(username string) (*User, error) {
    sql := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    // 攻击者输入：admin' --
    return db.Exec(sql)
}

// ✅ 正确：参数化查询
func GetUserByUsername(username string) (*User, error) {
    var user User
    db.Where("username = ?", username).First(&user) // 参数自动转义
    return &user, nil
}
```

---

## 6. A04:2021 不安全设计

### 6.1 风险说明

设计阶段未考虑安全因素，导致系统架构存在固有缺陷。

```
常见场景：
- 未进行威胁建模
- 缺少安全设计评审
- 信任客户端输入
- 缺少防御深度
```

### 6.2 防护措施

| 措施 | 说明 |
|------|------|
| **威胁建模** | 识别潜在威胁和攻击面 |
| **安全设计模式** | 使用经过验证的安全模式 |
| **防御深度** | 多层防护，不依赖单一防护 |
| **安全评审** | 设计阶段进行安全评审 |

### 6.3 示例：密码重置流程设计

```
❌ 不安全设计：
- 通过邮箱发送新密码（明文）
- 重置链接无过期时间
- 重置 Token 可重复使用

✅ 安全设计：
- 发送一次性重置链接
- 链接 30 分钟过期
- Token 使用后立即失效
- 记录重置日志
```

---

## 7. A05:2021 安全配置错误

### 7.1 风险说明

系统配置不当导致的安全漏洞。

```
常见场景：
- 默认密码未修改（admin/admin）
- 调试模式在生产环境开启
- 目录列表功能开启
- 详细的错误信息泄露
- 未使用的服务/端口未关闭
```

### 7.2 防护措施

| 措施 | 说明 |
|------|------|
| **最小化原则** | 只开启必要的功能和服务 |
| **自动化检查** | 使用工具扫描配置 |
| **环境隔离** | 开发/生产环境配置分离 |
| **定期审计** | 定期检查配置变更 |

### 7.3 Go 代码示例

```go
// ❌ 错误：生产环境开启调试模式
gin.SetMode(gin.DebugMode) // 暴露详细路由信息

// ✅ 正确：根据环境变量配置
mode := os.Getenv("GIN_MODE")
if mode == "" {
    mode = gin.ReleaseMode
}
gin.SetMode(mode)
```

---

## 8. A06:2021 易受攻击的组件

### 8.1 风险说明

使用的第三方库/框架存在已知漏洞。

```
常见场景：
- 使用有已知 CVE 的库
- 依赖版本过旧
- 未及时更新补丁
```

### 8.2 防护措施

| 措施 | 说明 |
|------|------|
| **依赖扫描** | 使用 SCA 工具扫描漏洞 |
| **固定版本** | 使用 go.mod 锁定依赖版本 |
| **及时更新** | 关注安全公告，及时升级 |
| **最小依赖** | 减少不必要的依赖 |

### 8.3 Go 工具

```bash
# 检查依赖漏洞
govulncheck ./...

# 更新依赖
go get -u
go mod tidy

# 查看依赖详情
go list -m all
```

---

## 9. A07:2021 认证和会话失效

### 9.1 风险说明

身份认证和会话管理缺陷导致身份冒用。

```
常见场景：
- 弱密码策略
- 会话固定攻击
- Token 未设置过期时间
- 登出后 Token 未撤销
- 凭证泄露（明文 Cookie）
```

### 9.2 防护措施

| 措施 | 说明 |
|------|------|
| **多因素认证** | 关键操作启用 MFA |
| **强密码策略** | 最小长度 8，复杂度要求 |
| **会话管理** | Token 设置合理过期时间 |
| **安全登出** | 登出后撤销 Token |
| **安全传输** | 强制 HTTPS，HttpOnly Cookie |

---

## 10. A08:2021 软件和数据完整性故障

### 10.1 风险说明

软件更新或数据处理过程被篡改。

```
常见场景：
- 依赖来源不可信
- CI/CD 管道被入侵
- 数据被恶意修改
- 反序列化未验证
```

### 10.2 防护措施

| 措施 | 说明 |
|------|------|
| **来源验证** | 从官方源下载依赖 |
| **签名验证** | 验证软件包签名 |
| **完整性校验** | 校验 Hash 值 |
| **安全 CI/CD** | 保护构建管道 |

---

## 11. A09:2021 安全日志和监控缺失

### 11.1 风险说明

缺少安全日志和监控，无法检测和响应攻击。

```
常见场景：
- 登录失败无日志
- 敏感操作无审计
- 异常行为无告警
- 日志未集中管理
```

### 11.2 防护措施

| 措施 | 说明 |
|------|------|
| **审计日志** | 记录敏感操作 |
| **实时监控** | 检测异常行为 |
| **告警机制** | 达到阈值触发告警 |
| **日志保护** | 日志防篡改、加密存储 |

### 11.3 IAM 审计日志要求

| 事件 | 记录内容 | 保留期 |
|------|----------|--------|
| 登录成功/失败 | 用户 ID、IP、设备、时间 | 180 天 |
| 密码修改 | 用户 ID、IP、时间 | 180 天 |
| 权限变更 | 操作人、目标用户、变更内容 | 180 天 |
| 数据导出 | 用户 ID、数据类型、记录数 | 180 天 |

---

## 12. A10:2021 服务端请求伪造（SSRF）

### 12.1 风险说明

攻击者诱导服务端发起恶意请求。

```
常见场景：
- 文件上传（从 URL 下载）
- Webhook 回调
- API 代理
- 文件包含
```

### 12.2 攻击示例

```
# 正常请求
POST /api/upload
{ "url": "https://example.com/file.pdf" }

# 恶意请求
POST /api/upload
{ "url": "http://169.254.169.254/latest/meta-data/" }
     ↑ AWS 元数据服务，可获取实例凭证
```

### 12.3 防护措施

| 措施 | 说明 |
|------|------|
| **白名单校验** | 只允许访问预定义域名 |
| **内网访问控制** | 禁止访问内网 IP |
| **禁用重定向** | 不跟随 HTTP 重定向 |
| **响应过滤** | 过滤敏感响应头 |

### 12.4 Go 代码示例

```go
// ❌ 错误：未校验 URL
func DownloadFile(url string) ([]byte, error) {
    resp, err := http.Get(url) // 可能访问内网
    return io.ReadAll(resp.Body)
}

// ✅ 正确：校验 URL
func DownloadFile(urlStr string) ([]byte, error) {
    // 1. 解析 URL
    parsed, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }

    // 2. 检查协议白名单
    if parsed.Scheme != "https" && parsed.Scheme != "http" {
        return nil, errors.New("invalid scheme")
    }

    // 3. 解析 IP，检查是否内网
    ips, err := net.LookupIP(parsed.Hostname())
    if err != nil {
        return nil, err
    }
    for _, ip := range ips {
        if isInternalIP(ip) {
            return nil, errors.New("internal IP not allowed")
        }
    }

    // 4. 发起请求
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return errors.New("redirect not allowed")
        },
    }
    resp, err := client.Get(urlStr)
    if err != nil {
        return nil, err
    }

    return io.ReadAll(resp.Body)
}

// 检查是否内网 IP
func isInternalIP(ip net.IP) bool {
    // 检查私有地址段：10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
    // 检查 AWS 元数据：169.254.169.254
    // ...
}
```

---

## 13. 总结

### 13.1 IAM 系统重点关注

| 风险 | IAM 相关性 | 优先级 |
|------|-----------|--------|
| A01 失效的访问控制 | 核心风险 | P0 |
| A02 加密机制失效 | 密码/Token 安全 | P0 |
| A03 注入攻击 | 基础防护 | P0 |
| A07 认证和会话失效 | 核心风险 | P0 |
| A09 安全日志和监控 | 审计合规 | P1 |
| A10 SSRF | 第三方集成 | P2 |

### 13.2 安全开发清单

- [ ] 所有 API 进行访问控制校验
- [ ] 密码使用 bcrypt/argon2 加密
- [ ] 所有查询使用参数化
- [ ] 强制 HTTPS
- [ ] Token 设置过期时间
- [ ] 敏感操作记录审计日志
- [ ] 依赖定期扫描漏洞
- [ ] 错误信息不泄露敏感数据

---

## 14. 参考链接

- OWASP Top 10 官方：https://owasp.org/www-project-top-ten/
- OWASP 中文社区：https://github.com/OWASP/owasp-asvs
- CWE 漏洞库：https://cwe.mitre.org/
