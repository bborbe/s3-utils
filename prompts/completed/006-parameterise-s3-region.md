---
status: completed
summary: 'Parameterised S3 region via a WithRegion functional option: CreateS3Client now defaults to the SDK''s unset-region behaviour (restoring pre-v0.0.28 backward compatibility), Garage callers opt in with WithRegion("garage"), and all four CLI commands read the region from S3_REGION'
execution_id: s3-utils-region-param-exec-006-parameterise-s3-region
dark-factory-version: dev
created: "2026-08-18T10:56:23Z"
queued: "2026-08-18T10:56:23Z"
started: "2026-08-18T10:57:22Z"
completed: "2026-08-18T11:00:57Z"
---

# Parameterise S3 region with a WithRegion option

<summary>
- CreateS3Client accepts an optional region via a functional option; the default is the SDK's unset-region behaviour (as before any region code was added)
- Callers targeting Garage pass WithRegion("garage") explicitly; the library no longer hardcodes a Seibert-specific region
- This restores backward compatibility: v0.0.28's hardcoded "garage" region changed behaviour for every consumer, including non-Garage targets
- The four in-repo CLI commands (upload, download, list-objects, list-buckets) read the region from a new S3_REGION env var with an empty default
- Existing tests still pass; new tests cover the default (empty) and the override
- Existing three-arg CreateS3Client calls compile and behave unchanged — full backward compatibility restored
</summary>

<objective>
Remove the hardcoded Garage region from the generic s3-utils library, restoring the original unset-region default so the library is generic and backward compatible, and let Garage consumers opt in via an explicit WithRegion("garage").
</objective>

<context>
Module: github.com/bborbe/s3-utils.

CreateS3Client currently hardcodes `Region: "garage"` in s3_client_creator.go (added in v0.0.28 to pass Garage's SigV4 validation). That was a **backward-incompatible change**: it altered the SigV4 credential scope for every consumer of this generic library, not just Garage. This prompt restores the library's original behaviour — region unset, exactly as it was before v0.0.28 — and moves the Garage-specific value to the call sites that need it.

Garage is Seibert-specific infrastructure (the IT replacement for the EOL MinIO server). It belongs at the Seibert call sites (the 7 Octopus services + the CLI commands), not in a generic library's default.

The library is consumed by 7 external services plus 4 in-repo CLI commands (cmd/upload, cmd/download, cmd/list-objects, cmd/list-buckets). The 7 external services must add `WithRegion("garage")` to their CreateS3Client calls when they migrate to Garage — that is the Seibert-side migration work, tracked in the Octopus S3 inventory (vault page [[Octopus S3 Inventory]] § Garage compatibility analysis).

Existing client creation pattern — s3_client_creator.go CreateS3Client(s3Url URL, s3AccessKey AccessKey, s3SecretKey SecretKey) *s3.Client. CLI callers build the client in their main.go Run method with those three positional args.

Follow the functional-options pattern per the coding plugin guide `go-functional-options-pattern.md` (singular option type + `With*` constructors, "Configuring External Types" variant).
</context>

<requirements>
1. In `s3_client_creator.go`, add a functional-options pattern. Define:

   ```go
   type CreateS3ClientOption func(*s3.Options)
   ```

   Change `CreateS3Client` to accept variadic `opts ...CreateS3ClientOption` as a fourth parameter. **Do NOT set a Region in the `s3.Options` literal** — leave it unset (the SDK's natural default, exactly as the code was before v0.0.28). Apply the options via a loop before returning the client:

   ```go
   for _, opt := range opts {
       opt(&options)
   }
   ```

   Provide the option constructor:

   ```go
   func WithRegion(region string) CreateS3ClientOption
   ```

   which sets `options.Region = region`. Follow the existing code style (see the surrounding file; there is no prior functional-options example in this repo, so keep it minimal and idiomatic — a plain variadic-option closure is sufficient).

2. Update the four CLI callers in `cmd/{upload,download,list-objects,list-buckets}/main.go`:
   - Add a `S3Region string` field to each `application` struct with the exact tag: `` `required:"false" arg:"s3-region" env:"S3_REGION" usage:"S3 region for SigV4 signing (empty = SDK default)"` `` (no `default:` — empty default, matching the sibling fields' style; the sibling S3Url/S3AccessKey/S3SecretKey carry `required:"false"` and a `usage:` line).
   - Pass it through only when non-empty: `s3utils.CreateS3Client(url, access, secret, s3utils.WithRegion(a.S3Region))` — note `WithRegion("")` sets an empty region, which is equivalent to the SDK default, so passing it unconditionally is also acceptable; prefer the unconditional form for simplicity.

3. Update `README.md`: in the "All commands share these common environment variables / flags" table (currently listing `--s3-url`/`S3_URL`, `--s3-access-key`/`S3_ACCESS_KEY`, `--s3-secret-key`/`S3_SECRET_KEY`), add a new row:
   ```
   | `--s3-region` | `S3_REGION` | Region for SigV4 signing (empty = SDK default) |
   ```

4. In `s3_client_creator_test.go`, keep the existing tests (they must still pass). Add a new `It` case:
   - Default: `client := s3utils.CreateS3Client("https://s3.example.com", "access", "secret")`, assert `client.Options().Region` equals `""` — the original pre-v0.0.28 behaviour (the SDK exposes `Options()` returning `s3.Options`; the unset region produces the empty-credential-scope signing that MinIO and AWS accept).
   - Override: `client := s3utils.CreateS3Client("https://s3.example.com", "access", "secret", s3utils.WithRegion("garage"))`, presign a GetObject and assert the signed URL's credential scope contains `garage` (mirror the existing boundary test that asserts `%2Fgarage%2Fs3%2Faws4_request`).
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- **Do NOT set a default region** — the library's default must be the unset-region behaviour from before v0.0.28 (backward compatibility is the point of this change)
- Garage is a Seibert-specific value: it belongs in call sites via WithRegion("garage"), NOT in the library
- Do NOT touch vendor/ (build-time artifact)
- Keep the API backward-compatible: existing three-arg CreateS3Client calls must compile unchanged (variadic option is optional) and behave as before v0.0.28
- Update CHANGELOG.md: add a bullet under `## Unreleased` (create the section if absent, keeping the preamble block intact)
- Existing tests must still pass
- Use repo-relative paths only
</constraints>

<verification>
Run `make precommit` — must pass.
</verification>
