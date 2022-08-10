package pool

import (
	"fmt"
	"time"
)

type Pool struct {
	items chan interface{}
}

func New(size int) *Pool {
	p := Pool{
		items: make(chan interface{}, size),
	}

	return &p
}

func (p *Pool) Put(item interface{}) {
	p.items <- item
}

func (p *Pool) PutTimeout(item interface{}, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case p.items <- item:
		return nil
	case <-timer.C:
		return fmt.Errorf("timed out waiting for pool not to be full after %v", timeout.String())
	}
}

func (p *Pool) Get() interface{} {
	return <-p.items
}

func (p *Pool) GetTimeout(timeout time.Duration) (interface{}, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case item := <-p.items:
		return item, nil
	case <-timer.C:
		return nil, fmt.Errorf("timed out waiting for pool not to be empty after %v", timeout.String())
	}
}
