// Author:         Wanghe
// Email:          wanghe@renrenche.com
// Author website: http://example.cn
//
// File: redis.go
// Create Date: 2016-06-29 20:18:57

package redis

import (
	"github.com/Sirupsen/logrus"
	toml "github.com/pelletier/go-toml"
	"gopkg.in/redis.v5"
	"log"
)

// RedisManager is managing redis clients.
type RedisManager map[string]*redis.Client

// New is redis manager.
func New(confpath string) map[string]*redis.Client {
	logrus.Debug(confpath)
	conf, err := toml.LoadFile(confpath)
	if err != nil {
		log.Fatalf("redis: %s", err.Error())
		return nil
	}
	var redisManager RedisManager
	redisManager = make(map[string]*redis.Client, 0)
	t := conf.Get("redis").(*toml.TomlTree)
	for _, value := range t.Keys() {
		instanceConf := t.Get(value).(*toml.TomlTree)
		address := instanceConf.Get("address").(string)
		password := instanceConf.Get("password").(string)
		db := instanceConf.Get("db").(int64)
		poolSize := instanceConf.Get("poolsize").(int64)

		client := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       int(db),
			PoolSize: int(poolSize),
		})

		redisManager[value] = client
	}
	return redisManager
}
