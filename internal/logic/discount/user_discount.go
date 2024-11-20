package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/internal/cmd/i18n"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	quantity2 "unibee/internal/logic/discount/quantity"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type UserDiscountApplyReq struct {
	MerchantId       uint64 `json:"merchantId"        description:"MerchantId"`
	UserId           uint64 `json:"userId"        description:"UserId"`
	PLanId           uint64 `json:"planId"        description:"PLanId"`
	DiscountCode     string `json:"discountCode"        description:"DiscountCode"`
	SubscriptionId   string `json:"subscriptionId"        description:"SubscriptionId"`
	PaymentId        string `json:"paymentId"        description:"PaymentId"`
	InvoiceId        string `json:"invoiceId"        description:"InvoiceId"`
	ApplyAmount      int64  `json:"applyAmount"        description:"ApplyAmount"`
	Currency         string `json:"currency"        description:"Currency"`
	TimeNow          int64  `json:"timeNow"        description:"TimeNow"`
	IsRecurringApply bool   `json:"isRecurringApply"        description:"IsRecurringApply"`
}

func UserDiscountApplyPreview(ctx context.Context, req *UserDiscountApplyReq) (canApply bool, isRecurring bool, message string) {
	if req.MerchantId == 0 {
		return false, false, "Invalid merchantId"
	}
	if len(req.DiscountCode) == 0 {
		return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}")
	}
	discountCode := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
	if discountCode == nil {
		return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}")
	}
	quantity2.GetDiscountQuantityUsedCount(ctx, discountCode.Id)
	if discountCode.DiscountType == consts.DiscountTypeFixedAmount && len(discountCode.Currency) == 0 {
		return false, false, "Code is fixed amount,but currency not set"
	}
	if discountCode.DiscountType == consts.DiscountTypeFixedAmount && strings.Compare(strings.ToUpper(req.Currency), strings.ToUpper(discountCode.Currency)) != 0 {
		return false, false, "Code currency not match plan"
	}
	compensateHistoryPlanId(ctx, req)
	{
		if !req.IsRecurringApply {
			if discountCode.Quantity > 0 {
				if discountCode.Quantity <= int64(quantity2.GetDiscountQuantityUsedCount(ctx, discountCode.Id)) {
					return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeReachLimitation}")
				}
			}
			if discountCode.Status != consts.DiscountStatusActive {
				return false, false, "Code not active"
			}
			if discountCode.StartTime > req.TimeNow {
				return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeNotStart}")
			}
			if discountCode.EndTime != 0 && discountCode.EndTime < req.TimeNow {
				return false, false, "Code expired"
			}
			if len(discountCode.PlanIds) > 0 {
				if req.PLanId <= 0 {
					return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodePlanError}")
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
					return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodePlanError}")
				}
			}
		}
	}
	if discountCode.UserLimit > 0 && req.UserId > 0 {
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
		if discountCode.UserLimit < count+1 {
			return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeReachLimitation}")
		}
	}
	if discountCode.SubscriptionLimit > 0 && req.UserId > 0 && len(req.SubscriptionId) > 0 {
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
		if discountCode.SubscriptionLimit < count+1 {
			return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeReachLimitation}")
		}
	}
	if discountCode.BillingType == consts.DiscountBillingTypeRecurring && discountCode.CycleLimit > 0 && req.UserId > 0 && len(req.SubscriptionId) > 0 {
		one := getLastNonRecurringPurchase(ctx, req)
		if one != nil && one.Status == 2 {
			return false, false, "First code purchase already rolled back"
		}
		var recurringId int64 = 0
		if one != nil {
			recurringId = one.RecurringId
		}
		//check user subscription limit
		count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
			Where(dao.MerchantUserDiscountCode.Columns().Code, req.DiscountCode).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.SubscriptionId).
			Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
			Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
			Where(dao.MerchantUserDiscountCode.Columns().RecurringId, recurringId).
			Count()
		if err != nil {
			g.Log().Error(ctx, "UserDiscountApplyPreview error:%s", err.Error())

			return false, false, "Server Error"
		}
		if discountCode.CycleLimit < count+1 {
			return false, false, i18n.LocalizationFormat(ctx, "{#DiscountCodeReachLimitation}")
		}
	} else if discountCode.BillingType == consts.DiscountBillingTypeRecurring && discountCode.CycleLimit == 0 && req.UserId > 0 && len(req.SubscriptionId) > 0 {
		one := getLastNonRecurringPurchase(ctx, req)
		if one != nil && one.Status == 2 {
			return false, false, "First code purchase already rolled back"
		}
	}

	return true, discountCode.BillingType == consts.DiscountBillingTypeRecurring, ""
}

