package service

import (
	"testing"

	"time"

	"fmt"

	"github.com/silentred/template/util"
	"github.com/labstack/echo"
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

	InitTestDBs()
	InitServices()

	dbSuite := new(DBTestSuite)
	Injector.Apply(dbSuite)
	suite.Run(t, dbSuite)
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type DBTestSuite struct {
	suite.Suite
	UserSvc *UserSV `inject`
}

// SetupTest runs before each test
func (suite *DBTestSuite) SetupTest() {
}

// All methods that begin with "Test" are run as tests within a suite.
func (suite *DBTestSuite) TestMysql() {
	suite.testUpdateDeceasePoints(20)

	suite.testCreatePlayer()
	suite.testGetPlayer()
}

// tested
func (suite *DBTestSuite) testCreatePlayer() {
	// insert player
	player := AffiliatePlayer{
		DeviceID:  fmt.Sprintf("test:%d", util.RandomCreateBytes(20)),
		BundleID:  "test",
		Points:    123,
		SDKVer:    "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userSV := suite.UserSvc

	err := userSV.UserRepo.createPlayer(player)
	suite.Assert().NoError(err)
}

func (suite *DBTestSuite) testGetPlayer() {
	userSV := suite.UserSvc

	player, err := userSV.UserRepo.getPlayerBy("sdf", "sdf")
	suite.Assert().Error(err, player.BundleID)
}

func (suite *DBTestSuite) testUpdateDeceasePoints(point int) {
	// update Points
	userSV := suite.UserSvc

	point, err := userSV.updatePlayerPoint("test", "test", point)
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, point >= 0)
}

func (suite *DBTestSuite) TestRedis() {
	userSV := UserSV{}
	Injector.Apply(&userSV)
	deviceID := "d123"
	// token, err := userSV.getPlayToken(deviceID)
	// suite.Assert().Equal(redis.Nil, err)
	// suite.Assert().Equal("", token)

	token, err := userSV.GetPlayTokenByDeviceID(deviceID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, len(token) > 0)
}
