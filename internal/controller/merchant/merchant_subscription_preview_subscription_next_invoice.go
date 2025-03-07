package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/billingcycle/cycle"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) PreviewSubscriptionNextInvoice(ctx context.Context, req *subscription.PreviewSubscriptionNextInvoiceReq) (res *subscription.PreviewSubscriptionNextInvoiceRes, err error) {
	utility.Assert(len(req.SubscriptionId) > 0, "Invalid SubscriptionId")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "sub not found")
	utility.Assert(sub.MerchantId == _interface.GetMerchantId(ctx), "no permission")
	invoice, one := cycle.PreviewSubscriptionNextInvoice(ctx, sub, utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd))
	var pendingUpdateDetail *detail.SubscriptionPendingUpdateDetail
	if one != nil {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("GetUnfinishedSubscriptionPendingUpdateDetailByPendingUpdateId Unmarshal Metadata error:%s", err.Error())
			}
		}
		pendingUpdateDetail = &detail.SubscriptionPendingUpdateDetail{
			MerchantId:      one.MerchantId,
			SubscriptionId:  one.SubscriptionId,
			PendingUpdateId: one.PendingUpdateId,
			GmtCreate:       one.GmtCreate,
			Amount:          one.Amount,
			Status:          one.Status,
			UpdateAmount:    one.UpdateAmount,
			Currency:        one.Currency,
			UpdateCurrency:  one.UpdateCurrency,
			PlanId:          one.PlanId,
			UpdatePlanId:    one.UpdatePlanId,
			Quantity:        one.Quantity,
			UpdateQuantity:  one.UpdateQuantity,
			AddonData:       one.AddonData,
			UpdateAddonData: one.UpdateAddonData,
			ProrationAmount: one.ProrationAmount,
			GatewayId:       one.GatewayId,
			UserId:          one.UserId,
			InvoiceId:       one.InvoiceId,
			GmtModify:       one.GmtModify,
			Paid:            one.Paid,
			Link:            one.Link,
			MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate: one.EffectImmediate,
			EffectTime:      one.EffectTime,
			Note:            one.Note,
			Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:        metadata,
		}
	}
	return &subscription.PreviewSubscriptionNextInvoiceRes{
		Subscription:              bean.SimplifySubscription(ctx, sub),
		Invoice:                   invoice,
		SubscriptionPendingUpdate: pendingUpdateDetail,
	}, nil
}
