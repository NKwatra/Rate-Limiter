package algorithms

import (
	"sync"
	"time"
)

type bucket struct {
	prevTime  time.Time
	prevCount int
	currCount int
}

type SlidingWindowCounter struct {
	lock        sync.Mutex
	ipCounts    map[string]*bucket
	maxRequests int
	duration    int
}

func NewSlidingWindowCounter(maxRequests, duration int) *SlidingWindowCounter {
	sw := SlidingWindowCounter{
		ipCounts:    make(map[string]*bucket),
		maxRequests: maxRequests,
		duration:    duration,
	}
	return &sw
}

func (sw *SlidingWindowCounter) Init() {
	go sw.reset()
}

func (sw *SlidingWindowCounter) ShouldAllow(ip string) bool {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	_, exists := sw.ipCounts[ip]
	if !exists {
		sw.ipCounts[ip] = &bucket{
			prevCount: 0,
			prevTime:  time.Now().Add(time.Duration(-sw.duration) * time.Second),
			currCount: 0,
		}
	}
	val := sw.ipCounts[ip]
	val.currCount++
	currWindowStart := time.Now().Add(time.Duration(-sw.duration) * time.Second)
	prevWindowEnd := val.prevTime.Add(time.Duration(sw.duration) * time.Second)
	fraction := float64(prevWindowEnd.Sub(currWindowStart)) / float64(time.Duration(sw.duration)*time.Second)
	return int(fraction*float64(val.prevCount))+val.currCount <= sw.maxRequests
}

func (sw *SlidingWindowCounter) reset() {
	sw.lock.Lock()
	for k, v := range sw.ipCounts {
		sw.ipCounts[k] = &bucket{
			prevTime:  v.prevTime.Add(time.Duration(sw.duration) * time.Second),
			prevCount: v.currCount,
			currCount: 0,
		}
	}
	sw.lock.Unlock()
}
