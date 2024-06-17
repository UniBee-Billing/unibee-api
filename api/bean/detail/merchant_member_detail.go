package detail

import (
	"context"
	"strings"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMemberDetail struct {
	Id          uint64                       `json:"id"         description:"userId"`          // userId
	MerchantId  uint64                       `json:"merchantId" description:"merchant id"`     // merchant id
	Email       string                       `json:"email"      description:"email"`           // email
	FirstName   string                       `json:"firstName"  description:"first name"`      // first name
	LastName    string                       `json:"lastName"   description:"last name"`       // last name
	CreateTime  int64                        `json:"createTime" description:"create utc time"` // create utc time
	Mobile      string                       `json:"mobile"     description:"mobile"`          // mobile
	IsOwner     bool                         `json:"isOwner" description:"Check Member is Owner" `
	MemberRoles []*bean.MerchantRoleSimplify `json:"MemberRoles" description:"The member's role list'" `
}

func ConvertMemberToDetail(ctx context.Context, one *entity.MerchantMember) *MerchantMemberDetail {
	isOwner, memberRoles := ConvertMemberRole(ctx, one)
	return &MerchantMemberDetail{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		Email:       one.Email,
		FirstName:   one.FirstName,
		LastName:    one.LastName,
		CreateTime:  one.CreateTime,
		Mobile:      one.Mobile,
		IsOwner:     isOwner,
		MemberRoles: memberRoles,
	}
}

func ConvertMemberRole(ctx context.Context, member *entity.MerchantMember) (isOwner bool, memberRoles []*bean.MerchantRoleSimplify) {
	memberRoles = make([]*bean.MerchantRoleSimplify, 0)
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
