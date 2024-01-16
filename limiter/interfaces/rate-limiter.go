package interfaces

type RateLimiter interface {
	ShouldAllow(string) *RateLimiterResponse
}
