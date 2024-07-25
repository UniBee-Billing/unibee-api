package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetMerchantWebhook(ctx context.Context, id uint64) (one *entity.MerchantWebhook) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantWebhook.Ctx(ctx).Where(dao.MerchantWebhook.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantWebhookByUrl(ctx context.Context, merchantId uint64, url string) (one *entity.MerchantWebhook) {
	if len(url) <= 0 {
		return nil
	}
	err := dao.MerchantWebhook.Ctx(ctx).Where(dao.MerchantWebhook.Columns().WebhookUrl, url).Where(dao.MerchantWebhook.Columns().MerchantId, merchantId).OmitEmpty().Scan(&one)
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
