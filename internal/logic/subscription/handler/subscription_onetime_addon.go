package handler

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
)

func UpdateSubscriptionAddonPurchasePaymentId(ctx context.Context, id uint64, paymentId string) error {
	one := query.GetSubscriptionOnetimeAddonById(ctx, id)
	if one == nil {
		return gerror.New("onetimeAddon not found")
	}
	if len(one.PaymentId) > 0 {
		return nil
	}
	_, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().PaymentId: paymentId,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionOnetimeAddonDetail(ctx context.Context, id uint64) *detail.SubscriptionOnetimeAddonDetail {
	one := query.GetSubscriptionOnetimeAddonById(ctx, id)
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			g.Log().Errorf(ctx, "SimplifySubscriptionOnetimeAddon Unmarshal Metadata error:%s", err.Error())
		}
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
	return &detail.SubscriptionOnetimeAddonDetail{
		Id:             one.Id,
		SubscriptionId: one.SubscriptionId,
		AddonId:        one.AddonId,
		Addon:          bean.SimplifyPlan(query.GetPlanById(ctx, one.AddonId)),
		Quantity:       one.Quantity,
		Status:         one.Status,
		CreateTime:     one.CreateTime,
		Payment:        bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, one.PaymentId)),
		User:           bean.SimplifyUserAccount(query.GetUserAccountById(ctx, sub.UserId)),
		Metadata:       metadata,
	}
}

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
