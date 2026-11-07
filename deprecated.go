package alipay

import "context"

// Type aliases for corrected spellings (legacy names remain usable).
type (
	SubscribeSendReq  = SubcribeSendReq
	SubscribeSendResp = SubcribeSendResp
)

// SendSubcribeMsg is deprecated; use SendSubscribeMessage.
func (c *Client) SendSubcribeMsg(ctx context.Context, req *SubcribeSendReq) (*SubcribeSendResp, error) {
	return c.SendSubscribeMessage(ctx, req)
}

// OrderSelfitemillBatchquery is deprecated; use OrderSelfItemBillBatchQuery.
func (c *Client) OrderSelfitemillBatchquery(ctx context.Context, req *OrderSelfitemillBatchqueryReq) (*OrderSelfitemillBatchqueryResp, error) {
	return c.OrderSelfItemBillBatchQuery(ctx, req)
}

// CreateQrcode is deprecated; use CreateQRCode.
func (c *Client) CreateQrcode(ctx context.Context, req *CreateQrcodeReq) (*CreateQrcodeResp, error) {
	return c.CreateQRCode(ctx, req)
}

// LocalitemQuery is deprecated; use LocalItemQuery.
func (c *Client) LocalitemQuery(ctx context.Context, req *LocalitemQueryReq) (*LocalitemQueryResp, error) {
	return c.LocalItemQuery(ctx, req)
}

// CertificationCertificationBatchquery is deprecated; use CertificationBatchQuery.
func (c *Client) CertificationCertificationBatchquery(ctx context.Context, req *CertificationBatchqueryReq) (*CertificationBatchqueryResp, error) {
	return c.CertificationBatchQuery(ctx, req)
}
