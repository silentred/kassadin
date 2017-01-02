package util

import (
	"time"

	"gopkg.in/redis.v5"
)

func RedisLock(cli *redis.Client, key string, expire int) bool {
	if expire == 0 {
		expire = 5
	}
	now := time.Now().Unix()
	err := cli.SetNX(key, now, time.Duration(expire)*time.Second).Err()

	return err == nil
}

func RedisUnlock(cli *redis.Client, key string) bool {
	err := cli.Del(key).Err()
	return err == nil
}
