package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	algos "github.com/NKwatra/Rate-Limiter/algorithms"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	bucket := algos.NewTokenBucket(&algos.TokenBucketArgs{
		Capacity:       2,
		RefillCount:    2,
		RefillInterval: time.Second,
		Redis:          rdb,
		Path:           "/limited",
	})
	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		response := bucket.ShouldAllow(ip)
		if response.Stats.Limit >= 0 {
			w.Header().Add("X-Ratelimit-Limit", fmt.Sprint(response.Stats.Limit))
		}
		if response.Stats.Remaining >= 0 {
			w.Header().Add("X-Ratelimit-Remaining", fmt.Sprint(response.Stats.Remaining))
		}
		if response.Stats.RetryAfter >= 0 {
			w.Header().Add("X-Ratelimit-Retry-After", fmt.Sprint(response.Stats.RetryAfter))
		}
		if !response.Allow {
			w.WriteHeader(response.StatusCode)
			io.WriteString(w, response.Message)
		} else {
			resp, err := http.Get("http://localhost:8080/limited")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Something went wrong")
				return
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		}

	})
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Server failed to start %v", err)
	}
	fmt.Println("Server started successfully!!")
}
