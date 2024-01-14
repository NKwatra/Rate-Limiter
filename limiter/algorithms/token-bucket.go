package algorithms

import (
	"sync"
	"time"
)

type TokenBucket struct {
	lock        sync.Mutex
	ipBucketMap map[string]int
}

func NewTokenBucket() *TokenBucket {
	bucket := TokenBucket{
		ipBucketMap: make(map[string]int),
	}
	return &bucket
}

func (t *TokenBucket) ShouldAllow(ip string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	val, exists := t.ipBucketMap[ip]
	if !exists {
		t.ipBucketMap[ip] = 9
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
	ticker := time.NewTicker(time.Second)
	for {
		_, ok := <-ticker.C
		if ok {
			t.lock.Lock()
			for k, v := range t.ipBucketMap {
				t.ipBucketMap[k] = v + 1
			}
			t.lock.Unlock()
		}
	}
}
