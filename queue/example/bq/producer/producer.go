package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuexiangyou/thor/queue/bq"
)

func main() {
	producer := bq.NewProducer([]bq.Beanstalk{
		{
			Endpoint: "localhost:11300",
			Tube: "tube",
		},
		{
			Endpoint: "localhost:11300",
			Tube: "tube",
		},
	})

	for i := 1006; i < 1007; i++ {
		fmt.Println(i)
		_, err := producer.Delay([]byte(strconv.Itoa(i)), time.Second*5)
		if err != nil {
			fmt.Println(err)
		}
	}
}


