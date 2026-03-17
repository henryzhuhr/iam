# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

IAM (身份认证与访问管理 / Identity and Access Management) - A multi-language project with Golang backend, designed to run in Docker containers.

## 技术栈

- **Golang**: 主要应用逻辑、API服务器和业务逻辑。
- **Python**: 业务接口测试，使用 uv 作为包管理器和运行环境。

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

## 测试规范

## Git 工作流

## 注意事项
