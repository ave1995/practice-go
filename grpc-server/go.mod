module github.com/ave1995/practice-go/grpc-server

go 1.25.3

require github.com/ave1995/practice-go/proto v0.0.0

replace github.com/ave1995/practice-go/proto => ../proto

require (
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10 // indirect
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.0
)

require (
	github.com/redis/go-redis/v9 v9.16.0
	github.com/segmentio/kafka-go v0.4.49
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
)