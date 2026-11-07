package alipay

import "context"

// API is the unified Alipay mini-program / open API client.
type API interface {
	AppID() string
	SystemOauthToken(ctx context.Context, req *OauthTokenReq) (*OauthTokenResp, error)
	DecryptMobile(ctx context.Context, ciphertext string) (string, error)
	CertificateUserBatchQuery(ctx context.Context, req *CertificateUserBatchqueryReq) (*CertificateUserBatchqueryResp, error)
	CertificationBatchQuery(ctx context.Context, req *CertificationBatchqueryReq) (*CertificationBatchqueryResp, error)
	OrderSelfItemBillBatchQuery(ctx context.Context, req *OrderSelfitemillBatchqueryReq) (*OrderSelfitemillBatchqueryResp, error)
	SendSubscribeMessage(ctx context.Context, req *SubcribeSendReq) (*SubcribeSendResp, error)
	InvoiceTitleInfo(ctx context.Context, req *TitleDynamicReq) (*TitleDynamicResp, error)
	CreateQRCode(ctx context.Context, req *CreateQrcodeReq) (*CreateQrcodeResp, error)
	LocalItemQuery(ctx context.Context, req *LocalitemQueryReq) (*LocalitemQueryResp, error)
}

// IAlipay is the internal OpenAPI request descriptor.
type IAlipay interface {
	ApiName() string
	ApiParam() map[string]string
}
