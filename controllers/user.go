package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	session "github.com/silentred/echo-session"
	"github.com/silentred/template/service"
)

type UserController struct {
	sv service.UserService
}

func NewUserController(sv service.UserService) *UserController {
	return &UserController{sv}
}

func (u *UserController) GetByID(c echo.Context) error {
	// 解析参数，验证参数
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	user := u.sv.GetByID(intID)

	// session test
	sess := session.Default(c)
	if sess != nil {
		var count int
		v := sess.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}
		sess.Set("count", count)
		sess.Save()
		c.Echo().Logger.Infof("count: %d \n", count)
	}

	c.JSON(http.StatusOK, user)

	return nil
}
