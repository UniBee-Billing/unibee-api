package model

type PspCustomerInfo struct {
	PspName             string `json:"pspName,omitempty"`
	PspCustomerId       string `json:"pspCustomerId,omitempty"`
	DisplayCustomerId   string `json:"displayCustomerId,omitempty"`
	DisplayCustomerName string `json:"displayCustomerName,omitempty"`
	Customer2088Id      string `json:"customer2088Id,omitempty"`
	ExtendInfo          string `json:"extendInfo,omitempty"`
}
