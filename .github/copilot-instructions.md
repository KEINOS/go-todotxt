---
applyTo: "**"
---
# Go-todotxt Project Instructions

## PREREQUISITE: Load Base Instructions First

**Before applying any go-todotxt-specific rules, you MUST:**

1. Fetch and read the base instructions from:
   - URL: https://raw.githubusercontent.com/KEINOS/copilot-instructions/refs/heads/main/common-instructions.md
   - Use `fetch_webpage` tool to retrieve this file
2. Apply the general Go development guidelines from `common-instructions.md`
3. Then apply go-todotxt specific overrides and extensions below

> **Note**: This file contains **only go-todotxt specific instructions**. All general Go development rules (TDD, naming, error handling, etc.) are in `common-instructions.md`. Both must be followed together with priority to this file for conflicts.

## Go-Todotxt-Specific Quick Reference

- **Linting Configuration**: See `.golangci.yml` for enabled/disabled linters
- **Approved Dependencies**: See `.golangci.yml` for allowed third-party packages
- **Error Package**: Use standard `errors` package
- **Mock Strategy**: Use `httptest` if needed for testing (no live API calls are allowed)
- **README Sections**: Update Features, Use Cases, API Reference when adding public APIs of this module
- **GoDoc Oriented**: Must document all public APIs with GoDoc comments.

## Basic Flow of Fixing Linting Issues

**IMPORTANT**: Always ensure all tests pass before and after each linting fixes.

1. Run `go test ./...` to ensure all tests pass before linting
2. Run `golangci-lint run` to identify linting issues
3. Apply `golangci-lint run --fix` to auto-fix issues where possible if any issues are found
4. Run `golangci-lint run` again to see remaining issues
5. Target one issue at a time and fix them manually
    - Use `golangci-lint run --enable-only=<linter_name>` to focus on a specific linter if needed
6. Run `go test ./...` again to ensure no tests are broken after fixes
7. Repeat steps 2-6 until `golangci-lint run` shows no issues
