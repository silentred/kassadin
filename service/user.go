package service

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/silentred/template/util"
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
}

func NewUserSV() *UserSV {
	return &UserSV{}
}

// GetPlayTokenByDeviceID gets user id (u123) by its device ID
func (u *UserSV) GetPlayTokenByDeviceID(deviceID string) (string, error) {
	token, err := u.getPlayToken(deviceID)
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

// mock
func (u *UserSV) getLocker() util.Locker {
	return GetRedisLocker()
}

// mock
func (u *UserSV) getPlayToken(deviceID string) (string, error) {
	return GetRedisClient().HGet(DeviceIDToPlayerToken, deviceID).Result()
}

// mock
func (u *UserSV) getNewUID() (int64, error) {
	return GetRedisClient().IncrBy(UidIncrease, 1).Result()
}

// mock
func (u *UserSV) setDeviceID(playerToken, deviceID string) error {
	return GetRedisClient().HSet(PlayerTokenToDeviceID, playerToken, deviceID).Err()
}

// mock
func (u *UserSV) setPlayToken(deviceID, playerToken string) error {
	return GetRedisClient().HSet(DeviceIDToPlayerToken, deviceID, playerToken).Err()
}

// mock
func (u *UserSV) increaeUserCount() error {
	return GetRedisClient().Incr(UserCount).Err()
}

func (u *UserSV) createNewUser(deviceID string) (User, error) {
	var user User
	var err error
	user.DeviceID = deviceID
	// get lock
	lockKey := fmt.Sprintf(DeviceLockFormat, deviceID)
	locker := u.getLocker()
	// here uses a locker object, which can be mocked to test the result of true and false
	if locker.Lock(lockKey, 3) {
		defer locker.Unlock(lockKey)
		// get new Uid
		user.ID, err = u.getNewUID()
		if err != nil {
			return user, err
		}
		user.ShowUid = fmt.Sprintf("u%d", user.ID)

		// ignore err
		u.setPlayToken(deviceID, user.ShowUid)
		u.setDeviceID(user.ShowUid, deviceID)
		// increate user count
		u.increaeUserCount()
		return user, nil
	}

	// cannot get lock, meaning there is another request of creating a new user
	// try to get the new created user 3 times
	for i := 0; i < 3; i++ {
		user.ShowUid, err = u.getPlayToken(deviceID)
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

// mock
func (u *UserSV) createPlayer(player AffiliatePlayer) error {
	orm := GetMysqlORM()
	_, err := orm.Insert(&player)
	return err
}

// mock
func (u *UserSV) getPlayerBy(deviceID, bundleID string) (AffiliatePlayer, error) {
	var player AffiliatePlayer

	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		util.Logger.Error(err)
	}
	qb.Select("id", "points").From("affi_player").Where("token = ?").And("app_id = ?").Limit(1)
	sql := qb.String()
	util.Logger.Debug(sql)

	o := GetMysqlORM()
	err = o.Raw(sql, deviceID, bundleID).QueryRow(&player)
	if err != nil {
		if err == orm.ErrNoRows {
			// no need to create new player
		}
		return player, err
	}

	return player, nil
}

// mock
func (u *UserSV) updatePlayerPoint(deviceID, bundleID string, point int) (int, error) {
	orm := GetMysqlORM()

	player, err := u.getPlayerBy(deviceID, bundleID)
	if err != nil {
		return 0, err
	}

	if player.Points+point >= 0 {
		player.Points += point
		_, err := orm.Update(&player, "points")
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

	player, err := u.getPlayerBy(deviceID, bundleID)
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
