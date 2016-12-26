package routers

import (
	"time"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/pingcap/tidb/util/hack"
)

type logInfo struct {
	URI      string        `json:"uri"`
	Method   string        `json:"mtd"`
	CostTime time.Duration `json:"time"`
	Code     int           `json:"code"`
}

func initLogMiddleware(ctx *context.Context) {
	ctx.Input.SetData("start_time", time.Now())
}

func afterLogMiddleware(ctx *context.Context) {
	if v, ok := ctx.Input.GetData("start_time").(time.Time); ok {
		costTime := time.Since(v)
		code := ctx.Output.Status
		if code == 0 {
			code = 200
		}
		log := logInfo{
			URI:      ctx.Input.URI(),
			Method:   ctx.Input.Method(),
			CostTime: costTime / time.Microsecond,
			Code:     code,
		}

		b, err := json.Marshal(log)
		if err != nil {
			beego.Trace(err)
		}
		beego.Info(hack.String(b))
	}
}

func RegisterLogMiddleware() {
	beego.InsertFilter("*", beego.BeforeRouter, initLogMiddleware, false)
	beego.InsertFilter("*", beego.AfterExec, afterLogMiddleware, false)
}
