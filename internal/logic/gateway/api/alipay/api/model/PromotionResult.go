package model

type PromotionResult struct {
	PromotionType PromotionType `json:"promotionType,omitempty"`
	Discount      Discount      `json:"discount,omitempty"`
}
