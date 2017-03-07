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

	"github.com/spf13/viper"
)

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
