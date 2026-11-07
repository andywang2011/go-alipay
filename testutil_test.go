package alipay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"sync"
)

var (
	testCfgOnce sync.Once
	testCfgData Config
	testCfgErr  error
)

func initTestConfig() {
	testCfgOnce.Do(func() {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			testCfgErr = err
			return
		}
		privDER := x509.MarshalPKCS1PrivateKey(key)
		pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			testCfgErr = err
			return
		}
		privPEM := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}))
		pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
		testCfgData = DefaultConfig()
		testCfgData.AppID = "test-app-id"
		testCfgData.PrivateKey = stripPEM(privPEM)
		testCfgData.PublicKey = stripPEM(pubPEM)
		testCfgData.AesKey = base64.StdEncoding.EncodeToString([]byte("1234567890123456"))
	})
}

func stripPEM(pemStr string) string {
	s := strings.ReplaceAll(pemStr, "-----BEGIN RSA PRIVATE KEY-----", "")
	s = strings.ReplaceAll(s, "-----END RSA PRIVATE KEY-----", "")
	s = strings.ReplaceAll(s, "-----BEGIN PUBLIC KEY-----", "")
	s = strings.ReplaceAll(s, "-----END PUBLIC KEY-----", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}

func testConfig() Config {
	initTestConfig()
	if testCfgErr != nil {
		panic(testCfgErr)
	}
	return testCfgData
}

func testClient(t testingT) *Client {
	t.Helper()
	initTestConfig()
	if testCfgErr != nil {
		t.Fatal(testCfgErr)
	}
	api, err := NewClient(testCfgData)
	if err != nil {
		t.Fatal(err)
	}
	c, ok := api.(*Client)
	if !ok {
		t.Fatal("expected *Client")
	}
	return c
}

type testingT interface {
	Helper()
	Fatal(...any)
}
