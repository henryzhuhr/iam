---
name: git-worktree
description: 检查、创建和管理 Git worktree，用于并行分支开发。适用于 Codex 需要按本地分支模式批量创建 worktree（例如 `skill/*`）、把 worktree 放到仓库同级目录、避免为已在其他目录检出的分支重复创建 worktree、核对分支与 worktree 路径映射，或在多个分支上修改前先准备隔离工作区的场景。
---

# Git Worktree

用这个 skill 以安全、可重复的方式管理本地 Git worktree。

## 工作流

### 1. 先检查当前状态

- 运行 `git branch --list '<pattern>'` 查看候选本地分支。
- 运行 `git worktree list --porcelain` 建立已有分支到路径的映射。
- 把当前工作区所在分支视为已有 worktree，不要为它重复创建 worktree。

### 2. 选择目标目录布局

- 默认批量布局：在仓库根目录同级创建一个名为 `<repo-name>-worktrees` 的目录。
- 每个 worktree 目录名由分支名推导而来，把 `/` 替换成 `-`。
- 如果用户指定了基础目录，优先使用用户指定值而不是默认值。

默认目录结构示例：

```text
/path/to/
├── repo/
└── repo-worktrees/
```

如果仓库是 `/path/to/repo`，分支模式是 `skill/*`，命中分支 `skill/git-worktree` 和 `skill/openai-docs`，则最终生成的 worktree 结构类似：

```text
/path/to/
├── repo/
└── repo-worktrees/
    ├── skill-git-worktree/
    └── skill-openai-docs/
```

路径映射示例：

- `skill/git-worktree` -> `repo-worktrees/skill-git-worktree`
- `feature/login` -> `repo-worktrees/feature-login`

### 3. 创建 worktree

- 批量创建时优先使用 `scripts/create_worktrees.sh`。
- 当范围超过一个分支时，先运行 `scripts/create_worktrees.sh --repo <repo> --pattern 'skill/*' --dry-run` 预览计划。
- 确认计划后，再去掉 `--dry-run` 正式执行。
- 如果只缺一个分支，直接使用 `git worktree add <path> <branch>` 也可以。

示例：

```bash
.agents/skills/git-worktree/scripts/create_worktrees.sh \
  --repo /path/to/repo \
  --pattern 'skill/*' \
  --dry-run

.agents/skills/git-worktree/scripts/create_worktrees.sh \
  --repo /path/to/repo \
  --pattern 'skill/*'
```

### 4. 验证结果

- 创建完成后运行 `git worktree list`。
- 如有需要，运行 `git -C <worktree-path> branch --show-current` 做抽样检查。
- 明确汇报哪些分支因为已有 worktree 被跳过，哪些路径被新建。

## 安全规则

- 不要为已经有 worktree 的分支再创建第二个 worktree。
- 除非用户明确要求，否则不要删除、`prune` 或移动 worktree。
- 不要为了管理本地 worktree 而推送远端变更。
- 如果目标路径已经存在，但没有注册为 worktree，立刻停止并把它作为人工冲突暴露出来。

## 脚本

使用 `scripts/create_worktrees.sh` 做确定性的批量创建。

- 必填：`--pattern`
- 可选：`--repo`，默认当前目录
- 可选：`--target-base`，默认 `<repo-parent>/<repo-name>-worktrees`
- 可选：`--dry-run`，只打印计划，不真正创建
