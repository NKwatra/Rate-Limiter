package algorithms

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	lock        sync.Mutex
	ipCounts    map[string]int
	duration    int
	maxRequests int
}

func NewFixedWindowCounter(duration int, maxRequest int) *FixedWindowCounter {
	counter := FixedWindowCounter{
		ipCounts:    make(map[string]int),
		duration:    duration,
		maxRequests: maxRequest,
	}
	return &counter
}

func (w *FixedWindowCounter) Init() {
	go w.reset()
}

func (w *FixedWindowCounter) ShouldAllow(ip string) bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	val, exists := w.ipCounts[ip]
	if !exists {
		w.ipCounts[ip] = 1
	} else {
		w.ipCounts[ip] = val + 1
	}
	return w.ipCounts[ip] <= w.maxRequests
}

func (w *FixedWindowCounter) reset() {
	timer := time.NewTicker(time.Duration(w.duration) * time.Second)
	for {
		<-timer.C
		w.lock.Lock()
		for k := range w.ipCounts {
			w.ipCounts[k] = 0
		}
		w.lock.Unlock()
	}
}
