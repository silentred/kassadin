package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/silentred/template/service"
	"github.com/silentred/template/util"
	"github.com/stretchr/testify/assert"
)

// TestGet is a sample to run an endpoint test
func Test_UserGetID(t *testing.T) {
	// Setup
	e := echo.New()
	req, err := util.NewHTTPReqeust(echo.GET, "/v1/user/123", nil, nil, nil)

	if assert.NoError(t, err) {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("123")

		mockUserService := &service.UserMockSV{}
		mockUserService.On("GetByID", 123).Return(&service.User{Id: 123, Username: "jason", Income: 1.2})
		controller := NewUserController(mockUserService)

		// Assertions
		if assert.NoError(t, controller.GetByID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	}
}
