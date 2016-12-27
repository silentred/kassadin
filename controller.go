package main

import (
	"beegotest/controllers"
	"beegotest/service"
)

func initControllers() *controllers.UserController {
	user := controllers.NewUserController(&service.UserSV{})

	return user
}
