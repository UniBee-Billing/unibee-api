package checkout

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/checkout/checkout"
)

func (c *ControllerCheckout) Get(ctx context.Context, req *checkout.GetReq) (res *checkout.GetRes, err error) {
	if req.CheckoutId <= 0 {
		return nil, gerror.New("checkout not found")
	}
	var one *entity.MerchantCheckout
	err = dao.MerchantCheckout.Ctx(ctx).
		Where(dao.MerchantCheckout.Columns().Id, req.CheckoutId).
		Where(dao.MerchantCheckout.Columns().IsDeleted, 0).
		Scan(&one)
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, gerror.New("checkout not found")
	}
	output := bean.SimplifyMerchantCheckout(ctx, one)
	output.StagingData = nil
	return &checkout.GetRes{MerchantCheckout: output}, nil
}
