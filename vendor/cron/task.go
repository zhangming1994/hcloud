package cron

import (
	"sync"
	"sync/atomic"
)

type Task struct {
	Name    string
	status  uint32
	running sync.Mutex
	wg      *sync.WaitGroup
}

func NewTask(name string, wg *sync.WaitGroup) *Task {
	return &Task{
		Name: name,
		wg:   wg,
	}
}

func (t *Task) Status() string {
	if atomic.LoadUint32(&t.status) > 0 {
		return "RUNNING"
	}
	return "IDLE"
}

func (t *Task) Run() {
	t.running.Lock()
	defer t.running.Unlock()

	atomic.StoreUint32(&t.status, 1)
	defer atomic.StoreUint32(&t.status, 0)
	t.wg.Done()
}
