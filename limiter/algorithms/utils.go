package algorithms

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func createRedisKey(ip, path string) string {
	return path + ":" + ip
}

func watchForTransaction(client *redis.Client, key string, tfn func(*redis.Tx) error, retries int) (err error) {
	ctx := context.TODO()
	for i := 0; i < retries; i++ {
		err := client.Watch(ctx, tfn, key)
		if err == redis.TxFailedErr {
			continue
		} else if err != nil {
			return err
		} else {
			return nil
		}
	}
	return errors.New("max retries reached")
}
