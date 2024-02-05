package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPayChannelById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPayChannelByChannel(ctx context.Context, channel string) (one *entity.MerchantGateway) {
	if len(channel) == 0 {
		return nil
	}
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{Channel: channel}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetPayChannelsGroupByEnumKey(ctx context.Context) []*entity.MerchantGateway {
	var data []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Group(dao.MerchantGateway.Columns().EnumKey).
		OmitEmpty().Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetPayChannelsGroupByEnumKey error:%s", err)
		return nil
	}
	return data
}

func GetPaymentTypePayChannelById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id)}).
		Where(m.Builder().
			Where(entity.MerchantGateway{ChannelType: consts.PayChannelTypePayment}).WhereOr("channel_type is null")).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionTypePayChannelById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id), ChannelType: consts.PayChannelTypeSubscription}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListSubscriptionTypePayChannels(ctx context.Context) (list []*entity.MerchantGateway) {
	var data []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{ChannelType: consts.PayChannelTypeSubscription}).
		OmitEmpty().Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetListSubscriptionTypePayChannels error:%s", err)
		return nil
	}
	return data
}

func SavePayChannelUniqueProductId(ctx context.Context, id int64, productId string) error {
	if len(productId) == 0 || id < 0 {
		return nil
	}
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().UniqueProductId: productId,
		dao.MerchantGateway.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, id).Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("savePayChannelUniqueProductId update err:%s", update)
	//}
	return nil
}

func UpdatePayChannelWebhookSecret(ctx context.Context, id int64, secret string) error {
	if id <= 0 {
		return gerror.New("invalid id")
	}
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().WebhookSecret: secret,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, id).Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("UpdatePayChannelWebhookSecret update err:%s", update)
	//}
	return nil
}
