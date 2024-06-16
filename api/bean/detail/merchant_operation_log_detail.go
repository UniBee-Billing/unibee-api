package detail

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type MerchantOperationLogDetail struct {
	Id             uint64                `json:"id"                 description:"id"`                        // id
	MerchantId     uint64                `json:"merchantId"         description:"merchant Id"`               // merchant Id
	MemberId       uint64                `json:"memberId"           description:"member_id"`                 // member_id
	OptTarget      string                `json:"optTarget"          description:"operation target"`          // operation target
	OptContent     string                `json:"optContent"         description:"operation content"`         // operation content
	CreateTime     int64                 `json:"createTime"         description:"operation create utc time"` // operation create utc time
	SubscriptionId string                `json:"subscriptionId"     description:"subscription_id"`           // subscription_id
	UserId         uint64                `json:"userId"             description:"user_id"`                   // user_id
	InvoiceId      string                `json:"invoiceId"          description:"invoice id"`                // invoice id
	PlanId         uint64                `json:"planId"             description:"plan id"`                   // plan id
	DiscountCode   string                `json:"discountCode"       description:"discount_code"`             // discount_code
	Member         *MerchantMemberDetail `json:"member"             description:"Member"`
	//User           *bean.UserAccountSimplify  `json:"user"               description:"User"`
	//Subscription   *bean.SubscriptionSimplify `json:"subscription"       description:"Subscription"`
}

func ConvertOperationLogToDetail(ctx context.Context, one *entity.MerchantOperationLog) *MerchantOperationLogDetail {
	if one == nil {
		return nil
	}
	return &MerchantOperationLogDetail{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		MemberId:       one.MemberId,
		OptTarget:      one.OptTarget,
		OptContent:     one.OptContent,
		CreateTime:     one.CreateTime,
		SubscriptionId: one.SubscriptionId,
		UserId:         one.UserId,
		InvoiceId:      one.InvoiceId,
		PlanId:         one.PlanId,
		DiscountCode:   one.DiscountCode,
		Member:         ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, one.MemberId)),
		//User:           bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		//Subscription:   bean.SimplifySubscription(query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)),
	}
}
