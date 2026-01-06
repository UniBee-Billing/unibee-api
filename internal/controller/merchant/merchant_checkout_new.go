package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/checkout"
)

func (c *ControllerCheckout) New(ctx context.Context, req *checkout.NewReq) (res *checkout.NewRes, err error) {
	utility.Assert(len(req.Name) > 0, "invalid name")
	utility.Assert(req.Name != bean.DefaultCheckoutName && req.Description != bean.DefaultCheckoutDescription, "Can't create default checkout")
	query.InitDefaultMerchantCheckout(ctx, _interface.GetMerchantId(ctx))
	one := &entity.MerchantCheckout{
		MerchantId:  _interface.GetMerchantId(ctx),
		Name:        req.Name,
		Description: req.Description,
		Data:        utility.MarshalToJsonString(req.Data),
		Staging:     utility.MarshalToJsonString(req.StagingData),
		IsDeleted:   0,
		CreateTime:  gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantCheckout.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`create merchant checkout record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = id
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Checkout(%d)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &checkout.NewRes{MerchantCheckout: bean.SimplifyMerchantCheckout(ctx, one)}, nil
}
