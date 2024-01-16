package algorithms

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/NKwatra/Rate-Limiter/interfaces"
	"github.com/redis/go-redis/v9"
)

// Structure of redis has for storing algorithm state
type TokenBucketRedisStructure struct {
	TokenLeft  int       `redis:"TokenLeft"`
	NextRefill time.Time `redis:"NextRefill"`
}

func (ds TokenBucketRedisStructure) MarshalBinary() ([]byte, error) {
	return json.Marshal(ds)
}

// Args to set up algorithm
type TokenBucketArgs struct {
	Capacity       int
	RefillCount    int
	RefillInterval time.Duration
	Redis          *redis.Client
	Path           string
}

// Token Bucket Algorithm
type TokenBucket struct {
	capacity       int
	refillCount    int
	refillInterval time.Duration
	redis          *redis.Client
	path           string
}

func NewTokenBucket(config *TokenBucketArgs) interfaces.RateLimiter {
	bucket := TokenBucket{
		redis:          config.Redis,
		refillCount:    config.RefillCount,
		capacity:       config.Capacity,
		refillInterval: config.RefillInterval,
		path:           config.Path,
	}
	return &bucket
}

func (t *TokenBucket) ShouldAllow(ip string) *interfaces.RateLimiterResponse {
	key := createRedisKey(ip, t.path)
	var response interfaces.RateLimiterResponse
	watchErr := watchForTransaction(t.redis, key, func(tx *redis.Tx) error {
		ctx := context.Background()
		exists, err := tx.Exists(ctx, key).Result()
		if err != nil {
			return err
		}
		fmt.Printf("Exists %d\n", exists)
		if exists == 0 {
			_, tnxError := tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
				p.HMSet(ctx, key, "TokenLeft", t.capacity-1, "NextRefill", time.Now().Add(t.refillInterval))
				p.Expire(ctx, key, time.Duration(t.capacity/t.refillCount)*t.refillInterval)
				return nil
			})
			if tnxError != nil {
				return tnxError
			} else {
				response = interfaces.RateLimiterResponse{
					Allow:      true,
					StatusCode: 200,
					Message:    "Allowed",
					Stats: &interfaces.RateLimiterStats{
						Remaining:  t.capacity - 1,
						Limit:      t.refillCount,
						RetryAfter: 0,
					},
				}
				return nil
			}
		} else {
			var value TokenBucketRedisStructure
			err = tx.HGetAll(ctx, key).Scan(&value)
			if err != nil {
				return err
			}
			if value.TokenLeft > 0 || time.Now().After(value.NextRefill) {
				newValue := value.TokenLeft
				nextRefill := value.NextRefill
				if time.Now().After(value.NextRefill) {
					diff := int(time.Since(value.NextRefill)/t.refillInterval) + 1
					newValue = int(math.Min(float64((diff)*t.refillCount+value.TokenLeft), float64(t.capacity)))
					nextRefill = value.NextRefill.Add(time.Duration(diff+1) * (t.refillInterval))
				}

				fmt.Printf("New data, value: %d, nextRefill: %v", newValue, nextRefill)

				_, tnxError := tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
					p.HMSet(ctx, key, "TokenLeft", newValue-1, "NextRefill", nextRefill)
					p.Expire(ctx, key, time.Duration(t.capacity/t.refillCount)*t.refillInterval)
					return nil
				})
				if tnxError != nil {
					return tnxError
				} else {
					response = interfaces.RateLimiterResponse{
						Allow:      true,
						StatusCode: 200,
						Message:    "Allowed",
						Stats: &interfaces.RateLimiterStats{
							Remaining:  t.capacity - 1,
							Limit:      t.refillCount,
							RetryAfter: 0,
						},
					}
					return nil
				}
			} else {
				response = interfaces.RateLimiterResponse{
					Allow:      false,
					StatusCode: 429,
					Message:    "Too Many Requests",
					Stats: &interfaces.RateLimiterStats{
						Remaining:  0,
						Limit:      t.refillCount,
						RetryAfter: time.Until(value.NextRefill).Seconds(),
					},
				}
				return nil
			}
		}
	}, 10)
	if watchErr != nil {
		return &interfaces.RateLimiterResponse{
			Allow:      false,
			StatusCode: 500,
			Message:    "Something went wrong",
			Stats: &interfaces.RateLimiterStats{
				Remaining:  -1,
				Limit:      -1,
				RetryAfter: -1,
			},
		}
	}
	return &response
}
