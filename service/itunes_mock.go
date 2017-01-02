package service

import (
	"github.com/stretchr/testify/mock"
)

type ItunesMockSV struct {
	mock.Mock
}

func (itune *ItunesMockSV) GenerateAdLink(bundleID, country, playerToken string) (string, int64, error) {
	args := itune.Called(bundleID, country, playerToken)
	return args.String(0), args.Get(1).(int64), args.Error(2)
}
