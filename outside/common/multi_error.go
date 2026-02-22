package common

import (
	"errors"
	"fmt"
	"strings"
)

// MultiError 複数の error を貯めて一つにする
type MultiError interface {
	Error() error
	DoMulti(callbacks ...func() error) error
	Do(callback func() error)
}

func FactoryMultiError() MultiError {
	return &multiErrorImpl{
		errors: []error{},
	}
}

type multiErrorImpl struct {
	errors []error
}

func (e *multiErrorImpl) Error() error {
	if len(e.errors) == 0 {
		return nil
	}
	sb := strings.Builder{}
	for i, err := range e.errors {
		if i > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString("- ")
		sb.WriteString(err.Error())
	}
	return errors.Join(e.errors...)
}

func (e *multiErrorImpl) DoMulti(callbacks ...func() error) error {
	for _, callback := range callbacks {
		e.Do(callback)
	}
	return e.Error()
}

func (e *multiErrorImpl) Do(callback func() error) {
	var err error
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New(fmt.Sprintf("recover from panic: %v", rec))
			e.errors = append(e.errors, err)
		}
	}()
	err = callback()
	if err == nil {
		return
	}
	e.errors = append(e.errors, err)
}
