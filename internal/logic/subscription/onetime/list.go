package onetime

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type SubscriptionOnetimeAddonListInternalReq struct {
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId"`
	SubscriptionId string `json:"subscriptionId"  dc:"SubscriptionId" `
	Page           int    `json:"page" dc:"Page, Start With 0" `
	Count          int    `json:"count" dc:"Count Of Page" `
}

func SubscriptionOnetimeAddonList(ctx context.Context, req *SubscriptionOnetimeAddonListInternalReq) (list []*detail.SubscriptionOnetimeAddonDetail) {
	var mainList []*entity.SubscriptionOnetimeAddon
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	baseQuery := dao.SubscriptionOnetimeAddon.Ctx(ctx).
		Where(dao.SubscriptionOnetimeAddon.Columns().SubscriptionId, req.SubscriptionId).WhereIn(dao.Subscription.Columns().Status, []int{1, 2})
	err := baseQuery.Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	for _, one := range mainList {
		var metadata = make(map[string]string)
		if len(one.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifySubscriptionOnetimeAddon Unmarshal Metadata error:%s", err.Error())
			}
		}
		list = append(list, &detail.SubscriptionOnetimeAddonDetail{
			Id:             one.Id,
			SubscriptionId: one.SubscriptionId,
			AddonId:        one.AddonId,
			Addon:          bean.SimplifyPlan(query.GetPlanById(ctx, one.AddonId)),
			Quantity:       one.Quantity,
			Status:         one.Status,
			CreateTime:     one.CreateTime,
			Payment:        bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, one.PaymentId)),
			Metadata:       metadata,
		})
	}
	return list
}