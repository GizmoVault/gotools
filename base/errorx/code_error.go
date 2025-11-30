package errorx

import "errors"

// CodeError is a structured error carrying:
//   - a business code
//   - an optional custom message
//   - an optional underlying cause (for error chaining)
type CodeError struct {
	code  Code   // Business error code
	msg   string // Optional override message
	cause error  // Wrapped error (for Unwrap)
}

// --- Constructors ---

// New creates a CodeError with the given code and default message.
func New(code Code) CodeError {
	return CodeError{code: code}
}

// NewEx creates a CodeError with a custom message.
func NewEx(code Code, msg string) CodeError {
	return CodeError{code: code, msg: msg}
}

// Wrap wraps an underlying error with a business code and optional message.
func Wrap(code Code, cause error, msg ...string) CodeError {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	return CodeError{
		code:  code,
		msg:   m,
		cause: cause,
	}
}

// FromError converts any error into a CodeError.
//   - nil → success
//   - CodeError in chain → returns it
//   - otherwise → CodeErrInternal with original error as cause
func FromError(err error) CodeError {
	if err == nil {
		return New(CodeSuccess)
	}

	var ce CodeError
	if errors.As(err, &ce) {
		return ce
	}

	return New(CodeErrInternal).WithCause(err)
}

// --- Chainable builders ---

// WithCause returns a copy with the given underlying error.
func (e CodeError) WithCause(cause error) CodeError {
	e.cause = cause

	return e
}

// WithMsg returns a copy with the given custom message.
func (e CodeError) WithMsg(msg string) CodeError {
	e.msg = msg

	return e
}

// --- error interface ---

// Error returns the error message.
// Uses custom msg if set; otherwise falls back to code.String().
func (e CodeError) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return e.code.String()
}

// --- Business methods ---

// GetCode returns the error code.
func (e CodeError) GetCode() Code { return e.code }

// GetMsg returns the custom message (empty if not set).
func (e CodeError) GetMsg() string { return e.msg }

// Success reports whether the operation succeeded.
func (e CodeError) Success() bool { return e.code == CodeSuccess }

// Cause returns the underlying error (for inspection).
func (e CodeError) Cause() error { return e.cause }

// Code returns the error code (enables type assertion).
func (e CodeError) Code() Code { return e.code }

// --- Error chaining support ---

// Unwrap returns the wrapped error to support errors.Unwrap and errors.Is.
func (e CodeError) Unwrap() error {
	return e.cause
}

// --- errors.Is support ---

// Is reports whether the current error matches the target by code.
// Uses errors.As to walk the entire error chain.
func (e CodeError) Is(target error) bool {
	var t CodeError

	return errors.As(target, &t) && e.code == t.code
}

// --- Utility functions ---

// CodeFromError extracts the business code and message from any error.
//   - nil → (CodeSuccess, "")
//   - CodeError found → (code, custom_msg or code.String())
//   - otherwise → (CodeErrInternal, original_error_string)
func CodeFromError(err error) (code Code, errorMsg string) {
	if err == nil {
		code = CodeSuccess

		return
	}

	var ce CodeError
	if errors.As(err, &ce) {
		code = ce.code

		errorMsg = ce.msg
		if errorMsg == "" {
			errorMsg = ce.code.String()
		}

		return
	}

	code = CodeErrUnknown
	errorMsg = err.Error()

	return
}

func TryGetCodeErrorFromError(err error) (CodeError, bool) {
	var ce CodeError

	if errors.As(err, &ce) {
		return ce, true
	}

	return ce, false
}

func CodeErrorFromError(err error, msg string) CodeError {
	var ce CodeError

	if errors.As(err, &ce) {
		if msg != "" {
			ce.msg = msg
		}

		return ce
	}

	return NewEx(CodeErrUnknown, msg)
}
