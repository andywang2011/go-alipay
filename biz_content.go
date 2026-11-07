package alipay

import "encoding/json"

// BizContent JSON-encodes v for the biz_content request field.
func BizContent(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func newReqWithBizContent[T any](p T) (string, error) {
	return BizContent(p)
}

// NewCertificateUserBatchqueryReq builds a request from business parameters.
func NewCertificateUserBatchqueryReq(p CertificateUserBatchqueryParam) (*CertificateUserBatchqueryReq, error) {
	s, err := newReqWithBizContent(p)
	if err != nil {
		return nil, err
	}
	return &CertificateUserBatchqueryReq{BizContent: s}, nil
}

// NewCertificationBatchqueryReq builds a request from business parameters.
func NewCertificationBatchqueryReq(p CertificationBatchqueryParam) (*CertificationBatchqueryReq, error) {
	s, err := newReqWithBizContent(p)
	if err != nil {
		return nil, err
	}
	return &CertificationBatchqueryReq{BizContent: s}, nil
}

// NewOrderSelfitemillBatchqueryReq builds a request from business parameters.
func NewOrderSelfitemillBatchqueryReq(p OrderSelfitemillBatchqueryParam) (*OrderSelfitemillBatchqueryReq, error) {
	s, err := newReqWithBizContent(p)
	if err != nil {
		return nil, err
	}
	return &OrderSelfitemillBatchqueryReq{BizContent: s}, nil
}

// NewSubcribeSendReq builds a subscribe message request from a template payload.
func NewSubcribeSendReq(p MessageTpl) (*SubcribeSendReq, error) {
	s, err := newReqWithBizContent(p)
	if err != nil {
		return nil, err
	}
	return &SubcribeSendReq{BizContent: s}, nil
}

// NewTitleDynamicReq builds a title request from arbitrary biz_content JSON.
func NewTitleDynamicReq(p any) (*TitleDynamicReq, error) {
	s, err := BizContent(p)
	if err != nil {
		return nil, err
	}
	return &TitleDynamicReq{BizContent: s}, nil
}

// NewCreateQrcodeReq builds a QR code request from arbitrary biz_content JSON.
func NewCreateQrcodeReq(p any) (*CreateQrcodeReq, error) {
	s, err := BizContent(p)
	if err != nil {
		return nil, err
	}
	return &CreateQrcodeReq{BizContent: s}, nil
}

// NewLocalitemQueryReq builds a local item query request from business parameters.
func NewLocalitemQueryReq(p LocalitemQueryParam) (*LocalitemQueryReq, error) {
	s, err := newReqWithBizContent(p)
	if err != nil {
		return nil, err
	}
	return &LocalitemQueryReq{BizContent: s}, nil
}
