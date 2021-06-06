package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xuexiangyou/thor/gateway/rest"
	"github.com/xuexiangyou/thor/gateway/rest/inter/context"
	"github.com/xuexiangyou/thor/gateway/rest/httpx"
)

type Steamer struct {
	source chan interface{}
}

var key string = "name"

type User struct {
	Id string `path:"id,default=4"`
	Name string `path:"name"`
}

func main() {
	var m map[string][]string
	m = map[string][]string{
		"a":[]string{"11"},
	}


	for key := range m {
		fmt.Println(key)
	}
	return

	var v = "abc"
	var a = "abc"

	fmt.Println(strings.HasPrefix(a, v))

	var user User

	err := httpx.ParseDemo(&user)
	fmt.Println(err)
	fmt.Println(user)

	return
	// http 实现route规则
	server := rest.MustNewServer()
	RegisterHandlers(server)
	server.Start()

	// fileds := strings.FieldsFunc("", func(r rune) bool {
	// 	return r == '.'
	// })
	// fmt.Println(fileds)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// // 附加值
	// valueCtx := context.WithValue(ctx, key, "【监控1】")
	// fmt.Println(valueCtx.Value(key))
	//
	// valueCtx = context.WithValue(valueCtx, key, "【监控1】233")
	//
	//
	// fmt.Println(valueCtx.Value(key))

	// Add(1,3,5)

	// ch := make(chan interface{})
	// p := From(ch)
	// p.FearchEcho()
	// p.source <- 2
}

func RegisterHandlers(engine *rest.Server) {
	engine.AddRoutes(
		[]rest.Route{
			{
				Method: http.MethodGet,
				Path: "/get",
				Handler: GetHandler(),
			},
			{
				Method: http.MethodGet,
				Path: "/list",
				Handler: ListHandler(),
			},
			{
				Method: http.MethodGet,
				Path: "/get/:id",
				Handler: GetIdHandler(),
			},
		},
	)
}

func GetIdHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// v := rest.Vars(request)
		v := context.Vars(request)
		id, ok := v["id"]
		if ok {
			httpx.OkJson(writer, id)
		} else {
			httpx.OkJson(writer, "不存在")
		}
	}
}

func GetHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		httpx.OkJson(writer, "nihao")
	}
}

func ListHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		httpx.OkJson(writer, "list")
	}
}

func Add(i ...int) {
	var v []int
	v = append(v, i...)
	fmt.Println(v)
}

func From(ch chan interface{}) Steamer {
	source := make(chan interface{}, 1000)
	go func() {
		defer close(source)
		for i := 0; i <= 100; i++ {
			source <- 1
			time.Sleep(10 * time.Millisecond)
			fmt.Println("111")
		}
	}()
	return Steamer{
		source: source,
	}
}

func (p Steamer) FearchEcho() {
	time.Sleep(1 * time.Second)
	for value := range p.source {
		fmt.Println(value)
	}
	// for  {
	// 	select {
	// 		case j := <- p.source :
	// 		fmt.Println(j)
	// 	}
	// 	// p.source <- 10
	// }
}
