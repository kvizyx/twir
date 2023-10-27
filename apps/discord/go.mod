module github.com/satont/twir/apps/discord

go 1.21.0

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/sentry => ../../libs/sentry
	github.com/satont/twir/libs/twitch => ../../libs/twitch
)

require (
	github.com/avast/retry-go/v4 v4.5.0
	github.com/diamondburned/arikawa/v3 v3.3.3
	github.com/google/uuid v1.3.1
	github.com/nicklaw5/helix/v2 v2.25.1
	github.com/redis/go-redis/v9 v9.2.1
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-20231015185800-07291c1491d4
	github.com/satont/twir/libs/gomodels v0.0.0-20231015185112-b1ddd14cbc8f
	github.com/satont/twir/libs/grpc v0.0.0-20231015185112-b1ddd14cbc8f
	github.com/satont/twir/libs/logger v0.0.0-20231015185800-07291c1491d4
	github.com/satont/twir/libs/sentry v0.0.0-20231015185800-07291c1491d4
	github.com/satont/twir/libs/twitch v0.0.0-20231015185112-b1ddd14cbc8f
	go.uber.org/fx v1.20.1
	golang.org/x/sync v0.4.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/getsentry/sentry-go v0.25.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
)