package alipay

import (
	"errors"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	cfg := testConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	cases := []Config{
		{PrivateKey: "p", PublicKey: "u"},
		{AppID: "a", PublicKey: "u"},
		{AppID: "a", PrivateKey: "p"},
		{AppID: "a", PrivateKey: "p", PublicKey: "u", SignType: "BAD"},
		{AppID: "a", PrivateKey: "p", PublicKey: "u", SignType: "RSA"},
	}
	for _, c := range cases {
		if err := c.Validate(); err == nil {
			t.Fatalf("expected error for %+v", c)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Gateway == "" || cfg.SignType == "" || cfg.Timeout <= 0 {
		t.Fatal(cfg)
	}
}

func TestNewAlipayAPI_InvalidConfig(t *testing.T) {
	if NewAlipayAPI(Config{}) != nil {
		t.Fatal("expected nil API for invalid config")
	}
}

func TestConfig_ValidateRejectsRSA1(t *testing.T) {
	c := Config{AppID: "a", PrivateKey: "p", PublicKey: "u", SignType: "RSA"}
	if err := c.Validate(); !errors.Is(err, ErrUnsupportedSignType) {
		t.Fatalf("err = %v, want ErrUnsupportedSignType", err)
	}
	c.SignType = SignTypeRSA2
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
}
