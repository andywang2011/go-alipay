// OAuth token exchange against a local mock gateway (httptest).
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	alipay "github.com/andywang2011/go-alipay"
)

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"alipay_system_oauth_token_response":{"code":"10000","user_id":"demo-user","access_token":"demo-tok"}}`))
	}))
	defer srv.Close()

	priv, pub, err := demoKeys()
	if err != nil {
		log.Fatal(err)
	}

	cfg := alipay.DefaultConfig()
	cfg.AppID = "demo-app"
	cfg.PrivateKey = priv
	cfg.PublicKey = pub
	cfg.AesKey = base64.StdEncoding.EncodeToString([]byte("1234567890123456"))
	cfg.Gateway = srv.URL

	api, err := alipay.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	out, err := api.SystemOauthToken(context.Background(), &alipay.OauthTokenReq{Code: "demo-code"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user_id:", out.AlipaySystemOauthToken.UserId)
}

func demoKeys() (privateKey string, publicKey string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	privDER := x509.MarshalPKCS1PrivateKey(key)
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	privPEM := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}))
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return stripPEM(privPEM), stripPEM(pubPEM), nil
}

func stripPEM(pemStr string) string {
	s := strings.ReplaceAll(pemStr, "-----BEGIN RSA PRIVATE KEY-----", "")
	s = strings.ReplaceAll(s, "-----END RSA PRIVATE KEY-----", "")
	s = strings.ReplaceAll(s, "-----BEGIN PUBLIC KEY-----", "")
	s = strings.ReplaceAll(s, "-----END PUBLIC KEY-----", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.TrimSpace(s)
}
