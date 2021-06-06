package trace

import "net/http"

const (
	HttpFormat = iota
	GrpcFormat
)

var (
	emptyHttpPropagator httpPropagator
	emptyGrpcPropagator grpcPropagator
)

type (
	Propagator interface {
		Extract(carrier interface{}) (Carrier, error)
		Inject(carrier interface{}) (Carrier, error)
	}

	httpPropagator struct {}
	grpcPropagator struct {}
)

func (h httpPropagator) Extract(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); ok {
		return httpCarrier(c), nil
	}
	return nil, ErrInvalidCarrier
}

func (h httpPropagator) Inject(carrier interface{}) (Carrier, error) {
	if c, ok := carrier.(http.Header); ok {
		return httpCarrier(c), nil
	}
	return nil, ErrInvalidCarrier
}

func Extract(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Extract(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}

func Inject(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Inject(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}
