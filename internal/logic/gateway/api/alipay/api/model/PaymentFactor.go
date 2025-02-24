package model

type InStorePaymentScenario string

const (
	PaymentCode InStorePaymentScenario = "PaymentCode"
	OrderCode   InStorePaymentScenario = "OrderCode"
	EntryCode   InStorePaymentScenario = "EntryCode"
)

type PresentmentMode string

const (
	BUNDLE  PresentmentMode = "BUNDLE"
	TILE    PresentmentMode = "TILE"
	UNIFIED PresentmentMode = "UNIFIED"
)

type PaymentFactor struct {
	IsPaymentEvaluation    bool                   `json:"isPaymentEvaluation,omitempty,omitempty"`
	InStorePaymentScenario InStorePaymentScenario `json:"inStorePaymentScenario,omitempty,omitempty"`
	PresentmentMode        PresentmentMode        `json:"presentmentMode,omitempty,omitempty"`
	CaptureMode            string                 `json:"captureMode,omitempty,omitempty"`
	IsAuthorization        bool                   `json:"isAuthorization,omitempty,omitempty"`
}
