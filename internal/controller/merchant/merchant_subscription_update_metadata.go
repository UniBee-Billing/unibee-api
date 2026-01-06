package merchant

import (
	"context"
	"fmt"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) UpdateMetadata(ctx context.Context, req *subscription.UpdateMetadataReq) (res *subscription.UpdateMetadataRes, err error) {
	utility.Assert(len(req.SubscriptionId) > 0, "subscription id should not be empty")
	utility.Assert(len(req.Metadata) > 0, "metadata should not be empty")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status: utility.MergeMetadata(sub.MetaData, &req.Metadata),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	utility.AssertError(err, "Update Subscription Metadata Error")
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("UpdateSubscriptionMetadata(%s)", sub.SubscriptionId),
		Content:        utility.MarshalToJsonString(req.Metadata),
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &subscription.UpdateMetadataRes{}, nil
}
