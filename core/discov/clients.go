package discov

import (
	"fmt"

	"github.com/xuexiangyou/thor/core/discov/internal"
)

const (
	_ = iota
	indexOfId
)

const timeToLive int64 = 10

var TimeToLive = timeToLive

func makeEtcdKey(key string, id int64) string {
	return fmt.Sprintf("%s%c%d", key, internal.Delimiter, id)
}


