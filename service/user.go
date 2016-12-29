package service

import "github.com/astaxie/beego/orm"

// UserService represents the service for user
type UserService interface {
	GetByID(int) *User
}

// User represents affi_user table
type User struct {
	Id       uint64  `orm:"pk;column(uid)" json:"uid"`
	Username string  `orm:"column(username)" json:"username"`
	Income   float64 `orm:"column(income)" json:"income"`
}

func (u *User) TableName() string {
	return "affi_user"
}

// UserSV is the implimentation of UserService
type UserSV struct {
	ormer orm.Ormer
}

// GetByID gets user by its ID
func (u *UserSV) GetByID(id int) *User {
	user := User{Id: uint64(id)}
	err := u.ormer.Read(&user)
	if err != nil {
		// NEED to log the error
		return nil
	}

	return &user
}
