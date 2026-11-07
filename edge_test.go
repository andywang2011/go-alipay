package alipay

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient_InvalidKeys(t *testing.T) {
	cfg := testConfig()
	cfg.PrivateKey = "not-a-key"
	if _, err := NewClient(cfg); err == nil {
		t.Fatal("expected private key error")
	}
	cfg = testConfig()
	cfg.PublicKey = "not-a-key"
	if _, err := NewClient(cfg); err == nil {
		t.Fatal("expected public key error")
	}
}

func TestNewClient_Defaults(t *testing.T) {
	cfg := Config{
		AppID:      testConfig().AppID,
		PrivateKey: testConfig().PrivateKey,
		PublicKey:  testConfig().PublicKey,
	}
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if api.(*Client).gateway != GatewayURL {
		t.Fatal("gateway default")
	}
}

func TestNewAlipayAPI_Valid(t *testing.T) {
	if NewAlipayAPI(testConfig()) == nil {
		t.Fatal("expected API")
	}
}

func TestParsePKCS8PrivateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if _, err := ParsePKCS1PrivateKey(pemBytes); err != nil {
		t.Fatal(err)
	}
	_, err = ParsePKCS1PrivateKey(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("bad")}))
	if err == nil {
		t.Fatal("expected pkcs8 parse error")
	}
	_, err = ParsePKCS1PrivateKey(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: []byte("bad")}))
	if err == nil {
		t.Fatal("expected pkcs1 parse error")
	}
}

func TestParsePublicKeyNotRSA(t *testing.T) {
	block := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte{1, 2, 3}})
	if _, err := ParsePKCS1PublicKey(block); err == nil {
		t.Fatal("expected error")
	}
}

func TestAlipayError_ErrorNoSub(t *testing.T) {
	e := &AlipayError{Code: "1", Msg: "m"}
	if e.Error() == "" {
		t.Fatal("expected message")
	}
	if _, ok := IsAlipayError(errors.New("other")); ok {
		t.Fatal("expected false")
	}
}

func TestVerifySignInvalidBase64(t *testing.T) {
	pub := FormatPublicKey(testConfig().PublicKey)
	publicKey, _ := ParsePKCS1PublicKey(pub)
	if err := verifySign([]byte("x"), SignTypeRSA2, "%%%", publicKey); err == nil {
		t.Fatal("expected decode error")
	}
	if err := verifySign([]byte("x"), "RSA", "%%%", publicKey); !errors.Is(err, ErrUnsupportedSignType) {
		t.Fatalf("err = %v, want ErrUnsupportedSignType", err)
	}
	if err := verifySign([]byte("x"), "", "%%%", publicKey); !errors.Is(err, ErrUnsupportedSignType) {
		t.Fatalf("err = %v, want ErrUnsupportedSignType", err)
	}
}

func TestDecryptMobileErrors(t *testing.T) {
	if _, err := decryptMobile("%%%", "x"); err == nil {
		t.Fatal("bad aes key")
	}
	if _, err := decryptMobile(testConfig().AesKey, "%%%"); err == nil {
		t.Fatal("bad cipher")
	}
	if _, err := decryptMobile(testConfig().AesKey, "AQIDBA=="); err == nil {
		t.Fatal("expected block size error")
	}
}

func TestExtractSignContentBadJSON(t *testing.T) {
	out := extractSignContent([]byte("{"))
	if len(out) == 0 {
		t.Fatal("expected fallback")
	}
}

func TestVerifyResponseSignFailures(t *testing.T) {
	c := testClient(t)
	if err := c.verifyResponseSign([]byte(`{"sign":"!!!","k":"v"}`)); err == nil {
		t.Fatal("expected verify error")
	}
	if err := c.verifyResponseSign([]byte(`{"sign":""}`)); !errors.Is(err, ErrMissingSign) {
		t.Fatalf("err = %v, want ErrMissingSign", err)
	}
	_ = c.verifyResponseSign([]byte(`not-json`))
}

func TestAPIErrorBranches(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error_response":{"code":"40004","msg":"fail"}}`))
	}))
	defer srv.Close()
	c := testClient(t)
	c.gateway = srv.URL
	ctx := context.Background()

	calls := []func() error{
		func() error { _, err := c.CertificateUserBatchQuery(ctx, &CertificateUserBatchqueryReq{}); return err },
		func() error { _, err := c.CertificationBatchQuery(ctx, &CertificationBatchqueryReq{}); return err },
		func() error {
			_, err := c.OrderSelfItemBillBatchQuery(ctx, &OrderSelfitemillBatchqueryReq{})
			return err
		},
		func() error { _, err := c.SendSubscribeMessage(ctx, &SubcribeSendReq{}); return err },
		func() error { _, err := c.InvoiceTitleInfo(ctx, &TitleDynamicReq{}); return err },
		func() error { _, err := c.CreateQRCode(ctx, &CreateQrcodeReq{}); return err },
		func() error { _, err := c.LocalItemQuery(ctx, &LocalitemQueryReq{}); return err },
	}
	for _, fn := range calls {
		if err := fn(); err == nil {
			t.Fatal("expected error")
		}
	}
}

