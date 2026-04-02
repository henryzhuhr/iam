---
name: git-worktree-helper
description: 检查、创建和管理 Git worktree（工作树），用于并行分支开发。适用于需要按本地分支模式批量创建 worktree、把 worktree 放到仓库同级目录、避免为已在其他目录检出的分支重复创建 worktree、核对分支与 worktree 路径映射，或在多个分支上修改前先准备隔离工作区的场景。
---

# 🌳 Git Worktree Helper

用这个 skill 以安全、可重复的方式管理本地 Git worktree。

## 标准工作流程

前置约束：

- 创建 worktree 时，必须由用户明确提供**目标分支名**或**分支模式**（如 `feature/*`）。
- 如果用户没有提供分支名或模式，只能执行现状检查。
- 禁止把“未指定范围”解释为“为所有分支批量创建 worktree”，避免目录数量失控和工作区膨胀。
- 绝对不能默认对所有本地分支创建 worktree。

### 1. 先检查当前状态

1. 执行命令 `git branch --list '<pattern>'` 查看本地分支，执行 `git branch -r --list 'origin/<pattern>'` 查看远程分支，汇总后告诉用户当前仓库的分支（按照字母顺序排序），输出格式示例：

    ```bash
    🌿 分支列表

    | 分支 | 本地 | 远程 | 工作树 |
    | ---- | ---- | ---- | --- |
    | feature/xxx | ✅ | ✅ | ✅ |
    | bugfix/typo | ✅ | ❌ | |
    | release/v1.0 | ❌ | ✅ | |
    ```

2. 执行命令 `git worktree list --porcelain` 查看当前仓库的 worktree 列表，建立已有分支到路径的映射，输出格式示例：

    ```bash
    🌳 worktree 列表

    | 分支 | 路径 |
    | ---- | ---- |
    | main | /path/to/worktree/ |
    | feature/xxx | /path/to/worktree-feature~xxx/ |
    ```

3. 询问用户目标范围，明确区分本地分支和远程分支，确认用户想要创建 worktree 的分支列表。你要提供一些可选项供用户选择，例如

    ```bash
    请问你想要做什么？

    1. 为特定分支创建 worktree - 告诉我分支名（如 <当前的分支> 或 feature/* 模式）
    2. 切换到某个分支的工作树 - 告诉我分支名（如 <当前的分支> 或 feature/* 模式）
    3. 仅查看映射关系 - 当前状态已显示
    ```

4. 如果用户选择的是远程分支，你需要先检查是否有对应的本地分支，如果没有，你需要自动创建本地分支（`git checkout -b <branch> origin/<branch>`），再继续后续步骤。你需要明确告诉用户哪些远程分支被转换成了本地分支。
5. 把当前工作区所在分支视为已有 worktree，不要为它重复创建 worktree。
6. 如果用户只要求查看映射关系、核对现状或确认哪些分支已检出，到这里即可，不进入创建步骤。
7. 如果用户要求创建 worktree，继续执行后续步骤。

### 2. 选择目标目录布局

- 默认批量布局：在 `/var/dev` 目录下为每个分支创建一个 worktree，路径由 `<project_name>-<branch-name>` 组成，其中 `<project_name>` 是当前 Git 仓库所在目录的名称，`<branch-name>` 是分支名经过特殊字符替换后的版本。
- 为了避免分支名中存在 `/` 导致的路径层级过深问题，默认使用 `_` 替代 `/`，例如 `feature/login` 会被映射成 `feature_login`。
- 如果用户指定了基础目录，优先使用用户指定值而不是默认值。
- 如果用户已经明确给出目标路径，直接使用该路径，不必再推导默认布局。

目录结构示例：

```bash
/var/dev/
├── <project_name>-skill~git-worktree/ # 对应[技能]开发分支 skill/git-worktree
└── <project_name>-feature~login/      # 对应[功能]开发分支 feature/login
```

> 默认使用 `_` 替代 `/`，以降低与原分支名中连字符 `-` 混用时的路径碰撞风险。

路径映射示例：

- `skill/git-worktree` -> `/var/dev/<project_name>-skill_git-worktree`
- `feature/login` -> `/var/dev/<project_name>-feature_login`

创建前检查：

- 批量处理时，先构造完整的 `branch -> target_path` 计划表。
- 在执行前检查目标路径唯一性；如果两个分支映射到同一路径，立刻停止，并明确报出冲突分支对。
- 如果目标路径已经存在，但没有注册为 worktree，立刻停止并把它作为人工冲突暴露出来。

### 3. 创建 worktree

- 批量创建时，先列出命中的本地分支，再逐个推导目标路径，形成明确的创建计划。
- 当范围超过一个分支时，先用自然语言汇总计划，至少说明命中的分支、对应路径、哪些分支会被跳过。
- 确认计划后，再逐个运行 `git worktree add <path> <branch>` 正式执行。
- 如果只缺一个分支，直接使用 `git worktree add <path> <branch>` 即可。

示例：

```bash
git branch --list 'skill/*'
git branch -r --list 'origin/skill/*'
git worktree list --porcelain
git worktree add /var/dev/<project_name>-skill_git-worktree skill/git-worktree
git worktree add /var/dev/<project_name>-feature_login feature/login
```

### 4. 验证结果

- 创建完成后运行 `git worktree list --porcelain`，确认新路径已注册。
- 对每个新建的 worktree 运行 `git -C <worktree-path> branch --show-current`，确认检出的分支正确。
- 核对计划与结果是否一致，尤其是 `created`、`skipped`、`failed/conflicted` 三类数量和明细。
- 明确汇报哪些分支因为已有 worktree 被跳过，哪些路径被新建，哪些分支因为冲突未处理。

### 5. 切换到工作树

1. 如果刚才用户明确要求创建 worktree，并且创建成功，你需要询问用户是否需要切换到新建的工作树目录：

   ```bash
   是否需要切换到新建的工作树目录？（y/n）
   ```

2. 如果用户确认需要切换，那么你需要自动完成切换目录的操作，例如：

   ```bash
   cd /var/dev/<project_name>-skill_git-worktree
   ```

3. 切换完成后，你需要再次确认**分支状态**和**当前所在目录**，必须告诉用户，例如：

   ```bash
   当前的分支是 skill/git-worktree
   当前目录是 /var/dev/<project_name>-skill_git-worktree
   ```

## 安全规则

- 不要为已经有 worktree 的分支再创建第二个 worktree。
- 除非用户明确要求，否则不要删除、`prune` 或移动 worktree。
- 不要为了管理本地 worktree 而推送远端变更。
- 如果检测到路径映射冲突，或者目标路径已经存在但没有注册为 worktree，立刻停止并把它作为人工冲突暴露出来。

## 执行要求

- 优先使用标准 Git 命令，不依赖额外脚本。
- 批量处理时，先汇总计划，再执行创建，避免一次性盲目落地。
- 无匹配分支：属于正常无操作结果，应明确说明“没有匹配到候选分支，因此未创建 worktree”。
- 已有 worktree：属于正常跳过，应在汇报中列出分支和现有路径。
- 路径映射冲突或目标路径冲突：属于失败，停止执行，等待人工处理。
- 非 Git 仓库、分支不存在、命令执行失败：属于失败，直接停止并暴露错误。
- 汇报结果时明确区分：已存在的 worktree、这次新建的 worktree、因冲突而未处理的分支。
