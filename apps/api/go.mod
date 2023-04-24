module onlineshop/api

go 1.18

require (
	domain v0.0.0
	github.com/brianvoe/gofakeit/v6 v6.20.2
	github.com/gofiber/fiber/v2 v2.44.0
	github.com/sonyamoonglade/s3-yandex-go v1.0.7
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.11.4
	go.uber.org/zap v1.24.0
	logger v0.0.0
	onlineshop/database v0.0.0
	redis v0.0.0
	repositories v0.0.0
	utils v0.0.0
)

require (
	github.com/aws/aws-sdk-go v1.44.44 // indirect
	github.com/aws/aws-sdk-go-v2 v1.16.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.15.11 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.7 // indirect
	github.com/aws/smithy-go v1.11.3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/redis/go-redis/v9 v9.0.3 // indirect
	github.com/sonyamoonglade/go_func v0.0.0-20230418180836-d7b9b025b11a // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	household_bot v0.0.0 // indirect
	nanoid v0.0.0 // indirect
)

require (
	dto v0.0.0
	functools v0.0.0 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.45.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	services v0.0.0
)

replace (
	domain v0.0.0 => ../../libs/domain
	dto v0.0.0 => ../../libs/dto
	functools v0.0.0 => ../../libs/functools

	household_bot v0.0.0 => ../household_bot
	logger v0.0.0 => ../../libs/logger
	nanoid => ../../libs/nanoid
	onlineshop/database v0.0.0 => ../../libs/database
	redis v0.0.0 => ../../libs/redis
	repositories v0.0.0 => ../../libs/repositories
	services => ../../libs/services
	usecase => ../../libs/usecase
	utils => ../../libs/utils
)
