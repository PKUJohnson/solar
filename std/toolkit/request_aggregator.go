package toolkit

import (
	std "github.com/PKUJohnson/solar/std"
	"runtime"
	"sync"
	"time"
)

type RequestAggregator struct {
	eventQueue chan interface{}
	batchSize  int
	workers    int

	quit chan struct{}
	wg   *sync.WaitGroup
	sendRequestFunc SendRequest
}

type SendRequest func([]interface{}) error

func NewRequestAggregator(batchSize, workers int, sendRequestFunc SendRequest) *RequestAggregator {
	return &RequestAggregator{
		eventQueue: make(chan interface{}, batchSize),
		batchSize:  batchSize,
		workers:    workers,
		quit:       make(chan struct{}),
		wg:         new(sync.WaitGroup),
		sendRequestFunc: sendRequestFunc,
	}
}

func (sa *RequestAggregator) Enqueue(req interface{} ) {
	select {
	case sa.eventQueue <- req:
	default:
		std.LogWarnLn("RequestAggregator is full and try call GoSched")
		runtime.Gosched()
		select {
		case sa.eventQueue <- req:
		default:
			std.LogWarnLn("RequestAggregator is still full then try use go routine")
			go func() {
				sa.eventQueue <- req
			}()
		}
	}
}

func (sa *RequestAggregator) Start() {
	for i:=0;i <sa.workers; i++{
		go sa.work()
	}
}

func (sa *RequestAggregator) Stop() {
	close(sa.quit)
	sa.wg.Wait()
}

func (sa *RequestAggregator) work() {
	sa.wg.Add(1)
	defer sa.wg.Done()

	reqs := make([]interface{}, 0, sa.batchSize)
	idleDelay := time.NewTimer(1 * time.Minute)
	defer idleDelay.Stop()

	loop:
	for {
		idleDelay.Reset(1 * time.Minute)
		select {
		case req := <-sa.eventQueue:
			reqs = append(reqs, req)
			if len(reqs) != sa.batchSize {
				break
			}

			sa.sendRequest(reqs)
			reqs = make([]interface{}, 0, sa.batchSize)
		case <-idleDelay.C:
			if len(reqs) == 0 {
				break
			}

			sa.sendRequest(reqs)
			reqs = make([]interface{}, 0, sa.batchSize)
		case <- sa.quit:
			break loop
		}
	}
}

func (sa *RequestAggregator) sendRequest(events []interface{}) {
	if err := sa.sendRequestFunc(events); err != nil {
		std.LogErrorc("RequestAggregator", err, "Fail to send requests")
	} else {
		std.LogDebugLn("RequestAggregator has sent ", len(events), " events")
	}
}
