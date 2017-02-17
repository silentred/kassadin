package controllers

import (
	"strconv"

	"fmt"

	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/labstack/echo"
)

type UserController struct {
	UserSV   service.UserService   `inject`
	ItunesSV service.ItunesService `inject`
}

func NewUserController() *UserController {
	u := &UserController{}
	service.Injector.Apply(u)
	return u
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
	ReqID   string `json:"req_id"`
}

func (e ErrorResp) Error() string {
	return e.Message
}

func (e *ErrorResp) Fill(code int, msg string) {
	e.Code = code
	e.Message = msg
}

func newErrResp(code int, msg string) ErrorResp {
	return ErrorResp{code, msg, ""}
}

// GenerateLink for user
func (u *UserController) GenerateLink(c echo.Context) error {
	var errResp ErrorResp
	var queryDTO GenLinkDTO
	errResp.ReqID = c.Request().Header.Get("X-Request-ID")

	queryDTO.BundleID = c.FormValue("bundleId")
	queryDTO.DeviceID = c.FormValue("playerToken")
	queryDTO.Country = c.FormValue("country")
	queryDTO.OSVersion = c.FormValue("os_version")

	if len(queryDTO.BundleID) == 0 {
		errResp.Fill(1, "app is empty")
		c.JSON(404, errResp)
		return errResp
	}

	token, err := u.UserSV.GetPlayTokenByDeviceID(queryDTO.DeviceID)
	if err != nil {
		c.Logger().Debug(err)
		errResp.Fill(1, "player token is empty")
		c.JSON(404, errResp)
		return err
	}

	urlStr, appID, err := u.ItunesSV.GenerateAdLink(queryDTO.BundleID, queryDTO.Country, token)
	if err != nil {
		c.Logger().Debug(err)
		errResp.Fill(1, "url is empty")
		c.JSON(404, errResp)
		return err
	}

	ret := map[string]interface{}{
		"url":                urlStr,
		"app_id":             appID,
		"bundleId":           queryDTO.BundleID,
		"error_code":         200,
		"switch":             0,
		"invalid_os_version": 0,
	}

	c.Logger().Info(fmt.Sprintf("reqID:%s app_id:%d bundleID:%s deviceID:%s os_version:%s country:%s retURL:%s", errResp.ReqID, appID, queryDTO.BundleID, queryDTO.DeviceID, queryDTO.OSVersion, queryDTO.Country, urlStr))
	return c.JSON(200, ret)
}

func (u *UserController) GetPoint(ctx echo.Context) error {
	var deviceID, bundleID string
	var errResp ErrorResp

	errResp.ReqID = ctx.Request().Header.Get("X-Request-ID")
	deviceID = ctx.FormValue("playerToken")
	bundleID = ctx.FormValue("bundleId")
	if deviceID == "" || bundleID == "" {
		errResp.Fill(1, "app is empty")
		ctx.JSON(404, errResp)
		return errResp
	}

	res, err := u.UserSV.HandleGetPlayerPoint(deviceID, bundleID)
	ctx.JSON(200, res)
	return err
}

func (u *UserController) UsePoint(ctx echo.Context) error {
	var deviceID, bundleID string
	var points int
	var errResp ErrorResp
	var err error

	errResp.ReqID = ctx.Request().Header.Get("X-Request-ID")
	deviceID = ctx.FormValue("playerToken")
	bundleID = ctx.FormValue("bundleId")
	points, err = strconv.Atoi(ctx.FormValue("points"))

	if deviceID == "" || bundleID == "" || points <= 0 {
		errResp.Fill(1, "app is empty")
		ctx.JSON(404, errResp)
		return errResp
	}
	if err != nil {
		points = 0
	}
	res, err := u.UserSV.HandleUpdatePlayerPoint(deviceID, bundleID, -1*points)
	ctx.JSON(200, res)
	return err
}

func (u *UserController) Log(ctx echo.Context) error {
	var typeVal int
	var osVer string
	var err error

	reqID := ctx.Request().Header.Get("X-Request-ID")
	osVer = ctx.FormValue("os_version")
	typeVal, err = strconv.Atoi(ctx.FormValue("type"))
	if err != nil {
		return err
	}
	if typeVal == 1 {
		ctx.Logger().Errorf("Invalid OSVersion: %s ReqID: %s", osVer, reqID)
		ctx.JSON(200, util.JSON{"error_code": 200})
	}

	return nil
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
