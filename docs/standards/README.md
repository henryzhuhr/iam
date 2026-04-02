# 项目规范目录

本目录收录 IAM 项目开发过程中的各类规范和标准，供研发团队遵循。

## 规范清单

| 编号 | 规范名称 | 说明 | 文档 |
|------|----------|------|------|
| 01 | 项目结构规范 | 目录结构、分层架构、文件命名 | [01-project-structure.md](./01-project-structure.md) |
| 02 | Go 编码规范 | Go 语言编码风格、命名约定、注释规范 | [02-go-coding-style.md](./02-go-coding-style.md) |
| 03 | Git 工作流规范 | 分支策略、提交规范、Code Review | [03-git-workflow.md](./03-git-workflow.md) |
| 04 | 术语表 | 项目统一术语和定义 | [04-glossary.md](./04-glossary.md) |
| 05 | API 设计规范 | RESTful API 设计规范、错误码规范 | [05-api-design.md](./05-api-design.md) |
| 06 | 数据库设计规范 | 表结构设计、索引、命名规范 | [06-database-design.md](./06-database-design.md) |

## 推荐阅读顺序

1. [项目结构规范](./01-project-structure.md) - 了解项目整体结构
2. [Go 编码规范](./02-go-coding-style.md) - 编码前必读
3. [Git 工作流规范](./03-git-workflow.md) - 提交代码前必读
4. [术语表](./04-glossary.md) - 统一项目语言

## 维护约定

- 新增规范时在 `standards/` 目录下按编号顺序创建
- 规范变更需要团队评审后更新
- 规范编号永久有效，删除后不重用
