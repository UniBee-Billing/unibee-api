package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type UserDiscountApplyReq struct {
	merchantId     uint64
	userId         uint64
	code           string
	subscriptionId string
	paymentId      string
	invoiceId      string
	applyAmount    int64
	currency       string
}

func UserDiscountApplyPreview(ctx context.Context, req *UserDiscountApplyReq) (canApply bool, message string) {
	if req.merchantId == 0 {
		return false, "invalid merchantId"
	}
	if req.userId == 0 {
		return false, "invalid userId"
	}
	if len(req.code) == 0 {
		return false, "invalid code"
	}
	discountCode := query.GetDiscountByCode(ctx, req.merchantId, req.code)
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
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.merchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.userId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.code).
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
		if len(req.subscriptionId) == 0 {
			return false, "invalid subscriptionId"
		}
		//check user subscription limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.merchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.userId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.code).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.subscriptionId).
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

func UserDiscountApply(ctx context.Context, req *UserDiscountApplyReq) (discountCode *entity.MerchantUserDiscountCode, err error) {
	one := &entity.MerchantUserDiscountCode{
		MerchantId:     req.merchantId,
		UserId:         req.userId,
		Code:           req.code,
		Status:         1,
		SubscriptionId: req.subscriptionId,
		PaymentId:      req.paymentId,
		InvoiceId:      req.invoiceId,
		UniqueId:       fmt.Sprintf("%d_%d_%s_%d_%s_%s_%s", req.merchantId, req.userId, req.code, 1, req.subscriptionId, req.paymentId, req.invoiceId),
		CreateTime:     gtime.Now().Timestamp(),
		ApplyAmount:    req.applyAmount,
		Currency:       req.currency,
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

func ComputeDiscountAmount(ctx context.Context, merchantId uint64, totalAmount int64, currency string, discountCode string, timeNow int64) int64 {
	if timeNow == 0 {
		timeNow = gtime.Now().Timestamp()
	}
	if merchantId <= 0 {
		return 0
	}
	if totalAmount <= 0 {
		return 0
	}
	if len(discountCode) == 0 {
		return 0
	}
	merchantDiscountCode := query.GetDiscountByCode(ctx, merchantId, discountCode)
	if merchantDiscountCode != nil {
		return 0
	}
	if merchantDiscountCode.Status != DiscountStatusActive {
		return 0
	}
	if merchantDiscountCode.EndTime < timeNow || merchantDiscountCode.StartTime > timeNow {
		return 0
	}
	if merchantDiscountCode.DiscountType == DiscountTypePercentage {
		return int64(float64(totalAmount) * utility.ConvertTaxScaleToInternalFloat(merchantDiscountCode.DiscountPercentage))
	} else if merchantDiscountCode.DiscountType == DiscountTypeFixedAmount &&
		strings.Compare(strings.ToUpper(currency), strings.ToUpper(merchantDiscountCode.Currency)) == 0 {
		return merchantDiscountCode.DiscountAmount
	}
	return 0
}