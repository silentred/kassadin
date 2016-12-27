package service

type UserService interface {
	GetByID(int) *User
}

type User struct {
	Id       string
	Username string
	Password string
}

type UserSV struct {
	// has db
}

func (*UserSV) GetByID(id int) *User {
	u := User{"user_11111", "astaxie", "11111"}
	return &u
}
