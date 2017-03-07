package db

import (
	"container/ring"
	"log"

	"github.com/go-xorm/xorm"
	"github.com/silentred/kassadin"
)

const (
	MaxIdle = 10
	MaxOpen = 20
)

type MysqlManager struct {
	Application  *kassadin.App `inject`
	Config       kassadin.MysqlConfig
	databases    map[string]*xorm.Engine
	readOnlyRing *ring.Ring
	master       *xorm.Engine
}

// NewMysqlManager returns a new MysqlManager
func NewMysqlManager(app *kassadin.App, config kassadin.MysqlConfig) *MysqlManager {
	var readOnlyLength int
	var err error
	var engine *xorm.Engine

	for _, instance := range config.Instances {
		if instance.ReadOnly {
			readOnlyLength++
		}
	}

	mm := &MysqlManager{
		Application:  app,
		Config:       config,
		databases:    make(map[string]*xorm.Engine),
		readOnlyRing: ring.New(readOnlyLength),
	}

	for _, instance := range config.Instances {
		if instance.ReadOnly {
			engine, err = mm.newORM(instance)
			if err != nil {
				mm.readOnlyRing.Value = engine
				mm.readOnlyRing = mm.readOnlyRing.Next()
			}
		} else {
			mm.master, err = mm.newORM(instance)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if mm.master == nil {
		log.Fatal("Mysql master is nil")
	}

	return mm
}

func (mm *MysqlManager) newORM(mysql kassadin.MysqlInstance) (*xorm.Engine, error) {
	orm, err := xorm.NewEngine("mysql", mysql.String())
	if err != nil {
		return nil, err
	}
	orm.SetMaxIdleConns(MaxIdle)
	orm.SetMaxOpenConns(MaxOpen)
	if mm.Application != nil {
		if mm.Application.Config.Mode == kassadin.ModeDev {
			orm.ShowSQL(true)
			orm.ShowExecTime(true)
		}
	}

	return nil, nil
}

// DB gets databases by name
func (mm *MysqlManager) DB(name string) *xorm.Engine {
	if engine, ok := mm.databases[name]; ok {
		return engine
	}
	return nil
}

func (mm *MysqlManager) R() *xorm.Engine {
	if mm.readOnlyRing.Len() == 0 {
		return mm.master
	}
	if e, ok := mm.readOnlyRing.Value.(*xorm.Engine); ok {
		return e
	}
	return nil
}

func (mm *MysqlManager) W() *xorm.Engine {
	return mm.master
}
