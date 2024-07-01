package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/logic/auth"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

type TaskUserImport struct {
}

func (t TaskUserImport) TaskName() string {
	return "UserImport"
}

func (t TaskUserImport) TemplateHeader() interface{} {
	return ImportUserEntity{}
}

func (t TaskUserImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, data map[string]string) (interface{}, error) {
	var err error
	target := &ImportUserEntity{
		FirstName:      fmt.Sprintf("%s", data["FirstName"]),
		LastName:       fmt.Sprintf("%s", data["LastName"]),
		Email:          fmt.Sprintf("%s", data["Email"]),
		Address:        fmt.Sprintf("%s", data["Address"]),
		VATNumber:      fmt.Sprintf("%s", data["VATNumber"]),
		CountryCode:    fmt.Sprintf("%s", data["CountryCode"]),
		ExternalUserId: fmt.Sprintf("%s", data["ExternalUserId"]),
		TaxPercentage:  fmt.Sprintf("%s", data["TaxPercentage"]),
		Type:           fmt.Sprintf("%s", data["Type"]),
	}
	one := query.GetUserAccountByEmail(ctx, task.MerchantId, target.Email)
	if one != nil {
		return target, gerror.New("Skip, same email user exist")
	}
	if len(target.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, task.MerchantId, target.ExternalUserId)
	}
	if one != nil {
		return target, gerror.New("Skip, same ExternalUserId user exist")
	}
	// todo mark
	_, err = auth.CreateUser(ctx, &auth.NewReq{
		ExternalUserId: target.ExternalUserId,
		Email:          target.Email,
		FirstName:      target.FirstName,
		LastName:       target.LastName,
		Address:        target.Address,
		CountryCode:    target.CountryCode,
		Custom:         "Import",
		MerchantId:     test.TestMerchant.Id,
	})
	return target, err
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
