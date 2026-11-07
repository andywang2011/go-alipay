package alipay

import (
	"fmt"
	"time"
)

const (
	GatewayURL    = "https://openapi.alipay.com/gateway.do"
	GateWayUrl    = GatewayURL // deprecated: use GatewayURL
	Format        = "JSON"
	Charset       = "utf-8"
	Version       = "1.0"
	ContentType   = "application/x-www-form-urlencoded;charset=utf-8"
	SignTypeRSA2  = "RSA2"
	ErrorResponse = "error_response"
	SuccessCode   = "10000"
)

// Config holds Alipay mini-program client configuration.
// Fill this struct in the application layer (not from viper inside this library).
type Config struct {
	// AppID is the Alipay application ID.
	AppID string
	// PrivateKey is the merchant RSA private key PEM.
	PrivateKey string
	// PublicKey is the Alipay platform RSA public key PEM.
	PublicKey string
	// Gateway is the OpenAPI gateway URL; defaults to GatewayURL.
	Gateway string
	// SignType must be RSA2 (default). The legacy RSA (SHA1) sign type is rejected.
	SignType string
	// AesKey is the optional AES key for encrypted responses.
	AesKey string
	// Timeout is the HTTP client timeout.
	Timeout time.Duration
	// VerifySign enables response signature verification. When enabled, responses
	// without a valid sign field are rejected.
	VerifySign bool
	// RetryConfig controls transient HTTP retry behavior.
	RetryConfig RetryConfig
}

// DefaultConfig returns defaults for tests and local development.
func DefaultConfig() Config {
	return Config{
		Gateway:  GatewayURL,
		SignType: SignTypeRSA2,
		Timeout:  5 * time.Second,
	}
}

// ConfigDefault is an alias for DefaultConfig (naming consistency across go-* libs).
func ConfigDefault() Config { return DefaultConfig() }

// Validate checks required fields and applies retry defaults.
func (c *Config) Validate() error {
	if c.AppID == "" {
		return fmt.Errorf("alipay config: app_id is required")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("alipay config: private_key is required")
	}
	if c.PublicKey == "" {
		return fmt.Errorf("alipay config: public_key is required")
	}
	if c.SignType != "" && c.SignType != SignTypeRSA2 {
		return fmt.Errorf("alipay config: sign_type %q: %w", c.SignType, ErrUnsupportedSignType)
	}
	if c.RetryConfig.MaxAttempts <= 0 {
		c.RetryConfig = DefaultRetryConfig()
	} else {
		normalizeRetry(&c.RetryConfig)
	}
	return nil
}
