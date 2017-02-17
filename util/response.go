package util

import "github.com/labstack/echo"
import . "github.com/HelloWorldDev/be-service/error"

type Resp struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func Response(c echo.Context, data interface{}) error{
    c.JSON(200, Resp{
        Code:    CodeSuccess,
        Message: "",
        Data:    data,
    })

    return nil
}

func ResponseError(c echo.Context, e error) error{
    if err, ok := e.(Error); ok {
        c.JSON(200, Resp{
            Code: err.Code,
            Message: err.Message,
            Data: nil,
        })
    } else {
        c.JSON(200, Resp{
            Code: CodeSystemError,
            Message: "response error type not matched",
            Data: nil,
        })
    }

    return nil
}