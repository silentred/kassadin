// Author:         Wanghe
// Email:          wanghe@renrenche.com
// Author website: http://example.cn
//
// File: db_test.go
// Create Date: 2016-06-29 20:18:57

package db

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewEngine(t *testing.T) {
	// normal
	_, err := xorm.NewEngine("mysql", "work:worktest@tcp(rds20d724tkol8vb44ky.mysql.rds.aliyuncs.com:3306)/rrc_user?charset=utf8")
	assert.Equal(t, err, nil)

	// not existing address
	_, err = xorm.NewEngine("mysql", "work:worktest@tcp(rdsd724tkol8vb44ky.mysql.rds.aliyuncs.com:3306)/rrc_user?charset=utf8")
	assert.Equal(t, err, nil)

	// wrong format
	_, err = xorm.NewEngine("mysql", "workworktest@rdsd724tkol8vb44ky.mysql.rds.aliyuncs.com:3306)/rrc_user?charset=utf8")
	assert.Equal(t, err, nil)
}

func TestInitDB(t *testing.T) {
	dbmap, err := InitDB("./config.toml")
	assert.NotNil(t, dbmap)
	assert.Nil(t, err)
}

func TestConfig(t *testing.T) {
	config := []byte(`
[[products]]
name = "Hammer"
sku = 738594937

[[products]]
name = "Nail"
sku = 284758393`)

	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer(config))

	obj := viper.Get("products")
	fmt.Printf("%#v", obj)

}
