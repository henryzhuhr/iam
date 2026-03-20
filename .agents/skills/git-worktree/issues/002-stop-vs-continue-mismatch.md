# 安全规则与脚本行为不一致

## 现状

`SKILL.md` 写的是：

- 如果目标路径已经存在，但没有注册为 worktree，立刻停止并把它作为人工冲突暴露出来

但脚本实际行为是：

- 记录一次失败
- 继续处理后续分支
- 最后统一以非零状态退出

## 问题

skill 承诺的是“立刻停止”，实现却是“继续跑完再失败”。

这会让使用者对安全边界产生错误预期，尤其是在批量创建多个 worktree 时。

## 影响

- 文档与实际行为不一致
- 调用方难以判断冲突后是否还会继续产生部分副作用
- 在保守场景下，当前实现不够严格

## 建议

二选一并保持一致：

- 保守方案：脚本在检测到未注册的同名路径后立刻退出
- 批处理方案：保留现有脚本行为，但把 skill 文案改成“记录冲突并在末尾失败退出”

如果倾向安全优先，建议采用前者。

## 相关位置

- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/SKILL.md`
- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/scripts/create_worktrees.sh`
