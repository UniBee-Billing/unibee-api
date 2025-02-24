package model

type SettlementDetail struct {
	SettleTo         SettleToType `json:"settleTo,omitempty"`
	SettlementAmount *Amount      `json:"settlementAmount,omitempty"`
}

type SettleToType string

const (
	SettleToType_SELLER      SettleToType = "SELLER"
	SettleToType_MARKETPLACE SettleToType = "MARKETPLACE"
)
