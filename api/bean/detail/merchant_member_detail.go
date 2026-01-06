package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/totp/client_activity"
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
	TotpType              int                                     `json:"totpType"   description:"0-Inactive, 1-General, Google Authenticator | 2-Microsoft Authenticator | 3-Authy | 4-1Password | 5-LastPass | 6-FreeOTP | 7-Other TOTP"`
	DeviceList            []*bean.MerchantMemberDevice            `json:"deviceList" description:"The devices list'"`
	MemberRoles           []*bean.MerchantRole                    `json:"MemberRoles" description:"The member role list'" `
	CurrentDeviceIdentity string                                  `json:"currentDeviceIdentity" description:"The Current DeviceIdentity'" `
	MemberGroupPermission map[string]*bean.MerchantRolePermission `json:"MemberGroupPermission" description:"The member group permission map'"`
	OAuthAccounts         []*bean.Oauth                           `json:"oauthAccounts" description:"List of connected OAuth accounts"`
}

func ConvertMemberToDetail(ctx context.Context, one *entity.MerchantMember) *MerchantMemberDetail {
	if ctx == nil || one == nil {
		return nil
	}
	isOwner, memberRoles := ConvertMemberRole(ctx, one)
	_, memberGroupPermission := ConvertMemberGroupPermissions(ctx, one)
	var currentDeviceIdentity string
	if _interface.Context() != nil &&
		_interface.Context().Get(ctx) != nil {
		currentDeviceIdentity = _interface.Context().Get(ctx).ClientIdentity
	}
	var identityData = &bean.OauthIdentity{
		OAuthAccountMap: make(map[string]*bean.Oauth),
	}
	_ = utility.UnmarshalFromJsonString(one.AuthJs, &identityData)
	oauthAccounts := make([]*bean.Oauth, 0)
	for _, account := range identityData.OAuthAccountMap {
		oauthAccounts = append(oauthAccounts, account)
	}
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
		TotpType:              one.TotpValidatorType,
		IsBlankPasswd:         len(one.Password) == 0,
		Status:                one.Status,
		MemberGroupPermission: memberGroupPermission,
		CurrentDeviceIdentity: currentDeviceIdentity,
		DeviceList:            GetClientIdentityList(ctx, one.Email),
		OAuthAccounts:         oauthAccounts,
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
								permissions = append(permissions, &bean.MerchantRolePermission{
									Group:       permission.Group,
									Permissions: permission.Permissions,
								})
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
			//for _, permission := range permissionGroupMap {
			//	permissions = append(permissions, &bean.MerchantRolePermission{
			//		Group:       permission.Group,
			//		Permissions: permission.Permissions,
			//	})
			//}
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

const MemberDeviceExpireDays = 30

func GetClientIdentityListCacheKey(email string) string {
	return fmt.Sprintf("Merchant#Totp#client#identity#list#%s", email)
}

func GetClientIdentityList(ctx context.Context, email string) []*bean.MerchantMemberDevice {
	var clientIdentityMap = map[string]int64{}
	data, err := g.Redis().Get(ctx, GetClientIdentityListCacheKey(email))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityMap)
	}
	var devices = make([]*bean.MerchantMemberDevice, 0)

	for key, value := range clientIdentityMap {
		var isCurrentDevice = false
		var name string
		if _interface.Context() != nil &&
			_interface.Context().Get(ctx) != nil &&
			_interface.Context().Get(ctx).ClientIdentity != "" &&
			_interface.Context().Get(ctx).ClientIdentity == key {
			if _interface.Context().Get(ctx).MerchantMember != nil && _interface.Context().Get(ctx).MerchantMember.Email == email {
				isCurrentDevice = true
			}
		}
		var lastLoginTime = value
		var lastActiveTime = value
		activity := client_activity.GetClientIdentityActivity(ctx, key)
		if activity != nil {
			if activity.LastActivityTime > 0 {
				lastLoginTime = activity.LastActivityTime
			}
			if activity.LastLoginTime > 0 {
				lastLoginTime = activity.LastLoginTime
			}
		}

		deviceTypes := strings.Split(key, "_")
		if len(deviceTypes) > 2 {
			name = fmt.Sprintf("%s_%s", deviceTypes[1], deviceTypes[2])
		} else {
			name = utility.ExtractBrowserOS(key)
		}
		devices = append(devices, &bean.MerchantMemberDevice{
			Name:                     name,
			Identity:                 key,
			LastLoginTime:            lastLoginTime,
			LastActiveTime:           lastActiveTime,
			LastTotpVerificationTime: value,
			Status:                   gtime.Now().Timestamp()-value < MemberDeviceExpireDays*24*60*60,
			IPAddress:                utility.ExtractFirstIPAddresses(key),
			CurrentDevice:            isCurrentDevice,
		})
	}
	return devices
}
