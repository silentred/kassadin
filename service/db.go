package service

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	redis "gopkg.in/redis.v5"
)

type MysqlInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Tags     []string
}

func (my MysqlInfo) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", my.User, my.Password, my.Host, my.Port, my.Database)
}

type RedisInfo struct {
	Host     string
	Port     int
	Database int
	Tags     []string
}

var (
	MysqlConfig MysqlInfo
	RedisConfig RedisInfo
	// mysql ORM
	MysqlORM orm.Ormer
	// redis client
	RedisClient *redis.Client
)

func InitDBInfo() {
	mysqlHost := viper.GetString("mysql.host")
	mysqlPort := viper.GetInt("mysql.port")
	mysqlUser := viper.GetString("mysql.user")
	mysqlPwd := viper.GetString("mysql.password")
	mysqlDB := viper.GetString("mysql.db")

	MysqlConfig = MysqlInfo{
		Host:     mysqlHost,
		Port:     mysqlPort,
		User:     mysqlUser,
		Password: mysqlPwd,
		Database: mysqlDB,
	}

	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")
	redisDB := viper.GetInt("redis.db")

	RedisConfig = RedisInfo{
		Host:     redisHost,
		Port:     redisPort,
		Database: redisDB,
	}
}

func InitMysqlORM(myInfo MysqlInfo) orm.Ormer {
	if MysqlORM == nil {
		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default", "mysql", myInfo.String())
		// register model
		orm.RegisterModel(new(AffiliatePlayer))

		MysqlORM = orm.NewOrm()
	}
	return MysqlORM
}

func InitRedisClient(redisInfo RedisInfo) *redis.Client {
	if RedisClient == nil {
		addr := fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port)
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			DB:       redisInfo.Database,
			Password: "", // no password set
		})
		RedisClient = client
	}
	return RedisClient
}

func GetMysqlORM() orm.Ormer {
	if MysqlORM == nil {
		panic("MysqlORM is not initialized")
	}
	return MysqlORM
}

func GetRedisClient() *redis.Client {
	if RedisClient == nil {
		panic("RedisClient is not initialized")
	}
	return RedisClient
}
