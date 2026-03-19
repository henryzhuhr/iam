#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
用法:
  create_worktrees.sh --pattern <branch-pattern> [--repo <repo>] [--target-base <dir>] [--dry-run]

示例:
  create_worktrees.sh --repo /path/to/repo --pattern 'skill/*' --dry-run
  create_worktrees.sh --pattern 'feature/*'
  create_worktrees.sh --pattern 'skill/*' --target-base /tmp/repo-worktrees
EOF
}

repo="."
pattern=""
target_base=""
dry_run=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      repo="${2:?missing value for --repo}"
      shift 2
      ;;
    --pattern)
      pattern="${2:?missing value for --pattern}"
      shift 2
      ;;
    --target-base)
      target_base="${2:?missing value for --target-base}"
      shift 2
      ;;
    --dry-run)
      dry_run=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "[错误] 未知参数: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ -z "$pattern" ]]; then
  echo "[错误] 必须提供 --pattern" >&2
  usage >&2
  exit 1
fi

repo_root="$(git -C "$repo" rev-parse --show-toplevel)"

if [[ -z "$target_base" ]]; then
  repo_parent="$(dirname "$repo_root")"
  repo_name="$(basename "$repo_root")"
  target_base="${repo_parent}/${repo_name}-worktrees"
fi

branches=()
while IFS= read -r branch; do
  if [[ -n "$branch" ]]; then
    branches+=("$branch")
  fi
done < <(git -C "$repo_root" for-each-ref --format='%(refname:short)' "refs/heads/${pattern}" | sort)

if [[ ${#branches[@]} -eq 0 ]]; then
  echo "[信息] 没有匹配分支模式的本地分支: ${pattern}"
  exit 0
fi

echo "[信息] 仓库根目录: ${repo_root}"
echo "[信息] 分支模式: ${pattern}"
echo "[信息] 目标根目录: ${target_base}"
if [[ "$dry_run" -eq 1 ]]; then
  echo "[信息] 已启用 dry-run 模式"
fi

if [[ "$dry_run" -eq 0 ]]; then
  mkdir -p "$target_base"
fi

created_count=0
skipped_count=0
failed_count=0

for branch in "${branches[@]}"; do
  existing_path="$(
    git -C "$repo_root" worktree list --porcelain | awk -v branch="$branch" '
      $1 == "worktree" { path = substr($0, 10) }
      index($0, "branch refs/heads/") == 1 {
        current = substr($0, 19)
        if (current == branch) {
          print path
          exit
        }
      }
    '
  )"

  if [[ -n "$existing_path" ]]; then
    echo "[跳过] ${branch} 已存在 worktree: ${existing_path}"
    skipped_count=$((skipped_count + 1))
    continue
  fi

  branch_dir="${branch//\//-}"
  worktree_path="${target_base}/${branch_dir}"

  if [[ -e "$worktree_path" ]]; then
    echo "[错误] 目标路径已存在，且未注册为 worktree: ${worktree_path}" >&2
    failed_count=$((failed_count + 1))
    continue
  fi

  if [[ "$dry_run" -eq 1 ]]; then
    echo "[计划] git -C ${repo_root} worktree add ${worktree_path} ${branch}"
    continue
  fi

  echo "[创建] ${branch} -> ${worktree_path}"
  git -C "$repo_root" worktree add "$worktree_path" "$branch"
  created_count=$((created_count + 1))
done

echo "[信息] 汇总: created=${created_count} skipped=${skipped_count} failed=${failed_count}"

if [[ "$failed_count" -gt 0 ]]; then
  exit 1
fi
