package bean

import (
	"fmt"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type MerchantRoleSimplify struct {
	Id          uint64                    `json:"id"             description:"id"`          // id
	MerchantId  uint64                    `json:"merchantId"     description:"merchant id"` // merchant id
	Role        string                    `json:"role"           description:"role"`        // role
	Permissions []*MerchantRolePermission `json:"permissions" description:"permissions"`
	CreateTime  int64                     `json:"createTime"     description:"create utc time"` // create utc time
}

func SimplifyMerchantRole(one *entity.MerchantRole) *MerchantRoleSimplify {
	if one == nil {
		return nil
	}
	var permissionData []*MerchantRolePermission
	err := utility.UnmarshalFromJsonString(one.PermissionData, &permissionData)
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	return &MerchantRoleSimplify{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		CreateTime:  one.CreateTime,
		Permissions: permissionData,
		Role:        one.Role,
	}
}

type MerchantRolePermission struct {
	Group       string   `json:"group"           description:"Group"`             // group
	Permissions []string `json:"permissions"           description:"Permissions"` // group
}
