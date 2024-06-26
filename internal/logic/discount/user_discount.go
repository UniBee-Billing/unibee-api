package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type UserDiscountApplyReq struct {
	MerchantId     uint64 `json:"merchantId"        description:"MerchantId"`
	UserId         uint64 `json:"userId"        description:"UserId"`
	PLanId         uint64 `json:"planId"        description:"PLanId"`
	DiscountCode   string `json:"discountCode"        description:"DiscountCode"`
	SubscriptionId string `json:"subscriptionId"        description:"SubscriptionId"`
	PaymentId      string `json:"paymentId"        description:"PaymentId"`
	InvoiceId      string `json:"invoiceId"        description:"InvoiceId"`
	ApplyAmount    int64  `json:"applyAmount"        description:"ApplyAmount"`
	Currency       string `json:"currency"        description:"Currency"`
	TimeNow        int64  `json:"timeNow"        description:"TimeNow"`
}

func UserDiscountApplyPreview(ctx context.Context, req *UserDiscountApplyReq) (canApply bool, isRecurring bool, message string) {
	if req.MerchantId == 0 {
		return false, false, "Invalid merchantId"
	}
	if req.UserId == 0 {
		return false, false, "Invalid userId"
	}
	if len(req.DiscountCode) == 0 {
		return false, false, "Invalid code"
	}
	discountCode := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
	if discountCode == nil {
		return false, false, "Code not found"
	}
	if discountCode.Status != consts.DiscountStatusActive {
		return false, false, "Code not active"
	}
	if discountCode.StartTime > req.TimeNow {
		return false, false, "Code not ready"
	}
	if discountCode.EndTime != 0 && discountCode.EndTime < req.TimeNow {
		return false, false, "Code expired"
	}
	if discountCode.DiscountType == consts.DiscountTypeFixedAmount && len(discountCode.Currency) == 0 {
		return false, false, "Code is fixed amount,but currency not set"
	}
	if discountCode.DiscountType == consts.DiscountTypeFixedAmount && strings.Compare(strings.ToUpper(req.Currency), strings.ToUpper(discountCode.Currency)) != 0 {
		return false, false, "Code currency not match plan"
	}
	if len(discountCode.PlanIds) > 0 {
		if req.PLanId <= 0 {
			return false, false, "Code not allow to use on this plan"
		}
		var match = false
		planIds := strings.Split(discountCode.PlanIds, ",")
		for _, s := range planIds {
			planId, err := strconv.ParseUint(s, 10, 64)
			if err == nil && planId == req.PLanId {
				match = true
				break
			}
		}
		if !match {
			return false, false, "Code not allow to use on this plan"
		}
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
			g.Log().Error(ctx, "UserDiscountApplyPreview error:%s", err.Error())
			return false, false, "Server Error"
		}
		if discountCode.UserLimit <= count+1 {
			return false, false, "Code reach out the limit"
		}
	}
	if discountCode.SubscriptionLimit > 0 && len(req.SubscriptionId) > 0 {
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
			g.Log().Error(ctx, "UserDiscountApplyPreview error:%s", err.Error())
			return false, false, "Server Error"
		}
		if discountCode.SubscriptionLimit <= count+1 {
			return false, false, "Code reach out the limit"
		}
	}
	if discountCode.BillingType == consts.DiscountBillingTypeRecurring && discountCode.CycleLimit > 0 && len(req.SubscriptionId) > 0 {
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
			g.Log().Error(ctx, "UserDiscountApplyPreview error:%s", err.Error())
			return false, false, "Server Error"
		}
		if discountCode.CycleLimit <= count+1 {
			return false, false, "Code reach out the limit"
		}
	}

	return true, discountCode.BillingType == consts.DiscountBillingTypeRecurring, ""
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
	if merchantDiscountCode == nil {
		return 0
	}
	if merchantDiscountCode.Status != consts.DiscountStatusActive {
		return 0
	}
	if (merchantDiscountCode.EndTime != 0 && merchantDiscountCode.EndTime < timeNow) || merchantDiscountCode.StartTime > timeNow {
		return 0
	}

	if merchantDiscountCode.DiscountType == consts.DiscountTypePercentage {
		return int64(float64(totalAmountExcludeTax) * utility.ConvertTaxPercentageToInternalFloat(merchantDiscountCode.DiscountPercentage))
	} else if merchantDiscountCode.DiscountType == consts.DiscountTypeFixedAmount &&
		strings.Compare(strings.ToUpper(currency), strings.ToUpper(merchantDiscountCode.Currency)) == 0 {
		return merchantDiscountCode.DiscountAmount
	}
	return 0
}
