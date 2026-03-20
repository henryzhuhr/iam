# 缺少使用边界说明

## 现状

当前 skill 主要描述了批量创建 worktree 的流程，但没有明确说明哪些情况应该直接使用原生命令，而不是走脚本。

## 问题

缺少“适用范围”和“非适用范围”会导致调用方在简单场景下也套用完整流程，增加不必要的操作复杂度。

## 影响

- 单分支场景被过度流程化
- 用户明确指定目标路径时，agent 仍可能先走 pattern 扫描和 dry-run
- skill 的使用成本偏高

## 建议

在 `SKILL.md` 补一段边界说明，例如：

- 单分支且目标路径已明确时，优先直接 `git worktree add <path> <branch>`
- 只有在批量按模式创建多个分支时，才优先使用 `scripts/create_worktrees.sh`
- 如果用户只要求查看映射关系，则只执行检查步骤，不进入创建步骤

## 相关位置

- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/SKILL.md`
