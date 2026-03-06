package funcB

import (
	"fmt"
	"net/http"
	"outcommon"

	"github.com/labstack/echo/v5"
)

func Bind(c *echo.Echo, root string) (outcommon.PreCondChecker, error) {
	c.GET(root, helloWorld)
	c.GET(root+"/:id", helloWorldId)
	return nil, nil
}

func helloWorld(c *echo.Context) error {
	return c.JSON(http.StatusOK, "FuncB: Hello World")
}

func helloWorldId(c *echo.Context) error {
	return c.JSON(http.StatusOK, fmt.Sprintf("FuncB: Hello World: %s", c.Param("id")))
}
