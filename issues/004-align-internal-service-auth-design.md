# 对齐内部服务认证需求与设计

## 问题描述

REQ-018 当前存在三个层面的不一致：用户故事映射错误、接口路径不一致、表结构设计不一致。继续开发会直接造成客户端认证模型返工。

## 背景

- REQ-018 关联的用户故事目前错误地指向第三方登录故事
- REQ 中的取 Token 接口与 TDD 中的接口路径不一致
- REQ 采用“客户端 + 凭证”两张表，SQL 目前是单表 `clients`

这类问题会直接影响凭证轮换、审计、禁用、scope 授权和中间件设计。

## 解决方案

### 已落定结论

- [x] 新增并绑定“平台管理员管理内部客户端”“内部服务换取 token”用户故事
- [x] 统一客户端 token 申请接口为 `/api/v1/clients/token`
- [x] 统一客户端数据模型为“客户端 + 凭证”双表
- [x] 明确 AK/SK 轮换后旧凭证按状态和失效时间淘汰
- [x] 明确 `scope` 与用户 RBAC 完全分离

### 需要回写的文档

- [x] 更新 `REQ-018` 的用户故事、API、数据库设计
- [x] 更新 TDD 中 `/auth` 与 `/clients` 的职责划分
- [x] 更新 SQL 中客户端相关表设计与字段命名
- [x] 更新产品摘要和系统架构中的主体模型说明

### 推荐补充内容

- [x] 补充客户端禁用、过期、轮换、审计查询的时序
- [x] 补充客户端凭证是否支持多套并存和灰度切换
- [x] 补充 `aud`、`subject_type`、`jti` 等 claim 的校验约束

### 完成标准

- [x] REQ、TDD、SQL 对内部客户端模型保持一致
- [x] 能明确回答客户端 token 入口、凭证轮换策略和 scope 校验位置
- [x] 用户故事和验收标准不再与 OAuth2 第三方登录混淆

## 处理结果

已将内部服务认证统一到 `clients + client_credentials` 双表模型，修正错误的用户故事引用，并对齐 `/clients/token` 路径与轮换策略。

## 相关链接

- [docs/PRD/05-functional-requirements/REQ-018-internal-service-authentication.md](../docs/PRD/05-functional-requirements/REQ-018-internal-service-authentication.md)
- [docs/PRD/10-appendix-user-stories.md](../docs/PRD/10-appendix-user-stories.md)
- [docs/TDD/001-iam-system-architecture-design.md](../docs/TDD/001-iam-system-architecture-design.md)
- [sql/001_init.sql](../sql/001_init.sql)
