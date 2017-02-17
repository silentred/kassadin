package service

import (
	"github.com/silentred/template/util/container"
	"github.com/astaxie/beego/orm"
	"github.com/silentred/template/util"
)

type Map map[string]interface{}

var (
	Store    *Map
	Injector container.Injector
)

func init() {
	Store = &Map{}
	Injector = container.NewInjector()
}

func (s *Map) Set(key string, val interface{}) {
	(*s)[key] = val
}

func (s *Map) Get(key string) interface{} {
	if val, ok := (*s)[key]; ok {
		return val
	}
	return nil
}

func (s *Map) Delete(key string) {
	delete(*s, key)
}

func InitDBs() {
	mysqlORM := initMysqlORM(mysqlConfig)
	Injector.Map(mysqlORM)
	Injector.MapTo(mysqlORM, new(orm.Ormer))
	Store.Set("mysql", mysqlORM)

	redisCli := initRedisClient(redisConfig)
	Injector.Map(redisCli)
	Store.Set("redis", redisCli)
}

func InitServices() {
	// from buttom to top
	//InitDBs()

	// cache
	cacheInMem := util.NewMemCache(3600)
	Store.Set("cache.mem", cacheInMem)
	Injector.Map(cacheInMem)
	Injector.MapTo(cacheInMem, new(util.Cache))

	// itunueService
	ituneSvc := NewItunesSV(AdToken)
	Store.Set("svc.itune", ituneSvc)
	Injector.Map(ituneSvc)
	Injector.MapTo(ituneSvc, new(ItunesService))

	// redis locker
	redisLocker := util.NewRedisLocker(nil, 3)
	Injector.Apply(redisLocker)
	Store.Set("locker.redisLocker", redisLocker)
	Injector.Map(redisLocker)
	Injector.MapTo(redisLocker, new(util.Locker))

	// userRepo
	userRepo := NewUserRepo()
	Store.Set("repo.user", userRepo)
	Injector.Map(userRepo)

	// userService
	userSvc := NewUserSV()
	Store.Set("svc.user", userSvc)
	Injector.Map(userSvc)
	Injector.MapTo(userSvc, new(UserService))

	// controller
}

func InitTestDBs() {
	mysqlInfo := MysqlInfo{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "jason",
		Password: "jason",
		Database: "fenda",
	}
	mysqlORM := initMysqlORM(mysqlInfo)
	Injector.Map(mysqlORM)
	Injector.MapTo(mysqlORM, new(orm.Ormer))
	Store.Set("mysql", mysqlORM)

	redisInfo := RedisInfo{
		Host:     "127.0.0.1",
		Port:     6379,
		Database: 0,
	}
	redisCli := initRedisClient(redisInfo)
	Injector.Map(redisCli)
	Store.Set("redis", redisCli)
}
