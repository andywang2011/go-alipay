package alipay

import "testing"

func TestAPI_ApiNameAndParam(t *testing.T) {
	tests := []struct {
		name   string
		req    IAlipay
		method string
		want   map[string]string
	}{
		{
			name:   "oauth",
			req:    &OauthTokenReq{Code: "c1"},
			method: "alipay.system.oauth.token",
			want:   map[string]string{"code": "c1", "grant_type": "authorization_code"},
		},
		{
			name:   "certificate_user",
			req:    &CertificateUserBatchqueryReq{BizContent: `{"user_id":"u"}`},
			method: "alipay.marketing.certificate.user.batchquery",
			want:   map[string]string{"biz_content": `{"user_id":"u"}`},
		},
		{
			name:   "certification",
			req:    &CertificationBatchqueryReq{BizContent: "{}"},
			method: "alipay.marketing.certificate.certification.batchquery",
			want:   map[string]string{"biz_content": "{}"},
		},
		{
			name:   "order_bill",
			req:    &OrderSelfitemillBatchqueryReq{BizContent: "{}"},
			method: "alipay.open.mini.order.selfitembill.batchquery",
			want:   map[string]string{"biz_content": "{}"},
		},
		{
			name:   "subscribe",
			req:    &SubcribeSendReq{BizContent: "{}"},
			method: "alipay.open.app.mini.templatemessage.send",
			want:   map[string]string{"biz_content": "{}"},
		},
		{
			name:   "title",
			req:    &TitleDynamicReq{BizContent: "{}"},
			method: "alipay.ebpp.invoice.title.dynamic.get",
			want:   map[string]string{"biz_content": "{}"},
		},
		{
			name:   "qrcode",
			req:    &CreateQrcodeReq{BizContent: "{}"},
			method: "alipay.open.app.qrcode.create",
			want:   map[string]string{"biz_content": "{}"},
		},
		{
			name:   "localitem",
			req:    &LocalitemQueryReq{BizContent: "{}"},
			method: "alipay.open.app.localitem.query",
			want:   map[string]string{"biz_content": "{}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.ApiName(); got != tt.method {
				t.Fatalf("ApiName=%q want %q", got, tt.method)
			}
			got := tt.req.ApiParam()
			for k, v := range tt.want {
				if got[k] != v {
					t.Fatalf("ApiParam[%q]=%q want %q", k, got[k], v)
				}
			}
		})
	}
}
