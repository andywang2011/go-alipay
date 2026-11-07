package alipay

import "context"

// TitleDynamicReq is the request to fetch dynamic invoice title info.
type TitleDynamicReq struct {
	BizContent string // biz_content JSON
}

// TitleDynamicResp is the response for dynamic invoice title.
type TitleDynamicResp struct {
	TitleResp struct {
		BaseResp
		Title TitleDynamic `json:"title"`
	} `json:"alipay_ebpp_invoice_title_dynamic_get_response"`
	Sign string `json:"sign"`
}

// TitleDynamic holds invoice title fields.
type TitleDynamic struct {
	UserId          string `json:"user_id"`           // Alipay user ID
	TitleName       string `json:"title_name"`        // invoice title name
	TitleType       string `json:"title_type"`        // title type
	IsDefault       bool   `json:"is_default"`        // default title flag
	LogonId         string `json:"logon_id"`          // Alipay login ID
	UserEmail       string `json:"user_email"`        // user email
	TaxRegisterNo   string `json:"tax_register_no"`   // tax registration number
	UserAddress     string `json:"user_address"`      // address
	OpenBankName    string `json:"open_bank_name"`    // bank name
	OpenBankAccount string `json:"open_bank_account"` // bank account
	TelePhoneNo     string `json:"tele_phone_no"`     // telephone
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/open/02auys
func (o *TitleDynamicReq) ApiName() string {
	return "alipay.ebpp.invoice.title.dynamic.get"
}

func (o *TitleDynamicReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) InvoiceTitleInfo(ctx context.Context, req *TitleDynamicReq) (*TitleDynamicResp, error) {
	out := &TitleDynamicResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.TitleResp.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
