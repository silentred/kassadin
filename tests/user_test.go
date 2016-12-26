package test

import (
	_ "beegotest/routers"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	//beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	assert.Equal(t, 200, w.Code)

	// Convey("Subject: Test Station Endpoint\n", t, func() {
	// 	Convey("Status Code Should Be 200", func() {
	// 		So(w.Code, ShouldEqual, 200)
	// 	})
	// 	Convey("The Result Should Not Be Empty", func() {
	// 		So(w.Body.Len(), ShouldBeGreaterThan, 0)
	// 	})
	// })
}
