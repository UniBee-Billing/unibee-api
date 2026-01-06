package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/checkout"
)

func (c *ControllerCheckout) Edit(ctx context.Context, req *checkout.EditReq) (res *checkout.EditRes, err error) {
	one := query.GetMerchantCheckoutById(ctx, _interface.GetMerchantId(ctx), uint64(req.CheckoutId))
	utility.Assert(one != nil, "Checkout not found, please setup first")
	utility.Assert(one.Name != bean.DefaultCheckoutName && one.Description != bean.DefaultCheckoutDescription, "Can't edit default checkout")
	_, err = dao.MerchantCheckout.Ctx(ctx).Data(g.Map{
		dao.MerchantCheckout.Columns().Name:        req.Name,
		dao.MerchantCheckout.Columns().Description: req.Description,
		dao.MerchantCheckout.Columns().Data:        utility.MarshalToJsonString(req.Data),
		dao.MerchantCheckout.Columns().Staging:     utility.MarshalToJsonString(req.StagingData),
		dao.MerchantCheckout.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.MerchantCheckout.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "Update Checkout Error:%s\n", err.Error())
		return nil, err
	}
	one = query.GetMerchantCheckoutById(ctx, _interface.GetMerchantId(ctx), uint64(req.CheckoutId))
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Checkout(%d)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &checkout.EditRes{MerchantCheckout: bean.SimplifyMerchantCheckout(ctx, one)}, nil
}
