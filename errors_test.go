package alipay

import (
	"errors"
	"testing"
)

func TestBaseResp_IsSuccess(t *testing.T) {
	if !(BaseResp{Code: SuccessCode}).IsSuccess() {
		t.Fatal("expected success")
	}
	if (BaseResp{Code: "40004"}).IsSuccess() {
		t.Fatal("expected failure")
	}
}

func TestAlipayError_ErrorAndIs(t *testing.T) {
	err := &AlipayError{Code: "40001", Msg: "invalid", SubCode: "isv", SubMsg: "bad"}
	if err.Error() == "" {
		t.Fatal("expected message")
	}
	var target *AlipayError
	if !errors.As(err, &target) {
		t.Fatal("IsAlipayError failed")
	}
	got, ok := IsAlipayError(err)
	if !ok || got.Code != "40001" {
		t.Fatalf("got=%+v ok=%v", got, ok)
	}
}

func TestBizErrFromBase(t *testing.T) {
	if err := bizErrFromBase(BaseResp{Code: SuccessCode}); err != nil {
		t.Fatal(err)
	}
	if err := bizErrFromBase(BaseResp{}); err != nil {
		t.Fatal("empty code should not error")
	}
	err := bizErrFromBase(BaseResp{Code: "40004", Msg: "fail"})
	if err == nil {
		t.Fatal("expected error")
	}
	ae, ok := IsAlipayError(err)
	if !ok || ae.Code != "40004" {
		t.Fatalf("err=%v", err)
	}
}

func TestNewAlipayErrorFromBase(t *testing.T) {
	e := NewAlipayErrorFromBase("1", "m", "s1", "s2")
	if e.Code != "1" || e.SubCode != "s1" {
		t.Fatal(e)
	}
}
