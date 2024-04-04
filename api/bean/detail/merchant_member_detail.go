package detail

import (
	"context"
	"fmt"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMemberDetail struct {
	Id          uint64                         `json:"id"         description:"userId"`          // userId
	MerchantId  uint64                         `json:"merchantId" description:"merchant id"`     // merchant id
	Email       string                         `json:"email"      description:"email"`           // email
	FirstName   string                         `json:"firstName"  description:"first name"`      // first name
	LastName    string                         `json:"lastName"   description:"last name"`       // last name
	CreateTime  int64                          `json:"createTime" description:"create utc time"` // create utc time
	Mobile      string                         `json:"mobile"     description:"mobile"`          // mobile
	Role        string                         `json:"role"       description:"role"`            // role
	Permissions []*bean.MerchantRolePermission `json:"permissions" description:"permissions"`
}

func ConvertMemberToDetail(ctx context.Context, one *entity.MerchantMember) *MerchantMemberDetail {
	role := query.GetRoleByName(ctx, one.MerchantId, one.Role)
	var permissionData = make([]*bean.MerchantRolePermission, 0)
	if role != nil {
		err := utility.UnmarshalFromJsonString(role.PermissionData, &permissionData)
		if err != nil {
			fmt.Printf("ConvertRolePermissions err:%s", err)
		}
	}
	return &MerchantMemberDetail{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		Email:       one.Email,
		FirstName:   one.FirstName,
		LastName:    one.LastName,
		CreateTime:  one.CreateTime,
		Mobile:      one.Mobile,
		Role:        one.Role,
		Permissions: permissionData,
	}
}
