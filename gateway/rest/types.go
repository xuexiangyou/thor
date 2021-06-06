package rest

import "net/http"

type (
	Route struct {
		Method string
		Path   string
		Handler http.HandlerFunc
	}

	RouteOption func(r *featureRoutes)

	featureRoutes struct {
		routes []Route
	}
)
