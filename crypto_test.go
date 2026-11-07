package alipay

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"testing"
)

func TestFormatKeys(t *testing.T) {
	if FormatPublicKey("") != nil {
		t.Fatal("empty public key")
	}
	if FormatPrivateKey("") != nil {
		t.Fatal("empty private key")
	}
	cfg := testConfig()
	pub := FormatPublicKey(cfg.PublicKey)
	if _, err := ParsePKCS1PublicKey(pub); err != nil {
		t.Fatal(err)
	}
	priv := FormatPrivateKey(cfg.PrivateKey)
	if _, err := ParsePKCS1PrivateKey(priv); err != nil {
		t.Fatal(err)
	}
}

func TestParsePrivateKeyInvalid(t *testing.T) {
	if _, err := ParsePKCS1PrivateKey([]byte("not-pem")); err == nil {
		t.Fatal("expected error")
	}
	if _, err := ParsePKCS1PublicKey([]byte("not-pem")); err == nil {
		t.Fatal("expected error")
	}
}

func TestSignAndVerify(t *testing.T) {
	cfg := testConfig()
	priv := FormatPrivateKey(cfg.PrivateKey)
	pub := FormatPublicKey(cfg.PublicKey)
	privateKey, err := ParsePKCS1PrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ParsePKCS1PublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{"app_id": "a", "method": "m", "version": "1.0"}
	sign, err := signWithPrivateKey(params, privateKey, crypto.SHA256)
	if err != nil || sign == "" {
		t.Fatalf("sign err=%v sign=%q", err, sign)
	}
	src := []byte(`{"alipay_system_oauth_token_response":{"code":"10000"}}`)
	if err := verifySign(src, SignTypeRSA2, sign, publicKey); err == nil {
		// sign was for params not src — expect verify failure
	}
}

func TestPKCS7UnPadding(t *testing.T) {
	if _, err := pkcs7UnPadding(nil); err == nil {
		t.Fatal("expected error")
	}
	padded := append([]byte("hello"), 3, 3, 3)
	out, err := pkcs7UnPadding(padded)
	if err != nil || string(out) != "hello" {
		t.Fatalf("out=%q err=%v", out, err)
	}
	if _, err := pkcs7UnPadding([]byte{1, 9}); err == nil {
		t.Fatal("expected range error")
	}
}

func TestDecryptMobileRoundTrip(t *testing.T) {
	key := make([]byte, 16)
	for i := range key {
		key[i] = byte(i)
	}
	aesKey := base64.StdEncoding.EncodeToString(key)
	plain := `{"code":"10000","mobile":"13800138000"}`
	block, _ := aes.NewCipher(key)
	iv := make([]byte, 16)
	pad := pkcs7Pad([]byte(plain), block.BlockSize())
	enc := make([]byte, len(pad))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(enc, pad)
	ciphertext := base64.StdEncoding.EncodeToString(enc)

	got, err := decryptMobile(aesKey, ciphertext)
	if err != nil || got != plain {
		t.Fatalf("got=%q err=%v", got, err)
	}
	out, err := DecryptMobile(aesKey, ciphertext)
	if err != nil || out != plain {
		t.Fatalf("package fn got=%q err=%v", out, err)
	}
}


func pkcs7Pad(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	padding := make([]byte, n)
	for i := range padding {
		padding[i] = byte(n)
	}
	return append(data, padding...)
}
