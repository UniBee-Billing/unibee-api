package discount

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	DiscountStatusEditable       = 1
	DiscountStatusActive         = 2
	DiscountStatusDeActive       = 3
	DiscountStatusExpired        = 4
	DiscountBillingTypeOnetime   = 1
	DiscountBillingTypeRecurring = 2
	DiscountTypePercentage       = 1
	DiscountTypeFixedAmount      = 2
)

type CreateDiscountCodeInternalReq struct {
	MerchantId         uint64 `json:"MerchantId"        description:"MerchantId"`
	Code               string `json:"Code"              description:"Code"`
	Name               string `json:"name"              description:"name"`                                                                        // name
	BillingType        int    `json:"billingType"       description:"billing_type, 1-one-time, 2-recurring"`                                       // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"      description:"discount_type, 1-percentage, 2-fixed_amount"`                                 // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"    description:"amount of discount, available when discount_type is fixed_amount"`            // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"Currency"          description:"Currency of discount, available when discount_type is fixed_amount"`          // Currency of discount, available when discount_type is fixed_amount
	UserLimit          int    `json:"userLimit"         description:"the limit of every user apply, 0-unlimited"`                                  // the limit of every user apply, 0-unlimited
	CycleLimit         int    `json:"cycleLimit"         description:"the count limitation of subscription cycle , 0-no limit"`                    // the count limitation of subscription cycle , 0-no limit
	SubscriptionLimit  int    `json:"subscriptionLimit" description:"the limit of every subscription apply, 0-unlimited"`                          // the limit of every subscription apply, 0-unlimited
	StartTime          int64  `json:"startTime"         description:"start of discount available utc time"`                                        // start of discount available utc time
	EndTime            int64  `json:"endTime"           description:"end of discount available utc time"`                                          // end of discount available utc time
}

func NewMerchantDiscountCode(ctx context.Context, req *CreateDiscountCodeInternalReq) error {
	utility.Assert(req.Code != "", "invalid Code")
	one := query.GetDiscountByCode(ctx, req.MerchantId, req.Code)
	utility.Assert(one == nil, "exist Code:"+req.Code)
	utility.Assert(req.BillingType == DiscountBillingTypeOnetime || req.BillingType == DiscountBillingTypeRecurring, "invalid billingType, 1-one-time, 2-recurring")
	utility.Assert(req.DiscountType == DiscountTypePercentage || req.DiscountType == DiscountTypeFixedAmount, "invalid billingType, 1-percentage, 2-fixed_amount")
	utility.Assert(req.UserLimit >= 0, "invalid UserLimit")
	utility.Assert(req.SubscriptionLimit >= 0, "invalid SubscriptionLimit")
	utility.Assert(req.StartTime >= gtime.Now().Timestamp(), "startTime should greater then time now")
	utility.Assert(req.EndTime >= req.StartTime, "startTime should lower then endTime")
	req.Currency = strings.ToUpper(req.Currency)
	if req.DiscountType == DiscountTypePercentage {
		utility.Assert(req.DiscountPercentage > 0 && req.DiscountPercentage < 10000, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount == 0, "invalid discountAmount")
		utility.Assert(len(req.Currency) == 0, "invalid Currency")
	} else if req.DiscountType == DiscountTypeFixedAmount {
		utility.Assert(req.DiscountPercentage == 0, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount >= 0, "invalid discountAmount")
		utility.Assert(len(req.Currency) >= 0, "invalid Currency")
	}

	one = &entity.MerchantDiscountCode{
		MerchantId:         req.MerchantId,
		Code:               req.Code,
		Name:               req.Name,
		Status:             DiscountStatusEditable,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		UserLimit:          req.UserLimit,
		CycleLimit:         req.CycleLimit,
		SubscriptionLimit:  req.SubscriptionLimit,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		CreateTime:         gtime.Now().Timestamp(),
	}
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(one).OmitNil().Insert(one)
	return err
}

