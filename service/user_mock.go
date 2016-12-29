package service

import (
	"github.com/stretchr/testify/mock"
)

type UserMockSV struct {
	mock.Mock
}

func (u *UserMockSV) GetByID(id int) *User {
	args := u.Called(id)
	return args.Get(0).(*User)
}
