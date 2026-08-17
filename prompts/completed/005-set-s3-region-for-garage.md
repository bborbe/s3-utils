---
status: completed
summary: Set S3 client region to garage in CreateS3Client and added a regression test pinning the SigV4 credential scope region
execution_id: s3-utils-region-exec-005-set-s3-region-for-garage
dark-factory-version: dev
created: "2026-08-17T14:38:51Z"
queued: "2026-08-17T14:38:51Z"
started: "2026-08-17T14:39:13Z"
completed: "2026-08-17T14:41:33Z"
---

# Set S3 client region to garage

<summary>
- The S3 client created by CreateS3Client now sends a SigV4 credential scope containing the region `garage` instead of an empty region
- Requests to Garage (the IT-hosted S3 replacement for the EOL MinIO server) authenticate successfully
- Existing behavior for MinIO is unchanged — MinIO accepts any region in the credential scope, so path-style requests still work
- A regression test pins the client's region so the value cannot silently regress to empty
- No API change: CreateS3Client keeps its existing signature
</summary>

<objective>
Make the S3 client produced by CreateS3Client compatible with Garage, which rejects SigV4 requests whose credential scope has an empty region. The client must sign with the region `garage`.
</objective>

<context>
The Go module is github.com/bborbe/s3-utils. It wraps aws-sdk-go-v2.

The production client is built by CreateS3Client in `s3_client_creator.go`. It currently sets UsePathStyle, EndpointResolver, EndpointOptions and Credentials — but NOT Region. When the SDK signs a request, the SigV4 credential scope contains an empty region, which Garage rejects with:

```
api error AuthorizationHeaderMalformed: Authorization header malformed,
unexpected scope: '20260817//s3/aws4_request',
expected: '20260817/garage/s3/aws4_request'
```

MinIO (the current target) accepts any region, so this only surfaced once testing against Garage. See the "Garage compatibility analysis" section of the vault page [[Octopus S3 Inventory]] for the full context.

This is the minimal "make it work" fix. Parameterising the region as a config surface is explicitly OUT OF SCOPE for this prompt.
</context>

<requirements>
1. In `s3_client_creator.go`, inside the `s3.Options{...}` literal of `CreateS3Client`, add a `Region` field with the value `"garage"` (place it alongside the existing `UsePathStyle: true` entry).

   Old:
   ```go
   s3.Options{
       UsePathStyle:     true,
       EndpointResolver: s3.EndpointResolverFromURL(s3Url.String()),
   ```
   New:
   ```go
   s3.Options{
       UsePathStyle:     true,
       Region:           "garage",
       EndpointResolver: s3.EndpointResolverFromURL(s3Url.String()),
   ```

2. In `s3_client_creator_test.go`, extend the existing `Describe("CreateS3Client", ...)` block with a new `It` case that signs a real request through the SigV4 boundary and asserts the credential scope contains the region `garage`:

   ```go
   It("signs requests with region garage in the credential scope", func() {
       client := s3utils.CreateS3Client("https://s3.example.com", "access", "secret")
       presigner := s3.NewPresignClient(client)
       req, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
           Bucket: aws.String("b"), Key: aws.String("k"),
       }, func(o *s3.PresignOptions) {})
       Expect(err).NotTo(HaveOccurred())
       Expect(req.SignedURI).To(ContainSubstring("X-Amz-Credential=access%2F%2Fgarage%2Fs3%2Faws4_request"))
   })
   ```

   This exercises the exact boundary that failed against Garage — an empty region in the SigV4 credential scope (`20260817//s3/aws4_request`) — and proves the fix produces `garage` in the scope. It is network-free (presigning signs locally; no server contact).

   Match the existing test style (Ginkgo v2 + Gomega, package `s3utils_test`). The SDK packages `github.com/aws/aws-sdk-go-v2/aws` and `github.com/aws/aws-sdk-go-v2/service/s3` are already direct dependencies of this module — import them as needed.
</requirements>

<constraints>
- Do NOT commit — dark-factory handles git
- Do NOT change the CreateS3Client signature or add new parameters/options
- Do NOT parameterise the region via config/env — that is explicitly out of scope for this change
- Existing tests must still pass
- Do NOT touch vendor/ (build-time artifact)
- Use repo-relative paths only
</constraints>

<verification>
Run `make precommit` — must pass.
</verification>
