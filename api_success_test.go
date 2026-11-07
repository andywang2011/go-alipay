package alipay

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISuccessPaths(t *testing.T) {
	responses := map[string]string{
		"alipay.marketing.certificate.user.batchquery":          `{"alipay_marketing_certificate_user_batchquery_response":{"code":"10000","total_size":1}}`,
		"alipay.marketing.certificate.certification.batchquery": `{"alipay_marketing_certificate_certification_batchquery_response":{"code":"10000"}}`,
		"alipay.open.mini.order.selfitembill.batchquery":        `{"alipay_open_mini_order_selfitembill_batchquery_response":{"code":"10000"}}`,
		"alipay.open.app.mini.templatemessage.send":             `{"alipay_open_app_mini_templatemessage_send_response":{"code":"10000"}}`,
		"alipay.ebpp.invoice.title.dynamic.get":                 `{"alipay_ebpp_invoice_title_dynamic_get_response":{"code":"10000"}}`,
		"alipay.open.app.qrcode.create":                         `{"alipay_open_app_qrcode_create_response":{"code":"10000","qr_code_url":"https://x"}}`,
		"alipay.open.app.localitem.query":                       `{"alipay_open_app_localitem_query_response":{"code":"10000","item_id":"i1"}}`,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		body, ok := responses[r.Form.Get("method")]
		if !ok {
			t.Errorf("unknown method %s", r.Form.Get("method"))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	ctx := context.Background()

	if _, err := c.CertificateUserBatchQuery(ctx, &CertificateUserBatchqueryReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CertificationBatchQuery(ctx, &CertificationBatchqueryReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.OrderSelfItemBillBatchQuery(ctx, &OrderSelfitemillBatchqueryReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.SendSubscribeMessage(ctx, &SubcribeSendReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.InvoiceTitleInfo(ctx, &TitleDynamicReq{}); err != nil {
		t.Fatal(err)
	}
	out, err := c.CreateQRCode(ctx, &CreateQrcodeReq{})
	if err != nil || out.AlipayQrcodeResponse.QrCodeUrl == "" {
		t.Fatalf("qrcode err=%v out=%+v", err, out)
	}
	li, err := c.LocalItemQuery(ctx, &LocalitemQueryReq{})
	if err != nil || li.AlipayOpenAppLocalitemQueryResponse.ItemId != "i1" {
		t.Fatalf("localitem err=%v out=%+v", err, li)
	}
}

func TestHTTPDoHook_ReadBodyError(t *testing.T) {
	orig := httpDo
	t.Cleanup(func() { httpDo = orig })

	httpDo = func(_ *http.Client, _ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(errReader{}),
		}, nil
	}
	c := testClient(t)
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil {
		t.Fatal("expected read error")
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (errReader) Close() error             { return nil }
