package member

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	go_redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/default"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMember struct {
	Id          uint64                    `json:"id"         description:"userId"`                // userId
	GmtCreate   *gtime.Time               `json:"gmtCreate"  description:"create time"`           // create time
	GmtModify   *gtime.Time               `json:"gmtModify"  description:"update time"`           // update time
	MerchantId  uint64                    `json:"merchantId" description:"merchant id"`           // merchant id
	IsDeleted   int                       `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Password    string                    `json:"password"   description:"password"`              // password
	UserName    string                    `json:"userName"   description:"user name"`             // user name
	Mobile      string                    `json:"mobile"     description:"mobile"`                // mobile
	Email       string                    `json:"email"      description:"email"`                 // email
	FirstName   string                    `json:"firstName"  description:"first name"`            // first name
	LastName    string                    `json:"lastName"   description:"last name"`             // last name
	CreateTime  int64                     `json:"createTime" description:"create utc time"`       // create utc time
	Role        string                    `json:"role"       description:"role"`                  // role
	Status      int                       `json:"status"     description:"0-Active, 2-Suspend"`   // 0-Active, 2-Suspend
	Permissions []*MerchantRolePermission `json:"permissions"       description:"Permissions"`    // Permissions
	IsOwner     bool                      `json:"isOwner"       description:"IsOwner"`            // role
}

type MerchantRolePermission struct {
	Group       string   `json:"group"           description:"Group"`             // group
	Permissions []string `json:"permissions"           description:"Permissions"` // group
}

func ReloadAllMembersCacheForSDKAuthBackground() {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		var list []*MerchantMember
		err = dao.MerchantMember.Ctx(ctx).
			Scan(&list)
		for _, member := range list {
			isOwner, permissions := ConvertMemberPermissions(ctx, member)
			member.IsOwner = isOwner
			member.Permissions = permissions
			_, _ = g.Redis().Set(ctx, fmt.Sprintf("UniBee#Member#%d", member.Id), utility.MarshalToJsonString(member))
			if isOwner {
				_, _ = g.Redis().Set(ctx, fmt.Sprintf("UniBee#Merchant#Owner#%d", member.MerchantId), utility.MarshalToJsonString(member))
				go func() {
					_ = go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
						Group:   "GID_UniBee_License",
						Method:  "GetLicenseByMerchantId",
						Request: member.MerchantId,
					}, 0)
				}()
			}
		}
	}()
}

func ReloadMemberCacheForSdkAuthBackground(id uint64) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		var member *MerchantMember
		err = dao.MerchantMember.Ctx(ctx).
			Where(dao.MerchantMember.Columns().Id, id).
			Scan(&member)
		if member != nil {
			isOwner, permissions := ConvertMemberPermissions(ctx, member)
			member.IsOwner = isOwner
			member.Permissions = permissions
			_, _ = g.Redis().Set(ctx, fmt.Sprintf("UniBee#Member#%d", member.Id), utility.MarshalToJsonString(member))
			if isOwner {
				_, _ = g.Redis().Set(ctx, fmt.Sprintf("UniBee#Merchant#Owner#%d", member.MerchantId), utility.MarshalToJsonString(member))
				go func() {
					_ = go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
						Group:   "GID_UniBee_License",
						Method:  "GetLicenseByMerchantId",
						Request: member.MerchantId,
					}, 0)
				}()
			}
		}
	}()
}

func ConvertMemberPermissions(ctx context.Context, member *MerchantMember) (isOwner bool, permissions []*MerchantRolePermission) {
	permissions = make([]*MerchantRolePermission, 0)
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
								permissions = append(permissions, &MerchantRolePermission{
									Group:       permission.Group,
									Permissions: permission.Permissions,
								})
							}
						}
					}
				}
			}
		}
	}
	return isOwner, permissions
}
