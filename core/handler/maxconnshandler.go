package handler

import (
	"fmt"
	"net/http"

	"github.com/xuexiangyou/thor/core/syncx"
)

func MaxConns(n int) func(next http.Handler) http.Handler {
	if n < 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	return func(next http.Handler) http.Handler {
		latch := syncx.NewLimit(n)
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if latch.TryBorrow() {
				defer func() {
					if err := latch.Return(); err != nil {
						// logx.Error(err)
						fmt.Println(err)
					}
				}()
				// fmt.Println("max middleware before")
				next.ServeHTTP(writer, request)
				// fmt.Println("max middleware last")
			} else {
				writer.WriteHeader(http.StatusServiceUnavailable)
			}
		})
	}
}
