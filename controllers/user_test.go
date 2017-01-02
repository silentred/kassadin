package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	ulog "github.com/silentred/echo-log"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/stretchr/testify/assert"
)

func Test_UserGenerateLink(t *testing.T) {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(ulog.DEBUG)

	query := map[string]string{
		"bundleId":    "com.nihao",
		"playerToken": "d123",
		"country":     "cn",
		"os_version":  "10.01",
	}

	query_empty_bundleid := map[string]string{
		"bundleId":    "",
		"playerToken": "d123",
		"country":     "cn",
		"os_version":  "10.01",
	}

	tests := []struct {
		query  map[string]string
		code   int
		hasErr bool
	}{
		{query, 200, false},
		{query_empty_bundleid, 404, true},
	}

	for _, test := range tests {
		req, err := util.NewHTTPReqeust(echo.POST, "/promotion/generatelink", test.query, nil, nil)
		if assert.NoError(t, err) {
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUserService := &service.UserMockSV{}
			mockUserService.On("GetPlayTokenByDeviceID", "d123").Return("u123", nil)

			mockItuneSV := &service.ItunesMockSV{}
			mockItuneSV.On("GenerateAdLink", "com.nihao", "cn", "u123").Return("http://sdf/id1233?at=123", int64(1233), nil)

			controller := NewUserController(mockUserService, mockItuneSV)
			err := controller.GenerateLink(c)

			// Assertions
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.code, rec.Code)
		}
	}
}
