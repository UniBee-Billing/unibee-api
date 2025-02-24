package model

type RedirectActionForm struct {
	Method         string `json:"method,omitempty"`
	Parameters     string `json:"parameters,omitempty"`
	RedirectURL    string `json:"redirectUrl,omitempty"`
	ActionFormType string `json:"actionFormType,omitempty"`
}
