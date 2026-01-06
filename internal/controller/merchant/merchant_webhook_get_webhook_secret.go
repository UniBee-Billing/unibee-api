package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	merchant2 "unibee/internal/logic/merchant"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) GetWebhookSecret(ctx context.Context, req *webhook.GetWebhookSecretReq) (res *webhook.GetWebhookSecretRes, err error) {
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "Server Error")
	if merchant.WebhookSecret != "" {
		return &webhook.GetWebhookSecretRes{Secret: merchant.WebhookSecret}, nil
	} else {
		secret := merchant2.GenerateMerchantWebHookSecret()
		_, err := dao.Merchant.Ctx(ctx).Data(g.Map{
			dao.Merchant.Columns().WebhookSecret: secret,
			dao.Merchant.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.Merchant.Columns().Id, merchant.Id).Update()
		merchant = query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     merchant.Id,
			Target:         fmt.Sprintf("WebhookSecret(%v)", utility.HideStar(merchant.WebhookSecret)),
			Content:        fmt.Sprintf("NewWebhookSecret(%v)", utility.HideStar(secret)),
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
		return &webhook.GetWebhookSecretRes{Secret: merchant.WebhookSecret}, nil
	}
}
