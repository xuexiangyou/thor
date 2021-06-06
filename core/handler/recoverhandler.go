package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if result := recover(); result != nil {
				fmt.Println(debug.Stack())
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		// fmt.Println("recover middleware before")
		next.ServeHTTP(writer, request)
		// fmt.Println("recover middleware last")
	})
}
