# Go 编码规范

> 最后更新：2026-03-29
> 适用范围：IAM 项目所有 Go 代码

---

## 1. 命名规范

### 1.1 包名

- 小写，无下划线
- 短小精悍
- 避免与标准库重名

```go
// ✅ 正确
package handler
package service
package repository

// ❌ 错误
package Handler      // 大写
package user_handler // 下划线
package string       // 与标准库重名
```

### 1.2 变量名

- 小写，驼峰式
- 简短但有意义
- 局部变量可缩写

```go
// ✅ 正确
var userName string
var ctx context.Context
var wg sync.WaitGroup

// ❌ 错误
var UserName string    // 导出变量才大写首字母
var user_name string   // 不用下划线
var thisIsAVeryLongVariableNameThatIsHardToRead string
```

### 1.3 常量名

- 大写首字母表示导出
- 驼峰式，无下划线

```go
// ✅ 正确
const MaxRetryCount = 3
const defaultTimeout = 30 * time.Second

// ❌ 错误
const MAX_RETRY_COUNT = 3  // 不用全大写
const default_timeout = 30 // 不用下划线
```

### 1.4 函数名

- 大写首字母表示导出
- 动词 + 名词结构

```go
// ✅ 正确
func CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*User, error)
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*User, error)

// ❌ 错误
func create_user(ctx context.Context)  // 不用下划线
func (s *UserService) getUser(ctx context.Context) // 导出函数首字母大写
```

---

## 2. 注释规范

### 2.1 包注释

每个包必须有包注释，位于 `package` 声明之前。

```go
// Package handler provides HTTP handlers for user management.
package handler
```

### 2.2 函数注释

导出的函数必须有注释。

```go
// CreateUser creates a new user with the given request.
// Returns error if user already exists or validation fails.
func CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*User, error)
```

### 2.3 注释格式

- 使用 `//` 单行注释
- 注释与代码之间空一行
- 注释首字母大写，句末加句号

```go
// ✅ 正确
// ValidateUser checks if the user data is valid.
func ValidateUser(user *User) error

// ❌ 错误
// validate user  // 格式不规范
func ValidateUser(user *User) error
```

---

## 3. 错误处理

### 3.1 错误检查

- 错误必须立即处理
- 不要忽略错误

```go
// ✅ 正确
result, err := doSomething()
if err != nil {
    return err
}

// ❌ 错误
result, err := doSomething()  // 忽略错误
_ = result
```

### 3.2 错误包装

- 使用 `fmt.Errorf` 或 `errors.Wrap` 包装错误
- 添加上下文信息

```go
// ✅ 正确
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ 错误
if err != nil {
    return err  // 丢失上下文
}
```

### 3.3 错误返回

- 错误消息小写开头
- 不要以标点符号结尾

```go
// ✅ 正确
return errors.New("user not found")

// ❌ 错误
return errors.New("User not found!")  // 不要大写和感叹号
```

---

## 4. 代码组织

### 4.1 导入顺序

```go
import (
    // 标准库
    "context"
    "fmt"
    "time"

    // 第三方库
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // 项目内部
    "github.com/iam/internal/dto"
    "github.com/iam/internal/service"
)
```

### 4.2 结构体组织

```go
type UserService struct {
    // 导出字段在前
    Logger *zap.Logger

    // 非导出字段在后
    repo *repository.UserRepository
    cache *cache.UserCache
}
```

---

## 5. 最佳实践

### 5.1 上下文使用

- 第一个参数传递 `context.Context`
- 不要存储 context 在结构体中

```go
// ✅ 正确
func (s *Service) DoSomething(ctx context.Context, id int64) error

// ❌ 错误
type Service struct {
    ctx context.Context  // 不要存储 context
}
```

### 5.2 指针使用

- 大结构体使用指针传递
- 基本类型不需要指针

```go
// ✅ 正确
func ProcessUser(user *User) error  // 结构体用指针
func Calculate(x int) int           // 基本类型不用指针

// ❌ 错误
func ProcessUser(user User) error   // 大结构体拷贝开销
func Calculate(x *int) int          // 基本类型不需要指针
```

### 5.3 接口设计

- 接口要小（1-2 个方法）
- 接受接口，返回结构体

```go
// ✅ 正确
type Reader interface {
    Read(p []byte) (n int, err error)
}

func Process(r io.Reader) error

// ❌ 错误
type BigInterface interface {
    Method1()
    Method2()
    Method3()
    Method4()
}
```

---

## 6. 工具

### 6.1 代码格式化

```bash
# 格式化代码
go fmt ./...

# 检查代码问题
go vet ./...

# 检查依赖
go mod tidy
```

### 6.2 静态检查

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行检查
golangci-lint run
```

---

## 7. 参考链接

- Effective Go: https://go.dev/doc/effective_go
- Go Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments
- Uber Go Style Guide: https://github.com/uber-go/guide/blob/master/style.md
