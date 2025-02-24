package model

type CreditPayPlan struct {
	InstallmentNum   int              `json:"installmentNum,omitempty"`
	Interval         string           `json:"interval,omitempty"`
	CreditPayFeeType CreditPayFeeType `json:"creditPayFeeType,omitempty"`
	FeePercentage    int              `json:"feePercentage,omitempty"`
}

type CreditPayFeeType string

const (
	PERCENTAGE CreditPayFeeType = "PERCENTAGE"
)
