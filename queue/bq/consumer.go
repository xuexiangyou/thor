package bq

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xuexiangyou/thor/core/hash"
	"github.com/xuexiangyou/thor/core/service"
)

const (
	expiration = 3600 // seconds
	tolerance  = time.Minute * 30
	guardValue = "1"
)

var maxCheckBytes = getMaxTimeLen()

type (
	Consume func(body []byte)

	Consumer interface {
		Consume(consume Consume)
	}

	consumerCluster struct {
		nodes []*consumerNode
		rdb   *redis.Client
	}
)

func NewConsumer(c BqConf) Consumer {
	var nodes []*consumerNode
	for _, node := range c.Beanstalk {
		nodes = append(nodes, newConsumerNode(node.Endpoint, node.Tube))
	}
	return &consumerCluster{
		nodes: nodes,
		rdb: redis.NewClient(&redis.Options{
			Addr: c.Redis.Host,
		}),
	}
}

func (c *consumerCluster) Consume(consume Consume) {
	guardedConsume := func(body []byte) {
		key := hash.Md5Hex(body)
		body, ok := c.unwrap(body)
		if !ok {
			return
		}
		err := c.rdb.SetNX(context.Background(), key, guardValue, time.Duration(expiration)*time.Second).Err()
		if err != nil {
			log.Println(err)
		} else {
			consume(body)
		}
	}

	group := service.NewServiceGroup()
	for _, node := range c.nodes {
		group.Add(consumerService{
			c:       node,
			consume: guardedConsume,
		})
	}

	group.Start()
}

func (c *consumerCluster) unwrap(body []byte) ([]byte, bool) {
	var pos = -1
	for i := 0; i < maxCheckBytes && i < len(body); i++ {
		if body[i] == timeSep {
			pos = i
			break
		}
	}

	if pos < 0 {
		return nil, false
	}

	val, err := strconv.ParseInt(string(body[:pos]), 10, 64)
	if err != nil {
		return nil, false
	}

	t := time.Unix(0, val)
	if t.Add(tolerance).Before(time.Now()) {
		return nil, false
	}
	return body[pos+1:], true
}

func getMaxTimeLen() int {
	return len(strconv.FormatInt(time.Now().UnixNano(), 10)) + 2
}
