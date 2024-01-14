package algorithms

import (
	"sync"
	"time"

	ts "github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
)

type SlidingWindowLog struct {
	lock        sync.Mutex
	ipCounts    map[string]*ts.Set
	duration    int
	maxRequests int
}

func NewSlidingWindowLog(duration int, maxRequest int) *SlidingWindowLog {
	window := SlidingWindowLog{
		ipCounts:    make(map[string]*ts.Set),
		duration:    duration,
		maxRequests: maxRequest,
	}
	return &window
}

func (sw *SlidingWindowLog) Init() {
}

func (sw *SlidingWindowLog) ShouldAllow(ip string) bool {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	_, exists := sw.ipCounts[ip]
	if !exists {
		sw.ipCounts[ip] = ts.NewWith(func(a, b interface{}) int {
			return utils.TimeComparator(a, b)
		})
	}
	val := sw.ipCounts[ip]
	now := time.Now()
	excluded := now.Add(time.Duration(-sw.duration) * time.Second)
	for it := val.Iterator(); it.Next(); {
		entry := it.Value()
		if entry.(time.Time).Compare(excluded) <= 0 {
			val.Remove(entry)
		} else {
			break
		}
	}
	val.Add(now)
	return val.Size() <= sw.maxRequests
}
