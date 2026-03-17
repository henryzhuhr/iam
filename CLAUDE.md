# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

IAM (身份认证与访问管理 / Identity and Access Management) - A multi-language project with Golang backend, designed to run in Docker containers.

## 技术栈

- **Golang**: 主要应用逻辑、API服务器和业务逻辑。
- **Python**: 业务接口测试，使用 uv 作为包管理器和运行环境。
- 后端微服务框架： go-zero
- 数据库： MySQL
- 缓存： Redis
- 消息队列： Kafka
- 配置管理： YAML 文件
- 容器化： Docker 和 Docker Compose

## 项目结构## 2. 项目结构

```bash
iam/
├── app/                    # 应用入口
│   └── main.go            # 主程序入口
├── etc/                    # 配置文件
│   └── dev.yaml            # 应用配置(开发环境配置)
├── infra/                  # 基础设施层（Infrastructure）
│   ├── cache/             # Redis 缓存封装
│   ├── database/          # MySQL 数据库连接
│   ├── executor/          # 批量执行器
│   └── queue/             # Kafka 消息队列
├── internal/               # 内部核心业务代码
│   ├── config/            # 配置结构体
│   ├── constant/          # 常量定义
│   ├── dto/               # 数据传输对象（Data Transfer Object）
│   ├── entity/            # 实体（Entity/Domain Model）
│   ├── handler/           # HTTP 处理器（Handler/Controller）
│   ├── repository/        # 数据访问层（Repository/DAO）
│   ├── routes/            # 路由注册
│   ├── service/           # 业务逻辑层（Service/Logic）
│   └── svc/               # 服务上下文（全局依赖注入容器）
├── sql/                    # SQL 脚本
├── scripts/                # 脚本文件
├── dockerfiles/            # Docker 相关文件
├── debug/                  # 调试脚本（Python）
├── docker-compose.yml      # Docker Compose 配置
├── go.mod                  # Go 模块依赖
└── README.md              # 项目说明
```

## 开发命令

```bash
go run app/main.go -f etc/dev.yaml
```

## 编码规范

### 目录结构说明

`internal/` 目录采用标准的分层架构，每个目录有明确的职责：

```bash
internal/
├── config/            # 配置结构体
│   └── config.go     # 定义 Config 结构体，包含所有配置项
├── constant/          # 常量定义
│   └── xxx.go        # 业务常量、错误码等
├── dto/               # 数据传输对象（Data Transfer Object）
│   └── <module>/     # 按业务模块划分子目录
│       └── xxx.go    # 定义 API 请求/响应结构
├── entity/            # 实体（Entity/Domain Model）
│   └── xxx.go        # 数据库表对应的实体结构
├── handler/           # HTTP 处理器（Handler/Controller）
│   └── <module>/     # 按业务模块划分子目录
│       └── xxx.go    # 接收 HTTP 请求、参数校验、调用 Service、返回响应
├── middleware/        # HTTP 中间件
│   └── xxx.go        # 用户代理、日志、认证等中间件
├── repository/        # 数据访问层（Repository/DAO）
│   └── xxx.go        # 数据库 CRUD 操作
├── routes/            # 路由注册
│   ├── routes.go     # 统一路由注册入口
│   └── <module>/     # 按业务模块划分子目录
│       ├── xxx.go    # 路由定义（使用 go-zero rest.Server）
│       └── xxx.swagger.yaml  # OpenAPI/Swagger 文档
├── service/           # 业务逻辑层（Service/Logic）
│   └── <module>/     # 按业务模块划分子目录
│       └── xxx.go    # 核心业务逻辑实现
└── svc/               # 服务上下文（全局依赖注入容器）
    └── servicecontext.go  # ServiceContext，包含所有依赖
```

### 新增标准接口的开发流程

1. **定义 DTO** (`internal/dto/<module>/`)
   - 定义请求和响应的数据结构

2. **实现 Service** (`internal/service/<module>/`)
   - 实现核心业务逻辑
   - 使用 `svcCtx` 访问依赖（数据库、缓存等）

3. **实现 Handler** (`internal/handler/<module>/`)
   - 接收 HTTP 请求
   - 参数校验
   - 调用 Service
   - 返回 JSON 响应

4. **注册路由** (`internal/routes/<module>/`)
   - 创建 `<module>.go` 定义路由
   - 创建 `<module>.swagger.yaml` 定义 OpenAPI 文档
   - 在 `internal/routes/routes.go` 中注册

5. **编写 Swagger 文档** (`internal/routes/<module>/<module>.swagger.yaml`)
   - 定义 API 路径、方法、参数
   - 定义请求/响应 Schema
   - 标注示例值

## 测试规范

## Git 工作流

## 注意事项
