package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	fnPre, fnEx, _, _ := getCode2MessageFn()
	assert.Nil(t, fnPre)
	assert.Nil(t, fnEx)
}

func TestCodeError1(t *testing.T) {
	// Simulate a DB timeout
	dbErr := errors.New("context deadline exceeded")
	err := Wrap(CodeErrCommunication, dbErr, "failed to query user")

	// Deep wrapping
	wrapped := fmt.Errorf("handler: %w", err)

	// THIS MUST MATCH THE SENTINEL NAME EXACTLY
	if errors.Is(wrapped, ErrCommunication) {
		fmt.Println("Communication error detected") // WILL PRINT
	} else {
		t.Fatal("errors.Is failed to detect CodeErrCommunication")
	}

	// Also verify extraction
	code, msg := CodeFromError(wrapped)
	if code != CodeErrCommunication {
		t.Fatalf("expected code %d, got %d", CodeErrCommunication, code)
	}

	if msg != "failed to query user" {
		t.Fatalf("expected msg %q, got %q", "failed to query user", msg)
	}
}

func utFn1(n int) error {
	if n >= 0 {
		return nil
	}

	return ErrDisabled
}

func TestFnR(t *testing.T) {
	err := utFn1(10)
	assert.Nil(t, err)

	err = utFn1(-11)
	assert.NotNil(t, err)

	c, _ := CodeFromError(err)
	assert.Equal(t, CodeErrDisabled, c)

	codeE, ok := TryGetCodeErrorFromError(err)
	assert.True(t, ok)
	assert.Equal(t, CodeErrDisabled, codeE.Code())
}
