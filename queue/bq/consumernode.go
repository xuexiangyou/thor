package bq

import (
	"log"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/xuexiangyou/thor/core/syncx"
)

type (
	consumerNode struct {
		conn *connection
		tube string
		on *syncx.AtomicBool
	}

	consumerService struct {
		c 	*consumerNode
		consume Consume
	}
)

func newConsumerNode(endpoint, tube string) *consumerNode {
	return &consumerNode{
		conn: newConnection(endpoint, tube),
		tube: tube,
		on: syncx.ForAtomicBool(true),
	}
}

func (c *consumerNode) dispose() {
	c.on.Set(false)
}

func (c *consumerNode) consumeEvents(consume Consume) {
	for c.on.True() {
		conn, err := c.conn.get()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		if !c.on.True() {
			break
		}

		conn.Tube.Name = c.tube
		conn.TubeSet.Name[c.tube] = true
		id, body, err := conn.Reserve(reserveTimeout)
		if err == nil {
			conn.Delete(id)
			consume(body)
			continue
		}

		// the error can only be beanstalk.NameError or beanstalk.ConnError
		switch cerr := err.(type) {
		case beanstalk.ConnError:
			switch cerr.Err {
			case beanstalk.ErrTimeout:
				// timeout error on timeout, just continue the loop
			case beanstalk.ErrBadChar, beanstalk.ErrBadFormat, beanstalk.ErrBuried, beanstalk.ErrDeadline,
				beanstalk.ErrDraining, beanstalk.ErrEmpty, beanstalk.ErrInternal, beanstalk.ErrJobTooBig,
				beanstalk.ErrNoCRLF, beanstalk.ErrNotFound, beanstalk.ErrNotIgnored, beanstalk.ErrTooLong:
				// won't reset
				// logx.Error(err)
				log.Println(err)
			default:
				// beanstalk.ErrOOM, beanstalk.ErrUnknown and other errors
				log.Println(err)
				c.conn.reset()
				time.Sleep(time.Second)
			}
		default:
			log.Println(err)
		}
	}

	if err := c.conn.Close(); err != nil {
		// logx.Error(err)
		log.Println(err)
	}
}

func (cs consumerService) Start() {
	cs.c.consumeEvents(cs.consume)
}

func (cs consumerService) Stop() {
	cs.c.dispose()
}