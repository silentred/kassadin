package util

import "github.com/labstack/echo"

// Get the form param
// Return value of key, and whether the key is exists
func GetParam(c echo.Context, key string) (string, bool) {

	r := c.Request()
	if r.Form == nil {
		r.ParseMultipartForm(32 << 20)
	}
	if vs := r.Form[key]; len(vs) > 0 {
		return vs[0], true
	}

	return "", false
}

func GetAppId(c echo.Context) string {
	r := c.Request()
	id := r.Header.Get("App-ID")
	return id
}

func GetAppSecret(c echo.Context) string {
	r := c.Request()
	secret := r.Header.Get("App-Secret")
	return secret
}
