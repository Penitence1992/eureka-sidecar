package runner

import "sync"

type RunFunc = func()

type Runner interface {
	Run()
}

type WaitGroupRunner struct {
	wg *sync.WaitGroup
	f  RunFunc
}

func NewRunner(wg *sync.WaitGroup, f RunFunc) *WaitGroupRunner {
	return &WaitGroupRunner{
		wg: wg,
		f:  f,
	}
}

func (r *WaitGroupRunner) Run() {
	go func() {
		r.wg.Add(1)
		defer r.wg.Done()
		r.f()
	}()
}
