package util

import "github.com/spf13/viper"

func InitConfig() error {
	viper.AddConfigPath(".")
	viper.AddConfigPath(SelfDir())
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
