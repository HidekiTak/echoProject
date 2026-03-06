package server_a

import (
	"context"
	"moduleA"
	"moduleB"
	"outcommon"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	checker, err := outcommon.MultiBind(
		e,
		map[string]func(c *echo.Echo, root string) (outcommon.PreCondChecker, error){
			"/module_a": moduleA.Bind,
			"/module_b": moduleB.Bind,
		},
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = checker.Check(ctx)
	if err != nil {
		panic(err)
	}

	err = e.Start(":1323")
	if err != nil {
		panic(err)
	}
}
