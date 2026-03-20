---
status: completed
summary: Modernized Makefile with .PHONY, golangci-lint v2, osv-scanner, gosec, and trivy targets; added .golangci.yml v2 config; updated tools.go with new dependencies; fixed CVE in modelcontextprotocol/go-sdk
container: s3-utils-003-add-golangci-lint-v2
dark-factory-version: v0.59.5-dirty
created: "2026-03-20T19:28:31Z"
queued: "2026-03-20T19:28:31Z"
started: "2026-03-20T19:28:39Z"
completed: "2026-03-20T19:36:12Z"
---

<summary>
- Add golangci-lint v2 to s3-utils with standard config
- Modernize Makefile to match other bborbe Go libraries
- Add `.golangci.yml` with standard linter configuration
- Update `tools.go` with current tool dependencies
</summary>

<objective>
Bring s3-utils Makefile, tools.go, and linting config to the same standard as other bborbe Go libraries (use kv as reference).
</objective>

<context>
Read CLAUDE.md for project conventions.
Read `docs/dod.md` for the Definition of Done criteria.

Reference: The `kv` library at `~/Documents/workspaces/kv/` has the standard Makefile pattern, `.golangci.yml`, and `tools.go`.

Current state:
- Makefile is outdated (missing `.PHONY`, missing lint/gosec/osv-scanner/trivy targets, old goimports-reviser invocation)
- No `.golangci.yml` exists
- `tools.go` imports deprecated `golang.org/x/lint/golint` and old `goimports-reviser` (not v3)
- No golangci-lint dependency at all
</context>

<requirements>
1. Replace Makefile with modern version matching kv pattern:
   - Add `.PHONY` declarations
   - Add `lint` target using `golangci-lint/v2/cmd/golangci-lint`
   - Add `osv-scanner`, `gosec`, `trivy` targets
   - Update `format` target to use goimports-reviser v3 with `-format -excludes vendor ./...` syntax
   - Add `golines` to format target
   - Add `go-modtool` to format target
   - Use `go mod tidy -e` instead of `go mod tidy`
   - Add `lint` to `check` target (not commented out like some other repos)
   - Keep project name as `github.com/bborbe/s3-utils`
   - Add `mkdir -p mocks` and `echo "package mocks" > mocks/mocks.go` to generate target
2. Create `.golangci.yml` matching kv's config but adapted for s3-utils (no repo-specific errname exclusions)
3. Update `tools.go`:
   - Remove deprecated `golang.org/x/lint/golint`
   - Add `github.com/golangci/golangci-lint/v2/cmd/golangci-lint`
   - Add `github.com/securego/gosec/v2/cmd/gosec`
   - Add `github.com/google/osv-scanner/v2/cmd/osv-scanner`
   - Add `github.com/segmentio/golines`
   - Add `github.com/shoenig/go-modtool`
   - Update `goimports-reviser` import if needed
4. Run `go get` for new dependencies and `go mod tidy`
5. Run `make precommit` — fix any lint issues found
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT refactor application code unrelated to lint fixes
- Keep the project name `github.com/bborbe/s3-utils` in all references
- Fix lint violations only in files that golangci-lint reports
</constraints>

<verification>
Run `make precommit` — must pass with exit code 0.
</verification>