func EditMerchantDiscountCode(ctx context.Context, req *CreateDiscountCodeInternalReq) error {
	utility.Assert(req.Code != "", "invalid Code")
	one := query.GetDiscountByCode(ctx, req.MerchantId, req.Code)
	utility.Assert(one != nil, "Code not found :"+req.Code)
	utility.Assert(one.Status == DiscountStatusEditable, "Code not editable :"+req.Code)
	utility.Assert(req.BillingType == DiscountBillingTypeOnetime || req.BillingType == DiscountBillingTypeRecurring, "invalid billingType, 1-one-time, 2-recurring")
	utility.Assert(req.DiscountType == DiscountTypePercentage || req.DiscountType == DiscountTypeFixedAmount, "invalid billingType, 1-percentage, 2-fixed_amount")
	utility.Assert(req.UserLimit >= 0, "invalid UserLimit")
	utility.Assert(req.SubscriptionLimit >= 0, "invalid SubscriptionLimit")
	utility.Assert(req.StartTime >= gtime.Now().Timestamp(), "startTime should greater then time now")
	utility.Assert(req.EndTime >= req.StartTime, "startTime should lower then endTime")
	req.Currency = strings.ToUpper(req.Currency)
	if req.DiscountType == DiscountTypePercentage {
		utility.Assert(req.DiscountPercentage > 0 && req.DiscountPercentage < 10000, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount == 0, "invalid discountAmount")
		utility.Assert(len(req.Currency) == 0, "invalid Currency")
	} else if req.DiscountType == DiscountTypeFixedAmount {
		utility.Assert(req.DiscountPercentage == 0, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount >= 0, "invalid discountAmount")
		utility.Assert(len(req.Currency) >= 0, "invalid Currency")
	}

	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().Name:               req.Name,
		dao.MerchantDiscountCode.Columns().BillingType:        req.BillingType,
		dao.MerchantDiscountCode.Columns().DiscountType:       req.DiscountType,
		dao.MerchantDiscountCode.Columns().DiscountAmount:     req.DiscountAmount,
		dao.MerchantDiscountCode.Columns().DiscountPercentage: req.DiscountPercentage,
		dao.MerchantDiscountCode.Columns().Currency:           req.Currency,
		dao.MerchantDiscountCode.Columns().UserLimit:          req.UserLimit,
		dao.MerchantDiscountCode.Columns().CycleLimit:         req.CycleLimit,
		dao.MerchantDiscountCode.Columns().SubscriptionLimit:  req.SubscriptionLimit,
		dao.MerchantDiscountCode.Columns().StartTime:          req.StartTime,
		dao.MerchantDiscountCode.Columns().EndTime:            req.EndTime,
		dao.MerchantDiscountCode.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	return err
}

func ActivateMerchantDiscountCode(ctx context.Context, merchantId uint64, code string) error {
	one := query.GetDiscountByCode(ctx, merchantId, code)
	utility.Assert(one != nil, "Code not found :"+code)
	if one.Status == DiscountStatusActive {
		return nil
	} else if one.Status == DiscountStatusExpired {
		return gerror.New("Code is expired")
	}
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().Status:    DiscountStatusActive,
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	return err
}

func DeactivateMerchantDiscountCode(ctx context.Context, merchantId uint64, code string) error {
	one := query.GetDiscountByCode(ctx, merchantId, code)
	utility.Assert(one != nil, "Code not found :"+code)
	if one.Status == DiscountStatusDeActive {
		return nil
	} else if one.Status != DiscountStatusActive {
		return gerror.New("Code is not active status")
	}
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().Status:    DiscountStatusDeActive,
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	return err
}

func DeleteMerchantDiscountCode(ctx context.Context, merchantId uint64, code string) error {
	one := query.GetDiscountByCode(ctx, merchantId, code)
	utility.Assert(one != nil, "Code not found :"+code)
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	return err
}

func HardDeleteMerchantDiscountCode(ctx context.Context, merchantId uint64, code string) error {
	utility.Assert(merchantId > 0, "invalid MerchantId")
	utility.Assert(len(code) > 0, "invalid Code")
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Where(dao.MerchantDiscountCode.Columns().Code, code).Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).Delete()
	return err
}
