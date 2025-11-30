package pjson

import (
	"context"

	"github.com/GizmoVault/gotools/base/errorx"
)

type ResponseWrapper struct {
	Code       errorx.Code     `json:"code"`
	Message    string          `json:"message"`
	Resp       interface{}     `json:"resp,omitempty"`
	RawMessage string          `json:"-" yaml:"-"`
	Ctx        context.Context `json:"-" yaml:"-"`
}

func (wr *ResponseWrapper) ApplyCodeError(ce errorx.CodeError) bool {
	return wr.ApplyCodeAndError(ce.GetCode(), ce.GetMsg())
}

func (wr *ResponseWrapper) ApplyCodeAndError(code errorx.Code, msg string) bool {
	wr.Code = code
	wr.RawMessage = msg
	wr.Message = errorx.CodeToMessageWithContext(wr.Ctx, code, msg)

	return wr.Code == errorx.CodeSuccess
}

func (wr *ResponseWrapper) Apply(err error) bool {
	code, msg := errorx.CodeFromError(err)

	return wr.ApplyCodeAndError(code, msg)
}

func (wr *ResponseWrapper) Clone(wro ResponseWrapper) bool {
	wr.Code = wro.Code
	wr.RawMessage = wro.RawMessage
	wr.Message = wro.Message

	return wr.Code == errorx.CodeSuccess
}
