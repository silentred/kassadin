package service

type UserMockSV struct {
}

func (*UserMockSV) GetByID(id int) *User {
	u := User{123, "Jason", 12}
	return &u
}
