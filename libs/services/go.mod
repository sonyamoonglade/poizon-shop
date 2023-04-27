module services

go 1.20

require (
	domain v0.0.0
	dto v0.0.0
	functools v0.0.0
	github.com/sonyamoonglade/go_func v0.0.0-20230418180836-d7b9b025b11a
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.11.4
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.24.0
	logger v0.0.0
	onlineshop/database v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	household_bot v0.0.0
	nanoid v0.0.0
	repositories v0.0.0
)

replace (
	domain => ../domain
	dto => ../dto
	functools => ../functools
	household_bot => ../../apps/household_bot
	logger => ../logger
	nanoid => ../nanoid
	onlineshop/database => ../database
	redis => ../redis
	repositories => ../repositories
	usecase => ../usecase
	utils => ../utils
)
