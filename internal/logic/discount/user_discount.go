package discount

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
)

func UserDiscountApplyPreview(ctx context.Context, merchantId uint64, userId uint64, code string, subscriptionId string) (canApply bool, message string) {
	if merchantId == 0 {
		return false, "invalid merchantId"
	}
	if userId == 0 {
		return false, "invalid userId"
	}
	if len(code) == 0 {
		return false, "invalid code"
	}
	discountCode := query.GetDiscountByCode(ctx, merchantId, code)
	if discountCode == nil {
		return false, "discount code not found"
	}
	if discountCode.Status != DiscountStatusActive {
		return false, "discount code not active"
	}
	if discountCode.StartTime > gtime.Now().Timestamp() {
		return false, "discount not start"
	}
	if discountCode.EndTime < gtime.Now().Timestamp() {
		return false, "discount expired"
	}
	if discountCode.UserLimit > 0 {
		//check user limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, merchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, userId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, code).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Count()
		if err != nil {
			return false, err.Error()
		}
		if discountCode.UserLimit <= count {
			return false, "reach out the limit"
		}
	}
	if discountCode.SubscriptionLimit > 0 {
		if len(subscriptionId) == 0 {
			return false, "invalid subscriptionId"
		}
		//check user subscription limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, merchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, userId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, code).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, subscriptionId).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Count()
		if err != nil {
			return false, err.Error()
		}
		if discountCode.SubscriptionLimit <= count {
			return false, "reach out the limit"
		}
	}

	return true, ""
}
