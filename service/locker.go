package service

import (
	"github.com/silentred/template/util"
)

var redisLocker util.Locker

func initRedisLocker() {
	if redisLocker == nil {
		redisLocker = util.NewRedisLocker(RedisClient, 3)
	}
}
