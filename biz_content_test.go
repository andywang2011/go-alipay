package alipay

import (
	"encoding/json"
	"testing"
)

func TestBizContent(t *testing.T) {
	s, err := BizContent(map[string]string{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil || m["k"] != "v" {
		t.Fatalf("s=%q err=%v", s, err)
	}
}

func TestBizContent_Invalid(t *testing.T) {
	_, err := BizContent(make(chan int))
	if err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestNewCertificateUserBatchqueryReq(t *testing.T) {
	req, err := NewCertificateUserBatchqueryReq(CertificateUserBatchqueryParam{
		UserId: "u1", PageNum: 1, PageSize: 10,
	})
	if err != nil || req.BizContent == "" {
		t.Fatalf("req=%+v err=%v", req, err)
	}
}

func TestNewCertificationBatchqueryReq(t *testing.T) {
	req, err := NewCertificationBatchqueryReq(CertificationBatchqueryParam{UserId: "u1"})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}

func TestNewOrderSelfitemillBatchqueryReq(t *testing.T) {
	req, err := NewOrderSelfitemillBatchqueryReq(OrderSelfitemillBatchqueryParam{
		MiniAppId: "app", PageNum: 1, PageSize: 20,
	})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}

func TestNewSubcribeSendReq(t *testing.T) {
	req, err := NewSubcribeSendReq(MessageTpl{ToUserId: "u", UserTemplateId: "t", Page: "p"})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}

func TestNewTitleDynamicReq(t *testing.T) {
	req, err := NewTitleDynamicReq(map[string]string{"user_id": "u"})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}

func TestNewCreateQrcodeReq(t *testing.T) {
	req, err := NewCreateQrcodeReq(map[string]string{"url_param": "x"})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}

func TestNewLocalitemQueryReq(t *testing.T) {
	req, err := NewLocalitemQueryReq(LocalitemQueryParam{ItemId: "i1"})
	if err != nil || req.BizContent == "" {
		t.Fatal(err)
	}
}
