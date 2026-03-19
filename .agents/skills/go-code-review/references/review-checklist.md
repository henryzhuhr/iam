# Go Review Checklist

## Core findings

- Broken behavior: wrong condition, wrong status code, wrong field mapping, stale interface contract, or silently ignored error.
- Reliability risk: panic path, nil dereference, partial write, leaked resource, missing timeout, or non-idempotent retry.
- Concurrency risk: race on shared state, lock order issue, blocked send or receive, runaway goroutine, or missing context cancellation.
- Data integrity risk: transactional gap, inconsistent validation, lost update, or schema mismatch.
- Security risk: authz bypass, unsafe path handling, injection surface, secret exposure, or unbounded request processing.

## Common Go-specific smells

- `defer` inside a loop that grows resource lifetime unexpectedly.
- Loop variable capture in goroutines or closures.
- Returning typed nil through an interface.
- Mutating shared slices or maps without synchronization.
- Comparing errors with `==` when wrapping is expected.
- Dropping `context.Context` or replacing it with `context.Background()`.
- Logging and swallowing an error instead of propagating it.
- Public function returning ambiguous zero values without documentation.
- Table-driven tests missing edge cases introduced by the change.

## Tool interpretation

- `gofmt -l`: formatting drift only; low severity unless generated files or CI would fail.
- `goimports -l`: missing import normalization; usually low severity.
- `go test`: failing tests are strong evidence of a regression, especially when they touch changed code paths.
- `go vet`: useful for suspicious constructs, copylocks, printf issues, and unreachable problems.
- `staticcheck`: high-value signal for correctness, dead code, and API misuse.
- `golangci-lint`: respect repo config; findings may include style and correctness, so separate the two.
- `ineffassign`: useful confirmation for ignored writes and dead assignments.
- `revive`: mainly maintainability unless the repo treats specific rules as required.
- `govulncheck`: use for dependency or call-path vulnerability evidence, not as a general security substitute.

## Evidence rules

- Prefer quoting the exact failing condition or tool output in summary form instead of dumping raw logs.
- When a tool is silent but the bug is still real, explain the reasoning path clearly.
- When tools disagree, prefer the code path and repo configuration over generic lint advice.
