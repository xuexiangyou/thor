package executors

import (
	"reflect"
	"sync"
	"time"

	"github.com/xuexiangyou/thor/core/lang"
	"github.com/xuexiangyou/thor/core/proc"
	"github.com/xuexiangyou/thor/core/syncx"
	"github.com/xuexiangyou/thor/core/threading"
	"github.com/xuexiangyou/thor/core/timex"
)

const idleRound = 10

type (
	TaskContainer interface {
		AddTask(task interface{}) bool
		Execute(tasks interface{})
		RemoveAll() interface{}
	}

	PeriodicalExecutor struct {
		commander chan interface{}
		interval time.Duration
		container TaskContainer
		waitGroup sync.WaitGroup

		wgBarrier syncx.Barrier
		confirmChan chan lang.PlaceholderType
		guarded bool
		newTicker func(duration time.Duration) timex.Ticker
		lock sync.Mutex
	}
)

func NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *PeriodicalExecutor {
	executor := &PeriodicalExecutor{
		commander: make(chan interface{}, 1),
		interval: interval,
		container: container,
		confirmChan: make(chan lang.PlaceholderType),
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(interval)
		},
	}

	proc.AddShutdownListener(func() {
		executor.Flush()
	})

	return executor
}

func (pe *PeriodicalExecutor) Add(task interface{}) {
	if vals, ok := pe.addAndCheck(task); ok {
		pe.commander <- vals
		<-pe.confirmChan
	}
}

func (pe *PeriodicalExecutor) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() interface{} {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll()
	} ())
}

func (pe *PeriodicalExecutor) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn()
}

func (pe *PeriodicalExecutor) Wait() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait()
	})
}

func (pe *PeriodicalExecutor) addAndCheck(taks interface{}) (interface{}, bool) {
	pe.lock.Lock()
	defer func() {
		var start bool
		if !pe.guarded {
			pe.guarded = true
			start = true
		}
		pe.lock.Unlock()
		if start {
			pe.backgroundFlush()
		}
	}()
	if pe.container.AddTask(taks) {
		return pe.container.RemoveAll(), true
	}
	return nil, false
}

func (pe *PeriodicalExecutor) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1)
	})
}

func (pe *PeriodicalExecutor) doneExecution() {
	pe.waitGroup.Done()
}

func (pe *PeriodicalExecutor) executeTasks(tasks interface{}) bool {
	defer pe.doneExecution()
	ok := pe.hasTasks(tasks)
	if ok {
		pe.container.Execute(tasks)
	}
	return ok
}

func (pe *PeriodicalExecutor) backgroundFlush() {
	threading.GoSafe(func() {
		ticker := pe.newTicker(pe.interval)
		defer ticker.Stop()

		var commanded bool
		last := timex.Now()
		for {
			select {
			case vals := <-pe.commander:
				commanded = true
				pe.enterExecution()
				pe.confirmChan <- lang.Placeholder
				pe.executeTasks(vals)
				last = timex.Now()
			case <- ticker.Chan():
				if commanded {
					commanded = false
				} else if pe.Flush() {
					last = timex.Now()
				} else if timex.Since(last) > pe.interval * idleRound {
					pe.lock.Lock()
					pe.guarded = false
					pe.lock.Unlock()

					pe.Flush()
					return
				}
			}
		}
	})
}

func (pe *PeriodicalExecutor) hasTasks(tasks interface{}) bool {
	if tasks == nil {
		return false
	}

	val := reflect.ValueOf(tasks)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() > 0
	default:
		// unknown type, let caller execute it
		return true
	}
}