package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	"unibee/internal/logic/auth"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskUser struct {
}

func (t TaskUser) TaskName() string {
	return "CustomerExport"
}

func (t TaskUser) Header() interface{} {
	return ExportUserEntity{}
}

func (t TaskUser) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	var mainList = make([]interface{}, 0)
	if task == nil && task.MerchantId <= 0 {
		return mainList, nil
	}
	merchant := query.GetMerchantById(ctx, task.MerchantId)
	result, _ := auth.UserList(ctx, &auth.UserListInternalReq{
		MerchantId: task.MerchantId,
		//UserId:        0,
		//Email:         "",
		//FirstName:     "",
		//LastName:      "",
		//Status:        nil,
		//DeleteInclude: false,
		//SortField:     "",
		//SortType:      "",
		Page:  page,
		Count: count,
	})
	if result != nil && result.UserAccounts != nil {
		for _, one := range result.UserAccounts {
			var userGateway = ""
			if one.Gateway != nil {
				userGateway = one.Gateway.GatewayName
			}
			mainList = append(mainList, &ExportUserEntity{
				Id:                 fmt.Sprintf("%v", one.Id),
				FirstName:          one.FirstName,
				LastName:           one.LastName,
				Email:              one.Email,
				MerchantName:       merchant.Name,
				AvatarUrl:          one.AvatarUrl,
				Phone:              one.Phone,
				Address:            one.Address,
				VATNumber:          one.VATNumber,
				CountryCode:        one.CountryCode,
				CountryName:        one.CountryName,
				SubscriptionName:   one.SubscriptionName,
				SubscriptionId:     one.SubscriptionId,
				SubscriptionStatus: consts.SubStatusToEnum(one.SubscriptionStatus).Description(),
				CreateTime:         gtime.NewFromTimeStamp(one.CreateTime),
				ExternalUserId:     one.ExternalUserId,
				Status:             consts.UserStatusToEnum(one.Status).Description(),
				TaxPercentage:      utility.ConvertTaxPercentageToPercentageString(one.TaxPercentage),
				Type:               consts.UserTypeToEnum(one.Type).Description(),
				Gateway:            userGateway,
				City:               one.City,
				ZipCode:            one.ZipCode,
			})
		}
	}
	return mainList, nil
}

type ExportUserEntity struct {
	Id                 string      `json:"Id"                 `
	FirstName          string      `json:"FirstName"          `
	LastName           string      `json:"LastName"           `
	Email              string      `json:"Email"              `
	MerchantName       string      `json:"MerchantName"       `
	AvatarUrl          string      `json:"AvatarUrl"          `
	Phone              string      `json:"Phone"              `
	Address            string      `json:"Address"            `
	VATNumber          string      `json:"VATNumber"          `
	CountryCode        string      `json:"CountryCode"        `
	CountryName        string      `json:"CountryName"        `
	SubscriptionName   string      `json:"SubscriptionName"   `
	SubscriptionId     string      `json:"SubscriptionId"     `
	SubscriptionStatus string      `json:"SubscriptionStatus" `
	CreateTime         *gtime.Time `json:"CreateTime"       layout:"2006-01-02 15:04:05"  `
	ExternalUserId     string      `json:"ExternalUserId"     `
	Status             string      `json:"Status"             `
	TaxPercentage      string      `json:"TaxPercentage"      `
	Type               string      `json:"Type"               `
	Gateway            string      `json:"Gateway"            `
	City               string      `json:"City"               `
	ZipCode            string      `json:"ZipCode"            `
}