package threading

import "log"

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer func() {
		if p := recover(); p != nil {
			log.Println(p)
		}
	}()
	fn()
}