package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/api/user/vat"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/internal/logic/subscription/handler"
	"go-oversea-pay/internal/logic/vat_gateway"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
	"time"
)

type SubscriptionCreatePrepareInternalRes struct {
	Plan              *entity.SubscriptionPlan           `json:"planId"`
	Quantity          int64                              `json:"quantity"`
	PlanChannel       *entity.ChannelPlan                `json:"planChannel"`
	PayChannel        *entity.MerchantChannelConfig      `json:"payChannel"`
	MerchantInfo      *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount       int64                              `json:"totalAmount"                `
	Currency          string                             `json:"currency"              `
	VatCountryCode    string                             `json:"vatCountryCode"              `
	VatCountryName    string                             `json:"vatCountryName"              `
	VatNumber         string                             `json:"vatNumber"              `
	VatNumberValidate *ro.ValidResult                    `json:"vatNumberValidate"              `
	TaxScale          int64                              `json:"taxScale"              `
	VatVerifyData     string                             `json:"vatVerifyData"              `
	Invoice           *ro.InvoiceDetailSimplify          `json:"invoice"`
	UserId            int64                              `json:"userId" `
	Email             string                             `json:"email" `
	VatCountryRate    *ro.VatCountryRate                 `json:"vatCountryRate" `
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo, channelId int64) []*ro.SubscriptionPlanAddonRo {
	var addons []*ro.SubscriptionPlanAddonRo
	var totalAddonIds []int64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.SubscriptionPlan
	if len(totalAddonIds) > 0 {
		//查询所有 Plan
		err := dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).Scan(&allAddonList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				//所有 Addon 项目必须要能查到
				//类型是 Addon
				//未删除
				//数量大于 0
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				planChannel := query.GetPlanChannel(ctx, int64(mapPlans[param.AddonPlanId].Id), channelId) // todo mark for 循环内调用 需做缓存，此数据基本不会变化,或者方案 2 使用 channelId 合并查询
				utility.Assert(len(planChannel.ChannelPlanId) > 0, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanId invalid", param.AddonPlanId, channelId))
				utility.Assert(planChannel.Status == consts.PlanChannelStatusActive, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanStatus not active", param.AddonPlanId, channelId))
				addons = append(addons, &ro.SubscriptionPlanAddonRo{
					Quantity:         param.Quantity,
					AddonPlan:        mapPlans[param.AddonPlanId],
					AddonPlanChannel: planChannel,
				})
			}
		}
	}
	return addons
}

func VatNumberValidate(ctx context.Context, req *vat.NumberValidateReq) (*vat.NumberValidateRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, req.MerchantId, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	if vatNumberValidate.Valid {
		vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, req.MerchantId, vatNumberValidate.CountryCode)
		utility.Assert(err == nil, fmt.Sprintf("vatNumber vatCountryCode check error:%s", err))
		utility.Assert(vatCountryRate != nil, fmt.Sprintf("vatNumber not found for countryCode:%v", vatNumberValidate.CountryCode))
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}

func SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (*SubscriptionCreatePrepareInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.ChannelId > 0, "ConfirmChannelId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	email := ""
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
		email = _interface.BizCtx().Get(ctx).User.Email
	}
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, "user not found")

	var err error
	utility.Assert(query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, req.UserId, merchantInfo.Id) == nil, "another active subscription find, only one subscription can create")

	//vat
	utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil, "Merchant Vat VATGateway not setup")
	var vatCountryCode = req.VatCountryCode
	var standardTaxScale int64 = 0
	var vatCountryName = ""
	var vatCountryRate *ro.VatCountryRate
	var vatNumberValidate *ro.ValidResult

	if len(req.VatNumber) > 0 {
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.VatNumber, "")
		if err != nil {
			return nil, err
		}
		utility.Assert(vatNumberValidate.Valid, fmt.Sprintf("vat number validate failure:"+req.VatNumber))
		vatCountryCode = vatNumberValidate.CountryCode
	}

	if len(vatCountryCode) == 0 && len(user.CountryCode) > 0 {
		//if not set , default user countryCode
		vatCountryCode = user.CountryCode
		req.VatCountryCode = user.CountryCode
	}

	if len(vatCountryCode) > 0 {
		vatCountryRate, err = vat_gateway.QueryVatCountryRateByMerchant(ctx, merchantInfo.Id, vatCountryCode)
		utility.Assert(err == nil, fmt.Sprintf("vat vatCountryCode check error:%s", err))
		utility.Assert(vatCountryRate != nil, fmt.Sprintf("vat not found for countryCode:%v", vatCountryCode))
		vatCountryName = vatCountryRate.CountryName
		if vatNumberValidate == nil || !vatNumberValidate.Valid {
			standardTaxScale = vatCountryRate.StandardTaxPercentage
		}
	}

	//设置Default值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	var currency = plan.Currency
	var TotalAmountExcludingTax = plan.Amount * req.Quantity

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
		TotalAmountExcludingTax = TotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		Currency:      currency,
		PlanId:        req.PlanId,
		Quantity:      req.Quantity,
		AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
		TaxScale:      standardTaxScale,
	})

	return &SubscriptionCreatePrepareInternalRes{
		Plan:              plan,
		Quantity:          req.Quantity,
		PlanChannel:       planChannel,
		PayChannel:        payChannel,
		MerchantInfo:      merchantInfo,
		AddonParams:       req.AddonParams,
		Addons:            addons,
		TotalAmount:       invoice.TotalAmount,
		Currency:          currency,
		VatCountryCode:    vatCountryCode,
		VatCountryName:    vatCountryName,
		VatNumber:         req.VatNumber,
		VatNumberValidate: vatNumberValidate,
		VatVerifyData:     utility.MarshalToJsonString(vatNumberValidate),
		TaxScale:          standardTaxScale,
		UserId:            req.UserId,
		Email:             email,
		Invoice:           invoice,
		VatCountryRate:    vatCountryRate,
	}, nil
}

func SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (*subscription.SubscriptionCreateRes, error) {

	prepare, err := SubscriptionCreatePreview(ctx, &subscription.SubscriptionCreatePreviewReq{
		PlanId:         req.PlanId,
		Quantity:       req.Quantity,
		ChannelId:      req.ChannelId,
		UserId:         req.UserId,
		AddonParams:    req.AddonParams,
		VatCountryCode: req.VatCountryCode,
		VatNumber:      req.VatNumber,
	})
	if err != nil {
		return nil, err
	}
	utility.Assert(len(prepare.VatCountryCode) > 0, "CountryCode Needed")

	//校验
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")

	var subType = consts.SubTypeDefault
	if consts.SubscriptionCycleUnderUniBeeControl {
		subType = consts.SubTypeUniBeeControl
	}

	var billingCycleAnchor = gtime.Now()
	var currentTimeStart = billingCycleAnchor
	var currentTimeEnd = billingCycleAnchor
	if strings.Compare(strings.ToLower(prepare.Plan.IntervalUnit), "day") == 0 {
		currentTimeEnd = currentTimeEnd.AddDate(0, 0, prepare.Plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(prepare.Plan.IntervalUnit), "week") == 0 {
		currentTimeEnd = currentTimeEnd.AddDate(0, 0, 7*prepare.Plan.IntervalCount)
	} else if strings.Compare(strings.ToLower(prepare.Plan.IntervalUnit), "month") == 0 {
		currentTimeEnd = currentTimeEnd.AddDate(0, prepare.Plan.IntervalCount, 0)
	} else if strings.Compare(strings.ToLower(prepare.Plan.IntervalUnit), "year") == 0 {
		currentTimeEnd = currentTimeEnd.AddDate(prepare.Plan.IntervalCount, 0, 0)
	}

	one := &entity.Subscription{
		MerchantId:         prepare.MerchantInfo.Id,
		Type:               subType,
		PlanId:             int64(prepare.Plan.Id),
		ChannelId:          prepare.PlanChannel.ChannelId,
		UserId:             prepare.UserId,
		Quantity:           prepare.Quantity,
		Amount:             prepare.TotalAmount,
		Currency:           prepare.Currency,
		AddonData:          utility.MarshalToJsonString(prepare.AddonParams),
		SubscriptionId:     utility.CreateSubscriptionId(),
		Status:             consts.SubStatusCreate,
		CustomerEmail:      prepare.Email,
		ReturnUrl:          req.ReturnUrl,
		Data:               "", //额外参数配置
		VatNumber:          prepare.VatNumber,
		VatVerifyData:      prepare.VatVerifyData,
		CountryCode:        prepare.VatCountryCode,
		TaxScale:           prepare.TaxScale,
		CurrentPeriodStart: currentTimeStart.Timestamp(),
		CurrentPeriodEnd:   currentTimeEnd.Timestamp(),
		BillingCycleAnchor: billingCycleAnchor.Timestamp(),
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	var createRes *ro.ChannelCreateSubscriptionInternalResp
	if consts.SubscriptionCycleUnderUniBeeControl {
		user := query.GetUserAccountById(ctx, uint64(one.UserId))
		var mobile = ""
		var firstName = ""
		var lastName = ""
		var gender = ""
		if user != nil {
			mobile = user.Mobile
			firstName = user.FirstName
			lastName = user.LastName
			gender = user.Gender
		}
		createPaymentResult, err := service.DoChannelPay(ctx, &ro.CreatePayContext{
			PayChannel: prepare.PayChannel,
			Pay: &entity.Payment{
				SubscriptionId: one.SubscriptionId,
				BizId:          one.SubscriptionId,
				BizType:        consts.BIZ_TYPE_SUBSCRIPTION,
				UserId:         prepare.UserId,
				ChannelId:      int64(prepare.PayChannel.Id),
				TotalAmount:    prepare.Invoice.TotalAmount,
				Currency:       prepare.Currency,
				CountryCode:    prepare.VatCountryCode,
				MerchantId:     prepare.MerchantInfo.Id,
				CompanyId:      prepare.MerchantInfo.CompanyId,
			},
			Platform:      "WEB",
			DeviceType:    "Web",
			ShopperUserId: strconv.FormatInt(one.UserId, 10),
			ShopperEmail:  prepare.Email,
			ShopperLocale: "en",
			Mobile:        mobile,
			Invoice:       prepare.Invoice,
			ShopperName: &v1.OutShopperName{
				FirstName: firstName,
				LastName:  lastName,
				Gender:    gender,
			},
			MediaData:              map[string]string{"BillingReason": "SubscriptionCreate"},
			MerchantOrderReference: one.SubscriptionId,
			PayMethod:              1, //automatic
			DaysUtilDue:            5, //one day expire
		})
		if err != nil {
			return nil, err
		}
		createRes = &ro.ChannelCreateSubscriptionInternalResp{
			ChannelSubscriptionId: createPaymentResult.PaymentId,
			Data:                  utility.MarshalToJsonString(createPaymentResult),
			Link:                  createPaymentResult.Link,
			Paid:                  createPaymentResult.Status == consts.PAY_SUCCESS,
		}
	} else {
		createRes, err = out.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionCreate(ctx, &ro.ChannelCreateSubscriptionInternalReq{
			Plan:           prepare.Plan,
			AddonPlans:     prepare.Addons,
			PlanChannel:    prepare.PlanChannel,
			VatCountryRate: prepare.VatCountryRate,
			Subscription:   one,
		})
		if err != nil {
			return nil, err
		}
	}

	//更新 Subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().ChannelSubscriptionId: createRes.ChannelSubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
		dao.Subscription.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return nil, gerror.Newf("SubscriptionCreate update err:%s", update)
	//}
	one.ChannelSubscriptionId = createRes.ChannelSubscriptionId
	one.Status = consts.PlanChannelStatusCreate
	one.Link = createRes.Link

	return &subscription.SubscriptionCreateRes{
		Subscription: one,
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}

type SubscriptionUpdatePrepareInternalRes struct {
	Subscription *entity.Subscription               `json:"subscription"`
	Plan         *entity.SubscriptionPlan           `json:"planId"`
	Quantity     int64                              `json:"quantity"`
	PlanChannel  *entity.ChannelPlan                `json:"planChannel"`
	PayChannel   *entity.MerchantChannelConfig      `json:"payChannel"`
	MerchantInfo *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams  []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons       []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount  int64                              `json:"totalAmount"                `
	Currency     string                             `json:"currency"              `
	UserId       int64                              `json:"userId" `
	OldPlan      *entity.SubscriptionPlan           `json:"oldPlan"`
	//OldPlanChannel    *entity.ChannelPlan    `json:"oldPlanChannel"`
	Invoice           *ro.InvoiceDetailSimplify `json:"invoice"`
	NextPeriodInvoice *ro.InvoiceDetailSimplify `json:"nextPeriodInvoice"`
	ProrationDate     int64                     `json:"prorationDate"`
	EffectImmediate   bool                      `json:"EffectImmediate"`
}

// SubscriptionUpdatePreview Default行为，升级订阅主方案不管总金额是否比之前高，都将按比例计算发票立即生效；降级订阅方案，次月生效；问题点，降级方案如果 addon 多可能的总金额可能比之前高
func SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq, prorationDate int64, merchantUserId int64) (res *SubscriptionUpdatePrepareInternalRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	//utility.Assert(sub.ChannelId == req.ConfirmChannelId, "channel not match")
	// todo mark addon binding check

	plan := query.GetPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.NewPlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	//试用期内不允许修改计划
	//utility.Assert(sub.TrialEnd < gtime.Now().Timestamp(), "subscription is in trial period ,should end trial before update")
	//设置了下周期取消不允许修改计划
	utility.Assert(sub.CancelAtPeriodEnd == 0, "subscription cannot be update as it will cancel at period end, should resume subscription first")

	//设置Default值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	var currency = sub.Currency
	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
	}
	oldPlan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")

	if req.NewPlanId != sub.PlanId {
		utility.Assert(oldPlan.IntervalUnit == plan.IntervalUnit, "newPlan must have same recurring interval to old")
		utility.Assert(oldPlan.IntervalCount == plan.IntervalCount, "newPlan must have same recurring interval to old")
	}
	//暂时不开放不同通道升级功能 todo mark
	//oldPlanChannel := query.GetPlanChannel(ctx, int64(oldPlan.Id), sub.ChannelId)
	//utility.Assert(oldPlanChannel != nil, "oldPlanChannel not found")

	var effectImmediate = false
	//升降级判断逻辑，升级设置payImmediate=true，保障马上能够生效；降级payImmediate=false,下周期生效
	//情况 1，NewPlan单价大于 OldPlan 单价，判断为升级，忽略Quantity 和 addon 变更
	//情况 2，NewPlan单价小于 OldPlan 单价，判断为降级，忽略Quantity 和 addon 变更
	//情况 3，NewPlan总价大于 OldPlan总价，判断为升级
	//情况 4，NewPlan总价小于 OldPlan总价，判断为降级
	//情况 5，NewPlan总价等于 OldPlan总价，则看 Addon 的变化，如果 addon 有数量增加情况或者新增 addon 情况为升级，否则降级

	if plan.Amount > oldPlan.Amount || plan.Amount*req.Quantity > oldPlan.Amount*sub.Quantity {
		effectImmediate = true
	} else if plan.Amount < oldPlan.Amount || plan.Amount*req.Quantity < oldPlan.Amount*sub.Quantity {
		effectImmediate = false
	} else {
		var oldAddonParams []*ro.SubscriptionPlanAddonParamRo
		err = utility.UnmarshalFromJsonString(sub.AddonData, &oldAddonParams)
		utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
		var oldAddonMap = make(map[int64]int64)
		for _, oldAddon := range oldAddonParams {
			if _, ok := oldAddonMap[oldAddon.AddonPlanId]; ok {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddonMap[oldAddon.AddonPlanId] + oldAddon.Quantity
			} else {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddon.Quantity
			}
		}
		var newAddonMap = make(map[int64]int64)
		for _, newAddon := range req.AddonParams {
			if _, ok := newAddonMap[newAddon.AddonPlanId]; ok {
				newAddonMap[newAddon.AddonPlanId] = newAddonMap[newAddon.AddonPlanId] + newAddon.Quantity
			} else {
				newAddonMap[newAddon.AddonPlanId] = newAddon.Quantity
			}
		}
		for newAddonPlanId, newAddonQuantity := range newAddonMap {
			if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
				if oldAddonQuantity < newAddonQuantity {
					//数量有增加,视为升级
					effectImmediate = true
					break
				}
			} else {
				//新增,视为升级
				effectImmediate = true
				break
			}
		}
		//如果是降级，校验是否有变化
		var changed = false
		if len(oldAddonMap) != len(newAddonMap) {
			changed = true
		} else {
			for newAddonPlanId, newAddonQuantity := range newAddonMap {
				if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
					if oldAddonQuantity != newAddonQuantity {
						//数量不等
						changed = true
						break
					}
				} else {
					//新增
					changed = true
					break
				}
			}
		}
		utility.Assert(changed, "subscription update should have plan or addons changed")
	}

	if req.WithImmediateEffect > 0 {
		utility.Assert(req.WithImmediateEffect == 1 || req.WithImmediateEffect == 2, "WithImmediateEffect should be 1 or 2")
		if req.WithImmediateEffect == 1 {
			effectImmediate = true
		} else {
			effectImmediate = false
		}
	}

	var nextPeriodStart = sub.CurrentPeriodEnd
	if sub.TrialEnd > sub.CurrentPeriodEnd {
		nextPeriodStart = sub.TrialEnd
	}
	var nextPeriodEnd = nextPeriodStart + (sub.CurrentPeriodEnd - sub.CurrentPeriodStart)

	var totalAmount int64
	var prorationInvoice *ro.InvoiceDetailSimplify
	//var nextPeriodTotalAmount int64
	var nextPeriodInvoice *ro.InvoiceDetailSimplify
	if effectImmediate {
		if consts.ProrationUsingUniBeeCompute {
			// Generate Proration Invoice Previe
			nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				Currency:      sub.Currency,
				PlanId:        req.NewPlanId,
				Quantity:      req.Quantity,
				AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
				TaxScale:      sub.TaxScale,
				PeriodStart:   nextPeriodStart,
				PeriodEnd:     nextPeriodEnd,
			})

			var oldAddonParams []*ro.SubscriptionPlanAddonParamRo
			err = utility.UnmarshalFromJsonString(sub.AddonData, &oldAddonParams)
			utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
			var oldProrationPlanParams []*invoice_compute.ProrationPlanParam
			oldProrationPlanParams = append(oldProrationPlanParams, &invoice_compute.ProrationPlanParam{
				PlanId:   sub.PlanId,
				Quantity: sub.Quantity,
			})
			for _, addonParam := range oldAddonParams {
				oldProrationPlanParams = append(oldProrationPlanParams, &invoice_compute.ProrationPlanParam{
					PlanId:   addonParam.AddonPlanId,
					Quantity: addonParam.Quantity,
				})
			}
			var newProrationPlanParams []*invoice_compute.ProrationPlanParam
			newProrationPlanParams = append(newProrationPlanParams, &invoice_compute.ProrationPlanParam{
				PlanId:   req.NewPlanId,
				Quantity: req.Quantity,
			})
			for _, addonParam := range req.AddonParams {
				newProrationPlanParams = append(newProrationPlanParams, &invoice_compute.ProrationPlanParam{
					PlanId:   addonParam.AddonPlanId,
					Quantity: addonParam.Quantity,
				})
			}

			if prorationDate == 0 {
				prorationDate = time.Now().Unix()
			}
			if prorationDate > sub.CurrentPeriodEnd || prorationDate < sub.CurrentPeriodStart {
				// after period end before trial end, also or sub data not sync todo mark
				prorationInvoice = &ro.InvoiceDetailSimplify{
					TotalAmount:                    0,
					TotalAmountExcludingTax:        0,
					Currency:                       sub.Currency,
					TaxAmount:                      0,
					SubscriptionAmount:             0,
					SubscriptionAmountExcludingTax: 0,
					Lines:                          make([]*ro.InvoiceItemDetailRo, 0),
					ProrationDate:                  prorationDate,
				}
			} else {
				prorationInvoice = invoice_compute.ComputeSubscriptionProrationInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					Currency:          sub.Currency,
					TaxScale:          sub.TaxScale,
					ProrationDate:     prorationDate,
					PeriodStart:       sub.CurrentPeriodStart,
					PeriodEnd:         sub.CurrentPeriodEnd,
					OldProrationPlans: oldProrationPlanParams,
					NewProrationPlans: newProrationPlanParams,
				})
			}
			prorationDate = prorationInvoice.ProrationDate
			totalAmount = prorationInvoice.TotalAmount
		} else {
			updatePreviewRes, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionUpdateProrationPreview(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
				Plan:            plan,
				Quantity:        req.Quantity,
				AddonPlans:      addons,
				PlanChannel:     planChannel,
				Subscription:    sub,
				ProrationDate:   prorationDate,
				EffectImmediate: effectImmediate,
			})
			if err != nil {
				return nil, err
			}
			utility.Assert(strings.Compare(updatePreviewRes.Currency, currency) == 0, fmt.Sprintf("preview currency not match for subscriptionId:%v preview currency:%s", sub.SubscriptionId, updatePreviewRes.Currency))

			prorationDate = updatePreviewRes.ProrationDate
			totalAmount = updatePreviewRes.TotalAmount
			if updatePreviewRes.Invoice.Lines == nil {
				updatePreviewRes.Invoice.Lines = make([]*ro.InvoiceItemDetailRo, 0)
			}
			prorationInvoice = &ro.InvoiceDetailSimplify{
				TotalAmount:                    updatePreviewRes.Invoice.TotalAmount,
				TotalAmountExcludingTax:        updatePreviewRes.Invoice.TotalAmountExcludingTax,
				Currency:                       updatePreviewRes.Invoice.Currency,
				TaxAmount:                      updatePreviewRes.Invoice.TaxAmount,
				SubscriptionAmount:             updatePreviewRes.Invoice.SubscriptionAmount,
				SubscriptionAmountExcludingTax: updatePreviewRes.Invoice.SubscriptionAmountExcludingTax,
				Lines:                          updatePreviewRes.Invoice.Lines,
			}

			utility.Assert(updatePreviewRes.NextPeriodInvoice.Lines != nil, "internal error, next_period_line is blank")

			//nextPeriodTotalAmount = updatePreviewRes.NextPeriodInvoice.TotalAmount
			nextPeriodInvoice = &ro.InvoiceDetailSimplify{
				TotalAmount:                    updatePreviewRes.NextPeriodInvoice.TotalAmount,
				TotalAmountExcludingTax:        updatePreviewRes.NextPeriodInvoice.TotalAmountExcludingTax,
				Currency:                       updatePreviewRes.NextPeriodInvoice.Currency,
				TaxAmount:                      updatePreviewRes.NextPeriodInvoice.TaxAmount,
				SubscriptionAmount:             updatePreviewRes.NextPeriodInvoice.SubscriptionAmount,
				SubscriptionAmountExcludingTax: updatePreviewRes.NextPeriodInvoice.SubscriptionAmountExcludingTax,
				Lines:                          updatePreviewRes.NextPeriodInvoice.Lines,
			}
		}
	} else {
		//Effect Next Period, Generate Invoice Preview
		var nextPeriodTotalAmountExcludingTax = plan.Amount * req.Quantity
		for _, addon := range addons {
			utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
			utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
			utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
			nextPeriodTotalAmountExcludingTax = nextPeriodTotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
		}

		var nextPeriodInvoiceItems []*ro.InvoiceItemDetailRo
		nextPeriodInvoiceItems = append(nextPeriodInvoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               currency,
			Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxScaleToInternalFloat(sub.TaxScale)),
			AmountExcludingTax:     req.Quantity * plan.Amount,
			Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxScaleToInternalFloat(sub.TaxScale)),
			UnitAmountExcludingTax: plan.Amount,
			Description:            plan.PlanName,
			Quantity:               req.Quantity,
		})
		for _, addon := range addons {
			nextPeriodInvoiceItems = append(nextPeriodInvoiceItems, &ro.InvoiceItemDetailRo{
				Currency:               currency,
				Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxScaleToInternalFloat(sub.TaxScale)),
				Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxScaleToInternalFloat(sub.TaxScale)),
				AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
				UnitAmountExcludingTax: addon.AddonPlan.Amount,
				Description:            addon.AddonPlan.PlanName,
				Quantity:               addon.Quantity,
			})
		}
		var nextPeriodTaxAmount = int64(float64(nextPeriodTotalAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(sub.TaxScale))
		nextPeriodInvoice = &ro.InvoiceDetailSimplify{
			TotalAmount:                    nextPeriodTotalAmountExcludingTax + nextPeriodTaxAmount,
			TotalAmountExcludingTax:        nextPeriodTotalAmountExcludingTax,
			Currency:                       currency,
			TaxAmount:                      nextPeriodTaxAmount,
			SubscriptionAmount:             nextPeriodTotalAmountExcludingTax + nextPeriodTaxAmount, // 在没有 discount 之前，保持于 Total 一致
			SubscriptionAmountExcludingTax: nextPeriodTotalAmountExcludingTax,                       // 在没有 discount 之前，保持于 Total 一致
			Lines:                          nextPeriodInvoiceItems,
		}

		if consts.ProrationUsingUniBeeCompute {
			selfComputeNextPeriodInvoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				Currency:      sub.Currency,
				PlanId:        req.NewPlanId,
				Quantity:      req.Quantity,
				AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
				TaxScale:      sub.TaxScale,
				PeriodStart:   nextPeriodStart,
				PeriodEnd:     nextPeriodEnd,
			})
			utility.Assert(selfComputeNextPeriodInvoice.TotalAmount == nextPeriodInvoice.TotalAmount, "System Error, Compute Error")
			nextPeriodInvoice = selfComputeNextPeriodInvoice
		}

		prorationInvoice = &ro.InvoiceDetailSimplify{
			TotalAmount:                    0,
			TotalAmountExcludingTax:        0,
			Currency:                       currency,
			TaxAmount:                      0,
			SubscriptionAmount:             0,
			SubscriptionAmountExcludingTax: 0,
			Lines:                          make([]*ro.InvoiceItemDetailRo, 0),
		}
		//nextPeriodTotalAmount = nextPeriodTotalAmountExcludingTax + nextPeriodTaxAmount
		prorationDate = sub.CurrentPeriodEnd

	}

	return &SubscriptionUpdatePrepareInternalRes{
		Subscription:      sub,
		Plan:              plan,
		Quantity:          req.Quantity,
		PlanChannel:       planChannel,
		PayChannel:        payChannel,
		MerchantInfo:      merchantInfo,
		AddonParams:       req.AddonParams,
		Addons:            addons,
		Currency:          currency,
		UserId:            sub.UserId,
		OldPlan:           oldPlan,
		TotalAmount:       totalAmount,
		Invoice:           prorationInvoice,
		NextPeriodInvoice: nextPeriodInvoice,
		ProrationDate:     prorationDate,
		EffectImmediate:   effectImmediate,
	}, nil

}

func SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq, merchantUserId int64, adminNote string) (*subscription.SubscriptionUpdateRes, error) {
	prepare, err := SubscriptionUpdatePreview(ctx, &subscription.SubscriptionUpdatePreviewReq{
		SubscriptionId:      req.SubscriptionId,
		NewPlanId:           req.NewPlanId,
		Quantity:            req.Quantity,
		AddonParams:         req.AddonParams,
		WithImmediateEffect: req.WithImmediateEffect,
	}, req.ProrationDate, merchantUserId)
	if err != nil {
		return nil, err
	}

	//subscription prepare 检查
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")

	var effectImmediate = 0
	var effectTime = prepare.Subscription.CurrentPeriodEnd
	if prepare.EffectImmediate {
		effectImmediate = 1
		effectTime = gtime.Now().Timestamp()
	}

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:           prepare.MerchantInfo.Id,
		ChannelId:            prepare.Subscription.ChannelId,
		UserId:               prepare.Subscription.UserId,
		SubscriptionId:       prepare.Subscription.SubscriptionId,
		UpdateSubscriptionId: utility.CreateSubscriptionUpdateId(),
		Amount:               prepare.Subscription.Amount,
		Currency:             prepare.Subscription.Currency,
		PlanId:               prepare.Subscription.PlanId,
		Quantity:             prepare.Subscription.Quantity,
		AddonData:            prepare.Subscription.AddonData,
		UpdateAmount:         prepare.NextPeriodInvoice.TotalAmount,
		ProrationAmount:      prepare.Invoice.TotalAmount,
		UpdateCurrency:       prepare.Currency,
		UpdatePlanId:         int64(prepare.Plan.Id),
		UpdateQuantity:       prepare.Quantity,
		UpdateAddonData:      utility.MarshalToJsonString(prepare.AddonParams), // addon 暂定不带上之前订阅
		Status:               consts.PendingSubStatusInit,
		Data:                 "", //额外参数配置
		MerchantUserId:       merchantUserId,
		AdminNote:            adminNote,
		ProrationDate:        req.ProrationDate,
		EffectImmediate:      effectImmediate,
		EffectTime:           effectTime,
	}

	result, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	var subUpdateRes *ro.ChannelUpdateSubscriptionInternalResp
	if consts.ProrationUsingUniBeeCompute {
		if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
			// createAndPayNewProrationInvoice
			merchantInfo := query.GetMerchantInfoById(ctx, one.MerchantId)
			utility.Assert(merchantInfo != nil, "merchantInfo not found")
			//utility.Assert(user != nil, "user not found")
			payChannel := query.GetSubscriptionTypePayChannelById(ctx, one.ChannelId)
			utility.Assert(payChannel != nil, "payChannel not found")
			user := query.GetUserAccountById(ctx, uint64(one.UserId))
			var mobile = ""
			var firstName = ""
			var lastName = ""
			var gender = ""
			if user != nil {
				mobile = user.Mobile
				firstName = user.FirstName
				lastName = user.LastName
				gender = user.Gender
			}

			createRes, err := service.DoChannelPay(ctx, &ro.CreatePayContext{
				PayChannel: payChannel,
				Pay: &entity.Payment{
					SubscriptionId: one.SubscriptionId,
					BizId:          one.UpdateSubscriptionId,
					BizType:        consts.BIZ_TYPE_SUBSCRIPTION,
					UserId:         one.UserId,
					ChannelId:      int64(payChannel.Id),
					TotalAmount:    prepare.Invoice.TotalAmount,
					Currency:       one.Currency,
					CountryCode:    prepare.Subscription.CountryCode,
					MerchantId:     one.MerchantId,
					CompanyId:      merchantInfo.CompanyId,
				},
				Platform:      "WEB",
				DeviceType:    "Web",
				ShopperUserId: strconv.FormatInt(one.UserId, 10),
				ShopperEmail:  prepare.Subscription.CustomerEmail,
				ShopperLocale: "en",
				Mobile:        mobile,
				Invoice:       prepare.Invoice,
				ShopperName: &v1.OutShopperName{
					FirstName: firstName,
					LastName:  lastName,
					Gender:    gender,
				},
				MediaData:              map[string]string{"BillingReason": "SubscriptionUpdate"},
				MerchantOrderReference: one.UpdateSubscriptionId,
				PayMethod:              1, //automatic
				DaysUtilDue:            5, //one day expire
				ChannelPaymentMethod:   prepare.Subscription.ChannelDefaultPaymentMethod,
			})
			if err != nil {
				return nil, err
			}
			// Upgrade
			subUpdateRes = &ro.ChannelUpdateSubscriptionInternalResp{
				ChannelUpdateId: createRes.PaymentId,
				Data:            utility.MarshalToJsonString(createRes),
				Link:            createRes.Link,
				Paid:            createRes.Status == consts.PAY_SUCCESS,
			}
			//subUpdateRes.ChannelUpdateId = createRes.ChannelInvoiceId
			//subUpdateRes.Paid = createRes.AlreadyPaid
			//subUpdateRes.Link = createRes.Link
			//subUpdateRes.Data = utility.MarshalToJsonString(createRes)
		} else {
			prepare.EffectImmediate = false
			subUpdateRes, err = out.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
				Plan:            prepare.Plan,
				Quantity:        prepare.Quantity,
				AddonPlans:      prepare.Addons,
				PlanChannel:     prepare.PlanChannel,
				Subscription:    prepare.Subscription,
				ProrationDate:   req.ProrationDate,
				EffectImmediate: false,
			})
		}
	} else {
		subUpdateRes, err = out.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
			Plan:            prepare.Plan,
			Quantity:        prepare.Quantity,
			AddonPlans:      prepare.Addons,
			PlanChannel:     prepare.PlanChannel,
			Subscription:    prepare.Subscription,
			ProrationDate:   req.ProrationDate,
			EffectImmediate: prepare.EffectImmediate,
		})
	}
	if err != nil {
		return nil, err
	}
	// 标记更新单
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PendingUpdateId: one.UpdateSubscriptionId,
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	// 当前更新单渠道执行成功，需要取消其他 PendingUpdate，保证只有一个在 Create 状态 todo mark need cancel payment、 invoice and send invoice email
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusCancelled,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, prepare.Subscription.SubscriptionId).
		WhereNot(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	var PaidInt = 0
	if subUpdateRes.Paid {
		PaidInt = 1
	}

	one.Link = subUpdateRes.Link
	one.Status = consts.PendingSubStatusCreate
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:          consts.PendingSubStatusCreate,
		dao.SubscriptionPendingUpdate.Columns().ResponseData:    subUpdateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:       gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Paid:            PaidInt,
		dao.SubscriptionPendingUpdate.Columns().Link:            subUpdateRes.Link,
		dao.SubscriptionPendingUpdate.Columns().ChannelUpdateId: subUpdateRes.ChannelUpdateId,
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	if prepare.EffectImmediate && subUpdateRes.Paid {
		//需要3DS校验的用户，在进行订阅更新，如果使用 PendingUpdate，经过验证也是需要 3DS 校验，如果不使用 PendingUpdate，下一周期再进行Invoice收款，可能面临发票自动收款失败，然后需要用户 3DS 校验的情况；使用了 PendingUpdate 提前收款只是把问题前置了
		one.Status = consts.PendingSubStatusFinished
		_, err = handler.FinishPendingUpdateForSubscription(ctx, prepare.Subscription, one)
		if err != nil {
			return nil, err
		}
	}

	return &subscription.SubscriptionUpdateRes{
		SubscriptionPendingUpdate: one,
		Paid:                      len(subUpdateRes.Link) == 0,
		Link:                      subUpdateRes.Link,
	}, nil
}

