module github.com/xuexiangyou/thor

go 1.14

require (
	github.com/beanstalkd/go-beanstalk v0.1.0
	github.com/go-redis/redis/v8 v8.8.2
	github.com/justinas/alice v1.2.0
	github.com/segmentio/kafka-go v0.4.14 // indirect
	github.com/spaolacci/murmur3 v1.1.0
	go.etcd.io/etcd v0.0.0-20200824191128-ae9734ed278b
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20210415220803-1a870aec2ff9 // indirect
	google.golang.org/protobuf v1.26.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
