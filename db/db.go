// Author:         Wanghe
// Email:          wanghe@renrenche.com
// Author website: http://example.cn
//
// File: users.go
// Create Date: 2016-06-29 20:18:57

// 一主多从的数据库处理， 从数据库用于读操作,
// 多个从数据库的连接对象存放一个环形链表中，按顺序循环使用连接实例
// 主数据库用于写操作， 目前只设置一个主数据库即可
package db

import (
	"container/ring"
	"fmt"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pelletier/go-toml"
)

type DBMap map[string]*DatabaseManager

func InitDB(f string) (DBMap, error) {
	conf, err := toml.LoadFile(f)
	if err != nil {
		return nil, err
	}

	dbMap := make(DBMap)

	if err := dbMap.New(conf); err != nil {
		return nil, err
	}

	return dbMap, nil
}

func (d DBMap) New(conf *toml.TomlTree) error {
	dbNames := conf.Get("db").(*toml.TomlTree).Keys()
	for _, name := range dbNames {
		key := fmt.Sprintf("db.%s", name)
		dbconf := conf.Get(key).(*toml.TomlTree)
		dm := new(DatabaseManager)
		err := dm.New(dbconf)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("fail to create db map.")
			return err
		}
		d[name] = dm
	}
	return nil
}

type DatabaseManager struct {
	r       *ring.Ring
	writeDB *xorm.Engine
}

func initOrm(conf *toml.TomlTree) (*xorm.Engine, error) {
	datasource := conf.Get("data_source").(string)

	orm, err := xorm.NewEngine("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("error connect to %s %s", datasource, err)
	}

	defaultMaxIdle := int64(10)
	maxIdle := conf.GetDefault("max_idle", defaultMaxIdle).(int64)
	orm.SetMaxIdleConns(int(maxIdle))

	defaultMaxOpen := int64(20)
	maxOpen := conf.GetDefault("max_open", defaultMaxOpen).(int64)
	orm.SetMaxOpenConns(int(maxOpen))

	isDebug := conf.GetDefault("is_debug", false).(bool)
	if isDebug {
		orm.ShowSQL(true)
		orm.ShowExecTime(true)
	}

	return orm, nil
}

func (dm *DatabaseManager) New(conf *toml.TomlTree) error {
	readconf := conf.Get("read").(*toml.TomlTree)
	slaveKeys := readconf.Keys()
	slaveLen := len(slaveKeys)
	slaveErrorCount := 0
	r := ring.New(slaveLen)
	dm.r = r
	for i := 0; i < slaveLen; i++ {
		// 读配置文件， 建立连接，建立数据库连接
		orm, err := initOrm(readconf.Get(slaveKeys[i]).(*toml.TomlTree))
		if err != nil {
			logrus.WithError(err).Warn("failed to connect db")

			// 连接数据出错的计数, 如果都报错，日志打印错误， 返回错误
			slaveErrorCount++
			if slaveErrorCount == slaveLen {
				e := fmt.Errorf("failed to connect all slave mysql")
				logrus.WithFields(logrus.Fields{
					"error": e.Error(),
				}).Error("mysql cluster ocurred exception.")
				return e
			}
			continue
		}

		dm.r.Value = orm
		dm.r = dm.r.Next()
	}

	// 写数据库连接实例
	masterconf := conf.Get("master").(*toml.TomlTree)
	writeOrm, err := initOrm(masterconf)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to connect master mysql")
		return err
	}
	dm.writeDB = writeOrm

	return nil
}

func (dm *DatabaseManager) db(typ string) *xorm.Engine {
	var orm *xorm.Engine
	switch typ {
	case "read":
		dm.r = dm.r.Next()
		orm = dm.r.Value.(*xorm.Engine)
		break
	case "write":
		orm = dm.writeDB
		break
	}
	return orm
}

func (dm *DatabaseManager) R() *xorm.Engine {
	return dm.db("read")
}

func (dm *DatabaseManager) W() *xorm.Engine {
	return dm.db("write")
}
