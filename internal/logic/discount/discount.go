package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type CreateDiscountCodeInternalReq struct {
	Id                 uint64                 `json:"id"                 description:"The discount's Id"`
	Type               int                    `json:"type"               description:"type, 1-external discount code"` // type, 1-external discount code
	MerchantId         uint64                 `json:"MerchantId"        description:"MerchantId"`
	Code               string                 `json:"Code"              description:"Code"`
	Name               string                 `json:"name"              description:"name"`                                                                        // name
	BillingType        int                    `json:"billingType"       description:"billing_type, 1-one-time, 2-recurring"`                                       // billing_type, 1-one-time, 2-recurring
	DiscountType       int                    `json:"discountType"      description:"discount_type, 1-percentage, 2-fixed_amount"`                                 // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64                  `json:"discountAmount"    description:"amount of discount, available when discount_type is fixed_amount"`            // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64                  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string                 `json:"Currency"          description:"Currency of discount, available when discount_type is fixed_amount"`          // Currency of discount, available when discount_type is fixed_amount
	CycleLimit         int                    `json:"cycleLimit"         description:"the count limitation of subscription cycle , 0-no limit"`                    // the count limitation of subscription cycle , 0-no limit
	SubscriptionLimit  int                    `json:"subscriptionLimit" description:"the limit of every subscription apply, 0-unlimited"`                          // the limit of every subscription apply, 0-unlimited
	StartTime          *int64                 `json:"startTime"         description:"start of discount available utc time"`                                        // start of discount available utc time
	EndTime            *int64                 `json:"endTime"           description:"end of discount available utc time"`                                          // end of discount available utc time
	Quantity           *uint64                `json:"quantity"           description:"Quantity of code"`
	PlanIds            []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Advance            *bool                  `json:"advance"            description:"AdvanceConfig, 0-false,1-true, will enable all advance config if set true"` // AdvanceConfig,  0-false,1-true, will enable all advance config if set 1
	UserScope          *int                   `json:"userScope"  dc:"AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew"`
	UpgradeOnly        *bool                  `json:"upgradeOnly"  dc:"AdvanceConfig, true or false, will forbid for all except upgrade action if set true" `
	UpgradeLongerOnly  *bool                  `json:"upgradeLongPlanOnly"  dc:"AdvanceConfig, true or false, will forbid for all except upgrade to longer plan if set true" `
	UserLimit          *int                   `json:"userLimit"         dc:"AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited"`
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func NewMerchantDiscountCode(ctx context.Context, req *CreateDiscountCodeInternalReq) (*entity.MerchantDiscountCode, error) {
	utility.Assert(req.Code != "", "invalid Code")
	one := query.GetDiscountByCode(ctx, req.MerchantId, req.Code)
	utility.Assert(one == nil, "exist Code:"+req.Code)
	utility.Assert(req.BillingType == consts.DiscountBillingTypeOnetime || req.BillingType == consts.DiscountBillingTypeRecurring, "invalid billingType, 1-one-time, 2-recurring")
	utility.Assert(req.DiscountType == consts.DiscountTypePercentage || req.DiscountType == consts.DiscountTypeFixedAmount, "invalid billingType, 1-percentage, 2-fixed_amount")
	utility.Assert(req.SubscriptionLimit >= 0, "invalid SubscriptionLimit")
	//utility.Assert(req.StartTime >= gtime.Now().Timestamp(), "startTime should greater than time now")
	utility.Assert(req.StartTime != nil && req.EndTime != nil, "startTime and endTime should not be nil")
	utility.Assert(*req.EndTime >= *req.StartTime, "startTime should lower than endTime")
	req.Currency = strings.ToUpper(req.Currency)
	if req.DiscountType == consts.DiscountTypePercentage {
		utility.Assert(req.DiscountPercentage >= 0 && req.DiscountPercentage <= 10000, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount == 0, "invalid discountAmount")
		//utility.Assert(len(req.Currency) == 0, "invalid Currency")
		req.Currency = ""
	} else if req.DiscountType == consts.DiscountTypeFixedAmount {
		utility.Assert(req.DiscountPercentage == 0, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount >= 0, "invalid discountAmount")
		utility.Assert(len(req.Currency) >= 0, "invalid Currency")
	}

	// advanceConfig
	userLimit := 0
	if req.UserLimit != nil {
		utility.Assert(*req.UserLimit >= 0, "invalid UserLimit")
		userLimit = *req.UserLimit
	}

	var quantity int64 = 0
	if req.Quantity != nil {
		quantity = int64(*req.Quantity)
	}
	advance := 0
	if req.Advance != nil && *req.Advance {
		advance = 1
	} else if req.Advance != nil && !*req.Advance {
		advance = 0
	}
	userScope := 0
	if req.UserScope != nil {
		userScope = *req.UserScope
	}
	upgradeOnly := 0
	if req.UpgradeOnly != nil && *req.UpgradeOnly {
		upgradeOnly = 1
	} else if req.UpgradeOnly != nil && !*req.UpgradeOnly {
		upgradeOnly = 0
	}
	upgradeLongerOnly := 0
	if req.UpgradeLongerOnly != nil && *req.UpgradeLongerOnly {
		upgradeLongerOnly = 1
	} else if req.UpgradeLongerOnly != nil && !*req.UpgradeLongerOnly {
		upgradeLongerOnly = 0
	}
	utility.Assert(upgradeOnly+upgradeLongerOnly <= 1, "invalid UpgradeOnly and UpgradeLongerOnly, You can only choose one of two")
	one = &entity.MerchantDiscountCode{
		MerchantId:         req.MerchantId,
		Code:               req.Code,
		Name:               req.Name,
		Status:             consts.DiscountStatusEditable,
		BillingType:        req.BillingType,
		DiscountType:       req.DiscountType,
		DiscountAmount:     req.DiscountAmount,
		Type:               req.Type,
		DiscountPercentage: req.DiscountPercentage,
		Currency:           req.Currency,
		CycleLimit:         req.CycleLimit,
		SubscriptionLimit:  req.SubscriptionLimit,
		StartTime:          *req.StartTime,
		EndTime:            *req.EndTime,
		Quantity:           quantity,
		MetaData:           utility.MarshalToJsonString(req.Metadata),
		PlanIds:            utility.IntListToString(req.PlanIds),
		CreateTime:         gtime.Now().Timestamp(),
		Advance:            advance,
		UserLimit:          userLimit,
		UserScope:          userScope,
		UpgradeOnly:        upgradeOnly,
		UpgradeLongerOnly:  upgradeLongerOnly,
	}
	result, err := dao.MerchantDiscountCode.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`NewMerchantDiscountCode insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return one, err
}

func EditMerchantDiscountCode(ctx context.Context, req *CreateDiscountCodeInternalReq) (*entity.MerchantDiscountCode, error) {
	utility.Assert(req.Id > 0, "invalid Id")
	one := query.GetDiscountById(ctx, req.Id)
	utility.Assert(one != nil, "Discount not found :"+strconv.FormatUint(req.Id, 10))
	utility.Assert(one.MerchantId == req.MerchantId, "Discount merchant not match :"+req.Code)

	advance := one.Advance
	if req.Advance != nil && *req.Advance {
		advance = 1
	} else if req.Advance != nil && !*req.Advance {
		advance = 0
	}
	upgradeOnly := one.UpgradeOnly
	if req.UpgradeOnly != nil && *req.UpgradeOnly {
		upgradeOnly = 1
	} else if req.UpgradeOnly != nil && !*req.UpgradeOnly {
		upgradeOnly = 0
	}
	upgradeLongerOnly := one.UpgradeLongerOnly
	if req.UpgradeLongerOnly != nil && *req.UpgradeLongerOnly {
		upgradeLongerOnly = 1
	} else if req.UpgradeLongerOnly != nil && !*req.UpgradeLongerOnly {
		upgradeLongerOnly = 0
	}
	utility.Assert(upgradeOnly+upgradeLongerOnly <= 1, "invalid UpgradeOnly and UpgradeLongerOnly, You can only choose one of two")
	//edit after activate
	if one.Status > consts.DiscountStatusEditable {
		utility.Assert((req.StartTime != nil && req.EndTime != nil) || req.PlanIds != nil, "startTime&endTime or planIds should not be nil")
		var planIdsStr *string = nil
		if req.PlanIds != nil {
			planIdsStr = unibee.String(utility.IntListToString(req.PlanIds))
		}
		_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
			dao.MerchantDiscountCode.Columns().Advance:           advance,
			dao.MerchantDiscountCode.Columns().UserLimit:         req.UserLimit,
			dao.MerchantDiscountCode.Columns().UserScope:         req.UserScope,
			dao.MerchantDiscountCode.Columns().UpgradeOnly:       upgradeOnly,
			dao.MerchantDiscountCode.Columns().UpgradeLongerOnly: upgradeLongerOnly,
			dao.MerchantDiscountCode.Columns().StartTime:         req.StartTime,
			dao.MerchantDiscountCode.Columns().EndTime:           req.EndTime,
			dao.MerchantDiscountCode.Columns().PlanIds:           planIdsStr,
			dao.MerchantDiscountCode.Columns().MetaData:          utility.MarshalToJsonString(utility.MergeMetadata(one.MetaData, &req.Metadata)),
			dao.MerchantDiscountCode.Columns().GmtModify:         gtime.Now(),
		}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			err = gerror.Newf(`EditMerchantDiscountCode update after activate failure %s`, err)
			return nil, err
		}
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
			Content:        "EditAfterActivate",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   one.Code,
		}, err)
		one = query.GetDiscountById(ctx, req.Id)
		return one, nil
	}

	utility.Assert(one.Type == 0, "Edit not available for external code :"+req.Code)
	utility.Assert(one.Status == consts.DiscountStatusEditable, "Code not editable :"+req.Code)
	utility.Assert(req.BillingType == consts.DiscountBillingTypeOnetime || req.BillingType == consts.DiscountBillingTypeRecurring, "invalid billingType, 1-one-time, 2-recurring")
	utility.Assert(req.DiscountType == consts.DiscountTypePercentage || req.DiscountType == consts.DiscountTypeFixedAmount, "invalid billingType, 1-percentage, 2-fixed_amount")
	if req.UserLimit != nil {
		utility.Assert(*req.UserLimit >= 0, "invalid UserLimit")
	}
	utility.Assert(req.SubscriptionLimit >= 0, "invalid SubscriptionLimit")
	//utility.Assert(req.StartTime >= gtime.Now().Timestamp(), "startTime should greater than time now")
	utility.Assert(req.StartTime != nil && req.EndTime != nil, "startTime and endTime should not be nil")
	utility.Assert(*req.EndTime >= *req.StartTime, "startTime should lower than endTime")
	req.Currency = strings.ToUpper(req.Currency)
	if req.DiscountType == consts.DiscountTypePercentage {
		utility.Assert(req.DiscountPercentage >= 0 && req.DiscountPercentage <= 10000, "invalid DiscountPercentage")
		utility.Assert(req.DiscountAmount == 0, "invalid discountAmount")
		//utility.Assert(len(req.Currency) == 0, "invalid Currency")
		req.Currency = ""
	} else if req.DiscountType == consts.DiscountTypeFixedAmount {
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
		dao.MerchantDiscountCode.Columns().CycleLimit:         req.CycleLimit,
		dao.MerchantDiscountCode.Columns().SubscriptionLimit:  req.SubscriptionLimit,
		dao.MerchantDiscountCode.Columns().StartTime:          req.StartTime,
		dao.MerchantDiscountCode.Columns().EndTime:            req.EndTime,
		dao.MerchantDiscountCode.Columns().Quantity:           req.Quantity,
		dao.MerchantDiscountCode.Columns().PlanIds:            utility.IntListToString(req.PlanIds),
		dao.MerchantDiscountCode.Columns().MetaData:           utility.MarshalToJsonString(req.Metadata),
		dao.MerchantDiscountCode.Columns().GmtModify:          gtime.Now(),
		dao.MerchantDiscountCode.Columns().Advance:            advance,
		dao.MerchantDiscountCode.Columns().UserLimit:          req.UserLimit,
		dao.MerchantDiscountCode.Columns().UserScope:          req.UserScope,
		dao.MerchantDiscountCode.Columns().UpgradeOnly:        upgradeOnly,
		dao.MerchantDiscountCode.Columns().UpgradeLongerOnly:  upgradeLongerOnly,
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		err = gerror.Newf(`EditMerchantDiscountCode update failure %s`, err)
		return nil, err
	}
	one = query.GetDiscountById(ctx, one.Id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return one, err
}

func ActivateMerchantDiscountCode(ctx context.Context, merchantId uint64, id uint64) error {
	utility.Assert(id > 0, "invalid Id")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	if one.Status == consts.DiscountStatusActive {
		return nil
	}
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().Status:    consts.DiscountStatusActive,
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        "Activate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return err
}

func DeactivateMerchantDiscountCode(ctx context.Context, merchantId uint64, id uint64) error {
	utility.Assert(id > 0, "invalid Id")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	if one.Status == consts.DiscountStatusDeactivate {
		return nil
	} else if one.Status != consts.DiscountStatusActive {
		return gerror.New("Code is not active status")
	}
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().Status:    consts.DiscountStatusDeactivate,
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        "Deactivate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return err
}

func DeleteMerchantDiscountCode(ctx context.Context, merchantId uint64, id uint64) error {
	utility.Assert(id > 0, "invalid Id")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	utility.Assert(one.Type == 0, "Delete not available for external code :"+strconv.FormatUint(id, 10))
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
		dao.MerchantDiscountCode.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return err
}

func QuantityIncrement(ctx context.Context, merchantId uint64, id uint64, amount uint64) error {
	utility.Assert(id > 0, "invalid Id")
	utility.Assert(amount > 0, "invalid Increase Value")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	utility.Assert(one.Type == 0, "Delete not available for external code :"+strconv.FormatUint(id, 10))
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Where(dao.MerchantDiscountCode.Columns().Id, id).Increment(dao.MerchantDiscountCode.Columns().Quantity, amount)
	if err != nil {
		return err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        fmt.Sprintf("QuantityIncrement(%d)", amount),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return err
}

func QuantityDecrement(ctx context.Context, merchantId uint64, id uint64, amount uint64) error {
	utility.Assert(id > 0, "invalid Id")
	utility.Assert(amount > 0, "invalid Decrease Value")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	utility.Assert(one.Type == 0, "Delete not available for external code :"+strconv.FormatUint(id, 10))
	utility.Assert(uint64(one.Quantity) >= amount, "decrease value should not greater than code quantity")
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Where(dao.MerchantDiscountCode.Columns().Id, id).Decrement(dao.MerchantDiscountCode.Columns().Quantity, amount)
	if err != nil {
		return err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("DiscountCode(%s)", one.Code),
		Content:        fmt.Sprintf("QuantityDecrement(%d)", amount),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   one.Code,
	}, err)
	return err
}

func HardDeleteMerchantDiscountCode(ctx context.Context, merchantId uint64, id uint64) error {
	utility.Assert(merchantId > 0, "invalid MerchantId")
	utility.Assert(id > 0, "invalid Id")
	one := query.GetDiscountById(ctx, id)
	utility.Assert(one != nil, "discount not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Discount merchant not match :"+strconv.FormatUint(id, 10))
	_, err := dao.MerchantDiscountCode.Ctx(ctx).Where(dao.MerchantDiscountCode.Columns().Id, id).Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).Delete()
	return err
}

func CreateExternalDiscount(ctx context.Context, merchantId uint64, userId uint64, source string, param *bean.ExternalDiscountParam, currency string, timeNow int64) *entity.MerchantDiscountCode {
	var cycleLimit = 0
	var endTime int64 = 0
	var BillingType = consts.DiscountBillingTypeOnetime
	if param.Recurring != nil && *param.Recurring {
		BillingType = consts.DiscountBillingTypeRecurring
		if param.CycleLimit != nil {
			cycleLimit = *param.CycleLimit
		}
		utility.Assert(cycleLimit >= 0, "invalid cycleLimit")
		if param.EndTime != nil {
			endTime = *param.EndTime
		}
		utility.Assert(endTime >= timeNow, "invalid endTime")
	} else {
		utility.Assert(param.CycleLimit == nil, "cycleLimit not available as recurring not enable")
		utility.Assert(param.EndTime == nil, "endTime not available as recurring not enable")
		endTime = timeNow + 600
	}
	var discountType = consts.DiscountTypePercentage
	var discountAmount int64 = 0
	var discountPercentage int64 = 0

	if param.DiscountAmount != nil && *param.DiscountAmount > 0 {
		discountType = consts.DiscountTypeFixedAmount
		discountAmount = *param.DiscountAmount

	} else if param.DiscountPercentage != nil && *param.DiscountPercentage > 0 {
		discountType = consts.DiscountTypePercentage
		discountPercentage = *param.DiscountPercentage
		utility.Assert(discountPercentage > 0 && discountPercentage <= 10000, "invalid discountPercentage")
	} else {
		utility.Assert(true, "one of discountAmount or discountPercentage should specified")
	}
	one, err := NewMerchantDiscountCode(ctx, &CreateDiscountCodeInternalReq{
		MerchantId:         merchantId,
		Code:               fmt.Sprintf("excode_%d_%d_%s_%d%s", merchantId, userId, source, utility.CurrentTimeMillis(), utility.GenerateRandomAlphanumeric(8)),
		Name:               fmt.Sprintf("excode_for_plan_%s_subscription", source),
		BillingType:        BillingType,
		DiscountType:       discountType,
		DiscountAmount:     discountAmount,
		Type:               1,
		DiscountPercentage: discountPercentage,
		Currency:           currency,
		CycleLimit:         cycleLimit,
		StartTime:          unibee.Int64(timeNow - 10),
		EndTime:            unibee.Int64(endTime),
		Metadata:           param.Metadata,
	})
	utility.AssertError(err, "Create discount error")
	err = ActivateMerchantDiscountCode(ctx, merchantId, one.Id)
	utility.AssertError(err, "Create discount error")
	return one
}
