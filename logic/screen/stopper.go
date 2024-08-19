// Copyright (c) 2024 by Marko Gaćeša

package screen

import "sync"

type Stopper struct {
	ch     chan error
	done   bool
	doneMx sync.Mutex
}

func NewStopper() *Stopper {
	return &Stopper{
		ch: make(chan error, 1),
	}
}

func (cc *Stopper) Done() <-chan error {
	return cc.ch
}

func (cc *Stopper) Stop() {
	cc.doneMx.Lock()
	defer cc.doneMx.Unlock()

	if cc.done {
		return
	}

	cc.done = true
	close(cc.ch)
}

func (cc *Stopper) Error(err error) {
	cc.doneMx.Lock()
	defer cc.doneMx.Unlock()

	if cc.done {
		return
	}

	cc.done = true
	cc.ch <- err
	close(cc.ch)
}
