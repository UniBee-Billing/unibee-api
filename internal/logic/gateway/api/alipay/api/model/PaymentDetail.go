package model

type PaymentDetail struct {
	Amount        *Amount        `json:"amount,omitempty"`
	PaymentMethod *PaymentMethod `json:"paymentMethod,omitempty"`
}

type AuthorizationPhase string

const (
	AuthorizationPhase_PRE_AUTHORIZATION  AuthorizationPhase = "PRE_AUTHORIZATION"
	AuthorizationPhase_POST_AUTHORIZATION AuthorizationPhase = "POST_AUTHORIZATION"
)
