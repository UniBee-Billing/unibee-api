package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetEarliestOneWebsocketUnreadMerchantWebhookMessage(ctx context.Context, merchantId uint64) (res *entity.MerchantWebhookMessage) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.MerchantWebhookMessage.Ctx(ctx).
		Where(dao.MerchantWebhookMessage.Columns().MerchantId, merchantId).
		Where(dao.MerchantWebhookMessage.Columns().WebsocketStatus, 10).
		WhereNotNull(dao.MerchantWebhookMessage.Columns().Data).
		OrderAsc(dao.MerchantWebhookMessage.Columns().CreateTime).
		Scan(&res)
	if err != nil {
		return nil
	}
	return res
}
