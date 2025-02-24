package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net/url"
	"unibee/internal/logic/gateway/api/alipay/api/exception"
)

func GenSign(httpMethod string, path string, clientId string, reqTime string, reqBody string, merchantPrivateKey string) (string, error) {
	block, _ := pem.Decode([]byte(getPkcsKey(merchantPrivateKey)))
	if block == nil {
		return "", &exception.AlipayLibraryError{Message: "Failed to decode private key"}
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", &exception.AlipayLibraryError{Message: "Failed to parse private key  " + err.Error()}
	}
	payload := genSignContent(httpMethod, path, clientId, reqTime, reqBody)
	signature, err := Sign(privateKey.(*rsa.PrivateKey), []byte(payload))
	if err != nil {
		return "", &exception.AlipayLibraryError{Message: "Failed to sign data  " + err.Error()}
	}

	return signature, nil

}

func Sign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	// 计算数据的SHA256哈希
	hashed := sha256.Sum256(data)

	// 使用私钥签名哈希值
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.Hash.HashFunc(crypto.SHA256), hashed[:])
	if err != nil {
		return "", err
	}
	//base64编码 如果signature直接传string会造成乱码
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	return Encode(base64Signature)
}

func Verify(httpMethod string, path string, clientId string, rspTimeStr string, rspBody string, signature string, alipayPublicKey string) (bool, error) {
	rspContent := genSignContent(httpMethod, path, clientId, rspTimeStr, rspBody)
	return verifySignatureWithSHA256RSA(rspContent, Decode(signature), alipayPublicKey)
}

func verifySignatureWithSHA256RSA(rspContent string, signature string, strPk string) (bool, error) {
	publicKey, err := getPublicKeyFromBase64String(strPk)
	if err != nil {
		return false, err
	}

	hash := sha256.New()
	hash.Write([]byte(rspContent))
	digest := hash.Sum(nil)

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, signatureBytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getPublicKeyFromBase64String(publicKeyString string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not an RSA public key")
	}
}

func genSignContent(httpMethod string, path string, clientId string, reqTime string, reqBody string) string {
	return httpMethod + " " + path + "\n" + clientId + "." + reqTime + "." + reqBody
}

func Encode(signature string) (string, error) {
	return url.QueryEscape(signature), nil
}

func Decode(originalStr string) string {
	unescape, err := url.QueryUnescape(originalStr)
	if err != nil {
		return ""
	}
	return unescape
}

func getPkcsKey(key string) string {
	return "-----BEGIN PRIVATE KEY-----\n" + key + "\n-----END PRIVATE KEY-----"
}
