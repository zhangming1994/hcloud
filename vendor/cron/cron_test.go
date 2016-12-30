package cron

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

type jobs struct {
	id  int
	bar int
	wg  *sync.WaitGroup
}

func NewJobs(id int, wg *sync.WaitGroup) *jobs {
	return &jobs{
		id: id,
		wg: wg,
	}
}

func (t *jobs) Run() {
	for i := 1; i <= 100; i++ {
		t.bar = i
	}
	fmt.Printf("The bar was: %d\n", t.bar)
	t.wg.Done()
}

// Test running the same job .
func TestRunningJob(t *testing.T) {
	fmt.Println("TestRunningJob....")

	wg := &sync.WaitGroup{}
	wg.Add(10)

	cron := New(1024)

	cron.Start()
	defer cron.Stop()

	go func() {
		for i := 0; i < 10; i++ {
			cmd := NewJobs(i, wg)
			cron.AddJob(cmd)
			fmt.Printf("add job %d\n", i)
		}
		fmt.Println("Entries:", cron.Entries())
		fmt.Println("add job Done")
	}()
	time.Sleep(1 * time.Second)
	select {
	case <-time.After(100 * time.Second):
		t.FailNow()
	case <-wait(wg):
	}
}
