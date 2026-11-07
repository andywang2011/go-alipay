package alipay

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Success(t *testing.T) {
	c := testClient(t)
	if c.AppID() != "test-app-id" {
		t.Fatal(c.AppID())
	}
}

func TestClient_DecryptMobile_NoAesKey(t *testing.T) {
	cfg := testConfig()
	cfg.AesKey = ""
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	_, err = api.DecryptMobile(context.Background(), "x")
	if err == nil || !strings.Contains(err.Error(), "aes_key") {
		t.Fatalf("err=%v", err)
	}
}

func TestClient_doCall_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("app_id") != "test-app-id" {
			t.Errorf("app_id=%s", r.Form.Get("app_id"))
		}
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"u1","access_token":"tok"}}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.Gateway = srv.URL
	c := testClient(t)
	c.gateway = srv.URL

	out, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "auth-code"})
	if err != nil {
		t.Fatal(err)
	}
	if out.AlipaySystemOauthToken.UserId != "u1" {
		t.Fatalf("out=%+v", out)
	}
}

func TestClient_BizCodeFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"40004","msg":"Business Failed"}}`))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	ae, ok := IsAlipayError(err)
	if !ok || ae.Code != "40004" || strings.Contains(ae.RawBody, "error_response") {
		t.Fatalf("err=%v", err)
	}
}

func TestClient_HTTPRetry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"u1"}}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.Gateway = srv.URL
	cfg.RetryConfig = RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 2}
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c := api.(*Client)
	c.gateway = srv.URL
	c.retry.sleep = func(context.Context, time.Duration) error { return nil }

	out, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err != nil || attempts != 2 || out.AlipaySystemOauthToken.UserId != "u1" {
		t.Fatalf("attempts=%d err=%v out=%+v", attempts, err, out)
	}
}

func TestClient_doCall_ErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error_response":{"code":"40001","msg":"Invalid","sub_code":"isv","sub_msg":"bad"}}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.Gateway = srv.URL
	c := testClient(t)
	c.gateway = srv.URL

	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	ae, ok := IsAlipayError(err)
	if !ok || ae.Code != "40001" {
		t.Fatalf("err=%v", err)
	}
}

func TestClient_doCall_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("err=%v", err)
	}
}

func TestClient_doCall_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	_, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil || !strings.Contains(err.Error(), "unmarshal") {
		t.Fatalf("err=%v", err)
	}
}

func TestParseErrorResponse_InvalidJSON(t *testing.T) {
	err := parseErrorResponse("not json")
	ae, ok := IsAlipayError(err)
	if !ok || ae.RawBody == "" {
		t.Fatalf("err=%v", err)
	}
}

func TestExtractSignContent(t *testing.T) {
	body := []byte(`{"sign":"abc","alipay_system_oauth_token_response":{"code":"10000"}}`)
	out := extractSignContent(body)
	if strings.Contains(string(out), "sign") {
		t.Fatalf("sign should be removed: %s", out)
	}
}

func TestClient_VerifySignEnabled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000"},"sign":"invalid"}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.Gateway = srv.URL
	cfg.VerifySign = true
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c := api.(*Client)
	c.gateway = srv.URL
	_, err = c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"})
	if err == nil {
		t.Fatal("expected verify error")
	}
}

func TestClient_RequestDeprecated(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000"}}`))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	var out OauthTokenResp
	if err := c.Request(context.Background(), &OauthTokenReq{Code: "c"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.AlipaySystemOauthToken.Code != SuccessCode {
		t.Fatal(out)
	}
	_, err := c.MakeRequestParams(&OauthTokenReq{Code: "c"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_RSA1SignTypeRejected(t *testing.T) {
	cfg := testConfig()
	cfg.SignType = "RSA"
	if _, err := NewClient(cfg); !errors.Is(err, ErrUnsupportedSignType) {
		t.Fatalf("NewClient err = %v, want ErrUnsupportedSignType", err)
	}
	if NewAlipayAPI(cfg) != nil {
		t.Fatal("NewAlipayAPI must return nil for RSA sign type")
	}

	// Defense in depth: even if signType is mutated after construction, signing must fail.
	c := testClient(t)
	c.signType = "RSA"
	if _, err := c.makeRequestParams(&OauthTokenReq{Code: "c"}); !errors.Is(err, ErrUnsupportedSignType) {
		t.Fatalf("makeRequestParams err = %v, want ErrUnsupportedSignType", err)
	}
}

func TestVerifyResponseSignEmptySignRejected(t *testing.T) {
	c := testClient(t)
	body, _ := json.Marshal(map[string]string{"code": "10000"})
	if err := c.verifyResponseSign(body); !errors.Is(err, ErrMissingSign) {
		t.Fatalf("err = %v, want ErrMissingSign", err)
	}
}

func TestClient_VerifySignEnabledRejectsMissingSign(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"forged"}}`))
	}))
	defer srv.Close()

	cfg := testConfig()
	cfg.VerifySign = true
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c := api.(*Client)
	c.gateway = srv.URL
	if _, err = c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"}); !errors.Is(err, ErrMissingSign) {
		t.Fatalf("err = %v, want ErrMissingSign", err)
	}
}

func TestClient_VerifySignDisabledAllowsMissingSign(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000"}}`))
	}))
	defer srv.Close()

	c := testClient(t)
	c.gateway = srv.URL
	if _, err := c.SystemOauthToken(context.Background(), &OauthTokenReq{Code: "x"}); err != nil {
		t.Fatal(err)
	}
}
