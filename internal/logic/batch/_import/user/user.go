package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/logic/auth"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type TaskUserImport struct {
}

func (t TaskUserImport) TaskName() string {
	return "UserImport"
}

func (t TaskUserImport) TemplateHeader() interface{} {
	return ImportUserEntity{}
}

func (t TaskUserImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportUserEntity{
		FirstName:      fmt.Sprintf("%s", row["FirstName"]),
		LastName:       fmt.Sprintf("%s", row["LastName"]),
		Email:          fmt.Sprintf("%s", row["Email"]),
		Address:        fmt.Sprintf("%s", row["Address"]),
		VATNumber:      fmt.Sprintf("%s", row["VATNumber"]),
		CountryCode:    fmt.Sprintf("%s", row["CountryCode"]),
		ExternalUserId: fmt.Sprintf("%s", row["ExternalUserId"]),
		TaxPercentage:  fmt.Sprintf("%s", row["TaxPercentage"]),
		Type:           fmt.Sprintf("%s", row["Type"]),
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
		MerchantId:     task.MerchantId,
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
