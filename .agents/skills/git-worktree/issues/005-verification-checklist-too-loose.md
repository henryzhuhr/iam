# 验证步骤不够具体

## 现状

`SKILL.md` 里的验证步骤目前只有：

- 创建后运行 `git worktree list`
- 如有需要再抽样检查分支

## 问题

这个验证流程偏松，执行者容易只做表面检查，没有确认：

- 新路径是否真的注册成 worktree
- worktree 实际检出的分支是否正确
- 汇总数字是否与预期一致

## 影响

- 创建后可能遗漏映射错误
- 批量执行结果不容易复核
- skill 的可重复性不够强

## 建议

把验证步骤改成固定 checklist：

- 运行 `git worktree list --porcelain`，确认新路径已注册
- 对新建路径执行 `git -C <path> branch --show-current`，确认分支正确
- 核对脚本汇总中的 `created`、`skipped`、`failed`
- 汇报中明确列出新建路径和被跳过分支

这样 skill 会更接近可执行 SOP，而不是仅提供方向。

## 相关位置

- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/SKILL.md`
