package main

import (
	"common"
	"moduleA"
	"moduleB"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	err := common.MultiBind(
		e,
		map[string]func(c *echo.Echo, root string) error{
			"/module_a": moduleA.Bind,
			"/module_b": moduleB.Bind,
		},
	)
	if err != nil {
		panic(err)
	}

	err = e.Start(":1323")
	if err != nil {
		panic(err)
	}
}
