package service

import (
	"log"

	"github.com/xuexiangyou/thor/core/proc"
	"github.com/xuexiangyou/thor/core/syncx"
	"github.com/xuexiangyou/thor/core/threading"
)

type (
	Starter interface {
		Start()
	}

	Stopper interface {
		Stop()
	}

	Service interface {
		Starter
		Stopper
	}

	ServiceGroup struct {
		services []Service
		stopOnce func()
	}
)

func NewServiceGroup() *ServiceGroup {
	sg := new(ServiceGroup)
	sg.stopOnce = syncx.Once(sg.stopOnce)
	return sg
}

func (sg *ServiceGroup) Add(service Service) {
	sg.services = append(sg.services, service)
}

func (sg *ServiceGroup) Start() {
	proc.AddShutdownListener(func() {
		log.Println("Shutting down...")
		sg.stopOnce()
	})

	sg.doStart()
}

func (sg *ServiceGroup) Stop() {
	sg.stopOnce()
}

func (sg *ServiceGroup) doStart() {
	routineGroup := threading.NewRoutineGroup()

	for i := range sg.services {
		service := sg.services[i]
		routineGroup.RunSafe(func() {
			service.Start()
		})
	}

	routineGroup.Wait()
}

func (sg *ServiceGroup) doStop() {
	for _, service := range sg.services {
		service.Stop()
	}
}


