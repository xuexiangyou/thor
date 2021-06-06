package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)
	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			// logx.Errorf("write response failed, error: %s", err)
			fmt.Printf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		// logx.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
		fmt.Printf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}