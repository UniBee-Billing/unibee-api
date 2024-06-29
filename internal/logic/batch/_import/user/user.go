package user

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type TaskUserImport struct {
}

func (t TaskUserImport) TaskName() string {
	return "UserImport"
}

func (t TaskUserImport) TemplateHeader() interface{} {
	return ImportUserEntity{}
}

func (t TaskUserImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, data interface{}) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

type ImportUserEntity struct {
	FirstName      string `json:"FirstName"          `
	LastName       string `json:"LastName"           `
	Email          string `json:"Email"              `
	Address        string `json:"Address"            `
	VATNumber      string `json:"VATNumber"          `
	CountryCode    string `json:"CountryCode"        `
	ExternalUserId string `json:"ExternalUserId"     `
	TaxPercentage  string `json:"TaxPercentage"      `
	Type           string `json:"Type"               `
}
