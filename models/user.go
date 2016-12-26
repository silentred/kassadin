package models

var (
	UserList map[string]*User
)

func init() {
	UserList = make(map[string]*User)
	u := User{"user_11111", "astaxie", "11111"}
	UserList["user_11111"] = &u
}

type User struct {
	Id       string
	Username string
	Password string
}

func GetAllUsers() map[string]*User {
	return UserList
}
