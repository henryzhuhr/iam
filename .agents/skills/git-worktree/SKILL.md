---
name: git-worktree
description: Inspect, create, and manage Git worktrees for parallel branch development. Use when Codex needs to batch-create one worktree per local branch pattern such as `skill/*`, place worktrees in a sibling directory, avoid duplicate worktrees for branches already checked out elsewhere, verify branch-to-worktree mappings, or prepare isolated workspaces before making changes on multiple branches.
---

# Git Worktree

Use this skill to manage local Git worktrees safely and repeatably.

## Workflow

### 1. Inspect the current state first

- Run `git branch --list '<pattern>'` to see candidate local branches.
- Run `git worktree list --porcelain` to build the existing branch-to-path mapping.
- Treat the current workspace branch as an existing worktree. Do not create a duplicate worktree for it.

### 2. Choose the target layout

- Default batch layout: create a sibling directory named `<repo-name>-worktrees` next to the repository root.
- Derive each worktree directory from the branch name by replacing `/` with `-`.
- If the user asks for a specific base directory, honor it instead of the default.

### 3. Create worktrees

- For batch creation, prefer `scripts/create_worktrees.sh`.
- Run `scripts/create_worktrees.sh --repo <repo> --pattern 'skill/*' --dry-run` first when the scope is larger than one branch.
- Re-run without `--dry-run` after confirming the plan.
- For a single missing branch, `git worktree add <path> <branch>` is fine.

Example:

```bash
.agents/skills/git-worktree/scripts/create_worktrees.sh \
  --repo /path/to/repo \
  --pattern 'skill/*' \
  --dry-run

.agents/skills/git-worktree/scripts/create_worktrees.sh \
  --repo /path/to/repo \
  --pattern 'skill/*'
```

### 4. Verify the result

- Run `git worktree list` after creation.
- If needed, run `git -C <worktree-path> branch --show-current` for spot checks.
- Report which branches were skipped because they already had worktrees and which paths were created.

## Safety Rules

- Never create a second worktree for a branch that already has one.
- Never remove, prune, or move worktrees unless the user explicitly asks.
- Never push remote changes just to manage local worktrees.
- If a target path already exists but is not registered as a worktree, stop and surface it as a manual conflict.

## Script

Use `scripts/create_worktrees.sh` for deterministic batch creation.

- Required: `--pattern`
- Optional: `--repo` defaults to the current directory
- Optional: `--target-base` defaults to `<repo-parent>/<repo-name>-worktrees`
- Optional: `--dry-run` prints the plan without creating anything
