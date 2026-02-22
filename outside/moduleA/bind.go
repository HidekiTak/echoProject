package moduleA

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Bind(c *echo.Echo, root string) error {
	c.GET(root, func(c *echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})
	return nil
}
