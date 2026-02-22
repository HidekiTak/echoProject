package common

import (
	"errors"

	"github.com/labstack/echo/v5"
)

func MultiBind(e *echo.Echo, binds map[string]func(e *echo.Echo, root string) error) error {
	if len(binds) == 0 {
		return errors.New("binds is empty")
	}
	multiError := FactoryMultiError()
	for root, bind := range binds {
		multiError.Do(func() error { return bind(e, root) })
	}
	return multiError.Error()
}
