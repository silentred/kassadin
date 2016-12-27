package filter

import "github.com/labstack/echo"
import "fmt"

func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			fmt.Println("metrics: before next")
			if err = next(c); err != nil {
				c.Error(err)
			}
			fmt.Println("metrics: after next")

			return nil
		}
	}
}
