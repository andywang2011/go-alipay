package alipay

import "context"

// OrderSelfitemillBatchqueryReq is the request to batch-query local order settlement bills.
type OrderSelfitemillBatchqueryReq struct {
	BizContent string // biz_content JSON; see OrderSelfitemillBatchqueryParam
}

const AlipayShopNotExist = "AlipayShopNotExist"

// OrderSelfitemillBatchqueryParam holds filters for settlement bill batch query.
type OrderSelfitemillBatchqueryParam struct {
	CertificateId      string `json:"certificate_id,omitempty"`       // certificate ID
	SettlementDate     string `json:"settlement_date,omitempty"`      // settlement date (yyyy-MM-dd)
	MiniAppId          string `json:"mini_app_id"`                    // mini-program app ID
	PageNum            int    `json:"page_num"`                       // page number, from 1
	SettleStatus       string `json:"settle_status,omitempty"`        // settlement status
	VerifyDate         string `json:"verify_date,omitempty"`          // verification date
	VerificationShopId string `json:"verification_shop_id,omitempty"` // verification shop ID
	SerialNo           string `json:"serial_no,omitempty"`            // verification sequence number
	CertificateStatus  string `json:"certificate_status,omitempty"`   // certificate status
	PageSize           int    `json:"page_size"`                      // page size
}

// OrderSelfitemillBatchqueryResp is the response for settlement bill batch query.
type OrderSelfitemillBatchqueryResp struct {
	AlipayOpenMiniOrderSelfitembillBatchqueryResponse struct {
		BaseResp
		PageNum   int       `json:"page_num"`   // current page
		PageSize  int       `json:"page_size"`  // page size
		Total     int       `json:"total"`      // total count
		OrderList OrderList `json:"order_list"` // settlement rows
	} `json:"alipay_open_mini_order_selfitembill_batchquery_response"`
	Sign string `json:"sign"`
}

// OrderList is a slice of settlement bill rows.
type OrderList []struct {
	CertificateId             string `json:"certificate_id"`               // certificate ID
	SerialNo                  string `json:"serial_no"`                    // verification sequence number
	MerchantVerificationNo    string `json:"merchant_verification_no"`     // merchant verification number
	VerificationTime          string `json:"verification_time"`            // verification time
	VerificationShopId        string `json:"verification_shop_id"`         // verification shop ID
	VerificationShopName      string `json:"verification_shop_name"`       // verification shop name
	SelfShopName              string `json:"self_shop_name"`               // self-operated shop name
	ItemType                  string `json:"item_type"`                    // item type
	FeeCategory               string `json:"fee_category"`                 // fee category
	OrderScene                string `json:"order_scene"`                  // order scene
	CertificateStatus         string `json:"certificate_status"`           // certificate status
	RealReceiptAmount         string `json:"real_receipt_amount"`          // actual receipt amount
	SettleAmount              string `json:"settle_amount"`                // settlement amount
	PlatformService           string `json:"platform_service"`             // platform service fee
	RefundFee                 string `json:"refund_fee"`                   // refund amount
	PlatformServiceRefund     string `json:"platform_service_refund"`      // refunded platform service fee
	SettleTime                string `json:"settle_time"`                  // settlement time
	SettleStatus              string `json:"settle_status"`                // settlement status
	TradeNo                   string `json:"trade_no"`                     // Alipay trade number
	SettleAccountType         string `json:"settle_account_type"`          // settlement account type
	VerificationTradeNo       string `json:"verification_trade_no"`        // verification trade number
	SettleAccountNo           string `json:"settle_account_no"`            // settlement account number
	MiniAppId                 string `json:"mini_app_id"`                  // mini-program app ID
	SettleBatchId             string `json:"settle_batch_id"`              // settlement batch ID
	PredictSettleTime         string `json:"predict_settle_time"`          // expected settlement time
	CertificatePrice          string `json:"certificate_price"`            // certificate price
	PayDiscountedPrice        string `json:"pay_discounted_price"`         // price after payment discounts
	PayCommission             string `json:"pay_commission"`               // payment commission
	RefundPayCommission       string `json:"refund_pay_commission"`        // refunded payment commission
	IsvAllocAmount            string `json:"isv_alloc_amount"`             // ISV revenue share amount
	PublicAllocAmount         string `json:"public_alloc_amount"`          // public-domain revenue share amount
	HuabeiTechnicalServiceFee string `json:"huabei_technical_service_fee"` // Huabei technical service fee
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/a1699288_alipay.open.mini.order.selfitembill.batchquery
func (o *OrderSelfitemillBatchqueryReq) ApiName() string {
	return "alipay.open.mini.order.selfitembill.batchquery"
}

func (o *OrderSelfitemillBatchqueryReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) OrderSelfItemBillBatchQuery(ctx context.Context, req *OrderSelfitemillBatchqueryReq) (*OrderSelfitemillBatchqueryResp, error) {
	out := &OrderSelfitemillBatchqueryResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipayOpenMiniOrderSelfitembillBatchqueryResponse.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
