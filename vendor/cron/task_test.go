package cron

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test running the same job twice.
func TestRunningTask(t *testing.T) {
	fmt.Println("TestRunningTask....")
	wg := &sync.WaitGroup{}
	wg.Add(10)

	cron := New(1024)

	cron.Start()
	defer cron.Stop()

	for i := 0; i < 10; i++ {
		cron.AddJob(NewTask(fmt.Sprintf("jobs id=%d", i), wg))
	}

	select {
	case <-time.After(100 * time.Second):
		t.FailNow()
	case <-wait(wg):
	}
}
