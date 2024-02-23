package ro

import entity "unibee-api/internal/model/entity/oversea_pay"

type UserAccountSimplify struct {
	Id             uint64 `json:"id"                 description:"userId"`             // userId
	MerchantId     uint64 `json:"merchantId"                description:"merchant id"` // merchant id
	Email          string `json:"email"              description:"email"`              // email
	Address        string `json:"address"            description:"address"`            // address
	CreateTime     int64  `json:"createTime"         description:"create utc time"`    // create utc time
	ExternalUserId string `json:"externalUserId"     description:"externalUserId"`     // external_user_id
	FirstName      string `json:"firstName"          description:"first name"`         // first name
	LastName       string `json:"lastName"           description:"last name"`          // last name
}

func SimplifyUserAccount(one *entity.UserAccount) *UserAccountSimplify {
	if one == nil {
		return nil
	}
	return &UserAccountSimplify{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		Email:          one.Email,
		Address:        one.Address,
		CreateTime:     one.CreateTime,
		ExternalUserId: one.ExternalUserId,
		FirstName:      one.FirstName,
		LastName:       one.LastName,
	}
}
