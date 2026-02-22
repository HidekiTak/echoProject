package moduleB

import (
	"common"
	"fmt"
	"moduleB/funcA"
	"moduleB/funcB"
	"net/http"

	"github.com/labstack/echo/v5"
)

func Bind(e *echo.Echo, root string) error {
	e.GET(root, helloWorld)
	e.GET(root+"/:id", helloWorldId)

	return common.MultiBind(
		e,
		map[string]func(e *echo.Echo, root string) error{
			root + "/func_a": funcA.Bind,
			root + "/func_b": funcB.Bind,
		},
	)
}

func helloWorld(c *echo.Context) error {
	return c.JSON(http.StatusOK, "Hello World")
}

func helloWorldId(c *echo.Context) error {
	return c.JSON(http.StatusOK, fmt.Sprintf("Hello World: %s", c.Param("id")))
}
