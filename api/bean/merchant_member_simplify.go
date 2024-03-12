package bean

import entity "unibee/internal/model/entity/oversea_pay"

type MerchantMemberSimplify struct {
	Id         uint64 `json:"id"         description:"userId"`          // userId
	MerchantId uint64 `json:"merchantId" description:"merchant id"`     // merchant id
	Email      string `json:"email"      description:"email"`           // email
	FirstName  string `json:"firstName"  description:"first name"`      // first name
	LastName   string `json:"lastName"   description:"last name"`       // last name
	CreateTime int64  `json:"createTime" description:"create utc time"` // create utc time
}

func SimplifyMerchantMember(one *entity.MerchantMember) *MerchantMemberSimplify {
	if one == nil {
		return nil
	}
	return &MerchantMemberSimplify{
		Id:         one.Id,
		MerchantId: one.MerchantId,
		Email:      one.Email,
		CreateTime: one.CreateTime,
		FirstName:  one.FirstName,
		LastName:   one.LastName,
	}
}
