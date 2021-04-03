package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xuexiangyou/thor/socket/server"
)

func main() {
	var svr *server.Server

	svr, err := server.NewServer("127.0.0.1:8089")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
		syscall.SIGUSR1,
	)
	go func() {
		for {
			sig := <- sc
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				fmt.Println("程序停止")
				svr.Close()
			}
		}
	}()

	svr.Run()
}
