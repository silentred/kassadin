package service

type UserMockSV struct {
}

func (*UserMockSV) GetByID(id int) *User {
	u := User{"user_11111", "Jason", "11111"}
	return &u
}
