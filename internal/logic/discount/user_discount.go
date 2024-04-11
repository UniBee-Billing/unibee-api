package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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

func UserDiscountApply(ctx context.Context, merchantId uint64, userId uint64, code string, subscriptionId string, paymentId string, invoiceId string) (discountCode *entity.MerchantUserDiscountCode, err error) {
	one := &entity.MerchantUserDiscountCode{
		MerchantId:     merchantId,
		UserId:         userId,
		Code:           code,
		Status:         1,
		SubscriptionId: subscriptionId,
		PaymentId:      paymentId,
		InvoiceId:      invoiceId,
		UniqueId:       fmt.Sprintf("%d_%d_%s_%d_%s_%s_%s", merchantId, userId, code, 1, subscriptionId, paymentId, invoiceId),
		CreateTime:     gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantUserDiscountCode.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = id

	return one, nil
}

func UserDiscountRollback(ctx context.Context, id int64) error {
	one := query.GetUserDiscountById(ctx, id)
	if one == nil {
		return gerror.New("not found")
	}
	if one.Status == 2 {
		return nil
	}
	_, err := dao.MerchantUserDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantUserDiscountCode.Columns().Status:    2,
		dao.MerchantUserDiscountCode.Columns().UniqueId:  fmt.Sprintf("%s_%d", one.UniqueId, gtime.Now().Timestamp()),
		dao.MerchantUserDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantUserDiscountCode.Columns().Id, id).OmitNil().Update()
	return err
}
