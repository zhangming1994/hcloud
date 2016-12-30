package utils

import (
	_list "container/list"
	"errors"
	"github.com/astaxie/beego"
	"sync"
)

type Queue struct {
	list    *_list.List
	lock    *sync.RWMutex
	sem     chan int
	running bool
}

func QueueNew() (*Queue, error) {
	return &Queue{list: _list.New(), lock: new(sync.RWMutex), sem: make(chan int, 2048)}, nil
}

func (this *Queue) Size() int {
	return this.list.Len()
}

func (this *Queue) Put(val interface{}) error {
	this.lock.Lock()
	e := this.list.PushFront(val)
	this.lock.Unlock()
	this.sem <- 1
	if e == nil {
		return errors.New("PushFront failed")
	}
	return nil
}

func (this *Queue) Get() (*_list.Element, error) {
	this.lock.Lock()
	e := this.list.Back()
	if e != nil {
		this.list.Remove(e)
	}
	this.lock.Unlock()
	if e == nil {
		return nil, errors.New("Back failed")
	}
	return e, nil
}

func (this *Queue) Poll() bool {
	state := <-this.sem
	switch state {
	case 0:
		return false
	case 1:
	default:
	}
	if this.Size() > 0 {
		return true
	}
	return false
}

func (this *Queue) Stop() {
	this.running = false
	this.sem <- 0
	beego.Debug("Queue Stop")
}

func (this *Queue) Start(f func(val interface{})) bool {
	this.running = true
	go func() {
		for {
			if this.running {
				if this.Poll() {
					//go f(this)
					n, err := this.Get()
					if err != nil {
						beego.Warn("Element is nil")
					} else {
						f(n.Value)
					}

				}
			} else {
				break
			}
		}
		beego.Debug("Queue Thread End")
	}()
	beego.Debug("Queue Start")
	return true
}
