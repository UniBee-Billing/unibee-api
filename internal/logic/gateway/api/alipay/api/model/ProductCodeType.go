package model

type ProductCodeType string

const (
	AGREEMENT_PAYMENT ProductCodeType = "AGREEMENT_PAYMENT"
	IN_STORE_PAYMENT  ProductCodeType = "IN_STORE_PAYMENT"
	CASHIER_PAYMENT   ProductCodeType = "CASHIER_PAYMENT"
)