func getLastNonRecurringPurchase(ctx context.Context, req *UserDiscountApplyReq) (one *entity.MerchantUserDiscountCode) {
	err := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
		Where(dao.MerchantUserDiscountCode.Columns().Code, req.DiscountCode).
		Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.SubscriptionId).
		Where(dao.MerchantUserDiscountCode.Columns().Recurring, 0).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		OrderDesc(dao.MerchantUserDiscountCode.Columns().Id).
		Scan(&one)
	if one == nil || err != nil {
		if err != nil {
			g.Log().Error(ctx, "getLastNonRecurringPurchase error:%s", err.Error())
		}
		return nil
	} else {
		return one
	}
}

func compensateHistoryPlanId(ctx context.Context, req *UserDiscountApplyReq) {
	if req.MerchantId <= 0 || req.UserId <= 0 || len(req.SubscriptionId) == 0 {
		return
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	if sub != nil {
		var list []*entity.MerchantUserDiscountCode
		err := dao.MerchantUserDiscountCode.Ctx(ctx).
			Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantUserDiscountCode.Columns().UserId, req.UserId).
			Where(dao.MerchantUserDiscountCode.Columns().SubscriptionId, req.SubscriptionId).
			Where(dao.MerchantUserDiscountCode.Columns().PlanId, 0).
			Scan(&list)
		if err == nil {
			for _, one := range list {
				_, err = dao.MerchantUserDiscountCode.Ctx(ctx).Data(g.Map{
					dao.MerchantUserDiscountCode.Columns().PlanId:    sub.PlanId,
					dao.MerchantUserDiscountCode.Columns().GmtModify: gtime.Now(),
				}).Where(dao.MerchantUserDiscountCode.Columns().Id, one.Id).Where(dao.MerchantUserDiscountCode.Columns().Recurring, 0).OmitNil().Update()
			}
		}
	}
}

func UserDiscountApply(ctx context.Context, req *UserDiscountApplyReq) (discountCode *entity.MerchantUserDiscountCode, err error) {
	if len(req.DiscountCode) == 0 {
		return nil, gerror.New("invalid discountCode")
	}
	if len(req.SubscriptionId) == 0 {
		return nil, gerror.New("invalid subscriptionId")
	}

	code := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)

	quantity2.GetDiscountQuantityUsedCount(ctx, code.Id)
	//should make sure cache loaded
	recurring := 0
	if req.IsRecurringApply {
		recurring = 1
	}
	if recurring == 0 {
		if code.Quantity > 0 {
			utility.Assert(code.Quantity > int64(quantity2.GetDiscountQuantityUsedCount(ctx, code.Id)), i18n.LocalizationFormat(ctx, "{#DiscountCodeReachLimitation}"))
		}
	}
	var recurringId int64 = 0
	if req.IsRecurringApply {
		one := getLastNonRecurringPurchase(ctx, req)
		if one == nil {
			return nil, gerror.New(fmt.Sprintf("getLastNonRecurringPurchase failed, request:%s", utility.MarshalToJsonString(req)))
		}
		recurringId = one.RecurringId
	}
	one := &entity.MerchantUserDiscountCode{
		MerchantId:     req.MerchantId,
		UserId:         req.UserId,
		Code:           req.DiscountCode,
		Status:         1,
		PlanId:         strconv.FormatUint(req.PLanId, 10),
		SubscriptionId: req.SubscriptionId,
		PaymentId:      req.PaymentId,
		InvoiceId:      req.InvoiceId,
		UniqueId:       fmt.Sprintf("%d_%d_%s_%d_%s_%s_%s", req.MerchantId, req.UserId, req.DiscountCode, 1, req.SubscriptionId, req.PaymentId, req.InvoiceId),
		CreateTime:     gtime.Now().Timestamp(),
		ApplyAmount:    req.ApplyAmount,
		Currency:       req.Currency,
		Recurring:      recurring,
		RecurringId:    recurringId,
	}
	result, err := dao.MerchantUserDiscountCode.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = id
	if !req.IsRecurringApply {
		//recurring not involved in quantity management
		_, _ = g.Redis().IncrBy(ctx, quantity2.GetDiscountQuantityUsedCountCacheKey(one.Code), 1)
		_, err = dao.MerchantUserDiscountCode.Ctx(ctx).Data(g.Map{
			dao.MerchantUserDiscountCode.Columns().RecurringId: one.Id,
			dao.MerchantUserDiscountCode.Columns().GmtModify:   gtime.Now(),
		}).Where(dao.MerchantUserDiscountCode.Columns().Id, one.Id).Where(dao.MerchantUserDiscountCode.Columns().Recurring, 0).OmitNil().Update()
	}
	return one, nil
}

