package handler

import (
	"fmt"
	"net/http"

	"github.com/xuexiangyou/thor/core/trace"
)

func TracingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		carrier, err := trace.Extract(trace.HttpFormat, request.Header)
		if err != nil && err != trace.ErrInvalidCarrier {
			fmt.Println(err)
		}
		ctx, span := trace.StartServerSpan(request.Context(), carrier, "hostName", request.RequestURI)
		defer span.Finish()

		// fmt.Println("Trace middleware before")
		request = request.WithContext(ctx)
		next.ServeHTTP(writer, request)
		// fmt.Println("Trace middleware last")
	})
}
