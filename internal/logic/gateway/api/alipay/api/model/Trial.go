package model

type Trial struct {
	TrialStartPeriod int     `json:"trialStartPeriod,omitempty"`
	TrialAmount      *Amount `json:"trialAmount,omitempty"`
	TrialEndPeriod   int     `json:"trialEndPeriod,omitempty"`
}
