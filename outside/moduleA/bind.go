package moduleA

import (
	"context"
	"net/http"
	"outcommon"

	"github.com/labstack/echo/v5"
)

func Bind(c *echo.Echo, root string) (outcommon.PreCondChecker, error) {
	c.GET(root, func(c *echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})

	return &checker{}, nil
}

type checker struct{}

func (c *checker) Check(e context.Context) error {
	// DB接続先チェック等
	// チェックミスはMultiError使ってまとめて送り返す
	return nil
}
