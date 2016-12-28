package service

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type UserService interface {
	GetByID(int) *User
}

type User struct {
	Id       uint64  `orm:"pk;column(uid)" json:"uid"`
	Username string  `orm:"column(username)" json:"username"`
	Income   float64 `orm:"column(income)" json:"income"`
}

func (u *User) TableName() string {
	return "affi_user"
}

type UserSV struct {
	ormer orm.Ormer
}

func (u *UserSV) GetByID(id int) *User {
	user := User{Id: uint64(id)}
	err := u.ormer.Read(&user)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &user
}
