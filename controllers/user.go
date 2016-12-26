package controllers

import (
	"beegotest/models"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	// 解析参数，验证参数

	users := models.GetAllUsers()
	u.Data["json"] = users

	// 需要再封装一层，因为需要记录 req 参数, resp 结果, 建议保存到 input.SetData() 中
	// 这样在 timerMiddleware 中可以一起记录

	// 返回错误时，需要手动设置 code
	u.Ctx.Output.SetStatus(403)
	u.ServeJSON()
}
