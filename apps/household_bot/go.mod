module household_bot

go 1.19

require (
	github.com/brianvoe/gofakeit/v6 v6.20.2
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.11.4
	go.uber.org/zap v1.24.0
	onlineshop/database v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elliotchance/pie/v2 v2.5.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/redis/go-redis/v9 v9.0.3 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nanoid v0.0.0 // indirect
)

require (
	domain v0.0.0
	dto v0.0.0
	functools v0.0.0
	github.com/klauspost/compress v1.16.3 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.11.0
	golang.org/x/text v0.8.0 // indirect
	logger v0.0.0
	redis v0.0.0
	repositories v0.0.0
	services v0.0.0
)

replace (
	domain v0.0.0 => ../../libs/domain
	dto v0.0.0 => ../../libs/dto
	functools v0.0.0 => ../../libs/functools
	logger v0.0.0 => ../../libs/logger
	nanoid => ../../libs/nanoid
	onlineshop/database v0.0.0 => ../../libs/database
	redis => ../../libs/redis
	repositories v0.0.0 => ../../libs/repositories
	services => ../../libs/services
	utils => ../../libs/utils
)
