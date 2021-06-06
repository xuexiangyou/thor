package rest

import "net/http"

type (
	runOptions struct {
		start func(*engine) error
	}

	RunOption func(*Server)

	Server struct {
		ngin *engine
		opts runOptions
	}
)

func MustNewServer(opts ...RunOption) *Server {
	server := &Server{
		ngin: newEngine(),
		opts: runOptions{
			start: func(e *engine) error {
				return e.Start()
			},
		},
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

func (s *Server) AddRoutes(rs []Route, opts ...RouteOption) {
	r := featureRoutes{
		routes: rs,
	}
	for _, opt := range opts {
		opt(&r)
	}
	s.ngin.AddRoutes(r)
}

func (s *Server) Start() {
	hanlderError(s.opts.start(s.ngin))
}

func hanlderError(err error) {
	// ErrServerClosed means the server is closed manually
	if err == nil || err == http.ErrServerClosed {
		return
	}
	// logx.Error(err)
	panic(err)
}
