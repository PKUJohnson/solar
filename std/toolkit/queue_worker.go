package toolkit

import (
	"sync"
)

type QueueWorker struct {
	workers  int
	quit     chan struct{}
	wg       *sync.WaitGroup
	jobsChan chan interface{}
	onWork   WorkAction
}

type WorkAction func(int)

func NewQueueWorker(workers int, onWork WorkAction) *QueueWorker {
	return &QueueWorker{
		workers:  workers,
		quit:     make(chan struct{}),
		wg:       new(sync.WaitGroup),
		jobsChan: make(chan interface{}, 4*workers),
		onWork:   onWork,
	}
}

func (sl *QueueWorker) Start() {
	for i := 0; i < sl.workers; i++ {
		go func(index int) {
			sl.Work(index)
		}(i)
	}
}

func (sl *QueueWorker) Stop() {
	close(sl.quit)
	sl.wg.Wait()
}

func (sl *QueueWorker) Enqueue(job interface{}) {
	sl.jobsChan <- job
}

func (sl *QueueWorker) Work(index int) {
	defer func() {
		if r := recover(); r != nil {
			// todo: need log here
			sl.Work(index)
		}
	}()

	sl.wg.Add(1)
	defer sl.wg.Done()

	sl.onWork(index)
}
