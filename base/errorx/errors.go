package errorx

// --- Global sentinel errors (recommended for errors.Is) ---
var (
	ErrSuccess           = New(CodeSuccess)
	ErrUnknown           = New(CodeErrUnknown)
	ErrCommunication     = New(CodeErrCommunication) // ← ADD THIS
	ErrInvalidArgs       = New(CodeErrInvalidArgs)
	ErrInternal          = New(CodeErrInternal)
	ErrFail              = New(CodeErrFail)
	ErrNoMoreData        = New(CodeErrNoMoreData)
	ErrSkip              = New(CodeErrSkip)
	NoErrSkip            = New(CodeSkip)
	ErrBadToken          = New(CodeErrBadToken)
	ErrInvalidToken      = New(CodeErrInvalidToken)
	ErrNeedAuth          = New(CodeErrNeedAuth)
	ErrVerify            = New(CodeErrVerify)
	ErrExists            = New(CodeErrExists)
	ErrNotExists         = New(CodeErrNotExists)
	ErrDisabled          = New(CodeErrDisabled)
	ErrConflict          = New(CodeErrConflict)
	ErrLogic             = New(CodeErrLogic)
	ErrResourceExhausted = New(CodeErrResourceExhausted)
	ErrPartSuccess       = New(CodeErrPartSuccess)
	ErrUnimplemented     = New(CodeErrUnimplemented)
	ErrCrashed           = New(CodeErrCrashed)
	ErrOverflow          = New(CodeErrOverflow)
)
