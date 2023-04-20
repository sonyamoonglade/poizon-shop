module repositories

go 1.18

require (
	go.mongodb.org/mongo-driver v1.11.4
	go.uber.org/zap v1.24.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.0.3 // indirect
)

require (
	domain v0.0.0
	dto v0.0.0
	functools v0.0.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	household_bot v0.0.0
	logger v0.0.0
	onlineshop/database v0.0.0
	redis v0.0.0
)

replace (
	domain v0.0.0 => ../domain
	dto v0.0.0 => ../dto
	functools v0.0.0 => ../functools
	household_bot v0.0.0 => ../../apps/household_bot
	logger v0.0.0 => ../logger
	nanoid v0.0.0 => ../nanoid
	onlineshop/database v0.0.0 => ../database
	redis v0.0.0 => ../redis
	services v0.0.0 => ../services
	utils => ../utils
)
