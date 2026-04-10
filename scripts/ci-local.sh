#!/bin/bash
# scripts/ci-local.sh
# Local CI validation script, run before pushing
set -e

echo "=== Step 1: Build ==="
go build ./...
echo "Build passed"

echo "=== Step 2: Unit tests ==="
go test ./... -race -cover -short
echo "Unit tests passed"

echo ""
echo "=== All checks passed ==="
echo "Run 'go test -tags=integration ./...' with Docker services for full integration tests."
