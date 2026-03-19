#!/usr/bin/env bash
set -u

status=0

tool_exists() {
  command -v "$1" >/dev/null 2>&1
}

print_section() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 2
}

append_unique() {
  local value
  for value in "$@"; do
    [ -n "$value" ] || continue
    case " $COLLECTED_TARGETS " in
      *" $value "*) ;;
      *) COLLECTED_TARGETS="${COLLECTED_TARGETS}${COLLECTED_TARGETS:+ }$value" ;;
    esac
  done
}

normalize_targets() {
  COLLECTED_TARGETS=""

  if [ "$#" -eq 0 ]; then
    return 0
  fi

  local target
  for target in "$@"; do
    if [ -f "$target" ]; then
      case "$target" in
        *.go) append_unique "$(dirname "$target")" ;;
        *) ;;
      esac
      continue
    fi

    append_unique "$target"
  done
}

collect_lint_targets() {
  LINT_TARGETS=()

  if [ -n "${COLLECTED_TARGETS:-}" ]; then
    local target
    for target in ${COLLECTED_TARGETS}; do
      LINT_TARGETS+=("$target")
    done
    return
  fi

  LINT_TARGETS=(./...)
}

collect_files() {
  if [ "${#PACKAGE_TARGETS[@]}" -eq 0 ]; then
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
      git ls-files '*.go'
    else
      find . -type f -name '*.go'
    fi
    return
  fi

  go list -f '{{range .GoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}{{range .CgoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}{{range .TestGoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}{{range .XTestGoFiles}}{{$.Dir}}/{{.}}{{"\n"}}{{end}}' "${PACKAGE_TARGETS[@]}" 2>/dev/null | sort -u
}

run_check() {
  local label="$1"
  shift

  print_section "$label"
  if "$@"; then
    return 0
  fi

  status=1
  return 1
}

run_optional_check() {
  local tool_name="$1"
  shift

  if tool_exists "$tool_name"; then
    run_check "$tool_name" "$tool_name" "$@"
  else
    print_section "$tool_name"
    printf 'SKIP: %s not found in PATH\n' "$tool_name"
  fi
}

tool_exists go || fail "go is required"
go env GOMOD >/dev/null 2>&1 || fail "go environment is unavailable"

normalize_targets "$@"
collect_lint_targets

PACKAGE_TARGETS=()
if [ -n "${COLLECTED_TARGETS:-}" ]; then
  while IFS= read -r line; do
    [ -n "$line" ] && PACKAGE_TARGETS+=("$line")
  done < <(go list ${COLLECTED_TARGETS} 2>/dev/null | sort -u)
fi

if [ "${#PACKAGE_TARGETS[@]}" -eq 0 ]; then
  PACKAGE_TARGETS=(./...)
fi

print_section "scope"
printf 'Packages: %s\n' "${PACKAGE_TARGETS[*]}"

GO_FILES=()
while IFS= read -r file; do
  [ -n "$file" ] && GO_FILES+=("$file")
done < <(collect_files)

if [ "${#GO_FILES[@]}" -gt 0 ]; then
  print_section "gofmt"
  gofmt_output="$(gofmt -l "${GO_FILES[@]}")"
  if [ -n "$gofmt_output" ]; then
    printf '%s\n' "$gofmt_output"
    status=1
  fi

  if tool_exists goimports; then
    print_section "goimports"
    goimports_output="$(goimports -l "${GO_FILES[@]}")"
    if [ -n "$goimports_output" ]; then
      printf '%s\n' "$goimports_output"
      status=1
    fi
  else
    print_section "goimports"
    printf 'SKIP: goimports not found in PATH\n'
  fi
else
  print_section "files"
  printf 'SKIP: no Go files found for the requested scope\n'
fi

run_check "go test" go test "${PACKAGE_TARGETS[@]}"
run_check "go vet" go vet "${PACKAGE_TARGETS[@]}"
run_optional_check staticcheck "${PACKAGE_TARGETS[@]}"

if tool_exists golangci-lint; then
  run_check "golangci-lint" golangci-lint run "${LINT_TARGETS[@]}"
else
  print_section "golangci-lint"
  printf 'SKIP: golangci-lint not found in PATH\n'
fi

run_optional_check ineffassign "${PACKAGE_TARGETS[@]}"
run_optional_check revive "${PACKAGE_TARGETS[@]}"
run_optional_check govulncheck "${PACKAGE_TARGETS[@]}"

exit "$status"
