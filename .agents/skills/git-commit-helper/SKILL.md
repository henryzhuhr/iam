---
name: git-commit-helper
description: Use When 提交代码到 git 仓库，使用 commit message
---

# Git Commit Helper

用于提交代码到 git 仓库，使用 commit message。

你需要自带判断当前项目的主要语言，并根据语言习惯生成合适的 commit message。

## 提交信息格式

提交信息应遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<类型>: <简短描述>

[可选的正文]

[可选的脚注]
```

### 类型说明

| 类型 | 说明 | Emoji |
|------|------|-------|
| `feat` | 新功能 | ✨ |
| `fix` | 修复 bug | 🐛 |
| `docs` | 文档更新 | 📝 |
| `style` | 代码格式调整 | 💄 |
| `refactor` | 代码重构 | ♻️ |
| `perf` | 性能优化 | ⚡ |
| `test` | 测试相关 | ✅ |
| `chore` | 构建/工具变动 | 🔧 |
| `ci` | CI/CD 配置 | 👷 |
| `revert` | 回滚提交 | ⏪ |

### 示例

```
/commit ✨ feat: 新增用户登录功能
/commit 🐛 fix: 修复密码重置接口的空指针异常
/commit 📝 docs: 更新 README.md 安装说明
/commit ♻️ refactor: 重构用户管理模块代码结构
```

## 提交流程

1. 检查 git 状态
2. 显示变更文件列表
3. 自动执行 `git commit` 提交代码
4. 显示提交结果
5. 询问用户：是否执行 `git push` 将提交推送到远程仓库

## 注意事项

- 提交前自动检查是否有未暂存的修改
- 如有多个不相关的修改，建议分批提交
- 提交信息应简洁明了，首行不超过 50 个字符
