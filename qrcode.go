package alipay

import "context"

// CreateQrcodeReq is the request to create a mini-program QR code.
type CreateQrcodeReq struct {
	BizContent string // biz_content JSON
}

// CreateQrcodeResp is the response for QR code creation.
type CreateQrcodeResp struct {
	AlipayQrcodeResponse struct {
		BaseResp
		QrCodeUrl            string `json:"qr_code_url"`              // square QR image URL
		QrCodeUrlCircleWhite string `json:"qr_code_url_circle_white"` // circular QR on white background
		QrCodeUrlCircleBlue  string `json:"qr_code_url_circle_blue"`  // circular QR on blue background
	} `json:"alipay_open_app_qrcode_create_response"`
	Sign string `json:"sign"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/02ovqp
func (o *CreateQrcodeReq) ApiName() string {
	return "alipay.open.app.qrcode.create"
}

func (o *CreateQrcodeReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) CreateQRCode(ctx context.Context, req *CreateQrcodeReq) (*CreateQrcodeResp, error) {
	out := &CreateQrcodeResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.AlipayQrcodeResponse.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
