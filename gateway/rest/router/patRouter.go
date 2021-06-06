package router

import (
	"errors"
	"net/http"
	"path"

	"github.com/xuexiangyou/thor/core/search"
	"github.com/xuexiangyou/thor/gateway/rest/httpx"
	"github.com/xuexiangyou/thor/gateway/rest/inter/context"
)

const (
	allowHeader = "Allow"
	allowMethodSeparator = ", "
)

var (
	ErrInvalidMethod = errors.New("not a valid http method")
	ErrInvalidPath = errors.New("path must begin with '/'")
)

type patRouter struct {
	trees map[string]*search.Tree
	notFound http.Handler
	notAllowed http.Handler
}

func NewRouter() httpx.Router {
	return &patRouter{
		trees: make(map[string]*search.Tree),
	}
}

func (pr *patRouter) Handle(method, reqPath string, handler http.Handler) error {
	if !validMethod(method) {
		return ErrInvalidMethod
	}
	if len(reqPath) == 0 || reqPath[0] != '/' {
		return ErrInvalidPath
	}
	clearPath := path.Clean(reqPath)
	tree, ok := pr.trees[method]
	if ok {
		return tree.Add(clearPath, handler)
	}

	tree = search.NewTree()
	pr.trees[method] = tree
	return tree.Add(clearPath, handler)
}

func (pr *patRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := path.Clean(r.URL.Path)
	if tree, ok := pr.trees[r.Method]; ok {
		if result, ok := tree.Search(reqPath); ok {
			if len(result.Params) > 0 {
				r = context.WithPathVars(r, result.Params)
			}
			result.Item.(http.Handler).ServeHTTP(w, r)
		}
	}
}

func (pr *patRouter) SetNotFoundHandler(handler http.Handler) {
	pr.notFound = handler
}

func (pr *patRouter) SetNotAllowedHandler(handler http.Handler) {
	pr.notAllowed = handler
}

func validMethod(method string) bool {
	return method == http.MethodDelete || method == http.MethodGet ||
		method == http.MethodHead || method == http.MethodOptions ||
		method == http.MethodPatch || method == http.MethodPut ||
		method == http.MethodPost
}

