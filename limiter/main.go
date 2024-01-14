package main

import (
	"fmt"
	"io"
	"net/http"
	algos "nkwatra/limiter/algorithms"
	"time"
)

func main() {
	bucket := algos.NewTokenBucket()
	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		fmt.Printf("Request at %v from ip:%s\n", time.Now(), ip)
		allowed := bucket.ShouldAllow(ip)
		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			io.WriteString(w, "Request limit reached!!")
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
	http.HandleFunc("/unlimited", func(w http.ResponseWriter, r *http.Request) {})
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Server failed to start %v", err)
	}
	bucket.Init()
}
