package controllers

import (
	"github.com/labstack/echo"
	"github.com/silentred/template/service"
)

type UserController struct {
	userSV   service.UserService
	itunesSV service.ItunesService
}

func NewUserController(sv service.UserService, itune service.ItunesService) *UserController {
	return &UserController{sv, itune}
}

type GenLinkDTO struct {
	BundleID  string `json:"bundle_id" form:"bundleId" query:"bundleId"`
	DeviceID  string `json:"device_id" form:"playerToken" query:"playerToken"`
	Country   string `json:"country" form:"country" query:"country"`
	OSVersion string `json:"os_version" form:"os_version" query:"os_version"`
}

type ErrorResp struct {
	Code    int    `json:"errcode"`
	Message string `json:"message"`
}

func (e ErrorResp) Error() string {
	return e.Message
}

func (u *UserController) GenerateLink(c echo.Context) error {
	var errResp ErrorResp
	var queryDTO GenLinkDTO
	queryDTO.BundleID = c.QueryParam("bundleId")
	queryDTO.DeviceID = c.QueryParam("playerToken")
	queryDTO.Country = c.QueryParam("country")
	queryDTO.OSVersion = c.QueryParam("os_version")

	c.Logger().Debug(queryDTO)

	if len(queryDTO.BundleID) == 0 {
		errResp.Code = 1
		errResp.Message = "app is empty"
		c.JSON(404, errResp)
		return errResp
	}

	token, err := u.userSV.GetPlayTokenByDeviceID(queryDTO.DeviceID)
	if err != nil {
		errResp.Code = 1
		errResp.Message = "player token is empty"
		c.JSON(404, errResp)
		return errResp
	}

	urlStr, appID, err := u.itunesSV.GenerateAdLink(queryDTO.BundleID, queryDTO.Country, token)
	if err != nil {
		errResp.Code = 1
		errResp.Message = "url is empty"
		c.JSON(404, errResp)
		return errResp
	}

	ret := map[string]interface{}{
		"url":                urlStr,
		"app_id":             appID,
		"bundleId":           queryDTO.BundleID,
		"error_code":         200,
		"switch":             0,
		"invalid_os_version": 0,
	}

	return c.JSON(200, ret)
}

// func (u *UserController) GetByID(c echo.Context) error {
// 	// 解析参数，验证参数
// 	id := c.Param("id")
// 	intID, _ := strconv.Atoi(id)

// 	user := u.sv.GetByID(intID)

// 	// session test
// 	sess := session.Default(c)
// 	if sess != nil {
// 		var count int
// 		v := sess.Get("count")
// 		if v == nil {
// 			count = 0
// 		} else {
// 			count = v.(int)
// 			count += 1
// 		}
// 		sess.Set("count", count)
// 		sess.Save()
// 		c.Echo().Logger.Infof("count: %d \n", count)
// 	}

// 	c.JSON(http.StatusOK, user)

// 	return nil
// }
