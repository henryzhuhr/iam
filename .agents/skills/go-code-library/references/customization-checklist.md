# Customization Checklist

Apply these changes after copying a template:

1. Replace the module path in `go.mod`.
2. Rename package paths under `internal/` to match the target repository layout.
3. Change environment variable names if the project already has conventions.
4. Tune HTTP timeouts, database pool sizes, and Redis address/DB settings.
5. Replace the sample `/healthz` endpoint with real liveness/readiness checks if needed.
6. Add migrations, repositories, and business services instead of putting SQL or Redis logic in handlers.
7. Run `gofmt` and `go test ./...` after integrating the copied code.
