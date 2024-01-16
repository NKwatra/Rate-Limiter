package algorithms

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestAlgorithms(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	t.Cleanup(func() {
		client.Close()
	})
	t.Run("Token Bucket", func(t1 *testing.T) {
		bucket := NewTokenBucket(&TokenBucketArgs{
			Capacity:       10,
			Path:           "/limited",
			RefillCount:    2,
			RefillInterval: time.Second,
			Redis:          client,
		})
		resp := bucket.ShouldAllow("127.0.0.1:3000")
		if !resp.Allow {
			t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 1)
		}
		resp = bucket.ShouldAllow("127.0.0.1:3000")
		if !resp.Allow {
			t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 2)
		}
		resp = bucket.ShouldAllow("127.0.0.1:3000")
		if resp.Allow {
			t.Fatalf("For domain %s, request %d should not be allowed", "127.0.0.1:3000", 3)
		}
		resp = bucket.ShouldAllow("127.0.0.1:3001")
		if !resp.Allow {
			t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3001", 1)
		}
		time.Sleep(time.Second)
		resp = bucket.ShouldAllow("127.0.0.1:3000")
		if !resp.Allow {
			t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 4)
		}
	})

}

// func TestFixedWindowCounter(t *testing.T) {

// 	window := NewFixedWindowCounter(3, 1)
// 	window.Init()
// 	allowed := window.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 1)
// 	}
// 	time.Sleep(time.Second)
// 	allowed = window.ShouldAllow("127.0.0.1:3000")
// 	if allowed {
// 		t.Fatalf("For domain %s, request %d should not be allowed", "127.0.0.1:3000", 2)
// 	}
// 	allowed = window.ShouldAllow("127.0.0.1:3001")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3001", 1)
// 	}
// 	time.Sleep(2 * time.Second)
// 	allowed = window.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 3)
// 	}
// }

// func TestSlidingWindowLog(t *testing.T) {
// 	sw := NewSlidingWindowLog(2, 1)
// 	sw.Init()
// 	allowed := sw.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 1)
// 	}
// 	allowed = sw.ShouldAllow("127.0.0.1:3000")
// 	if allowed {
// 		t.Fatalf("For domain %s, request %d should not be allowed", "127.0.0.1:3000", 2)
// 	}
// 	allowed = sw.ShouldAllow("127.0.0.1:3001")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3001", 1)
// 	}
// 	time.Sleep(2 * time.Second)
// 	allowed = sw.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 3)
// 	}
// }

// func TestSlidingWindowCounter(t *testing.T) {
// 	sw := NewSlidingWindowCounter(2, 2)
// 	sw.Init()
// 	allowed := sw.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 1)
// 	}
// 	allowed = sw.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should  be allowed", "127.0.0.1:3000", 2)
// 	}
// 	time.Sleep(2 * time.Second)
// 	allowed = sw.ShouldAllow("127.0.0.1:3000")
// 	if !allowed {
// 		t.Fatalf("For domain %s, request %d should be allowed", "127.0.0.1:3000", 3)
// 	}
// 	allowed = sw.ShouldAllow("127.0.0.1:3000")
// 	if allowed {
// 		t.Fatalf("For domain %s, request %d should not be allowed", "127.0.0.1:3000", 4)
// 	}
// }