func SubscriptionCancel(ctx context.Context, subscriptionId string, proration bool, invoiceNow bool) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusCancelled, "subscription already cancelled")
	plan := query.GetPlanById(ctx, sub.PlanId)
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payment channel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if !consts.GetConfigInstance().IsServerDev() || !consts.GetConfigInstance().IsLocal() {
		// todo mark will support proration invoiceNow later
		invoiceNow = false
		proration = false
		// todo mark will support proration invoiceNow later
		// only local env can cancel immediately invoice_compute proration invoice
		utility.Assert(invoiceNow == false && proration == false, "cancel subscription with proration invoice immediate not support for this version")
	}
	var nextStatus = consts.SubStatusPendingInActive
	if sub.Type == consts.SubTypeDefault {
		_, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancel(ctx, &ro.ChannelCancelSubscriptionInternalReq{
			Plan:         plan,
			PlanChannel:  planChannel,
			Subscription: sub,
			InvoiceNow:   invoiceNow,
			Prorate:      proration,
		})
		if err != nil {
			return err
		}
	} else {
		nextStatus = consts.SubStatusCancelled
	}
	// cancel will generate proration invoice need compute todo mark
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:    nextStatus,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool, merchantUserId int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 1 {
		//已经设置未周期结束取消
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if sub.Type == consts.SubTypeDefault {
		_, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx, plan, planChannel, sub)
		if err != nil {
			return err
		}
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1, // todo mark 如果您在计费周期结束时取消订阅（即设置cancel_at_period_end为true），customer.subscription.updated则会立即触发事件。该事件反映了订阅值的变化cancel_at_period_end。当订阅在期限结束时实际取消时，customer.subscription.deleted就会发生一个事件
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
	// SendEmail
	if merchantUserId > 0 {
		//merchant Cancel
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, email.TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin, "", &email.TemplateVariable{
			UserName:            user.UserName,
			MerchantProductName: plan.ChannelProductName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd).Layout("2006-01-02"),
		})
		if err != nil {
			fmt.Printf("SendTemplateEmail err:%s", err.Error())
		}
	} else {
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, email.TemplateSubscriptionCancelledAtPeriodEndByUser, "", &email.TemplateVariable{
			UserName:            user.UserName,
			MerchantProductName: plan.ChannelProductName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd).Layout("2006-01-02"),
		})
		if err != nil {
			fmt.Printf("SendTemplateEmail err:%s", err.Error())
		}
	}
	return nil
}

func SubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 0 {
		//已经设置未周期结束取消
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if sub.Type == consts.SubTypeDefault {
		_, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx, plan, planChannel, sub)
		if err != nil {
			return err
		}
	}

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("SubscriptionCancel subscription err:%s", update)
	//}
	return nil
}

func SubscriptionAddNewTrialEnd(ctx context.Context, subscriptionId string, AppendNewTrialEndByHour int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")

	if sub.Type == consts.SubTypeDefault {
		details, err := out.GetPayChannelServiceProvider(ctx, sub.ChannelId).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, sub)
		utility.Assert(err == nil, fmt.Sprintf("SubscriptionDetail Fetch error%s", err))
		err = handler.UpdateSubWithChannelDetailBack(ctx, sub, details)
		sub = query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		utility.Assert(err == nil, fmt.Sprintf("UpdateSubWithChannelDetailBack Fetch error%s", err))
	}
	utility.Assert(AppendNewTrialEndByHour > 0, "invalid AppendNewTrialEndByHour , should > 0")
	newTrialEnd := sub.CurrentPeriodEnd + AppendNewTrialEndByHour*3600
	if sub.Type == consts.SubTypeDefault {
		_, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionNewTrialEnd(ctx, plan, planChannel, sub, newTrialEnd)
		if err != nil {
			return err
		}
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().TrialEnd:  newTrialEnd,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("SubscriptionAddNewTrialEnd subscription err:%s", update)
	//}
	return nil
}

func SubscriptionEndTrial(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	var newTrialEnd = sub.TrialEnd
	if sub.Type == consts.SubTypeDefault {
		details, err := out.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionEndTrial(ctx, plan, planChannel, sub)
		if err != nil {
			return err
		}
		newTrialEnd = details.TrialEnd
	} else {
		newTrialEnd = sub.CurrentPeriodStart - 1
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().TrialEnd:  newTrialEnd,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	// end trail may cause payment immediately todo mark
	return nil
}
