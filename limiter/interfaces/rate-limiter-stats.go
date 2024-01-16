package interfaces

type RateLimiterStats struct {
	Remaining  int
	Limit      int
	RetryAfter float64
}

type RateLimiterResponse struct {
	Stats      *RateLimiterStats
	Allow      bool
	StatusCode int
	Message    string
}
