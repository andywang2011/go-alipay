package alipay

import "context"

// CertificationBatchqueryReq is the request to batch-query certificate info.
type CertificationBatchqueryReq struct {
	BizContent string // biz_content JSON; see CertificationBatchqueryParam
}

// CertificationBatchqueryParam holds biz_content for certificate batch query.
type CertificationBatchqueryParam struct {
	CertificateIdList []string `json:"certificate_id_list"` // certificate IDs
	UserId            string   `json:"user_id"`             // Alipay user ID
	OpenId            string   `json:"open_id,omitempty"`   // open ID of certificate owner
	OrderId           string   `json:"order_id,omitempty"`  // purchase order ID
}

// CertificationBatchqueryResp is the response for certificate batch query.
type CertificationBatchqueryResp struct {
	AlipayMarketingCertificateCertificationBatchqueryResponse struct {
		BaseResp
		CertificateInfoList []CertificateInfo `json:"certificate_info_list"` // certificate list
	} `json:"alipay_marketing_certificate_certification_batchquery_response"`
	Sign string `json:"sign"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/2da99f0e_alipay.marketing.certificate.certification.batchquery
func (o *CertificationBatchqueryReq) ApiName() string {
	return "alipay.marketing.certificate.certification.batchquery"
}

func (o *CertificationBatchqueryReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) CertificationBatchQuery(ctx context.Context, req *CertificationBatchqueryReq) (*CertificationBatchqueryResp, error) {
	out := &CertificationBatchqueryResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipayMarketingCertificateCertificationBatchqueryResponse.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
