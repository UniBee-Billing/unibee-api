package model

type AuthCodeForm struct {
	CodeDetails []CodeDetail `json:"codeDetails,omitempty"`
}
