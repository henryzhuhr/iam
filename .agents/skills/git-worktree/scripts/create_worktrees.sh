#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  create_worktrees.sh --pattern <branch-pattern> [--repo <repo>] [--target-base <dir>] [--dry-run]

Examples:
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
      echo "[ERROR] Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ -z "$pattern" ]]; then
  echo "[ERROR] --pattern is required" >&2
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
  echo "[INFO] No local branches matched pattern: ${pattern}"
  exit 0
fi

echo "[INFO] Repo root: ${repo_root}"
echo "[INFO] Pattern: ${pattern}"
echo "[INFO] Target base: ${target_base}"
if [[ "$dry_run" -eq 1 ]]; then
  echo "[INFO] Dry run mode enabled"
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
    echo "[SKIP] ${branch} already has worktree: ${existing_path}"
    skipped_count=$((skipped_count + 1))
    continue
  fi

  branch_dir="${branch//\//-}"
  worktree_path="${target_base}/${branch_dir}"

  if [[ -e "$worktree_path" ]]; then
    echo "[ERROR] Target path already exists and is not registered as a worktree: ${worktree_path}" >&2
    failed_count=$((failed_count + 1))
    continue
  fi

  if [[ "$dry_run" -eq 1 ]]; then
    echo "[PLAN] git -C ${repo_root} worktree add ${worktree_path} ${branch}"
    continue
  fi

  echo "[CREATE] ${branch} -> ${worktree_path}"
  git -C "$repo_root" worktree add "$worktree_path" "$branch"
  created_count=$((created_count + 1))
done

echo "[INFO] Summary: created=${created_count} skipped=${skipped_count} failed=${failed_count}"

if [[ "$failed_count" -gt 0 ]]; then
  exit 1
fi
