package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

type Ewma struct {
	uncounted int64
	rate      float64
	alpha     float64
	interval  float64
	init      bool
	mutex     sync.Mutex
}

func (e *Ewma) Update(n int64) {
	atomic.AddInt64(&e.uncounted, n)
}
func (e *Ewma) Tick() {
	count := atomic.LoadInt64(&e.uncounted)
	atomic.AddInt64(&e.uncounted, -count)
	instantRate := float64(count) / e.interval
	if e.init {
		oldRate := e.rate
		e.rate = oldRate + (e.alpha * (instantRate - oldRate))
	} else {
		e.rate = instantRate
		e.init = true
	}
}

func (e *Ewma) Rate(duration time.Duration) float64 {
	return e.rate * float64(duration.Nanoseconds())
}
