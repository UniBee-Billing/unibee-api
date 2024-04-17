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
	MerchantId     uint64 `json:"merchantId"        description:"MerchantId"`
	UserId         uint64 `json:"userId"        description:"UserId"`
	DiscountCode   string `json:"discountCode"        description:"DiscountCode"`
	SubscriptionId string `json:"subscriptionId"        description:"SubscriptionId"`
	PaymentId      string `json:"paymentId"        description:"PaymentId"`
	InvoiceId      string `json:"invoiceId"        description:"InvoiceId"`
	ApplyAmount    int64  `json:"applyAmount"        description:"ApplyAmount"`
	Currency       string `json:"currency"        description:"Currency"`
}

func UserDiscountApplyPreview(ctx context.Context, req *UserDiscountApplyReq) (canApply bool, isRecurring bool, message string) {
	if req.MerchantId == 0 {
		return false, false, "invalid MerchantId"
	}
	if req.UserId == 0 {
		return false, false, "invalid UserId"
	}
	if len(req.DiscountCode) == 0 {
		return false, false, "invalid Code"
	}
	discountCode := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
	if discountCode == nil {
		return false, false, "discount Code not found"
	}
	if discountCode.Status != DiscountStatusActive {
		return false, false, "discount Code not active"
	}
	if discountCode.StartTime > gtime.Now().Timestamp() {
		return false, false, "discount not start"
	}
	if discountCode.EndTime != 0 && discountCode.EndTime < gtime.Now().Timestamp() {
		return false, false, "discount expired"
	}
	if discountCode.DiscountType == DiscountTypeFixedAmount && strings.Compare(strings.ToUpper(req.Currency), strings.ToUpper(discountCode.Currency)) != 0 {
		return false, false, "currency not match"
	}
	if discountCode.UserLimit > 0 {
		//check user limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.DiscountCode).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Count()
		if err != nil {
			return false, false, err.Error()
		}
		if discountCode.UserLimit <= count {
			return false, false, "reach out the limit"
		}
	}
	if discountCode.SubscriptionLimit > 0 {
		if len(req.SubscriptionId) == 0 {
			return false, false, "invalid SubscriptionId"
		}
		//check user subscription limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.DiscountCode).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.SubscriptionId).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Count()
		if err != nil {
			return false, false, err.Error()
		}
		if discountCode.SubscriptionLimit <= count {
			return false, false, "reach out the limit"
		}
	}
	if discountCode.CycleLimit > 0 {
		if len(req.SubscriptionId) == 0 {
			return false, false, "invalid SubscriptionId"
		}
		//check user subscription limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.DiscountCode).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.SubscriptionId).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Count()
		if err != nil {
			return false, false, err.Error()
		}
		if discountCode.CycleLimit <= count {
			return false, false, "reach out the limit"
		}
	}

	return true, discountCode.BillingType == DiscountBillingTypeRecurring, ""
}

func UserDiscountApply(ctx context.Context, req *UserDiscountApplyReq) (discountCode *entity.MerchantUserDiscountCode, err error) {
	if len(req.DiscountCode) == 0 {
		return nil, gerror.New("invalid discountCode")
	}
	one := &entity.MerchantUserDiscountCode{
		MerchantId:     req.MerchantId,
		UserId:         req.UserId,
		Code:           req.DiscountCode,
		Status:         1,
		SubscriptionId: req.SubscriptionId,
		PaymentId:      req.PaymentId,
		InvoiceId:      req.InvoiceId,
		UniqueId:       fmt.Sprintf("%d_%d_%s_%d_%s_%s_%s", req.MerchantId, req.UserId, req.DiscountCode, 1, req.SubscriptionId, req.PaymentId, req.InvoiceId),
		CreateTime:     gtime.Now().Timestamp(),
		ApplyAmount:    req.ApplyAmount,
		Currency:       req.Currency,
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

func UserDiscountRollbackFromPayment(ctx context.Context, paymentId string) error {
	var one *entity.MerchantUserDiscountCode
	err := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().PaymentId, paymentId).
		Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		Scan(&one)
	if one != nil {
		return UserDiscountRollback(ctx, one.Id)
	}
	return err
}

func ComputeDiscountAmount(ctx context.Context, merchantId uint64, totalAmountExcludeTax int64, currency string, discountCode string, timeNow int64) int64 {
	if timeNow == 0 {
		timeNow = gtime.Now().Timestamp()
	}
	if merchantId <= 0 {
		return 0
	}
	if totalAmountExcludeTax <= 0 {
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
	if (merchantDiscountCode.EndTime != 0 && merchantDiscountCode.EndTime < timeNow) || merchantDiscountCode.StartTime > timeNow {
		return 0
	}

	if merchantDiscountCode.DiscountType == DiscountTypePercentage {
		return int64(float64(totalAmountExcludeTax) * utility.ConvertTaxPercentageToInternalFloat(merchantDiscountCode.DiscountPercentage))
	} else if merchantDiscountCode.DiscountType == DiscountTypeFixedAmount &&
		strings.Compare(strings.ToUpper(currency), strings.ToUpper(merchantDiscountCode.Currency)) == 0 {
		return merchantDiscountCode.DiscountAmount
	}
	return 0
}
