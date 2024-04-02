package handler

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
)

func HandleOnetimeAddonPaymentFailure(ctx context.Context, id uint64) (bool, error) {
	one := query.GetSubscriptionOnetimeAddonById(ctx, id)
	if one == nil {
		return false, gerror.New("HandleOnetimeAddonPaymentFailure Id Not Found:" + strconv.FormatUint(id, 10))
	}
	if one.Status > 1 {
		return true, nil
	}
	_, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().Status:    4,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).Where(dao.SubscriptionOnetimeAddon.Columns().Status, 1).OmitNil().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func HandleOnetimeAddonPaymentCancel(ctx context.Context, id uint64) (bool, error) {
	one := query.GetSubscriptionOnetimeAddonById(ctx, id)
	if one == nil {
		return false, gerror.New("HandleOnetimeAddonPaymentFailure Id Not Found:" + strconv.FormatUint(id, 10))
	}
	if one.Status > 1 {
		return true, nil
	}
	_, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().Status:    3,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).Where(dao.SubscriptionOnetimeAddon.Columns().Status, 1).OmitNil().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func HandleOnetimeAddonPaymentSuccess(ctx context.Context, id uint64) (bool, error) {
	one := query.GetSubscriptionOnetimeAddonById(ctx, id)
	if one == nil {
		return false, gerror.New("HandleOnetimeAddonPaymentFailure Id Not Found:" + strconv.FormatUint(id, 10))
	}
	if one.Status > 1 {
		return true, nil
	}
	_, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().Status:    2,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).Where(dao.SubscriptionOnetimeAddon.Columns().Status, 1).OmitNil().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}
