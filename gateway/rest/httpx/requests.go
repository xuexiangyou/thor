package httpx

import (
	"net/http"

	"github.com/xuexiangyou/thor/core/mapping"
	"github.com/xuexiangyou/thor/gateway/rest/inter/context"
)

const (
	formKey           = "form"
	pathKey           = "path"
)

var (
	formUnmarshaler = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
)

// func Parse(r *http.Request, v interface{}) error {
// 	if err := ParsePath(r, v); err != nil {
//
// 	}
// }

func ParsePath(r *http.Request, v interface{}) error {
	vars := context.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func ParseDemo(v interface{}) error {
	m := make(map[string]interface{})
	m["id"] = "1"
	m["name"] = "lilonggen"
	return pathUnmarshaler.Unmarshal(m, v)
}
