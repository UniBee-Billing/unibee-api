package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

//func GetSubscriptionById(ctx context.Context, id int64) (one *entity.Subscription) {
//	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{Id: uint64(id)}).OmitEmpty().Scan(&one)
//	if err != nil {
//		one = nil
//	}
//	return
//}

func GetLatestActiveOrIncompleteSubscriptionByUserId(ctx context.Context, userId int64, merchantId uint64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusIncomplete, consts.SubStatusActive}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx context.Context, userId uint64, merchantId uint64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusCreate, consts.SubStatusActive, consts.SubStatusIncomplete}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionBySubscriptionId(ctx context.Context, subscriptionId string) (one *entity.Subscription) {
	if len(subscriptionId) == 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{SubscriptionId: subscriptionId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionPendingUpdateByInvoiceId(ctx context.Context, invoiceId string) *entity.SubscriptionPendingUpdate {
	if len(invoiceId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().InvoiceId, invoiceId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionUpgradePendingUpdateByInvoiceId(ctx context.Context, invoiceId string) *entity.SubscriptionPendingUpdate {
	if len(invoiceId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().InvoiceId, invoiceId).
		Where(dao.SubscriptionPendingUpdate.Columns().EffectImmediate, 1).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionPendingUpdateByPendingUpdateId(ctx context.Context, pendingUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx context.Context, pendingUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUnfinishedSubscriptionPendingUpdateByInvoiceId(ctx context.Context, invoiceId string) *entity.SubscriptionPendingUpdate {
	if len(invoiceId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().InvoiceId, invoiceId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionTimeLineByUniqueId(ctx context.Context, uniqueId string) (one *entity.SubscriptionTimeline) {
	if len(uniqueId) == 0 {
		return nil
	}
	err := dao.SubscriptionTimeline.Ctx(ctx).Where(entity.SubscriptionTimeline{UniqueId: uniqueId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
