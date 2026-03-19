---
name: go-code-library
description: Reusable Go code template library for backend scenarios such as PostgreSQL access, Redis access, configuration loading, HTTP server bootstrap, logging, graceful shutdown, and repository/service layering. Use when Codex needs to scaffold a Go service, provide reference snippets, or assemble common infrastructure modules for new Go code.
---

# Go Code Library

## Overview

Use the bundled templates to assemble production-oriented Go backend code quickly. Prefer copying from `assets/` and adapting names, configuration keys, and domain logic instead of drafting infrastructure code from scratch.

## Workflow

1. Identify the target shape:
   - Use `assets/basic-service-template/` when the user needs a small service skeleton.
   - Use `assets/snippets/` when the user only needs one focused module or code fragment.
2. Read [references/catalog.md](references/catalog.md) to locate the closest template.
3. Copy the relevant file or directory into the working tree.
4. Rename packages, adjust module path, and replace placeholder config values.
5. Keep business logic thin in handlers and place integration details in `internal/platform` or `internal/repository`.

## Guidance

- Prefer `pgx/v5` with `pgxpool` for PostgreSQL connection pooling.
- Prefer `go-redis/v9` for Redis clients.
- Keep configuration in a dedicated package and load it once at process startup.
- Return infrastructure clients as concrete wrappers or interfaces only when a test seam is needed.
- Use `context.Context` on all I/O paths.
- Add health/readiness handlers before adding business endpoints.

## Resources

- [references/catalog.md](references/catalog.md): Index of available templates and the intended use for each one.
- [references/customization-checklist.md](references/customization-checklist.md): What to rename or adjust after copying a template.
- `assets/basic-service-template/`: Minimal service skeleton with config, logging, PG, Redis, HTTP server, and graceful shutdown.
- `assets/snippets/`: Standalone code snippets for common modules.
