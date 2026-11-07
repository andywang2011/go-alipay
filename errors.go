package alipay

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingSign is returned when response signature verification is enabled
	// but the gateway response carries no sign field.
	ErrMissingSign = errors.New("alipay response sign is missing")
	// ErrUnsupportedSignType is returned for any sign type other than RSA2 (SHA256withRSA).
	// The legacy RSA (SHA1withRSA) sign type is rejected as cryptographically weak.
	ErrUnsupportedSignType = errors.New("alipay unsupported sign_type: only RSA2 is allowed")
)

// AlipayError represents a structured Alipay gateway or business error.
type AlipayError struct {
	Code    string
	Msg     string
	SubCode string
	SubMsg  string
	RawBody string
}

func (e *AlipayError) Error() string {
	if e.SubCode != "" || e.SubMsg != "" {
		return fmt.Sprintf("alipay error code=%s msg=%s sub_code=%s sub_msg=%s", e.Code, e.Msg, e.SubCode, e.SubMsg)
	}
	return fmt.Sprintf("alipay error code=%s msg=%s", e.Code, e.Msg)
}

// IsAlipayError reports whether err is an *AlipayError.
func IsAlipayError(err error) (*AlipayError, bool) {
	var ae *AlipayError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// IsSuccess reports whether the Alipay business code indicates success.
func (b BaseResp) IsSuccess() bool {
	return b.Code == SuccessCode
}

// NewAlipayErrorFromBase builds an AlipayError from BaseResp fields.
func NewAlipayErrorFromBase(code, msg, subCode, subMsg string) *AlipayError {
	return &AlipayError{
		Code:    code,
		Msg:     msg,
		SubCode: subCode,
		SubMsg:  subMsg,
	}
}

// bizErrFromBase returns an AlipayError when the business code is present and not success.
func bizErrFromBase(b BaseResp) error {
	if b.Code == "" || b.IsSuccess() {
		return nil
	}
	return NewAlipayErrorFromBase(b.Code, b.Msg, b.SubCode, b.SubMsg)
}
