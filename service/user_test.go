package service

import (
	"testing"

	redis "gopkg.in/redis.v5"

	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUserSV(t *testing.T) {
	testUserMock(t)
}

func testUserMock(t *testing.T) {
	// create an instance of our test object
	testObj := new(UserMockSV)

	// setup expectations
	testObj.On("GetByID", 123).Return(&User{123, "jason", 1.2})

	// call the code we are testing
	user := testObj.GetByID(123)
	assert.Equal(t, "jason", user.Username)

	// assert that the expectations were met
	testObj.AssertExpectations(t)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func Test_DBSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type DBTestSuite struct {
	suite.Suite
	ormer    orm.Ormer
	redisCli *redis.Client
}

// SetupTest runs before each test
func (suite *DBTestSuite) SetupTest() {
	if suite.ormer == nil {
		mysqlInfo := MysqlInfo{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "jason",
			Password: "jason",
			Database: "fenda",
		}
		suite.ormer = InitMysqlORM(mysqlInfo)
	}

	if suite.redisCli == nil {
		redisInfo := RedisInfo{
			Host:     "127.0.0.1",
			Port:     6379,
			Database: 0,
		}
		suite.redisCli = InitRedisClient(redisInfo)
	}
}

// All methods that begin with "Test" are run as tests within a suite.
func (suite *DBTestSuite) TestMysql() {
	usv := UserSV{suite.ormer}
	user := usv.GetByID(10010026)
	suite.Equal(true, len(user.Username) > 0)
}

func (suite *DBTestSuite) TestRedis() {
	err := suite.redisCli.Set("foo", "bar", 0).Err()
	suite.Nil(err)

	res, err := suite.redisCli.Get("foo").Result()
	suite.Nil(err)
	suite.Equal("bar", res)
}
