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

func (c *ControllerCheckout) Archive(ctx context.Context, req *checkout.ArchiveReq) (res *checkout.ArchiveRes, err error) {
	one := query.GetMerchantCheckoutById(ctx, _interface.GetMerchantId(ctx), uint64(req.CheckoutId))
	utility.Assert(one != nil, "Checkout not found")
	if one.IsDeleted != 0 {
		return &checkout.ArchiveRes{}, nil
	}
	utility.Assert(one.Name != bean.DefaultCheckoutName && one.Description != bean.DefaultCheckoutDescription, "Can't archive default checkout")
	_, err = dao.MerchantCheckout.Ctx(ctx).Data(g.Map{
		dao.MerchantCheckout.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantCheckout.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantCheckout.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "Archive Checkout Error:%s\n", err.Error())
		return nil, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Checkout(%d)", one.Id),
		Content:        "Archive",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &checkout.ArchiveRes{}, nil
}
