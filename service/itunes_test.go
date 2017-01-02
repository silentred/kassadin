package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	redis "gopkg.in/redis.v5"
)

func TestItunesSV(t *testing.T) {
	sv := ItunesSV{}

	tests := []struct {
		bundleID string
		country  string
	}{
		{"com.meituan.imeituan", ""},
		{"com.nintendo.zara", ""},
		{"com.aladinfun.petsisland.alipay", "CN"},
		{"com.digipixie.piggy.app", "SG"},
		{"com.aok.dbp", ""},
	}

	for _, test := range tests {
		app, err := sv.searchAllCountryByBundleID(test.bundleID, test.country)
		assert.NoError(t, err)
		assert.Equal(t, true, len(app.TrackViewUrl) > 0)
	}
}

func TestGenLink(t *testing.T) {
	sv := ItunesSV{}
	tests := []struct {
		bundleID string
		country  string
	}{
		{"com.meituan.imeituan", ""},
		// {"com.nintendo.zara", ""},
		// {"com.aladinfun.petsisland.alipay", "CN"},
		// {"com.digipixie.piggy.app", "SG"},
		// {"com.aok.dbp", ""},
	}

	for _, test := range tests {
		urlStr, appID, err := sv.GenerateAdLink(test.bundleID, test.country, "u123")
		assert.NoError(t, err)
		assert.Equal(t, true, len(urlStr) > 0)
		assert.Equal(t, true, appID > 0)
	}
}

func Test_ItunesSuite(t *testing.T) {
	suite.Run(t, new(ItunesTestSuite))
}

type ItunesTestSuite struct {
	suite.Suite
	redisCli *redis.Client
}

// SetupTest runs before each test
func (suite *ItunesTestSuite) SetupTest() {
	if suite.redisCli == nil {
		redisInfo := RedisInfo{
			Host:     "127.0.0.1",
			Port:     6379,
			Database: 0,
		}
		suite.redisCli = InitRedisClient(redisInfo)
	}
}

func (suite *ItunesTestSuite) TestRedis() {
	sv := ItunesSV{"token", suite.redisCli}

	tests := []struct {
		bundleID string
		country  string
	}{
		{"com.meituan.imeituan", ""},
		{"com.nintendo.zara", ""},
		{"com.aladinfun.petsisland.alipay", "CN"},
		{"com.digipixie.piggy.app", "SG"},
		{"com.aok.dbp", ""},
	}

	for _, test := range tests {
		app, err := sv.searchAllCountryByBundleID(test.bundleID, test.country)
		suite.Assert().NoError(err)
		suite.Assert().Equal(true, len(app.TrackViewUrl) > 0)
	}

}
