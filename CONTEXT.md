# CONTEXT.md

## Build/Test Commands
- `just golangci` (alias: `just ci`) - Run golangci-lint
- `just test-units` (alias: `just tu`) - Run unit tests with `-tags=unit`
- `just test-integrations` (alias: `just ti`) - Run integration tests with `-tags=integration`
- `just test-e2e` - Run end-to-end tests with `-tags=e2e`
- `just test-all` - Run all tests
- `go test -tags=unit ./path/to/package` - Run single package unit tests
- `go test -tags=integration ./path/to/package` - Run single package integration tests
- NEVER run `go build`

## Code Style Guidelines
- **Imports**: Standard library first, blank line, third-party, blank line, local packages
- **Naming**: PascalCase for exported types/functions, camelCase for unexported, Entity suffix for domain models
- **Error Handling**: Return errors explicitly, use `errors.Join()` for wrapping, define package-level error vars
- **Database**: Use sqlc for type-safe queries, transactions with defer rollback pattern
- **Testing**: Use build tags (`//go:build integration`), table-driven tests, testify/assert
- **Handlers**: Return `error`, use Echo context, validate payloads, render templates with `Render(renderArgs(ctx))`
- **Services**: Accept context first, return domain entities, handle business logic
- **Comments**: Minimal comments, focus on why not what, use TODO for future improvements
- **Styling**: Use vanilla CSS for creating styles. NEVER use inline styles.
- **Views**: Only edit the .templ files and run `just ct` to compile the result
