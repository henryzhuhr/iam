# 路径名映射存在碰撞风险

## 现状

`SKILL.md` 规定 worktree 目录名由分支名推导，把 `/` 替换成 `-`。

脚本也按这个规则实现：

- `branch_dir="${branch//\//-}"`

## 问题

不同分支可能映射到同一路径，例如：

- `a/b-c` -> `a-b-c`
- `a-b/c` -> `a-b-c`

当前实现下，第二个分支只会在运行时遇到“目标路径已存在”的错误，无法区分这是：

- 分支名映射冲突
- 用户目录里本来就有同名路径

## 影响

- dry-run 结果不够可信，不能提前暴露全部冲突
- 批量创建时错误信息不够明确
- 使用者可能误以为是本地目录脏状态，而不是命名规则有缺陷

## 建议

- 在 skill 文档里明确要求“创建前检查目标路径唯一性”
- 在脚本里先构造 `branch -> target_path` 计划表，再检测重复路径
- 如果检测到两个分支映射到同一路径，直接报出冲突分支对并退出

## 相关位置

- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/SKILL.md`
- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/scripts/create_worktrees.sh`
