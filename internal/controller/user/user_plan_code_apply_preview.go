package user

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/internal/cmd/i18n"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/plan"
)

func (c *ControllerPlan) CodeApplyPreview(ctx context.Context, req *plan.CodeApplyPreviewReq) (res *plan.CodeApplyPreviewRes, err error) {
	if req.PlanId == 0 && len(req.ExternalPlanId) > 0 {
		one := query.GetPlanByExternalPlanId(ctx, _interface.GetMerchantId(ctx), req.ExternalPlanId)
		if one != nil {
			req.PlanId = int64(one.Id)
		}
	}
	utility.Assert(req.PlanId > 0, "Invalid planId")
	one := query.GetPlanById(ctx, uint64(req.PlanId))
	utility.Assert(one != nil, "Plan Not Found")
	utility.Assert(one.Status == consts.PlanStatusActive, "Plan Not Active")
	utility.Assert(one.Type == consts.PlanTypeMain, "Not Main Plan")
	utility.Assert(len(req.Code) > 0, "Invalid Code")
	oneDiscount := query.GetDiscountByCode(ctx, _interface.GetMerchantId(ctx), req.Code)
	if oneDiscount == nil {
		return &plan.CodeApplyPreviewRes{
			Valid:          false,
			DiscountAmount: 0,
			DiscountCode:   nil,
			FailureReason:  i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}"),
		}, nil
	}
	//utility.Assert(oneDiscount != nil, i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}"))
	canApply, _, message := discount2.UserDiscountApplyPreview(ctx, &discount2.UserDiscountApplyReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		PLanId:             uint64(req.PlanId),
		DiscountCode:       req.Code,
		Currency:           one.Currency,
		TimeNow:            gtime.Now().Timestamp(),
		IsUpgrade:          req.IsUpgrade,
		IsChangeToLongPlan: req.IsChangeToLongPlan,
		IsRenew:            req.IsRenew,
		IsNewUser:          service.IsNewSubscriptionUser(ctx, _interface.GetMerchantId(ctx), strings.ToLower(req.Email)),
	})
	discountAmount := utility.MinInt64(discount2.ComputeDiscountAmount(ctx, one.MerchantId, one.Amount, one.Currency, req.Code, gtime.Now().Timestamp()), one.Amount)
	return &plan.CodeApplyPreviewRes{
		Valid:          canApply,
		DiscountAmount: discountAmount,
		DiscountCode:   bean.SimplifyMerchantDiscountCode(oneDiscount),
		FailureReason:  message,
	}, nil
}
