package router

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		method string
		path   string
		result bool
	}{
		{echo.GET, "/v1/user/123", true},
		{echo.PUT, "/v1/user/123", true},
		{echo.POST, "/v1/user/123", true},
		{echo.DELETE, "/v1/user/123", true},
	}

	viper.Set("app.sessionEnable", true)
	viper.Set("app.sessionProvider", "file")

	e := echo.New()
	InitRoutes(e)
	InitMiddleware(e)
	r := e.Router()

	for _, test := range tests {
		c := e.NewContext(nil, nil)
		r.Find(test.method, test.path, c)
		assert.Equal(t, test.result, len(c.ParamValues()) > 0)
	}
}
