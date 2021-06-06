package inter

import (
	"fmt"
	"net/http"
)

func StartHttp(host string, port int, handler http.Handler) error {
	return start(host, port, handler, func(srv *http.Server) error {
		return srv.ListenAndServe()
	})
}

func start(host string, port int, handler http.Handler, run func(*http.Server) error) error {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}

	return run(server)
}