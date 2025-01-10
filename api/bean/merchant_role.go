package bean

import (
	"fmt"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type MerchantRole struct {
	Id          uint64                    `json:"id"             description:"id"`          // id
	MerchantId  uint64                    `json:"merchantId"     description:"merchant id"` // merchant id
	Role        string                    `json:"role"           description:"role"`        // role
	Permissions []*MerchantRolePermission `json:"permissions" description:"permissions"`
	CreateTime  int64                     `json:"createTime"     description:"create utc time"` // create utc time
}

func removeByValue(arr []*MerchantRolePermission, value *MerchantRolePermission) []*MerchantRolePermission {
	for i, v := range arr {
		if v == value {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func SimplifyMerchantRole(one *entity.MerchantRole) *MerchantRole {
	if one == nil {
		return nil
	}
	var permissionData []*MerchantRolePermission
	err := utility.UnmarshalFromJsonString(one.PermissionData, &permissionData)
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	for _, permission := range permissionData {
		if len(permission.Permissions) == 0 {
			permissionData = removeByValue(permissionData, permission)
		}
	}
	return &MerchantRole{
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
