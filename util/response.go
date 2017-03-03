package util

import "github.com/labstack/echo"

const CodeSuccess = 200
const CodeSystemError = 500

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(c echo.Context, data interface{}) error {
	c.JSON(200, Resp{
		Code:    CodeSuccess,
		Message: "",
		Data:    data,
	})

	return nil
}

func ResponseError(c echo.Context, e error) error {
	c.JSON(200, Resp{
		Code:    CodeSystemError,
		Message: "response error type not matched",
		Data:    nil,
	})
	return nil
}
