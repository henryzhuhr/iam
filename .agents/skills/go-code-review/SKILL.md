---
name: go-code-review
description: Review Go code for correctness, concurrency, API design, error handling, tests, performance, security, and maintainability, and confirm suspicions with tool-based checks. Use when Codex is asked to review Go diffs, PRs, files, packages, or repos; find bugs or regressions; audit code quality; or verify findings with go test, go vet, gofmt, goimports, staticcheck, golangci-lint, ineffassign, revive, or govulncheck.
---

# Go Code Review

## Overview

Review Go changes with a bug-finding mindset, not a style-summary mindset. Combine direct reasoning over the diff and surrounding code with tool-backed evidence so findings are concrete and defensible.

Read `references/review-checklist.md` when you need a broader checklist or help interpreting tool output.

## Workflow

1. Identify review scope.
Determine whether the user wants a diff review, a file review, or a package or repo audit. Prefer the smallest scope that still supports a correct conclusion.

2. Build context before judging.
Read the changed code, surrounding call sites, interfaces, tests, and configuration. Check how data flows through handlers, services, goroutines, and error paths before reporting an issue.

3. Run tool-based verification.
Use `scripts/run-go-review-checks.sh` from the repo root. Pass changed directories or packages when the scope is narrow; otherwise let it inspect the whole module. Treat tool output as evidence, not as a substitute for reasoning.

4. Separate confirmed issues from suspicions.
If a tool proves the issue, say so explicitly. If the issue is based on code reasoning and tools are silent, label it as a reasoning-based finding instead of overstating certainty.

5. Report findings first.
Order findings by severity and include file references. Focus on bugs, regressions, race risks, broken error handling, contract mismatches, test gaps, and operational risk. Keep summaries brief.

## Review Priorities

- Check correctness first: nil handling, error propagation, wrong branches, ignored return values, shadowed variables, and broken invariants.
- Check concurrency next: goroutine lifetime, channel ownership, locking, context cancellation, and data races around shared state.
- Check public behavior: API compatibility, config defaults, HTTP status codes, serialization tags, database boundaries, and backward compatibility.
- Check tests and observability: missing test coverage for changed logic, brittle tests, silent failures, and unhelpful logs or metrics.
- Check maintainability after correctness: duplicated logic, misleading names, dead code, impossible branches, and confusing abstractions.

## Tooling Guidance

- Prefer `go test` and `go vet` as the baseline.
- Use `staticcheck` for correctness and API misuse signals that `go vet` misses.
- Use `golangci-lint` when present because repo-specific config often encodes team rules and extra analyzers.
- Use formatting checks such as `gofmt -l` and `goimports -l` to catch generated or unformatted edits, but do not elevate pure formatting to high-severity findings.
- Run optional tools such as `ineffassign`, `revive`, and `govulncheck` when available; mention when they were unavailable.

## Reporting Format

- Present findings first, ordered by severity.
- Include a short explanation of impact and the concrete code path.
- Include file references for each finding.
- Note which checks ran and which were unavailable or skipped.
- State explicitly when no findings were discovered, then mention residual risk or missing verification.

## Commands

Run the bundled checker from the repository root:

```bash
./.codex/skills/go-code-review/scripts/run-go-review-checks.sh
./.codex/skills/go-code-review/scripts/run-go-review-checks.sh ./internal/... ./app
```

If the user asked for a review of a specific diff, inspect the diff first and then run the checker on the affected directories or packages rather than on unrelated parts of the module.
