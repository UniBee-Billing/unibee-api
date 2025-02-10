package detail

import (
	"context"
	"strings"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMemberDetail struct {
	Id                    uint64                                  `json:"id"         description:"userId"`          // userId
	MerchantId            uint64                                  `json:"merchantId" description:"merchant id"`     // merchant id
	Email                 string                                  `json:"email"      description:"email"`           // email
	FirstName             string                                  `json:"firstName"  description:"first name"`      // first name
	LastName              string                                  `json:"lastName"   description:"last name"`       // last name
	CreateTime            int64                                   `json:"createTime" description:"create utc time"` // create utc time
	Mobile                string                                  `json:"mobile"     description:"mobile"`          // mobile
	IsOwner               bool                                    `json:"isOwner" description:"Check Member is Owner" `
	Status                int                                     `json:"status"             description:"0-Active, 2-Suspend"`
	IsBlankPasswd         bool                                    `json:"isBlankPasswd" description:"is blank password"`
	MemberRoles           []*bean.MerchantRole                    `json:"MemberRoles" description:"The member role list'" `
	MemberGroupPermission map[string]*bean.MerchantRolePermission `json:"MemberGroupPermission" description:"The member group permission map'"`
}

func ConvertMemberToDetail(ctx context.Context, one *entity.MerchantMember) *MerchantMemberDetail {
	if ctx == nil || one == nil {
		return nil
	}
	isOwner, memberRoles := ConvertMemberRole(ctx, one)
	_, memberGroupPermission := ConvertMemberGroupPermissions(ctx, one)
	return &MerchantMemberDetail{
		Id:                    one.Id,
		MerchantId:            one.MerchantId,
		Email:                 one.Email,
		FirstName:             one.FirstName,
		LastName:              one.LastName,
		CreateTime:            one.CreateTime,
		Mobile:                one.Mobile,
		IsOwner:               isOwner,
		MemberRoles:           memberRoles,
		IsBlankPasswd:         len(one.Password) == 0,
		Status:                one.Status,
		MemberGroupPermission: memberGroupPermission,
	}
}

func ConvertMemberRole(ctx context.Context, member *entity.MerchantMember) (isOwner bool, memberRoles []*bean.MerchantRole) {
	memberRoles = make([]*bean.MerchantRole, 0)
	if member != nil {
		if strings.Contains(member.Role, "Owner") {
			isOwner = true
		} else {
			var roleIdList = make([]uint64, 0)
			_ = utility.UnmarshalFromJsonString(member.Role, &roleIdList)
			for _, roleId := range roleIdList {
				if roleId > 0 {
					role := query.GetRoleById(ctx, roleId)
					if role != nil {
						memberRoles = append(memberRoles, bean.SimplifyMerchantRole(role))
					}
				}
			}
		}
	}
	return isOwner, memberRoles
}

func ConvertMemberPermissions(ctx context.Context, member *entity.MerchantMember) (isOwner bool, permissions []*bean.MerchantRolePermission, groupPermissionMap map[string]*bean.MerchantRolePermission) {
	permissions = make([]*bean.MerchantRolePermission, 0)
	permissionGroupMap := make(map[string]*bean.MerchantRolePermission)
	if member != nil {
		if strings.Contains(member.Role, "Owner") {
			isOwner = true
		} else {
			var roleIdList = make([]uint64, 0)
			_ = utility.UnmarshalFromJsonString(member.Role, &roleIdList)
			for _, roleId := range roleIdList {
				if roleId > 0 {
					role := query.GetRoleById(ctx, roleId)
					if role != nil {
						roleDetail := bean.SimplifyMerchantRole(role)
						if roleDetail != nil {
							for _, permission := range roleDetail.Permissions {
								permissions = append(permissions, permission)
								if groupPermission, ok := permissionGroupMap[permission.Group]; ok {
									for _, p := range permission.Permissions {
										if groupPermission.Permissions == nil {
											groupPermission.Permissions = make([]string, 0)
										}
										if len(p) > 0 && !utility.IsStringInArray(groupPermission.Permissions, p) {
											groupPermission.Permissions = append(groupPermission.Permissions, p)
										}
									}
								} else {
									permissionGroupMap[permission.Group] = permission
								}
							}
						}
					}
				}
			}
		}
	}
	return isOwner, permissions, permissionGroupMap
}

func ConvertMemberGroupPermissions(ctx context.Context, member *entity.MerchantMember) (isOwner bool, groupPermissionMap map[string]*bean.MerchantRolePermission) {
	permissionGroupMap := make(map[string]*bean.MerchantRolePermission)
	if member != nil {
		if strings.Contains(member.Role, "Owner") {
			isOwner = true
		} else {
			var roleIdList = make([]uint64, 0)
			_ = utility.UnmarshalFromJsonString(member.Role, &roleIdList)
			for _, roleId := range roleIdList {
				if roleId > 0 {
					role := query.GetRoleById(ctx, roleId)
					if role != nil {
						roleDetail := bean.SimplifyMerchantRole(role)
						if roleDetail != nil {
							for _, permission := range roleDetail.Permissions {
								if groupPermission, ok := permissionGroupMap[permission.Group]; ok {
									for _, p := range permission.Permissions {
										if groupPermission.Permissions == nil {
											groupPermission.Permissions = make([]string, 0)
										}
										if len(p) > 0 && !utility.IsStringInArray(groupPermission.Permissions, p) {
											groupPermission.Permissions = append(groupPermission.Permissions, p)
										}
									}
								} else {
									permissionGroupMap[permission.Group] = permission
								}
							}
						}
					}
				}
			}
		}
	}
	return isOwner, permissionGroupMap
}
