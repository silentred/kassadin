package service

import "gopkg.in/redis.v5"
import "fmt"
import "github.com/silentred/template/util"
import "time"

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
}

// User represents user in redis
type User struct {
	ID       int64
	ShowUid  string // "u"+Id
	DeviceID string
}

// UserSV is the implimentation of UserService
type UserSV struct {
	redisCli *redis.Client
}

func NewUserSV(redisCli *redis.Client) *UserSV {
	return &UserSV{redisCli}
}

// GetByID gets user by its ID
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

func (u *UserSV) getPlayToken(deviceID string) (string, error) {
	return u.redisCli.HGet(DeviceIDToPlayerToken, deviceID).Result()
}

func (u *UserSV) getNewUID() (int64, error) {
	return u.redisCli.IncrBy(UidIncrease, 1).Result()
}

func (u *UserSV) setDeviceID(playerToken, deviceID string) error {
	return u.redisCli.HSet(PlayerTokenToDeviceID, playerToken, deviceID).Err()
}

func (u *UserSV) setPlayToken(deviceID, playerToken string) error {
	return u.redisCli.HSet(DeviceIDToPlayerToken, deviceID, playerToken).Err()
}

func (u *UserSV) increaeUserCount() error {
	return u.redisCli.Incr(UserCount).Err()
}

func (u *UserSV) createNewUser(deviceID string) (User, error) {
	var user User
	var err error
	user.DeviceID = deviceID
	// get lock
	lockKey := fmt.Sprintf(DeviceLockFormat, deviceID)
	if util.RedisLock(u.redisCli, lockKey, 3) {
		defer util.RedisUnlock(u.redisCli, lockKey)
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
