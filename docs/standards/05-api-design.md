# API 设计规范

> 最后更新：2026-03-29
> 适用范围：IAM 项目所有 RESTful API

---

## 1. API 风格

### 1.1 RESTful 原则

- 使用 HTTP 方法表示操作
- 使用名词复数表示资源
- 使用 HTTP 状态码表示结果

| 操作 | HTTP 方法 | 路径示例 |
|------|----------|----------|
| 查询列表 | GET | `/api/v1/users` |
| 查询详情 | GET | `/api/v1/users/{id}` |
| 创建资源 | POST | `/api/v1/users` |
| 更新资源 | PUT | `/api/v1/users/{id}` |
| 部分更新 | PATCH | `/api/v1/users/{id}` |
| 删除资源 | DELETE | `/api/v1/users/{id}` |

### 1.2 API 版本

- 版本号放在 URL 路径中
- 格式：`/api/v{version}/`

```
✅ 正确：/api/v1/users
❌ 错误：/api/users/v1
```

---

## 2. 请求规范

### 2.1 请求头

|  header | 说明 | 示例 |
|--------|------|------|
| `Authorization` | 认证 Token | `Bearer eyJhbGciOiJSUzI1NiIs...` |
| `Content-Type` | 请求体类型 | `application/json` |
| `X-Request-ID` | 请求唯一标识 | `uuid-12345-abcde` |
| `X-Tenant-ID` | 租户 ID | `tenant-67890` |

### 2.2 查询参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `page` | 页码 | `?page=2` |
| `page_size` | 每页数量 | `?page_size=20` |
| `sort` | 排序字段 | `?sort=created_at` |
| `order` | 排序方向 | `?order=desc` |
| `search` | 搜索关键词 | `?search=zhangsan` |

### 2.3 请求体

```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "张三"
}
```

---

## 3. 响应规范

### 3.1 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 12345,
    "email": "user@example.com",
    "name": "张三",
    "created_at": "2026-03-29T10:00:00Z"
  }
}
```

### 3.2 列表响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {"id": 1, "name": "用户 1"},
      {"id": 2, "name": "用户 2"}
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 3.3 错误响应

```json
{
  "code": 1001,
  "message": "用户不存在",
  "data": null,
  "details": {
    "field": "user_id",
    "reason": "record not found"
  }
}
```

---

## 4. HTTP 状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|----------|
| 200 OK | 成功 | GET/PUT/PATCH 成功 |
| 201 Created | 已创建 | POST 创建资源成功 |
| 204 No Content | 无内容 | DELETE 成功 |
| 400 Bad Request | 请求错误 | 参数校验失败 |
| 401 Unauthorized | 未授权 | Token 无效或缺失 |
| 403 Forbidden | 禁止访问 | 权限不足 |
| 404 Not Found | 未找到 | 资源不存在 |
| 409 Conflict | 冲突 | 资源已存在 |
| 429 Too Many Requests | 请求过多 | 触发限流 |
| 500 Internal Server Error | 服务器错误 | 系统异常 |

---

## 5. 错误码规范

### 5.1 错误码格式

```
错误码 = 模块码 (2 位) + 子模块码 (2 位) + 错误序号 (2 位)
```

### 5.2 模块码分配

| 模块码 | 模块 |
|--------|------|
| 01 | 认证模块 |
| 02 | 用户模块 |
| 03 | 租户模块 |
| 04 | 角色模块 |
| 05 | 权限模块 |
| 06 | MFA 模块 |
| 07 | 审计日志 |
| 08 | 应用模块 |
| 09 | 客户端模块 |
| 99 | 系统错误 |

### 5.3 错误码示例

| 错误码 | 说明 |
|--------|------|
| 010101 | 用户名或密码错误 |
| 010102 | Token 已过期 |
| 010103 | Token 无效 |
| 020101 | 用户不存在 |
| 020102 | 用户已存在 |
| 020103 | 用户已被禁用 |
| 030101 | 租户不存在 |
| 999999 | 系统繁忙，请稍后重试 |

---

## 6. 安全规范

### 6.1 认证要求

- 所有 API 默认需要认证（公开 API 除外）
- 使用 Bearer Token 认证
- Token 放在 `Authorization` 请求头

### 6.2 敏感数据

- 密码不在响应中返回
- 密码在请求中使用加密传输（HTTPS）
- 敏感字段脱敏展示

### 6.3 限流

- 单用户限流：100 次/分钟
- 单 IP 限流：1000 次/分钟
- 敏感接口限流：10 次/分钟（登录、改密）

---

## 7. 文档规范

### 7.1 Swagger 注释

```go
// @Summary 创建用户
// @Description 创建一个新的用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "用户信息"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/users [post]
func CreateUser(c *gin.Context) {
    // ...
}
```

### 7.2 API 文档生成

```bash
# 使用 swag 生成 Swagger 文档
swag init -g app/main.go -o docs/api

# 访问 Swagger UI
# http://localhost:8080/swagger/index.html
```

---

## 8. 参考链接

- REST API Best Practices: https://restfulapi.net/
- OpenAPI Specification: https://swagger.io/specification/
- Google API Design Guide: https://cloud.google.com/apis/design/
