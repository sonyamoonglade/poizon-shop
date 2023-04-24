module usecase

go 1.20

require repositories v0.0.0

require (
	domain v0.0.0 // indirect
	dto v0.0.0 // indirect
	functools v0.0.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sonyamoonglade/go_func v0.0.0-20230418180836-d7b9b025b11a // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.11.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	household_bot v0.0.0 // indirect
	logger v0.0.0 // indirect
	onlineshop/database v0.0.0 // indirect
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
	services => ../services
	utils => ../utils
)
