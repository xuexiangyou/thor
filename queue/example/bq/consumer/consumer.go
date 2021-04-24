package main

import (
	"fmt"

	"github.com/xuexiangyou/thor/queue/bq"
)

func main() {
	consumer := bq.NewConsumer(bq.BqConf{
		Beanstalk: []bq.Beanstalk{
			{
				Endpoint: "localhost:11300",
				Tube: "tube",
			},
			{
				Endpoint: "localhost:11300",
				Tube: "tube",
			},
		},
		Redis: bq.Redis{
			Host: "127.0.0.1:6379",
		},
	})
	consumer.Consume(func(body []byte) {
		fmt.Println(string(body))
	})
}
