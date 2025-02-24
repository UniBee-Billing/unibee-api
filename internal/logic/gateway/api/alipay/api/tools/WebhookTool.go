package tools

import (
	"errors"
	"strings"
)

func CheckSignature(requestUri, httpMethod, clientId, requestTime, responseBody, signature, alipayPublicKey string) (bool, error) {
	realSignature := ""

	// Check if signature is nil or empty
	if signature == "" {
		return false, errors.New("empty notify signature")
	}

	// Get valid part from raw signature
	parts := strings.Split(signature, "signature=")
	if len(parts) > 1 {
		realSignature = parts[1]
	}

	// Verify signature
	isValid, _ := Verify(httpMethod, requestUri, clientId, requestTime, responseBody, realSignature, alipayPublicKey)
	if !isValid {
		return false, errors.New("signature verification failed")
	}

	return true, nil
}