func TestSignWithEmptyValues(t *testing.T) {
	priv := FormatPrivateKey(testConfig().PrivateKey)
	key, _ := ParsePKCS1PrivateKey(priv)
	sign, err := signWithPrivateKey(map[string]string{"empty": ""}, key, crypto.SHA256)
	if err != nil || sign == "" {
		t.Fatalf("err=%v sign=%q", err, sign)
	}
}

func TestMakeRequestParams_SignError(t *testing.T) {
	c := testClient(t)
	c.privateKey = nil
	if _, err := c.makeRequestParams(&OauthTokenReq{Code: "c"}); err == nil {
		t.Fatal("expected sign error")
	}
}

func TestDoCall_InvalidGateway(t *testing.T) {
	c := testClient(t)
	c.gateway = "://bad-url"
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "c"})
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestDoCall_ReadBodyError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer srv.Close()
	c := testClient(t)
	c.gateway = srv.URL
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "c"})
	if err == nil {
		t.Fatal("expected read error")
	}
}

func TestParsePKCS8NonRSA(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(ecKey)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	_, err = ParsePKCS1PrivateKey(pemBytes)
	if err == nil || !strings.Contains(err.Error(), "RSA") {
		t.Fatalf("err=%v", err)
	}
}

func TestParsePublicKeyInvalidDER(t *testing.T) {
	block := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte("invalid")})
	_, err := ParsePKCS1PublicKey(block)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestVerifyPKCS1v15Errors(t *testing.T) {
	pub := FormatPublicKey(testConfig().PublicKey)
	publicKey, _ := ParsePKCS1PublicKey(pub)
	if err := verifyPKCS1v15([]byte("x"), []byte{1, 2, 3}, publicKey, crypto.SHA256); err == nil {
		t.Fatal("expected verify error")
	}
}

func TestExtractSignContentMarshalFail(t *testing.T) {
	old := jsonMarshalResponse
	defer func() { jsonMarshalResponse = old }()
	jsonMarshalResponse = func(any) ([]byte, error) { return nil, errors.New("marshal fail") }
	out := extractSignContent([]byte(`{"sign":"x","k":"v"}`))
	if len(out) == 0 {
		t.Fatal("expected fallback body")
	}
}

func TestDoCallVerifySignSuccess(t *testing.T) {
	c := testClient(t)
	body := []byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"u1"}}`)
	content := extractSignContent(body)
	priv := c.privateKey
	h := crypto.SHA256.New()
	_, _ = h.Write(content)
	hashed := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		t.Fatal(err)
	}
	signed := `{"alipay_system_oauth_token_response":{"code":"10000","user_id":"u1"},"sign":"` +
		base64.StdEncoding.EncodeToString(sig) + `"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(signed))
	}))
	defer srv.Close()
	cfg := testConfig()
	cfg.VerifySign = true
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	cl := api.(*Client)
	cl.gateway = srv.URL
	_, err = cl.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "c"})
	if err != nil {
		t.Fatalf("verify sign call: %v", err)
	}
}

func TestVerifyPKCS1v15Success(t *testing.T) {
	cfg := testConfig()
	priv := FormatPrivateKey(cfg.PrivateKey)
	pub := FormatPublicKey(cfg.PublicKey)
	privateKey, _ := ParsePKCS1PrivateKey(priv)
	publicKey, _ := ParsePKCS1PublicKey(pub)
	src := []byte("content-bytes")
	h := crypto.SHA256.New()
	_, _ = h.Write(src)
	hashed := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyPKCS1v15(src, sig, publicKey, crypto.SHA256); err != nil {
		t.Fatal(err)
	}
}

func TestParsePublicKeyPKIXParseError(t *testing.T) {
	block := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte{0x30, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01}})
	_, err := ParsePKCS1PublicKey(block)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParsePublicKeyEmptyBlock(t *testing.T) {
	_, err := ParsePKCS1PublicKey([]byte("-----BEGIN PUBLIC KEY-----\n-----END PUBLIC KEY-----"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDecryptMobileInvalidAESKeySize(t *testing.T) {
	_, err := decryptMobile(base64.StdEncoding.EncodeToString([]byte("short-key")), "AQIDBAUGB0g=")
	if err == nil {
		t.Fatal("expected cipher error")
	}
}

func TestDecryptMobilePKCS7Error(t *testing.T) {
	key := testConfig().AesKey
	// valid base64, full block, invalid padding
	ct := base64.StdEncoding.EncodeToString([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	_, err := decryptMobile(key, ct)
	if err == nil {
		t.Fatal("expected unpadding error")
	}
}
