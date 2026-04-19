# 项目结构规范

> 最后更新：2026-03-29
> 适用范围：IAM 项目

---

## 1. 项目概述

IAM (身份认证与访问管理 / Identity and Access Management) 是一个以 Golang 后端为主、运行于 Docker 容器环境的多语言项目，包含 Vue 3 前端控制台。

## 2. 技术栈

| 技术 | 用途 |
|------|------|
| **Golang** | 主要应用逻辑、API 服务器和业务逻辑 |
| **Python** | 业务接口测试，使用 `uv` 作为包管理器和运行环境 |
| **Vue 3 / TypeScript** | 前端 Web 控制台，Vite 构建 |
| **go-zero** | 后端微服务框架 |
| **MySQL** | 数据库 |
| **Redis** | 缓存 |
| **Kafka** | 消息队列 |
| **YAML** | 配置管理 |
| **Docker/Docker Compose** | 容器化 |

## 3. 项目结构

```bash
iam/
├── app/                    # 应用入口
│   └── main.go             # 主程序入口
├── etc/                    # 配置文件
│   └── dev.yaml            # 开发环境配置
├── infra/                  # 基础设施层（Infrastructure）
│   ├── cache/              # Redis 缓存封装
│   ├── database/           # MySQL 数据库连接
│   ├── executor/           # 批量执行器
│   └── queue/              # Kafka 消息队列
├── internal/               # 内部核心业务代码
│   ├── config/             # 配置结构体
│   ├── constant/           # 常量定义
│   ├── dto/                # 数据传输对象
│   ├── entity/             # 实体模型
│   ├── handler/            # HTTP 处理器
│   ├── repository/         # 数据访问层
│   ├── routes/             # 路由注册
│   ├── service/            # 业务逻辑层
│   └── svc/                # 服务上下文
├── web/                    # 前端源码目录
│   ├── public/             # 静态资源
│   ├── index.html          # HTML 入口
│   └── src/                # 前端源码（Vue 3 + TypeScript）
├── sql/                    # SQL 脚本
├── scripts/                # 脚本文件
├── dockerfiles/            # Docker 相关文件
├── debug/                  # 调试脚本（Python）
├── docker-compose.yml      # Docker Compose 配置
├── go.mod                  # Go 模块依赖
├── package.json            # 前端构建配置
├── vite.config.ts          # Vite 构建配置
├── tsconfig.json           # TypeScript 配置
├── node_modules/           # 前端依赖（.gitignore）
└── README.md               # 项目说明
```

## 4. internal/ 分层约定

`internal/` 目录采用标准分层架构：

```bash
internal/
├── config/            # 配置结构体
│   └── config.go      # 定义 Config 结构体，包含所有配置项
├── constant/          # 常量定义
│   └── xxx.go         # 业务常量、错误码等
├── dto/               # 数据传输对象
│   └── <module>/      # 按业务模块划分子目录
│       └── xxx.go     # 定义 API 请求/响应结构
├── entity/            # 实体模型
│   └── xxx.go         # 数据库表对应的实体结构
├── handler/           # HTTP 处理器
│   └── <module>/      # 按业务模块划分子目录
│       └── xxx.go     # 接收请求、参数校验、调用 Service、返回响应
├── middleware/        # HTTP 中间件
│   └── xxx.go         # 用户代理、日志、认证等中间件
├── repository/        # 数据访问层
│   └── xxx.go         # 数据库 CRUD 操作
├── routes/            # 路由注册
│   ├── routes.go      # 统一路由注册入口
│   └── <module>/      # 按业务模块划分子目录
│       ├── xxx.go     # 路由定义（使用 go-zero rest.Server）
│       └── xxx.swagger.yaml  # OpenAPI/Swagger 文档
├── service/           # 业务逻辑层
│   └── <module>/      # 按业务模块划分子目录
│       └── xxx.go     # 核心业务逻辑实现
└── svc/               # 服务上下文
    └── servicecontext.go  # ServiceContext，包含所有依赖
```

## 5. 开发命令

```bash
# 启动后端服务
go run app/main.go -f etc/dev.yaml

# 启动前端开发服务器（构建文件在根目录，无需切换目录）
npm run dev

# 前端构建
npm run build

# 前端单元测试 + 组件测试
npm run test:unit

# 前端 E2E 测试
npm run test:e2e

# 全部测试（单元 + 组件 + E2E）
npm run test
```

## 6. 新增标准接口的开发流程

1. **定义 DTO**：在 `internal/dto/<module>/` 中定义请求和响应结构。
2. **实现 Service**：在 `internal/service/<module>/` 中实现核心业务逻辑，通过 `svcCtx` 访问依赖。
3. **实现 Handler**：在 `internal/handler/<module>/` 中接收请求、校验参数、调用 Service，并返回 JSON 响应。
4. **注册路由**：在 `internal/routes/<module>/` 中定义路由，并在 `internal/routes/routes.go` 注册。
5. **编写 Swagger 文档**：在 `internal/routes/<module>/<module>.swagger.yaml` 中维护接口路径、参数和请求/响应 Schema。

## 7. Issues 规范

### 7.1 目录和文件命名规范

- 整个项目统一使用名为 `issues/` 的目录记录 issue，不限定于 skill。
- issue 文件名使用三位递增编号开头，格式为 `NNN-short-kebab-case.md`，例如 `001-path-name-collision.md`。
- 新 issue 必须延续当前最大编号，不能跳号，也不要重用已有编号。
- 新增 issue 时，需要同步在对应层级的 `README.md` 中维护 index，方便按编号查阅。

## 8. 相关规范

| 规范 | 说明 | 文档 |
|------|------|------|
| Go 编码规范 | Go 语言编码风格、命名约定、注释规范 | [02-go-coding-style.md](./02-go-coding-style.md) |
| Git 工作流规范 | 分支策略、提交规范、Code Review | [03-git-workflow.md](./03-git-workflow.md) |
| 术语表 | 项目统一术语和定义 | [04-glossary.md](./04-glossary.md) |
| API 设计规范 | RESTful API 设计规范、错误码规范 | [05-api-design.md](./05-api-design.md) |
| 数据库设计规范 | 表结构设计、索引、命名规范 | [06-database-design.md](./06-database-design.md) |
