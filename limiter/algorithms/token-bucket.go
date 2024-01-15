package algorithms

import (
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	lock           sync.Mutex
	ipBucketMap    map[string]int
	capacity       int
	refillCount    int
	refillInterval int
}

func NewTokenBucket(capacity int, refillCount int, refillInterval int) *TokenBucket {
	bucket := TokenBucket{
		ipBucketMap:    make(map[string]int),
		capacity:       capacity,
		refillCount:    refillCount,
		refillInterval: refillInterval,
	}
	return &bucket
}

func (t *TokenBucket) ShouldAllow(ip string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	val, exists := t.ipBucketMap[ip]
	if !exists {
		t.ipBucketMap[ip] = t.capacity - 1
		return true
	} else if val == 0 {
		return false
	} else {
		t.ipBucketMap[ip] = val - 1
		return true
	}
}

func (t *TokenBucket) Init() {
	go t.refill()
}

func (t *TokenBucket) refill() {
	ticker := time.NewTicker(time.Duration(t.refillInterval) * time.Second)
	for {
		_, ok := <-ticker.C
		if ok {
			t.lock.Lock()
			for k, v := range t.ipBucketMap {
				t.ipBucketMap[k] = int(math.Min(float64(v+t.refillCount), float64(t.capacity)))
			}
			t.lock.Unlock()
		}
	}
}
