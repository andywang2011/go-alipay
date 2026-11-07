package alipay

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"sort"
	"strings"
)

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----")
}

func FormatPrivateKey(raw string) []byte {
	return formatKey(raw, "-----BEGIN RSA PRIVATE KEY-----", "-----END RSA PRIVATE KEY-----")
}

func formatKey(raw, prefix, suffix string) []byte {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.ReplaceAll(raw, " ", "")
	raw = strings.ReplaceAll(raw, "\n", "")
	raw = strings.ReplaceAll(raw, "\r", "")
	raw = strings.ReplaceAll(raw, "\t", "")

	const chunk = 64
	sl := len(raw)
	c := sl / chunk
	if sl%chunk > 0 {
		c++
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		b := i * chunk
		e := b + chunk
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func ParsePKCS1PrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return rsaKey, nil
}

func ParsePKCS1PublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}
	return key, nil
}

func signWithPrivateKey(param map[string]string, privateKey *rsa.PrivateKey, hash crypto.Hash) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is nil")
	}
	pList := make([]string, 0, len(param))
	for key, value := range param {
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	src := strings.Join(pList, "&")

	h := hash.New()
	_, _ = h.Write([]byte(src))
	hashed := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func verifyPKCS1v15(src, sig []byte, publicKey *rsa.PublicKey, hash crypto.Hash) error {
	h := hash.New()
	_, _ = h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, hash, hashed, sig)
}

// verifySign verifies sign over data. Only RSA2 (SHA256withRSA) is accepted.
func verifySign(data []byte, signType, sign string, publicKey *rsa.PublicKey) error {
	if signType != SignTypeRSA2 {
		return ErrUnsupportedSignType
	}
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return verifyPKCS1v15(data, sig, publicKey, crypto.SHA256)
}

// decryptMobile decrypts Alipay mini-program encrypted mobile payload.
func decryptMobile(aesKey, ciphertext string) (string, error) {
	kb, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return "", err
	}
	cb, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(kb)
	if err != nil {
		return "", err
	}
	if len(cb)%block.BlockSize() != 0 {
		return "", errors.New("ciphertext not full blocks")
	}
	iv := make([]byte, 16)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(cb))
	blockMode.CryptBlocks(plain, cb)
	b, err := pkcs7UnPadding(plain)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func pkcs7UnPadding(plainText []byte) ([]byte, error) {
	length := len(plainText)
	if length <= 0 {
		return nil, errors.New("plainText length = 0")
	}
	unpadding := int(plainText[length-1])
	if length < unpadding {
		return nil, errors.New("index out of range")
	}
	return plainText[:(length - unpadding)], nil
}

// DecryptMobile decrypts ciphertext using the package-level helper (deprecated: use Client.DecryptMobile).
func DecryptMobile(aesKey, ciphertext string) (string, error) {
	return decryptMobile(aesKey, ciphertext)
}
