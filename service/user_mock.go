package service

import (
	"github.com/silentred/template/util"
	"github.com/stretchr/testify/mock"
)

type UserMockSV struct {
	mock.Mock
}

func (u *UserMockSV) GetPlayTokenByDeviceID(deviceID string) (string, error) {
	args := u.Called(deviceID)
	return args.String(0), args.Error(1)
}

func (u *UserMockSV) HandleGetPlayerPoint(deviceID, bundleID string) (util.JSON, error) {
	args := u.Called()
	return args.Get(0).(util.JSON), args.Error(1)
}

func (u *UserMockSV) HandleUpdatePlayerPoint(deviceID, bundleID string, point int) (util.JSON, error) {
	args := u.Called()
	return args.Get(0).(util.JSON), args.Error(1)
}
