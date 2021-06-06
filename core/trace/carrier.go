package trace

import (
	"errors"
	"net/http"
)

var ErrInvalidCarrier = errors.New("invalid carrier")

type (
	Carrier interface {
		Get(key string) string
		Set(key, value string)
	}

	httpCarrier http.Header

	grpcCarrier map[string][]string
)

func (h httpCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}

func (h httpCarrier) Set(key string, val string) {
	http.Header(h).Set(key, val)
}


