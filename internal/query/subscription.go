package query

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetUserAllActiveOrIncompleteSubscriptions(ctx context.Context, userId uint64, merchantId uint64) (list []*entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusIncomplete, consts.SubStatusActive}).
		Scan(&list)
	if err != nil {
		list = make([]*entity.Subscription, 0)
	}
	return
}

func GetLatestSubscriptionByUserId(ctx context.Context, userId uint64, merchantId uint64, productId int64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		WhereIn(dao.Subscription.Columns().PlanId, GetPlanIdsByProductId(ctx, merchantId, productId)).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetLatestActiveOrIncompleteSubscriptionByUserId(ctx context.Context, userId uint64, merchantId uint64, productId int64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		WhereIn(dao.Subscription.Columns().PlanId, GetPlanIdsByProductId(ctx, merchantId, productId)).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusIncomplete, consts.SubStatusActive}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetLatestCreateOrProcessingSubscriptionByUserId(ctx context.Context, userId uint64, merchantId uint64, productId int64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		WhereIn(dao.Subscription.Columns().PlanId, GetPlanIdsByProductId(ctx, merchantId, productId)).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusPending, consts.SubStatusProcessing}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx context.Context, userId uint64, merchantId uint64, productId int64) (one *entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return nil
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		WhereIn(dao.Subscription.Columns().PlanId, GetPlanIdsByProductId(ctx, merchantId, productId)).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusPending, consts.SubStatusProcessing, consts.SubStatusActive, consts.SubStatusIncomplete}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetLatestActiveOrIncompleteOrCreateSubscriptionsByUserId(ctx context.Context, userId uint64, merchantId uint64) (list []*entity.Subscription) {
	if userId <= 0 || merchantId <= 0 {
		return make([]*entity.Subscription, 0)
	}
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().UserId, userId).
		Where(dao.Subscription.Columns().MerchantId, merchantId).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusPending, consts.SubStatusProcessing, consts.SubStatusActive, consts.SubStatusIncomplete}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		Scan(&list)
	if err != nil || list == nil {
		list = make([]*entity.Subscription, 0)
	}
	return
}

func GetSubscriptionByExternalSubscriptionId(ctx context.Context, externalSubscriptionId string) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{ExternalSubscriptionId: externalSubscriptionId}).OmitEmpty().Scan(&one)
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

func GetSubscriptionEffectImmediatePendingUpdateByInvoiceId(ctx context.Context, invoiceId string) *entity.SubscriptionPendingUpdate {
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
		Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdateId).
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
		Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdateId).
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

func GetLatestSubscriptionTimeLine(ctx context.Context, subscriptionId string) (one *entity.SubscriptionTimeline) {
	if len(subscriptionId) == 0 {
		return nil
	}
	err := dao.SubscriptionTimeline.Ctx(ctx).
		Where(dao.SubscriptionTimeline.Columns().SubscriptionId, subscriptionId).
		OrderDesc(dao.SubscriptionTimeline.Columns().Id).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
