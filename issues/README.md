# Issue 跟踪目录

本目录用于记录 IAM 项目开发过程中的所有 Issue。

## Issue 编号规则

- 文件名格式：`NNN-short-kebab-case.md`
- 编号从 001 开始，依次递增
- 新 Issue 延续当前最大编号，不跳号、不重用

## Issue 列表

| 编号 | 标题 | 状态 | 创建日期 |
|------|------|------|----------|
| 001 | [需求缺口收敛清单](./001-requirements-gap-checklist.md) | 已完成 | 2026-04-21 |
| 002 | [明确注册模型与租户归属规则](./002-clarify-registration-and-tenant-binding.md) | 已完成 | 2026-04-21 |
| 003 | [确认 REQ-017 应用隔离最终方案](./003-finalize-application-isolation-scope.md) | 已完成 | 2026-04-21 |
| 004 | [对齐内部服务认证需求与设计](./004-align-internal-service-auth-design.md) | 已完成 | 2026-04-21 |
| 005 | [冻结 v1 权限模型边界](./005-freeze-permission-model-v1.md) | 已完成 | 2026-04-21 |

## 新增 Issue 流程

1. 查看当前最大编号
2. 创建新文件 `NNN-title.md`
3. 在本文件中添加索引记录
4. 在 Issue 文件中编写详细内容

## Issue 模板

```markdown
# {Issue 标题}

## 问题描述
<!-- 简要描述问题 -->

## 背景
<!-- 问题产生的背景 -->

## 解决方案
<!-- 解决方案或讨论 -->

## 相关链接
<!-- 关联的 PR、文档等 -->
```
