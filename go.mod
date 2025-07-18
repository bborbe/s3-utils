module github.com/bborbe/s3-utils

go 1.24.4

exclude (
	github.com/go-logr/glogr v1.0.0-rc1
	github.com/go-logr/glogr v1.0.0
	github.com/go-logr/logr v1.0.0-rc1
	github.com/go-logr/logr v1.0.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10
	github.com/aws/aws-sdk-go-v2/config => github.com/aws/aws-sdk-go-v2/config v1.18.39
	github.com/aws/aws-sdk-go-v2/credentials => github.com/aws/aws-sdk-go-v2/credentials v1.13.37
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11
	github.com/aws/aws-sdk-go-v2/feature/s3/manager => github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.50
	github.com/aws/aws-sdk-go-v2/internal/configsources => github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35
	github.com/aws/aws-sdk-go-v2/internal/ini => github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42
	github.com/aws/aws-sdk-go-v2/internal/v4a => github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.28
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11
	github.com/aws/aws-sdk-go-v2/service/internal/checksum => github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.31
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.5
	github.com/aws/aws-sdk-go-v2/service/s3 => github.com/aws/aws-sdk-go-v2/service/s3 v1.37.1
	github.com/aws/aws-sdk-go-v2/service/sso => github.com/aws/aws-sdk-go-v2/service/sso v1.13.6
	github.com/aws/aws-sdk-go-v2/service/ssooidc => github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.6
	github.com/aws/aws-sdk-go-v2/service/sts => github.com/aws/aws-sdk-go-v2/service/sts v1.21.5
)

require (
	github.com/actgardner/gogen-avro/v9 v9.2.0
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/credentials v1.13.37
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/s3 v1.30.1
	github.com/aws/smithy-go v1.22.3
	github.com/bborbe/errors v1.3.0
	github.com/bborbe/sentry v1.7.1
	github.com/bborbe/service v1.6.2
	github.com/golang/glog v1.2.5
	github.com/google/addlicense v1.1.1
	github.com/incu6us/goimports-reviser/v3 v3.9.1
	github.com/kisielk/errcheck v1.9.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.11.2
	github.com/onsi/ginkgo/v2 v2.23.4
	github.com/onsi/gomega v1.37.0
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067
	golang.org/x/vuln v1.1.4
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.5 // indirect
	github.com/bborbe/argument/v2 v2.3.0 // indirect
	github.com/bborbe/collection v1.9.0 // indirect
	github.com/bborbe/math v1.2.0 // indirect
	github.com/bborbe/parse v1.7.0 // indirect
	github.com/bborbe/run v1.7.0 // indirect
	github.com/bborbe/time v1.15.1 // indirect
	github.com/bborbe/validation v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/getsentry/sentry-go v0.31.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250501235452-c0086092b71a // indirect
	github.com/incu6us/goimports-reviser v0.1.6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.63.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/telemetry v0.0.0-20250506010939-b15a553ce495 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
