module github.com/haqury/user-service

go 1.24.0

toolchain go1.24.4

replace github.com/haqury/user-service => ./

require (
	github.com/gorilla/mux v1.8.1
	github.com/haqury/helpy v0.0.7
	github.com/joho/godotenv v1.5.1
	github.com/uptrace/bun/driver/pgdriver v1.2.5
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.46.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.4.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun v1.2.5 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	mellium.im/sasl v0.3.2 // indirect
)
