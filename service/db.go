package service

import (
	"github.com/spf13/viper"
)

type MysqlInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Tags     []string
}

type RedisInfo struct {
	Host     string
	Port     int
	Database int
	Tags     []string
}

var (
	mysqlInfo MysqlInfo
	redisInfo RedisInfo
)

func init() {
	mysqlHost := viper.GetString("mysql.host")
	mysqlPort := viper.GetInt("mysql.port")
	mysqlUser := viper.GetString("mysql.user")
	mysqlPwd := viper.GetString("mysql.password")
	mysqlDB := viper.GetString("mysql.db")

	mysqlInfo = MysqlInfo{
		Host:     mysqlHost,
		Port:     mysqlPort,
		User:     mysqlUser,
		Password: mysqlPwd,
		Database: mysqlDB,
	}

	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetInt("redis.port")
	redisDB := viper.GetInt("redis.db")

	redisInfo = RedisInfo{
		Host:     redisHost,
		Port:     redisPort,
		Database: redisDB,
	}
}
