package errorx

import (
	"context"
	"fmt"
	"sync/atomic"
)

// Code represents a business-level error code.
// All codes are positive integers. Zero is reserved for success.
type Code int

const (
	CodeSuccess Code = iota

	CodeErrUnknown  // Unknown error (catch-all)
	CodeErrInternal // Server internal error
	CodeErrFail     // General operation failure (use when no better code fits)

	CodeErrNoMoreData
	CodeErrSkip
	CodeSkip

	CodeErrCommunication // Network / RPC / timeout / connection issues

	CodeErrInvalidArgs // Invalid or malformed arguments

	CodeErrBadToken     // Token is malformed
	CodeErrInvalidToken // Token is valid but expired / revoked
	CodeErrNeedAuth     // Authentication required
	CodeErrVerify       // Verification failed (e.g. CAPTCHA, email)

	CodeErrExists    // Resource already exists
	CodeErrNotExists // Resource does not exist
	CodeErrDisabled  // Resource is disabled
	CodeErrConflict  // State conflict (e.g. concurrent update)

	CodeErrLogic // Business rule violation

	CodeErrResourceExhausted // Rate limit, quota, storage full, etc.

	CodeErrPartSuccess // Operation partially succeeded

	CodeErrUnimplemented // Feature not implemented

	CodeErrCrashed
	CodeErrOverflow

	//
	// Custom code range (user-defined)
	//

	CodeErrCustomStart Code = 1000
	CodeErrCustomEnd   Code = 8000
)

func (c Code) Key() string {
	switch c {
	case CodeSuccess:
		return "success"
	case CodeErrUnknown:
		return "unknown error"
	case CodeErrInternal:
		return "internal server error"
	case CodeErrFail:
		return "operation failed"
	case CodeErrNoMoreData:
		return "no more data"
	case CodeErrSkip:
		return "skip"
	case CodeErrCommunication:
		return "communication error"
	case CodeErrInvalidArgs:
		return "invalid arguments"
	case CodeErrBadToken:
		return "bad token"
	case CodeErrInvalidToken:
		return "invalid token"
	case CodeErrNeedAuth:
		return "authentication required"
	case CodeErrVerify:
		return "verification failed"
	case CodeErrExists:
		return "already exists"
	case CodeErrNotExists:
		return "not found"
	case CodeErrDisabled:
		return "disabled"
	case CodeErrConflict:
		return "conflict"
	case CodeErrLogic:
		return "business logic error"
	case CodeErrResourceExhausted:
		return "resource exhausted"
	case CodeErrPartSuccess:
		return "partial success"
	case CodeErrUnimplemented:
		return "not implemented"
	case CodeErrCrashed:
		return "crashed"
	case CodeErrOverflow:
		return "overflow"
	default:
		return ""
	}
}

type FNCode2Message func(code Code) (msg string, ok bool)
type FNCode2MessageWithContext func(ctx context.Context, code Code) (msg string, ok bool)

type fnCode2MessageWrapper struct {
	fnPre            FNCode2Message
	fnEx             FNCode2Message
	fnPreWithContext FNCode2MessageWithContext
	fnExWithContext  FNCode2MessageWithContext
}

var (
	_exCode2Message atomic.Pointer[fnCode2MessageWrapper]
)

// InstallCode2Message warning, not thread safe
func InstallCode2Message(fnPre, fnEx FNCode2Message) {
	InstallCode2MessageEx(fnPre, fnEx, nil, nil)
}

// InstallCode2MessageEx warning, not thread safe
func InstallCode2MessageEx(fnPre, fnEx FNCode2Message, fnPreWithContext, fnExWithContext FNCode2MessageWithContext) {
	_exCode2Message.Store(&fnCode2MessageWrapper{
		fnPre:            fnPre,
		fnEx:             fnEx,
		fnPreWithContext: fnPreWithContext,
		fnExWithContext:  fnExWithContext,
	})
}

func getCode2MessageFn() (fnPre, fnEx FNCode2Message, fnPreWithContext, fnExWithContext FNCode2MessageWithContext) {
	wrapper := _exCode2Message.Load()
	if wrapper == nil {
		return
	}

	fnPre = wrapper.fnPre
	fnEx = wrapper.fnEx
	fnPreWithContext = wrapper.fnPreWithContext
	fnExWithContext = wrapper.fnExWithContext

	return
}

func (c Code) String() string {
	return c.StringWithContext(context.Background())
}

func (c Code) StringWithContext(ctx context.Context) string {
	fnPre, fnEx, fnPreWithContext, fnExWithContext := getCode2MessageFn()

	if fnPreWithContext != nil {
		t, ok := fnPreWithContext(ctx, c)
		if ok {
			return t
		}
	}

	if fnPre != nil {
		t, ok := fnPre(c)
		if ok {
			return t
		}
	}

	t := c.Key()
	if t != "" {
		return t
	}

	if fnExWithContext != nil {
		msg, ok := fnExWithContext(ctx, c)
		if ok {
			return msg
		}
	}

	if fnEx != nil {
		msg, ok := fnEx(c)
		if ok {
			return msg
		}
	}

	return fmt.Sprintf("Unknow error: %d", c)
}

func CodeToMessage(code Code, msg string) string {
	return CodeToMessageWithContext(context.Background(), code, msg)
}

func CodeToMessageWithContext(ctx context.Context, code Code, msg string) string {
	if msg != "" {
		return msg
	}

	return code.StringWithContext(ctx)
}
