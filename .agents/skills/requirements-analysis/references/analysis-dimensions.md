# Analysis Dimensions

Use this as a checklist. Do not force every item into the final answer; include only relevant dimensions.

## Business

- 要解决什么问题
- 成功标准是什么
- 谁是主要用户，谁是次要用户
- 哪些能力是刚需，哪些是增强项

## Scope

- 本次交付边界
- 明确不做的内容
- 与现有系统的关系
- 是否涉及迁移、兼容、替换

## Functional

- 核心流程
- CRUD 或业务动作
- 状态机 / 生命周期
- 搜索、筛选、排序、分页
- 批量操作
- 导入导出
- 通知与消息

## Authorization

- 哪些角色可以访问
- 谁能查看、创建、修改、删除、审批
- 是否存在租户隔离、数据隔离、字段级权限

## Data

- 核心实体
- 关键字段
- 唯一性约束
- 关联关系
- 审计字段
- 历史记录 / 软删除 / 版本化

## Integration

- 第三方登录 / 支付 / 消息 / 存储 / 搜索
- 回调、重试、幂等
- 上下游系统依赖

## Non-Functional

- 性能目标
- 安全要求
- 稳定性 / 可用性
- 监控告警
- 审计与合规
- 国际化 / 本地化

## Delivery

- MVP 切分
- 依赖前置项
- 技术风险
- 测试重点
- 发布与回滚

## Clarification Questions

Ask only the questions that materially change scope or implementation. Typical examples:

- 目标用户是谁，是否区分管理员与普通用户？
- 本期必须上线的最小能力是什么？
- 是否需要兼容现有数据、旧接口或旧流程？
- 是否有明确的性能、安全或合规要求？
- 是否会接入第三方系统或外部回调？
