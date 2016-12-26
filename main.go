package main

import (
	_ "beegotest/routers"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }

	testConfig()
	testLog()
	beego.Run()
}

func testConfig() {
	user := beego.AppConfig.String("mysqluser")
	fmt.Println(user)
}

func testLog() {
	logs.Async()
	// need to set location by workingDir
	logs.SetLogger(logs.AdapterFile, `{"filename":"storage/log/test.log"}`)
}
