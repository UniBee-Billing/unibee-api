package model

type PeriodRule struct {
	PeriodType  PeriodType `json:"periodType,omitempty"`
	PeriodCount int        `json:"periodCount,omitempty"`
}
