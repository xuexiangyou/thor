package discov

import (
	"fmt"

	"github.com/xuexiangyou/thor/core/discov/internal"
	"github.com/xuexiangyou/thor/core/lang"
	"github.com/xuexiangyou/thor/core/syncx"
	"github.com/xuexiangyou/thor/core/threading"
	"go.etcd.io/etcd/client/v3"
)

type (
	PublisherOption func(client *Publisher)

	Publisher struct {
		endpoints  []string
		key        string
		fullKey    string
		id         int64
		value      string
		lease      clientv3.LeaseID
		quit       *syncx.DoneChan
		pauseChan  chan lang.PlaceholderType
		resumeChan chan lang.PlaceholderType
	}
)

func NewPublisher(endpoints []string, key, value string, opts ...PublisherOption) *Publisher {
	publisher := &Publisher{
		endpoints: endpoints,
		key: key,
		value: value,
		quit: syncx.NewDoneChan(),
		pauseChan:  make(chan lang.PlaceholderType),
		resumeChan: make(chan lang.PlaceholderType),
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

// Pause pauses the renewing of key:value.
func (p *Publisher) Pause() {
	p.pauseChan <- lang.Placeholder
}

// Resume resumes the renewing of key:value.
func (p *Publisher) Resume() {
	p.resumeChan <- lang.Placeholder
}

// Stop stops the renewing and revokes the registration.
func (p *Publisher) Stop() {
	p.quit.Close()
}

func (p *Publisher) KeepAlive() error {
	cli, err := internal.GetRegistry().GetConn(p.endpoints)
	if err != nil {
		return err
	}

	p.lease, err = p.register(cli)
	if err != nil {
		return err
	}
	// todo 添加平滑退出的逻辑


	return p.keepAliveAsync(cli)
}

func (p *Publisher) keepAliveAsync(cli internal.EtcdClient) error {
	ch, err := cli.KeepAlive(cli.Ctx(), p.lease)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		for {
			select {
			case _, ok := <-ch :
				if !ok {
					p.revoke(cli)
					if err := p.KeepAlive(); err != nil {
						fmt.Println("KeepAlive:", err.Error())
					}
					return
				}
			case <-p.pauseChan:
				fmt.Sprintf("paused etcd renew, key: %s, value: %s", p.key, p.value)
				p.revoke(cli)
				select {
				case <-p.resumeChan:
					if err := p.KeepAlive(); err != nil {
						fmt.Println("KeepAlive:", err.Error())
					}
					return
				case <-p.quit.Done():
					return
				}
			case <-p.quit.Done():
				p.revoke(cli)
				return
			}
		}
	})

	return nil
}

func (p *Publisher) revoke(cli internal.EtcdClient) {
	if _, err := cli.Revoke(cli.Ctx(), p.lease); err != nil {
		fmt.Println(err)
	}
}

func (p *Publisher) register(client internal.EtcdClient) (clientv3.LeaseID, error) {
	resp, err := client.Grant(client.Ctx(), TimeToLive)
	if err != nil {
		return clientv3.NoLease, err
	}

	lease := resp.ID
	if p.id > 0 {
		p.fullKey = makeEtcdKey(p.key, p.id)
	} else {
		p.fullKey = makeEtcdKey(p.key, int64(lease))
	}

	_, err = client.Put(client.Ctx(), p.fullKey, p.value, clientv3.WithLease(lease))

	return lease, err
}

// WithId customizes a Publisher with the id.
func WithId(id int64) PublisherOption {
	return func(publisher *Publisher) {
		publisher.id = id
	}
}

