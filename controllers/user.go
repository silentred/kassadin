package controllers

import (
	"beegotest/service"
	"net/http"

	"strconv"

	"fmt"

	"github.com/labstack/echo"
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

	fmt.Println("userController: id is ", id)

	user := u.sv.GetByID(intID)
	c.JSON(http.StatusOK, user)

	return nil
}
