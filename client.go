package alipay

import (
	"context"
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andywang2011/go-httplib"
)

// jsonMarshalResponse is swappable in tests to cover marshal error paths.
var jsonMarshalResponse = json.Marshal

// httpDo is swappable in tests (e.g. simulate transport failures).
var httpDo = func(c *http.Client, req *http.Request) (*http.Response, error) {
	return c.Do(req)
}

type Client struct {
	cfg        Config
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	gateway    string
	signType   string
	aesKey     string
	httpClient *http.Client
	retry      *retryManager
}

// NewClient builds an Alipay API client from config.
func NewClient(cfg Config) (API, error) {
	if cfg.Gateway == "" {
		cfg.Gateway = GatewayURL
	}
	if cfg.SignType == "" {
		cfg.SignType = SignTypeRSA2
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if err := (&cfg).Validate(); err != nil {
		return nil, err
	}

	priv, err := ParsePKCS1PrivateKey(FormatPrivateKey(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("alipay private key: %w", err)
	}
	pub, err := ParsePKCS1PublicKey(FormatPublicKey(cfg.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("alipay public key: %w", err)
	}

	retryCfg := cfg.RetryConfig
	normalizeRetry(&retryCfg)

	return &Client{
		cfg:        cfg,
		privateKey: priv,
		publicKey:  pub,
		gateway:    cfg.Gateway,
		signType:   cfg.SignType,
		aesKey:     cfg.AesKey,
		httpClient: httplib.NewClient(httplib.WithTimeout(cfg.Timeout)),
		retry:      newRetryManager(retryCfg),
	}, nil
}

// NewAlipayAPI is the Wire-friendly constructor (returns nil API when config/keys invalid).
func NewAlipayAPI(cfg Config) API {
	api, err := NewClient(cfg)
	if err != nil {
		return nil
	}
	return api
}

func (c *Client) AppID() string {
	return c.cfg.AppID
}

func (c *Client) DecryptMobile(_ context.Context, ciphertext string) (string, error) {
	if c.aesKey == "" {
		return "", fmt.Errorf("alipay config: aes_key is required for DecryptMobile")
	}
	return decryptMobile(c.aesKey, ciphertext)
}

type BaseResp struct {
	Code    string `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

type errorResponseBody struct {
	ErrorResponse struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
	} `json:"error_response"`
}

func (c *Client) doCall(ctx context.Context, req IAlipay, resp any) error {
	return c.retry.execute(ctx, func(ctx context.Context) error {
		return c.doCallOnce(ctx, req, resp)
	})
}

func (c *Client) doCallOnce(ctx context.Context, req IAlipay, resp any) error {
	params, err := c.makeRequestParams(req)
	if err != nil {
		return err
	}

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.gateway, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", ContentType)

	httpResp, err := httpDo(c.httpClient, httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("alipay http status %d: %s", httpResp.StatusCode, string(body))
	}

	raw := string(body)
	if strings.Contains(raw, ErrorResponse) {
		return parseErrorResponse(raw)
	}

	if c.cfg.VerifySign {
		if err := c.verifyResponseSign(body); err != nil {
			return err
		}
	}

	if err := json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("alipay unmarshal: %w", err)
	}
	return nil
}

func parseErrorResponse(raw string) error {
	var er errorResponseBody
	if err := json.Unmarshal([]byte(raw), &er); err != nil {
		return &AlipayError{RawBody: raw, Msg: raw}
	}
	return &AlipayError{
		Code:    er.ErrorResponse.Code,
		Msg:     er.ErrorResponse.Msg,
		SubCode: er.ErrorResponse.SubCode,
		SubMsg:  er.ErrorResponse.SubMsg,
		RawBody: raw,
	}
}

func (c *Client) verifyResponseSign(body []byte) error {
	var envelope struct {
		Sign string `json:"sign"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("alipay verify sign: %w", err)
	}
	if envelope.Sign == "" {
		return ErrMissingSign
	}
	content := extractSignContent(body)
	return verifySign(content, c.signType, envelope.Sign, c.publicKey)
}

// extractSignContent returns JSON bytes used for response signature verification.
func extractSignContent(body []byte) []byte {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return body
	}
	delete(raw, "sign")
	b, err := json.Marshal(raw)
	if err != nil {
		return body
	}
	return b
}

func (c *Client) makeRequestParams(req IAlipay) (map[string]string, error) {
	data := map[string]string{
		"app_id":    c.cfg.AppID,
		"method":    req.ApiName(),
		"format":    Format,
		"charset":   Charset,
		"sign_type": c.signType,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"version":   Version,
	}
	for k, v := range req.ApiParam() {
		data[k] = v
	}
	if c.signType != SignTypeRSA2 {
		return nil, ErrUnsupportedSignType
	}
	sign, err := signWithPrivateKey(data, c.privateKey, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	data["sign"] = sign
	return data, nil
}

// Request is deprecated; use doCall via API methods.
func (c *Client) Request(ctx context.Context, alipay IAlipay, result interface{}) error {
	return c.doCall(ctx, alipay, result)
}

// MakeRequestParams is deprecated; use makeRequestParams.
func (c *Client) MakeRequestParams(alipay IAlipay) (map[string]string, error) {
	return c.makeRequestParams(alipay)
}