func userDiscountRollback(ctx context.Context, one *entity.MerchantUserDiscountCode) error {
	discount := query.GetDiscountByCode(ctx, one.MerchantId, one.Code)
	if discount == nil {
		return gerror.New("not found")
	}
	if one.Status == 2 {
		return nil
	}
	_, err := dao.MerchantUserDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantUserDiscountCode.Columns().Status:    2,
		dao.MerchantUserDiscountCode.Columns().UniqueId:  fmt.Sprintf("%s_%d", one.UniqueId, gtime.Now().Timestamp()),
		dao.MerchantUserDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantUserDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	if err == nil && quantity2.GetDiscountQuantityUsedCount(ctx, discount.Id) > 0 && one.Recurring == 0 {
		quantityRollbackCount, quantityRollbackErr := g.Redis().DecrBy(ctx, quantity2.GetDiscountQuantityUsedCountCacheKey(one.Code), 1)
		g.Log().Infof(ctx, "Discount_Quantity_Rollback, count:%d, error:%s", quantityRollbackCount, quantityRollbackErr)
		if quantityRollbackErr != nil {
			g.Log().Errorf(ctx, "Discount_Quantity_Rollback contain error, count:%d, error:%s", quantityRollbackCount, quantityRollbackErr.Error())
		}
	} else if err != nil {
		g.Log().Errorf(ctx, "Discount_Quantity_Rollback last error, error:%s", err.Error())
	} else if quantity2.GetDiscountQuantityUsedCount(ctx, discount.Id) == 0 {
		g.Log().Errorf(ctx, "Discount_Quantity_Rollback last error, GetDiscountQuantityUsedCount <= 0")
	} else if one.Recurring != 0 {
		g.Log().Errorf(ctx, "Discount_Quantity_Rollback last error, Recurring != 0")
	}
	return err
}

// UserDiscountRollbackFromPayment payment total refund, partial refund is not involved
func UserDiscountRollbackFromPayment(ctx context.Context, invoiceId string, paymentId string) error {
	if len(paymentId) == 0 {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invalid paymentId:%s", paymentId)
		return nil
	}
	invoice := query.GetInvoiceByPaymentId(ctx, paymentId)
	if invoice == nil {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invoice not found, paymentId:%s", paymentId)
		return nil
	}
	if len(invoiceId) == 0 {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invalid invoiceId:%s", invoiceId)
		return nil
	}
	var one *entity.MerchantUserDiscountCode
	err := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().InvoiceId, invoiceId).
		Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		Scan(&one)
	if one != nil {
		return userDiscountRollback(ctx, one)
	}
	return err
}

// UserDiscountRollbackFromInvoice invoice create failed|cancel|failed, partial refund is not involved
func UserDiscountRollbackFromInvoice(ctx context.Context, invoiceId string) error {
	if len(invoiceId) == 0 {
		g.Log().Error(ctx, "UserDiscountRollbackFromPayment invalid invoiceId:%s", invoiceId)
		return nil
	}
	var one *entity.MerchantUserDiscountCode
	err := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().InvoiceId, invoiceId).
		Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		Scan(&one)
	if one != nil {
		return userDiscountRollback(ctx, one)
	}
	return err
}

func UpdateUserDiscountPaymentIdWhenInvoicePaid(ctx context.Context, invoiceId string, paymentId string) {
	_, _ = dao.MerchantUserDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantUserDiscountCode.Columns().PaymentId: paymentId,
		dao.MerchantUserDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantUserDiscountCode.Columns().InvoiceId, invoiceId).OmitNil().Update()
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
	//if merchantDiscountCode.Status != consts.DiscountStatusActive {
	//	return 0
	//}
	//if (merchantDiscountCode.EndTime != 0 && merchantDiscountCode.EndTime < timeNow) || merchantDiscountCode.StartTime > timeNow {
	//	return 0
	//}

	if merchantDiscountCode.DiscountType == consts.DiscountTypePercentage {
		return int64(float64(totalAmountExcludeTax) * utility.ConvertTaxPercentageToInternalFloat(merchantDiscountCode.DiscountPercentage))
	} else if merchantDiscountCode.DiscountType == consts.DiscountTypeFixedAmount &&
		strings.Compare(strings.ToUpper(currency), strings.ToUpper(merchantDiscountCode.Currency)) == 0 {
		return merchantDiscountCode.DiscountAmount
	}
	return 0
}

func ComputeHistoryDiscountAmount(ctx context.Context, merchantId uint64, totalAmountExcludeTax int64, currency string, discountCode string, timeNow int64) int64 {
	if timeNow == 0 {
		timeNow = gtime.Now().Timestamp()
	}
	if merchantId <= 0 {
		return 0
	}
	if len(discountCode) == 0 {
		return 0
	}
	merchantDiscountCode := query.GetDiscountByCode(ctx, merchantId, discountCode)
	if merchantDiscountCode == nil {
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
