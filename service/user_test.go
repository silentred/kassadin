package service

import (
	"testing"

	redis "gopkg.in/redis.v5"

	"time"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo"
	"github.com/silentred/template/util"
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
	testObj.On("GetPlayTokenByDeviceID", "d123").Return("u123", nil)

	// call the code we are testing
	token, err := testObj.GetPlayTokenByDeviceID("d123")
	assert.NoError(t, err)
	assert.Equal(t, "u123", token)

	// assert that the expectations were met
	testObj.AssertExpectations(t)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func Test_DBSuite(t *testing.T) {
	e := echo.New()
	util.InitLogger(e)
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
	// suite.updateDeceasePoints(20)
	// suite.updateDeceasePoints(-10)
	suite.testGetPlayer()
}

// tested
func (suite *DBTestSuite) testCreatePlayer() {
	// insert player
	player := AffiliatePlayer{
		DeviceID:  "test",
		BundleID:  "test",
		Points:    123,
		SDKVer:    "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userSV := UserSV{}
	err := userSV.createPlayer(player)
	suite.Assert().NoError(err)
}

func (suite *DBTestSuite) testGetPlayer() {
	userSV := UserSV{}
	player, err := userSV.getPlayerBy("sdf", "sdf")
	suite.Assert().Error(err, player)
}

func (suite *DBTestSuite) updateDeceasePoints(point int) {
	// update Points
	userSV := UserSV{}
	point, err := userSV.updatePlayerPoint("test", "test", point)
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, point >= 0)
}

func (suite *DBTestSuite) TestRedis() {
	userSV := UserSV{}
	deviceID := "d123"
	// token, err := userSV.getPlayToken(deviceID)
	// suite.Assert().Equal(redis.Nil, err)
	// suite.Assert().Equal("", token)

	token, err := userSV.GetPlayTokenByDeviceID(deviceID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, len(token) > 0)
}
