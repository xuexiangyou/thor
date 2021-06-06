package rest

import (
	"time"

	"github.com/justinas/alice"
	"github.com/xuexiangyou/thor/core/handler"
	"github.com/xuexiangyou/thor/gateway/rest/httpx"
	"github.com/xuexiangyou/thor/gateway/rest/inter"
	"github.com/xuexiangyou/thor/gateway/rest/router"
)

type engine struct {
	routes []featureRoutes
}

func newEngine() *engine {
	srv := &engine{}
	return srv
}

func (s *engine) AddRoutes(r featureRoutes) {
	s.routes = append(s.routes, r)
}

func (s *engine) Start() error {
	return s.StartWithRouter(router.NewRouter())
}

func (s *engine) StartWithRouter(router httpx.Router) error {
	if err := s.bindRoutes(router); err != nil {
		return err
	}

	return inter.StartHttp("localhost", 8000, router)
}

func (s *engine) bindRoutes(router httpx.Router) error {
	for _, fr := range s.routes {
		if err := s.bindFeaturedRoutes(router, fr); err != nil {
			return err
		}
	}

	return nil
}

func (s *engine) bindFeaturedRoutes(router httpx.Router, fr featureRoutes) error {
	for _, route := range fr.routes {
		if err := s.bindRoute(fr, router, route); err != nil {
			return err
		}
	}
	return nil
}

func (s *engine) bindRoute(fr featureRoutes, router httpx.Router, route Route) error {
	chian := alice.New(
		handler.TracingHandler,
		handler.MaxConns(100),
		handler.TimeoutHandler(1 * time.Second),
		handler.RecoverHandler,
	)
	handle := chian.ThenFunc(route.Handler)
	return router.Handle(route.Method, route.Path, handle)
}
