package model

type PaymentVerificationData struct {
	VerifyRequestId    string `json:"verifyRequestId,omitempty"`
	AuthenticationCode string `json:"authenticationCode,omitempty"`
}
