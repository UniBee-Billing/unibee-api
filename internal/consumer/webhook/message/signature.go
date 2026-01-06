package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// VerifyHMACSignature Verify HMAC-SHA256 Signature（Base64）
func VerifyHMACSignature(body string, secret string, signatureTarget string) error {
	if signatureTarget == "" {
		return errors.New("missing X-Signature header")
	}
	if secret == "" {
		return errors.New("missing secret key")
	}
	if body == "" {
		return errors.New("missing body")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signatureTarget)) {
		return errors.New("invalid HMAC signatureTarget")
	}
	return nil
}

func SignHMACWebhook(body string, secret string) (signature string, algorithm string) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, "hmac"
}
