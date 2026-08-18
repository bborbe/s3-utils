# S3 Utils

Go library and CLI tools for interacting with S3-compatible object storage (MinIO, AWS S3, etc.).

## Installation

```bash
go get github.com/bborbe/s3-utils
```

## Library Usage

```go
import s3utils "github.com/bborbe/s3-utils"

client := s3utils.CreateS3Client(
    s3utils.URL("http://localhost:9000"),
    s3utils.AccessKey("myaccesskey"),
    s3utils.SecretKey("mysecretkey"),
)
```

## CLI Commands

All commands share these common environment variables / flags:

| Flag | Env | Description |
|------|-----|-------------|
| `--s3-url` | `S3_URL` | URL of the S3 server |
| `--s3-access-key` | `S3_ACCESS_KEY` | Access key |
| `--s3-secret-key` | `S3_SECRET_KEY` | Secret key |
| `--s3-region` | `S3_REGION` | Region for SigV4 signing (empty = SDK default) |

### download

Downloads an object from S3 to stdout.

```bash
S3_URL=http://localhost:9000 \
S3_ACCESS_KEY=access \
S3_SECRET_KEY=secret \
BUCKET=mybucket \
KEY=myfile.txt \
./download > myfile.txt
```

### upload

Uploads stdin to an S3 object.

```bash
cat myfile.txt | \
S3_URL=http://localhost:9000 \
S3_ACCESS_KEY=access \
S3_SECRET_KEY=secret \
BUCKET=mybucket \
KEY=myfile.txt \
./upload
```

### list-buckets

Lists all buckets.

```bash
S3_URL=http://localhost:9000 \
S3_ACCESS_KEY=access \
S3_SECRET_KEY=secret \
./list-buckets
```

### list-objects

Lists all objects in a bucket.

```bash
S3_URL=http://localhost:9000 \
S3_ACCESS_KEY=access \
S3_SECRET_KEY=secret \
BUCKET=mybucket \
./list-objects
```

## Development

```bash
make test        # run tests
make precommit   # format, vet, lint, test
```

## License

BSD-2-Clause — see [LICENSE](LICENSE)
