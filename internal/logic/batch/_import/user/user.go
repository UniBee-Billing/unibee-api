package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
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
	return &ImportUserEntity{
		Email:          "accounts.unibee@unibee.dev",
		ExternalUserId: "example_id",
		FirstName:      "unibee",
		LastName:       "unibee",
		Address:        "Harju maakond, Tallinn, Pirita linnaosa, Supluse pst 1-201a, 11911",
		Phone:          "5586977",
	}
}

func (t TaskUserImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportUserEntity{
		Email:          fmt.Sprintf("%s", row["Email"]),
		ExternalUserId: fmt.Sprintf("%s", row["ExternalUserId"]),
		FirstName:      fmt.Sprintf("%s", row["FirstName"]),
		LastName:       fmt.Sprintf("%s", row["LastName"]),
		Address:        fmt.Sprintf("%s", row["Address"]),
		//VATNumber:      fmt.Sprintf("%s", row["VATNumber"]),
		//CountryCode:    fmt.Sprintf("%s", row["CountryCode"]),
		Phone: fmt.Sprintf("%s", row["Phone"]),
		//TaxPercentage:  fmt.Sprintf("%s", row["TaxPercentage"]),
		//Type:           fmt.Sprintf("%s", row["Type"]),
	}
	if len(target.Email) == 0 {
		return target, gerror.New("Error, Email is blank")
	}
	one := query.GetUserAccountByEmail(ctx, task.MerchantId, target.Email)
	if one != nil {
		//enter update process
		if len(target.ExternalUserId) > 0 {
			otherOne := query.GetUserAccountByExternalUserId(ctx, task.MerchantId, target.ExternalUserId)
			if otherOne != nil && one.Id != otherOne.Id {
				return target, gerror.New("Error, same ExternalUserId user exist")
			}
		}
		_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().ExternalUserId: target.ExternalUserId,
			dao.UserAccount.Columns().FirstName:      target.FirstName,
			dao.UserAccount.Columns().LastName:       target.LastName,
			dao.UserAccount.Columns().Address:        target.Address,
			dao.UserAccount.Columns().Phone:          target.Phone,
			dao.UserAccount.Columns().GmtModify:      gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
		return target, err
	}
	if len(target.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, task.MerchantId, target.ExternalUserId)
	}
	if one != nil {
		return target, gerror.New("Error, same ExternalUserId user exist")
	}
	_, err = auth.CreateUser(ctx, &auth.NewReq{
		ExternalUserId: target.ExternalUserId,
		Email:          target.Email,
		FirstName:      target.FirstName,
		LastName:       target.LastName,
		Address:        target.Address,
		Phone:          target.Phone,
		Custom:         "Import",
		MerchantId:     task.MerchantId,
	})
	return target, err
}

type ImportUserEntity struct {
	Email          string `json:"Email"              `
	ExternalUserId string `json:"ExternalUserId"     `
	FirstName      string `json:"FirstName"          `
	LastName       string `json:"LastName"           `
	Address        string `json:"Address"            `
	Phone          string `json:"Phone"            `
	//VATNumber      string `json:"VATNumber"          `
	//CountryCode    string `json:"CountryCode"        `
	//TaxPercentage  string `json:"TaxPercentage"      `
	//Type           string `json:"Type"               `
}
