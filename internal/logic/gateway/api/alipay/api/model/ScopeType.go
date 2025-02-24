package model

type ScopeType string

const (
	//BASE_USER_INFO, AGREEMENT_PAY, USER_INFO, USER_LOGIN_ID, HASH_LOGIN_ID, SEND_OTP;
	ScopeTypeBaseUserInfo ScopeType = "BASE_USER_INFO"
	ScopeTypeAgreementPay ScopeType = "AGREEMENT_PAY"
	ScopeTypeUserInfo     ScopeType = "USER_INFO"
	ScopeTypeUserLoginId  ScopeType = "USER_LOGIN_ID"
	ScopeTypeHashLoginId  ScopeType = "HASH_LOGIN_ID"
	ScopeTypeSendOtp      ScopeType = "SEND_OTP"
	ScopeTypeTAOBAOREBIND ScopeType = "TAOBAO_REBIND"
)
