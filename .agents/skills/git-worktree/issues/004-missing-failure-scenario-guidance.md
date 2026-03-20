# 缺少失败场景和输出解释

## 现状

脚本已经覆盖了一些典型场景，例如：

- 没有匹配到本地分支时打印信息并退出 0
- 参数缺失时退出 1
- 目标路径冲突时累计失败并最终退出 1

但 `SKILL.md` 没有说明这些场景该如何解释，也没有告诉调用方哪些属于正常无操作结果，哪些属于真正失败。

## 问题

调用方如果只看 skill 文档，容易把“无匹配分支”误判为异常，或者忽略“部分成功但最终失败”的情况。

## 影响

- agent 汇报可能不一致
- 自动化流程里可能错误处理退出状态或输出
- 用户难以理解 dry-run、skip、failed 的语义差异

## 建议

在 `SKILL.md` 增加常见结果说明：

- 无匹配分支：正常结束，无需创建
- 已有 worktree：正常跳过，应在汇报中列出
- 路径冲突：失败，需要人工处理
- 参数错误或非 Git 仓库：失败，直接停止

最好再补一条：汇报时必须区分 `created`、`skipped`、`failed`。

## 相关位置

- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/SKILL.md`
- `/Users/henryzhuhr/project/iam/.agents/skills/git-worktree/scripts/create_worktrees.sh`
