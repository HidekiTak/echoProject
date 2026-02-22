package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactoryMultiError(t *testing.T) {
	me := FactoryMultiError()
	err := me.DoMulti(
		func() error { panic("some error") },
	)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "recover from panic: some error")
}
