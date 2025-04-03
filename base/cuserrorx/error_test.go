package cuserrorx_test

import (
	"errors"
	"testing"

	"github.com/GizmoVault/gotools/base/cuserrorx"
	"github.com/stretchr/testify/assert"
)

var errI = cuserrorx.NewWithCode(1)

func errIReturn() error {
	return errI
}

func TestError1(t *testing.T) {
	err1 := cuserrorx.NewWithCode(1)
	assert.NotNil(t, cuserrorx.As(err1))

	err2 := errors.New("xx")
	assert.Nil(t, cuserrorx.As(err2))

	assert.False(t, cuserrorx.Is(err1, errI))
	assert.True(t, cuserrorx.Is(errIReturn(), errI))

	err3 := cuserrorx.NewWithError(2, errI)
	assert.True(t, cuserrorx.Is(err3, errI))
}
