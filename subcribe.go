package alipay

import "context"

// SubcribeSendReq is the request to send a mini-program template message.
type SubcribeSendReq struct {
	BizContent string // biz_content JSON; see MessageTpl
}

// SubcribeSendResp is the response for template message send.
type SubcribeSendResp struct {
	SendResp struct {
		BaseResp
	} `json:"alipay_open_app_mini_templatemessage_send_response"`
	Sign string `json:"sign"`
}

// ApiName returns the Alipay OpenAPI method name.
// Doc: https://opendocs.alipay.com/mini/02cth2
func (o *SubcribeSendReq) ApiName() string {
	return "alipay.open.app.mini.templatemessage.send"
}

func (o *SubcribeSendReq) ApiParam() map[string]string {
	return map[string]string{"biz_content": o.BizContent}
}

func (c *Client) SendSubscribeMessage(ctx context.Context, req *SubcribeSendReq) (*SubcribeSendResp, error) {
	out := &SubcribeSendResp{}
	if err := c.doCall(ctx, req, out); err != nil {
		return nil, err
	}
	if err := bizErrFromBase(out.SendResp.BaseResp); err != nil {
		return nil, err
	}
	return out, nil
}
