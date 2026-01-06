package query

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetMerchantCheckoutById(ctx context.Context, merchantId uint64, id uint64) (one *entity.MerchantCheckout) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.MerchantCheckout.Ctx(ctx).
		Where(dao.MerchantCheckout.Columns().MerchantId, merchantId).
		Where(dao.MerchantCheckout.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantCheckoutList(ctx context.Context, merchantId uint64) (list []*entity.MerchantCheckout) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.MerchantCheckout.Ctx(ctx).
		Where(dao.MerchantCheckout.Columns().MerchantId, merchantId).
		Where(dao.MerchantCheckout.Columns().IsDeleted, 0).
		OrderAsc(dao.MerchantCheckout.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil
	}
	return list
}

func InitDefaultMerchantCheckout(ctx context.Context, merchantId uint64) {
	if merchantId <= 0 {
		return
	}
	var defaultOne *entity.MerchantCheckout
	err := dao.MerchantCheckout.Ctx(ctx).
		Where(dao.MerchantCheckout.Columns().MerchantId, merchantId).
		Where(dao.MerchantCheckout.Columns().Name, bean.DefaultCheckoutName).
		Where(dao.MerchantCheckout.Columns().Description, bean.DefaultCheckoutDescription).
		Where(dao.MerchantCheckout.Columns().IsDeleted, 0).
		Scan(&defaultOne)
	if err != nil {
		return
	}
	if defaultOne == nil {
		one := &entity.MerchantCheckout{
			MerchantId:  merchantId,
			Name:        bean.DefaultCheckoutName,
			Description: bean.DefaultCheckoutDescription,
			Data:        utility.MarshalToJsonString(nil),
			Staging:     utility.MarshalToJsonString(nil),
			IsDeleted:   0,
			CreateTime:  gtime.Now().Timestamp(),
		}
		_, _ = dao.MerchantCheckout.Ctx(ctx).Data(one).OmitNil().Insert(one)
	}
}
