# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v0.0.18

- chore: Migrate to tools.env + Makefile @version pattern; remove tools.go and obsolete replace block. go.mod reduced from 483 to 58 lines.

## v0.0.17

- Update bborbe/* dependencies (errors, sentry, service, argument, collection, etc.)
- Update golangci-lint to v2.11.4 with revised linter config (depguard, timeout 15m)
- Update osv-scanner to v2.3.5, add new CVE ignores to .trivyignore and .osv-scanner.toml
- Update numerous transitive dependencies (otel, prometheus, charmbracelet, google, etc.)
- Improve vulncheck Makefile target to filter known-ignored vulnerabilities

## v0.0.16

- Update go version to 1.26.2
- Bump aws/smithy-go to v1.24.3
- Add vulnerability ignores for GHSA-xmrv-pmrh-hhx2 and CVE-2026-33817
- Pin anthropic-sdk-go and diskfs/go-diskfs in replace directives

## v0.0.15

- Update AWS SDK v2 dependencies (v1.41.4, credentials v1.19.12)
- Update Docker/containerd/moby dependencies
- Update OpenTelemetry to v1.40.0
- Enable parallel golangci-lint runners
- Clean up go.mod excludes and replace directives

## v0.0.14

- fix: install trivy in CI workflow to resolve missing binary error

## v0.0.13

- chore: add golangci-lint v2 with `.golangci.yml` config and modernize Makefile with `.PHONY`, `lint`, `osv-scanner`, `gosec`, `trivy` targets
- chore: update `tools.go` to add golangci-lint/v2, gosec, osv-scanner, golines, go-modtool and remove deprecated golint
- fix: upgrade `modelcontextprotocol/go-sdk` to v1.4.1 to resolve two high-severity CVEs

## v0.0.12

- chore: verified project health — all tests pass, linting succeeds, precommit exits 0

## v0.0.11

- chore: verified all tests pass, linting and precommit checks succeed

## v0.0.10

- Improve README with installation, usage, CLI command docs, and examples

## v0.0.9

- Fix DisableHTTPS detection: swap strings.HasPrefix args so HTTPS URLs correctly enable TLS

## v0.0.8

- Split s3-client-creator.go into one file per type (URL, AccessKey, SecretKey, CreateS3Client)
- Add Ginkgo tests for URL, AccessKey, SecretKey, and CreateS3Client

## v0.0.7

- Update Go from 1.25.5 to 1.26.1
- Update AWS SDK v2 to v1.41.1
- Update bborbe/* dependencies
- Update golang.org/x/* packages
- Remove tracked upload binary

## v0.0.6

- Update Go to 1.25.7
- Update AWS SDK v2 dependencies
- Update testing dependencies (ginkgo, gomega)
- Add .update-logs/ and .mcp-* to .gitignore

## v0.0.5

- Update Go to 1.25.5
- Update golang.org/x/crypto to v0.47.0
- Update dependencies

## v0.0.4

- Fix stdin seeking issue in upload command by buffering stdin data before upload

## v0.0.3

- Modify upload command to read content from stdin instead of static string

## v0.0.2

- Add typed S3 credentials (URL, AccessKey, SecretKey) with String() methods
- Change package name from `s3` to `s3utils` to avoid conflicts with AWS SDK
- Update all commands to use new typed credentials

## v0.0.1

- Initial Version
