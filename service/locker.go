package service

import (
	"github.com/silentred/template/util"
)

var redisLocker util.Locker

func GetRedisLocker() util.Locker {
	if redisLocker == nil {
		redisLocker = util.NewRedisLocker(RedisClient, 3)
	}

	return redisLocker
}
