package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetMerchantWebhook(ctx context.Context, id int64) (one *entity.MerchantWebhook) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantWebhook.Ctx(ctx).Where(entity.MerchantWebhook{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantWebhookByUrl(ctx context.Context, url string) (one *entity.MerchantWebhook) {
	if len(url) <= 0 {
		return nil
	}
	err := dao.MerchantWebhook.Ctx(ctx).Where(entity.MerchantWebhook{WebhookUrl: url}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantWebhooksByMerchantId(ctx context.Context, merchantId uint64) (list []*entity.MerchantWebhook) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.MerchantWebhook.Ctx(ctx).Where(dao.MerchantWebhook.Columns().MerchantId, merchantId).Where(dao.MerchantWebhook.Columns().IsDeleted, 0).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
