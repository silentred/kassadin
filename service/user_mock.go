package service

import (
	"github.com/stretchr/testify/mock"
)

type UserMockSV struct {
	mock.Mock
}

func (u *UserMockSV) GetPlayTokenByDeviceID(deviceID string) (string, error) {
	args := u.Called(deviceID)
	return args.String(0), args.Error(1)
}
