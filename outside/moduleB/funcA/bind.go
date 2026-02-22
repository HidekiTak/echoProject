package funcA

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func Bind(c *echo.Echo, root string) error {
	c.GET(root, helloWorld)
	c.GET(root+"/:id", helloWorldId)
	return nil
}

func helloWorld(c *echo.Context) error {
	return c.JSON(http.StatusOK, "FuncA: Hello World")
}

func helloWorldId(c *echo.Context) error {
	return c.JSON(http.StatusOK, fmt.Sprintf("FuncA: Hello World: %s", c.Param("id")))
}
