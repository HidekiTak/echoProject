package moduleB

import (
	"fmt"
	"moduleB/funcA"
	"moduleB/funcB"
	"net/http"
	"outcommon"

	"github.com/labstack/echo/v5"
)

func Bind(e *echo.Echo, root string) (outcommon.PreCondChecker, error) {
	e.GET(root, helloWorld)
	e.GET(root+"/:id", helloWorldId)

	// bind children
	return outcommon.MultiBind(
		e, nil,
		map[string]func(e *echo.Echo, root string) (outcommon.PreCondChecker, error){
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
