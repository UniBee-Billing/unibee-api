package merchant

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"

	"unibee/api/merchant/checkout"
)

func (c *ControllerCheckout) List(ctx context.Context, req *checkout.ListReq) (res *checkout.ListRes, err error) {
	// init merchant default checkout
	query.InitDefaultMerchantCheckout(ctx, _interface.GetMerchantId(ctx))
	list := make([]*entity.MerchantCheckout, 0)
	q := dao.MerchantCheckout.Ctx(ctx).
		Where(dao.MerchantCheckout.Columns().MerchantId, _interface.GetMerchantId(ctx)).
		Where(dao.MerchantCheckout.Columns().IsDeleted, 0)
	if len(req.SearchKey) > 0 {
		q = q.Where(
			q.Builder().WhereOrLike(dao.MerchantCheckout.Columns().Id, "%"+req.SearchKey+"%").
				WhereOrLike(dao.MerchantCheckout.Columns().Name, "%"+req.SearchKey+"%").
				WhereOrLike(dao.MerchantCheckout.Columns().Description, "%"+req.SearchKey+"%"))
	}
	err = q.OrderDesc(dao.MerchantCheckout.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return &checkout.ListRes{MerchantCheckouts: bean.SimplifyMerchantCheckoutList(ctx, list)}, nil
}
