package service

import "testing"
import "fmt"

func TestUserSV(t *testing.T) {
	mysqlInfo := MysqlInfo{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "jason",
		Password: "jason",
		Database: "fenda",
	}
	InitMysqlORM(mysqlInfo)

	usv := UserSV{ormer}
	user := usv.GetByID(10010026)
	fmt.Println(*user)
}
