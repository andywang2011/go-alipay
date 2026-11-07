package alipay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeprecatedClientMethods(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		switch r.Form.Get("method") {
		case "alipay.open.mini.order.selfitembill.batchquery":
			_, _ = w.Write([]byte(`{"alipay_open_mini_order_selfitembill_batchquery_response":{"code":"10000"}}`))
		case "alipay.open.app.qrcode.create":
			_, _ = w.Write([]byte(`{"alipay_open_app_qrcode_create_response":{"code":"10000"}}`))
		case "alipay.open.app.localitem.query":
			_, _ = w.Write([]byte(`{"alipay_open_app_localitem_query_response":{"code":"10000"}}`))
		case "alipay.marketing.certificate.certification.batchquery":
			_, _ = w.Write([]byte(`{"alipay_marketing_certificate_certification_batchquery_response":{"code":"10000"}}`))
		default:
			_, _ = w.Write([]byte(`{"alipay_open_app_mini_templatemessage_send_response":{"code":"10000"}}`))
		}
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	ctx := context.Background()

	if _, err := c.SendSubcribeMsg(ctx, &SubcribeSendReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.OrderSelfitemillBatchquery(ctx, &OrderSelfitemillBatchqueryReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CreateQrcode(ctx, &CreateQrcodeReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.LocalitemQuery(ctx, &LocalitemQueryReq{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CertificationCertificationBatchquery(ctx, &CertificationBatchqueryReq{}); err != nil {
		t.Fatal(err)
	}
}
