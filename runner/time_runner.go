package runner

import (
	"fmt"
	"os"
	"time"
)

type TimeRunner struct {
	sign   <-chan os.Signal
	period time.Duration
	f      RunFunc
	doneF  RunFunc
}

func CreateTimeRunner(sign <-chan os.Signal, period time.Duration, f RunFunc, doneF RunFunc) *TimeRunner {
	return &TimeRunner{
		sign:   sign,
		period: period,
		f:      f,
		doneF:  doneF,
	}
}

func (r *TimeRunner) Start() {
	ticker := time.NewTicker(r.period)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.f()
			break
		case sg := <-r.sign:
			fmt.Printf("accept signal %s\n", sg)
			r.doneF()
			return
		}
	}
}
