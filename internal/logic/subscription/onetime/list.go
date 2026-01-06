package onetime

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type SubscriptionOnetimeAddonPurchaseListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId"`
	UserId     uint64 `json:"userId"  dc:"UserId" `
	Page       int    `json:"page" dc:"Page, Start With 0" `
	Count      int    `json:"count" dc:"Count Of Page" `
}

func SubscriptionOnetimeAddonPurchaseList(ctx context.Context, req *SubscriptionOnetimeAddonPurchaseListInternalReq) (list []*detail.SubscriptionOnetimeAddonDetail) {
	var mainList []*entity.SubscriptionOnetimeAddon
	if req.Count <= 0 {
		req.Count = 100
	}
	if req.Page < 0 {
		req.Page = 0
	}
	baseQuery := dao.SubscriptionOnetimeAddon.Ctx(ctx).
		Where(dao.SubscriptionOnetimeAddon.Columns().UserId, req.UserId).
		WhereIn(dao.SubscriptionOnetimeAddon.Columns().Status, []int{1, 2})
	err := baseQuery.Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil
	}
	for _, one := range mainList {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifySubscriptionOnetimeAddon Unmarshal Metadata error:%s", err.Error())
			}
		}
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		var invoiceId = ""
		if len(one.InvoiceId) > 0 {
			invoiceId = one.InvoiceId
		} else {
			if payment != nil {
				invoiceId = payment.InvoiceId
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
			Payment:        bean.SimplifyPayment(payment),
			Invoice:        bean.SimplifyInvoice(query.GetInvoiceByInvoiceId(ctx, invoiceId)),
			Metadata:       metadata,
		})
	}
	return list
}
