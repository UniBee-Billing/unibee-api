package merchant

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/cmd/i18n"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) PlanApplyPreview(ctx context.Context, req *discount.PlanApplyPreviewReq) (res *discount.PlanApplyPreviewRes, err error) {
	if req.PlanId == 0 && len(req.ExternalPlanId) > 0 {
		plan := query.GetPlanByExternalPlanId(ctx, _interface.GetMerchantId(ctx), req.ExternalPlanId)
		if plan != nil {
			req.PlanId = int64(plan.Id)
		}
	}
	utility.Assert(req.PlanId > 0, "Invalid planId")
	plan := query.GetPlanById(ctx, uint64(req.PlanId))
	utility.Assert(plan != nil, "Plan Not Found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "Plan Not Active")
	utility.Assert(plan.Type == consts.PlanTypeMain, "Not Main Plan")
	utility.Assert(len(req.Code) > 0, "Invalid Code")
	oneDiscount := query.GetDiscountByCode(ctx, _interface.GetMerchantId(ctx), req.Code)
	if oneDiscount == nil || oneDiscount.IsDeleted > 0 {
		return &discount.PlanApplyPreviewRes{
			Valid:          false,
			DiscountAmount: 0,
			DiscountCode:   nil,
			FailureReason:  i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}"),
		}, nil
	}
	//utility.Assert(oneDiscount != nil, i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}"))
	canApply, _, message := discount2.UserDiscountApplyPreview(ctx, &discount2.UserDiscountApplyReq{
		MerchantId:   _interface.GetMerchantId(ctx),
		PLanId:       uint64(req.PlanId),
		DiscountCode: req.Code,
		Currency:     plan.Currency,
		TimeNow:      gtime.Now().Timestamp(),
	})
	discountAmount := utility.MinInt64(discount2.ComputeDiscountAmount(ctx, plan.MerchantId, plan.Amount, plan.Currency, req.Code, gtime.Now().Timestamp()), plan.Amount)
	return &discount.PlanApplyPreviewRes{
		Valid:          canApply,
		DiscountAmount: discountAmount,
		DiscountCode:   bean.SimplifyMerchantDiscountCode(oneDiscount),
		FailureReason:  message,
	}, nil
}
