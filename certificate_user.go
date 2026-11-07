package alipay

import "context"

// CertificateUserBatchqueryReq is the request to query user certificates by filter.
type CertificateUserBatchqueryReq struct {
	BizContent string // biz_content JSON; see CertificateUserBatchqueryParam
}

const (
	// CertificateStatusUnUse is the unused certificate status.
	CertificateStatusUnUse = "UNUSE"
)

// CertificateUserBatchqueryParam holds biz_content for user certificate batch query.
type CertificateUserBatchqueryParam struct {
	ShopId            string `json:"shop_id,omitempty"`            // Alipay shop ID
	UserId            string `json:"user_id"`                      // Alipay user ID
	BelongMerchantId  string `json:"belong_merchant_id,omitempty"` // merchant ID that owns the certificate
	PageNum           int    `json:"page_num"`                     // page number (1-based)
	CertificateStatus string `json:"certificate_status,omitempty"` // certificate status
	PageSize          int    `json:"page_size"`                    // page size
}

// CertificateUserBatchqueryResp is the response for user certificate batch query.
type CertificateUserBatchqueryResp struct {
	AlipayMarketingCertificateUserBatchqueryResponse struct {
		BaseResp
		CertificateInfoList []CertificateInfo `json:"certificate_info_list"` // certificate list
		TotalSize           int               `json:"total_size"`            // total count
		PageNum             int               `json:"page_num"`              // current page
		PageSize            int               `json:"page_size"`             // page size
	} `json:"alipay_marketing_certificate_user_batchquery_response"`
	Sign string `json:"sign"`
}

// CertificateInfo describes one marketing certificate.
type CertificateInfo struct {
	CertificateId  string     `json:"certificate_id"`   // user certificate ID
	Code           string     `json:"code"`             // third-party voucher code
	Status         string     `json:"status"`           // UNUSE / USED / EXPIRED / INVALID
	ValidBeginTime string     `json:"valid_begin_time"` // valid-from time
	ValidEndTime   string     `json:"valid_end_time"`   // expiry time
	OutOrderId     string     `json:"out_order_id"`     // third-party order ID
	CanUse         string     `json:"can_use"`          // true if usable, false otherwise
	AmountInfo     AmountInfo `json:"amount_info"`
	SkuInfo        SkuInfo    `json:"sku_info"`
}

// AmountInfo holds certificate amount breakdown.
type AmountInfo struct {
	OriginalPrice          string `json:"original_price"`           // original item price
	SalePrice              string `json:"sale_price"`               // sale price
	PayAmount              string `json:"pay_amount"`               // amount paid by user
	ReceiptAmount          string `json:"receipt_amount"`           // amount received by merchant
	MerchantDiscountAmount string `json:"merchant_discount_amount"` // merchant discount
	PlatformDiscountAmount string `json:"platform_discount_amount"` // Alipay discount
}

// SkuInfo is the SKU linked to the certificate.
type SkuInfo struct {
	SkuId     string `json:"sku_id"`      // Alipay SKU ID
	OutSkuId  string `json:"out_sku_id"`  // merchant SKU ID
	ItemId    string `json:"item_id"`     // Alipay item ID
	OutItemId string `json:"out_item_id"` // merchant item ID
	Title     string `json:"title"`       // item title
	ItemType  string `json:"item_type"`   // template type: 1 group buy, 2 voucher
	Commodity string `json:"commodity"`   // Alipay M-for-N group-buy payload
}

// Commodity is M-for-N group-buy bundle metadata.
type Commodity struct {
	HasOutCommodityId string          `json:"has_out_commodity_id"`
	GroupName         string          `json:"group_name"`
	TotalCount        int             `json:"total_count"`
	OptionCount       int             `json:"option_count"`
	ItemList          []CommodityItem `json:"item_list"`
}

// CommodityItem is one item in an M-for-N group-buy bundle.
type CommodityItem struct {
	OutCommodityId string `json:"out_commodity_id"`
	Name           string `json:"name"`
	Count          int    `json:"count"`
	CountUnit      string `json:"count_unit"`
	Price          int    `json:"price"`
	Unit           string `json:"unit"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/09695q
func (o *CertificateUserBatchqueryReq) ApiName() string {
	return "alipay.marketing.certificate.user.batchquery"
}

func (o *CertificateUserBatchqueryReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) CertificateUserBatchQuery(ctx context.Context, req *CertificateUserBatchqueryReq) (*CertificateUserBatchqueryResp, error) {
	out := &CertificateUserBatchqueryResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipayMarketingCertificateUserBatchqueryResponse.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
