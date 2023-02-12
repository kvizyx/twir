module github.com/satont/tsuwari/apps/emotes-cacher

go 1.19

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

replace github.com/satont/tsuwari/libs/config => ../../libs/config

require (
	github.com/getsentry/sentry-go v0.18.0
	github.com/redis/go-redis/v9 v9.0.2
	github.com/samber/do v1.6.0
	github.com/samber/lo v1.37.0
	github.com/satont/tsuwari/apps/parser v0.0.0-20230201225635-782be9343513
	github.com/satont/tsuwari/libs/config v0.0.0
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.7
	gorm.io/gorm v1.24.5
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20221028150844-83b7d23a625f // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
)