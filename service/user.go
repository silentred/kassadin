package service

import (
	"fmt"
	"time"

	"github.com/silentred/template/util"
	"github.com/astaxie/beego/orm"
	"gopkg.in/redis.v5"
)

const (
	DeviceIDToPlayerToken = "fame:affiliate:device_id:player_token"
	PlayerTokenToDeviceID = "fame:affiliate:player_token:device_id"
	UidIncrease           = "fame:affiliate:user:id"
	DeviceLockFormat      = "qingchang:lock:affiliate:lock:register:deviceId:%s"
	UserCount             = "fame:affiliate:bundle:id:user:count"
)

// UserService represents the service for user
type UserService interface {
	GetPlayTokenByDeviceID(string) (string, error)
	HandleGetPlayerPoint(deviceID, bundleID string) (util.JSON, error)
	HandleUpdatePlayerPoint(deviceID, bundleID string, point int) (util.JSON, error)
}

// User represents user in redis
type User struct {
	ID       int64
	ShowUid  string // "u"+Id
	DeviceID string
}

// UserSV is the implimentation of UserService
type UserSV struct {
	UserRepo *UserRepo   `inject`
	Locker   util.Locker `inject`
}

func NewUserSV() *UserSV {
	svc := &UserSV{}
	Injector.Apply(svc)
	return svc
}

// GetPlayTokenByDeviceID gets user id (u123) by its device ID
func (u *UserSV) GetPlayTokenByDeviceID(deviceID string) (string, error) {
	token, err := u.UserRepo.getPlayToken(deviceID)
	if err != nil {
		if err == redis.Nil {
			// create new user
			user, err := u.createNewUser(deviceID)
			return user.ShowUid, err
		}
		return "", err
	}

	return token, nil
}

func (u *UserSV) createNewUser(deviceID string) (User, error) {
	var user User
	var err error
	user.DeviceID = deviceID
	// get lock
	lockKey := fmt.Sprintf(DeviceLockFormat, deviceID)
	locker := u.Locker
	// here uses a locker object, which can be mocked to test the result of true and false
	if locker.Lock(lockKey, 3) {
		defer locker.Unlock(lockKey)
		// get new Uid
		user.ID, err = u.UserRepo.getNewUID()
		if err != nil {
			return user, err
		}
		user.ShowUid = fmt.Sprintf("u%d", user.ID)

		// ignore err
		u.UserRepo.setPlayToken(deviceID, user.ShowUid)
		u.UserRepo.setDeviceID(user.ShowUid, deviceID)
		// increate user count
		u.UserRepo.increaeUserCount()
		return user, nil
	}

	// cannot get lock, meaning there is another request of creating a new user
	// try to get the new created user 3 times
	for i := 0; i < 3; i++ {
		user.ShowUid, err = u.UserRepo.getPlayToken(deviceID)
		if err == nil {
			_, err := fmt.Sscanf(user.ShowUid, "u%d", &user.ID)
			if err == nil {
				return user, nil
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	return user, fmt.Errorf("cannot get playToken by device")
}

type AffiliatePlayer struct {
	ID        int64     `orm:"column(id);pk"`
	DeviceID  string    `orm:"column(token)"`
	BundleID  string    `orm:"column(app_id)"`
	Points    int       `orm:"column(points)"`
	SDKVer    string    `orm:"column(sdk_version)"`
	ClientID  string    `orm:"column(client_id)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp)"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp)"`
}

func (player *AffiliatePlayer) TableName() string {
	return "affi_player"
}

// tested
func (u *UserSV) updatePlayerPoint(deviceID, bundleID string, point int) (int, error) {
	player, err := u.UserRepo.getPlayerBy(deviceID, bundleID)
	fmt.Printf("player: %#v", player)
	if err != nil {
		return 0, err
	}

	if player.Points+point >= 0 {
		player.Points += point
		err = u.UserRepo.updateMysqlPlayerPt(&player)
		if err != nil {
			return player.Points, err
		}
		return player.Points, nil
	}

	return player.Points, fmt.Errorf("player_id %d no enough point: having %d, delta %d", player.ID, player.Points, point)
}

func (u *UserSV) HandleGetPlayerPoint(deviceID, bundleID string) (util.JSON, error) {
	var result = util.JSON{
		"error_code":  200,
		"bundleId":    bundleID,
		"playerToken": deviceID,
		"points":      0,
	}

	player, err := u.UserRepo.getPlayerBy(deviceID, bundleID)
	if err != nil {
		return result, err
	}
	result["points"] = player.Points

	return result, nil
}

func (u *UserSV) HandleUpdatePlayerPoint(deviceID, bundleID string, point int) (util.JSON, error) {
	var result = util.JSON{
		"error_code":  200,
		"bundleId":    bundleID,
		"playerToken": deviceID,
		"points":      0,
	}
	resultPoint, err := u.updatePlayerPoint(deviceID, bundleID, point)
	if err != nil {
		result["msg"] = err.Error()
		result["error_code"] = 404
		return result, err
	}
	result["points"] = resultPoint
	return result, nil
}

//UserRepo containers all DB operation in here
type UserRepo struct {
	Mysql orm.Ormer     `inject`
	Redis *redis.Client `inject`
}

func NewUserRepo() *UserRepo {
	u := UserRepo{}
	Injector.Apply(&u)
	return &u
}

// mock
func (u *UserRepo) getPlayToken(deviceID string) (string, error) {
	return u.Redis.HGet(DeviceIDToPlayerToken, deviceID).Result()
}

// mock
func (u *UserRepo) getNewUID() (int64, error) {
	return u.Redis.IncrBy(UidIncrease, 1).Result()
}

// mock
func (u *UserRepo) setDeviceID(playerToken, deviceID string) error {
	return u.Redis.HSet(PlayerTokenToDeviceID, playerToken, deviceID).Err()
}

// mock
func (u *UserRepo) setPlayToken(deviceID, playerToken string) error {
	return u.Redis.HSet(DeviceIDToPlayerToken, deviceID, playerToken).Err()
}

// mock
func (u *UserRepo) increaeUserCount() error {
	return u.Redis.Incr(UserCount).Err()
}

// Mysql Mock
func (u *UserRepo) updateMysqlPlayerPt(player *AffiliatePlayer) error {
	_, err := u.Mysql.Update(player, "points")
	if err != nil {
		return err
	}
	return nil
}

// mock
func (u *UserRepo) createPlayer(player AffiliatePlayer) error {
	orm := u.Mysql
	_, err := orm.Insert(&player)
	return err
}

// mock
func (u *UserRepo) getPlayerBy(deviceID, bundleID string) (AffiliatePlayer, error) {
	var player AffiliatePlayer

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		util.Logger.Error(err)
	}
	qb.Select("id", "points").From("affi_player").Where("token = ?").And("app_id = ?").Limit(1)
	sql := qb.String()
	util.Logger.Debug(sql)

	o := u.Mysql
	err = o.Raw(sql, deviceID, bundleID).QueryRow(&player)
	if err != nil {
		if err == orm.ErrNoRows {
			// no need to create new player
		}
		return player, err
	}

	return player, nil
}
