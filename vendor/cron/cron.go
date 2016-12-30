// This library implements a cron spec parser and runner.  See the README for
// more details.
package cron

import (
	"container/list"
	"log"
	"runtime"
	"time"
)

// Cron keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Cron struct {
	entries  *list.List
	stop     chan struct{}
	add      chan *Entry
	snapshot chan []*Entry
	running  bool
	ErrorLog *log.Logger
}

// Job is an interface for submitted cron jobs.
type Job interface {
	Run()
}

// Entry consists of a schedule and the func to execute on that schedule.
type Entry struct {
	// The Job to run.
	Job Job
}

// New returns a new Cron job runner.
func New(size int) *Cron {
	return &Cron{
		entries:  list.New(),
		add:      make(chan *Entry, size),
		stop:     make(chan struct{}),
		snapshot: make(chan []*Entry),
		running:  false,
		ErrorLog: nil,
	}

}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (c *Cron) AddJob(cmd Job) error {
	c.Schedule(cmd)
	return nil
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *Cron) Schedule(cmd Job) {
	entry := &Entry{
		Job: cmd,
	}
	if !c.running {
		c.entries.PushBack(entry)
		return
	}

	c.add <- entry
}

// Entries returns a snapshot of the cron entries.
func (c *Cron) Entries() []*Entry {
	if c.running {
		c.snapshot <- nil
		x := <-c.snapshot
		return x
	}
	return c.entrySnapshot()
}

// Start the cron scheduler in its own go-routine.
func (c *Cron) Start() {
	c.running = true
	go c.run()
}

func (c *Cron) runWithRecovery(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.logf("cron: panic running job: %v\n%s", r, buf)
		}
	}()
	j.Run()
}

// Run the scheduler.. this is private just due to the need to synchronize
// access to the 'running' state variable.
func (c *Cron) run() {
	for {
		select {
		case newEntry := <-c.add:
			c.entries.PushBack(newEntry)
		case <-c.snapshot:
			c.snapshot <- c.entrySnapshot()
		case <-c.stop:
			return
		default:
			if e := c.entries.Front(); e != nil {
				value := c.entries.Remove(e)
				entry, ok := value.(*Entry)
				if ok {
					c.runWithRecovery(entry.Job)
				}
			}
		}
		time.Sleep(0)
	}
}

// Logs an error to stderr or to the configured error log
func (c *Cron) logf(format string, args ...interface{}) {
	if c.ErrorLog != nil {
		c.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

// entrySnapshot returns a copy of the current cron entry list.
func (c *Cron) entrySnapshot() []*Entry {
	entries := []*Entry{}
	for e := c.entries.Front(); e != nil; e = e.Next() {
		if entry, ok := e.Value.(*Entry); ok {
			entries = append(entries, &Entry{
				Job: entry.Job,
			})
		}
	}
	return entries
}
