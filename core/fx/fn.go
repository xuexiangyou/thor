package fx

import (
	"sync"

	"github.com/xuexiangyou/thor/core/lang"
	"github.com/xuexiangyou/thor/core/threading"
)

const(
	defaultWorkers = 16
	minWorks = 1
)

type PlaceholderType struct {}

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers 	int
	}

	FilterFunc func(item interface{}) bool
	ForAllFunc func(pipe <-chan interface{})
	ForEachFunc func(item interface{})
	GenerateFunc func(source chan<- interface{})
	KeyFunc func(item interface{}) interface{}
	LessFunc func(a, b interface{}) bool
	MapFunc func(item interface{}) interface{}
	Option func(opts *rxOptions)
	ParalleFunc func(item interface{})
	ReduceFunc func(pipe <-chan interface{}) (interface{}, error)
	WalkFunc     func(item interface{}, pipe chan<- interface{})

	Stream struct {
		source <-chan interface{}
	}
)

func From(generate GenerateFunc) Stream {
	source := make(chan interface{})

	go func() {
		defer close(source)
		generate(source)
	} ()

	return Range(source)
}

func (p Stream) ForEach(fn ForEachFunc) {
	for item := range p.source {
		fn(item)
	}
}

func (p Stream) Map(fn MapFunc, opts ...Option) Stream {
	return p.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opts...)
}

func (p Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return p.walkUnLimited(fn, option)
	} else {
		return p.walkLimited(fn, option)
	}
}

func (p Stream) walkUnLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, defaultWorkers)

	go func() {
		var wg sync.WaitGroup

		for {
			item, ok := <- p.source // 判断channel是否关闭了
			if !ok {
				break
			}
			wg.Add(1)
			threading.GoSafe(func() {
				defer wg.Done()
				fn(item, pipe)
			})
		}
	}()
	return Range(pipe)
}

func (p Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan lang.PlaceholderType, option.workers)
		for {
			pool <- lang.Placeholder
			item, ok := <-p.source // todo 关闭channel时ok为false
			if !ok {
				<-pool
				break
			}
			wg.Add(1)
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				fn(item, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

func Range(source <-chan interface{}) Stream {
	return Stream{
		source: source,
	}
}

func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}