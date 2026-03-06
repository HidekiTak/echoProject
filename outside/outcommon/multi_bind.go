package outcommon

import (
	"errors"

	"github.com/labstack/echo/v5"
)

func MultiBind(e *echo.Echo, checker PreCondChecker, binds map[string]func(e *echo.Echo, root string) (PreCondChecker, error)) (PreCondChecker, error) {
	if len(binds) == 0 {
		return nil, errors.New("binds is empty")
	}
	multiError := FactoryMultiError()
	checkers := &preCondCheckers{}
	if nil != checker {
		checkers.Append(checker)
	}
	for root, bind := range binds {
		multiError.Do(func() error {
			checker, err := bind(e, root)
			if nil != checker {
				checkers.Append(checker)
			}
			return err
		})
	}
	return checkers, multiError.Error()
}
