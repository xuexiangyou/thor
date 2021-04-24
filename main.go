package main

import (
	"fmt"
	"time"
)

type Steamer struct {
	source chan interface{}
}


func main() {
	ch := make(chan interface{})
	p := From(ch)
	p.FearchEcho()
	p.source <- 2
}

func From(ch chan interface{}) Steamer {
	source := make(chan interface{}, 1000)
	go func() {
		defer close(source)
		for i:= 0; i <= 100; i++{
			source <- 1
			time.Sleep(10 * time.Millisecond)
			fmt.Println("111")
		}
	}()
	return Steamer{
		source: source,
	}
}

func (p Steamer) FearchEcho () {
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
