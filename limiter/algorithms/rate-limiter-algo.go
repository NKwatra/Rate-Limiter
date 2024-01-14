package algorithms

type RateLimiterAlgorithm interface {
	Init()
	ShouldAllow(ip string) bool
}
