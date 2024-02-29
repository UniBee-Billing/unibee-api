package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"time"
	v1 "unibee/api/onetime/payment"
	"unibee/api/user/subscription"
	"unibee/api/user/vat"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/ro"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionCreatePrepareInternalRes struct {
	Plan              *entity.SubscriptionPlan           `json:"plan"`
	Quantity          int64                              `json:"quantity"`
	Gateway           *entity.MerchantGateway            `json:"gateway"`
	MerchantInfo      *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.PlanAddonVo                  `json:"addons"`
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

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo) []*ro.PlanAddonVo {
	var addons []*ro.PlanAddonVo
	var totalAddonIds []uint64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.SubscriptionPlan
	if len(totalAddonIds) > 0 {
		//query all plan
		err := dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//add to list
			mapPlans := make(map[uint64]*entity.SubscriptionPlan)
			for _, pair := range allAddonList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				addons = append(addons, &ro.PlanAddonVo{
					Quantity:  param.Quantity,
					AddonPlan: ro.SimplifyPlan(mapPlans[param.AddonPlanId]),
				})
			}
		}
	}
	return addons
}

func VatNumberValidate(ctx context.Context, req *vat.NumberValidateReq, userId int64) (*vat.NumberValidateRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), userId, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	if vatNumberValidate.Valid {
		vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, _interface.GetMerchantId(ctx), vatNumberValidate.CountryCode)
		utility.Assert(err == nil, fmt.Sprintf("vatNumber vatCountryCode check error:%s", err))
		utility.Assert(vatCountryRate != nil, fmt.Sprintf("vatNumber not found for countryCode:%v", vatNumberValidate.CountryCode))
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}

func MerchantGatewayCheck(ctx context.Context, merchantId uint64, reqGatewayId uint64) *entity.MerchantGateway {
	if reqGatewayId > 0 {
		gateway := query.GetGatewayById(ctx, reqGatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		utility.Assert(gateway.MerchantId == merchantId, "gateway not match")
		return gateway
	} else {
		list := query.GetMerchantGatewayList(ctx, merchantId)
		utility.Assert(len(list) > 0, "merchant gateway need setup")
		utility.Assert(len(list) == 1, "gateway need specify")
		return list[0]
	}
}

func SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (*SubscriptionCreatePrepareInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.GatewayId > 0, "Id invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	email := ""
	if !consts.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
		email = _interface.BizCtx().Get(ctx).User.Email
	}
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := MerchantGatewayCheck(ctx, plan.MerchantId, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
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
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.UserId, req.VatNumber, "")
		if err != nil {
			return nil, err
		}
		utility.Assert(vatNumberValidate.Valid, fmt.Sprintf("vat number validate failure:"+req.VatNumber))
		vatCountryCode = vatNumberValidate.CountryCode
	}

	if len(vatCountryCode) == 0 && len(user.CountryCode) > 0 {
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

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	var currency = plan.Currency
	var TotalAmountExcludingTax = plan.Amount * req.Quantity

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
		TotalAmountExcludingTax = TotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var billingCycleAnchor = gtime.Now()
	var currentTimeStart = billingCycleAnchor
	var currentTimeEnd = subscription2.GetPeriodEndFromStart(ctx, billingCycleAnchor.Timestamp(), req.PlanId)

	invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:   "SubscriptionCreate",
		Currency:      currency,
		PlanId:        req.PlanId,
		Quantity:      req.Quantity,
		AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
		TaxScale:      standardTaxScale,
		PeriodStart:   currentTimeStart.Timestamp(),
		PeriodEnd:     currentTimeEnd,
	})

	return &SubscriptionCreatePrepareInternalRes{
		Plan:              plan,
		Quantity:          req.Quantity,
		Gateway:           gateway,
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
		GatewayId:      req.GatewayId,
		UserId:         req.UserId,
		AddonParams:    req.AddonParams,
		VatCountryCode: req.VatCountryCode,
		VatNumber:      req.VatNumber,
	})
	if err != nil {
		return nil, err
	}
	utility.Assert(len(prepare.VatCountryCode) > 0, "CountryCode Needed")
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")

	var subType = consts.SubTypeDefault
	if consts.SubscriptionCycleUnderUniBeeControl {
		subType = consts.SubTypeUniBeeControl
	}

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, prepare.Invoice.PeriodEnd, prepare.Plan.Id)

	one := &entity.Subscription{
		MerchantId:         prepare.MerchantInfo.Id,
		Type:               subType,
		PlanId:             prepare.Plan.Id,
		GatewayId:          prepare.Gateway.Id,
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
		CurrentPeriodStart: prepare.Invoice.PeriodStart,
		CurrentPeriodEnd:   prepare.Invoice.PeriodEnd,
		DunningTime:        dunningTime,
		BillingCycleAnchor: prepare.Invoice.PeriodStart,
		CreateTime:         gtime.Now().Timestamp(),
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	var createRes *ro.GatewayCreateSubscriptionInternalResp
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
	invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, one)
	utility.AssertError(err, "System Error")
	createPaymentResult, err := service.GatewayPaymentCreate(ctx, &ro.CreatePayContext{
		CheckoutMode: true,
		Gateway:      prepare.Gateway,
		Pay: &entity.Payment{
			SubscriptionId: one.SubscriptionId,
			BizId:          one.SubscriptionId,
			BizType:        consts.BizTypeSubscription,
			UserId:         prepare.UserId,
			GatewayId:      prepare.Gateway.Id,
			TotalAmount:    prepare.Invoice.TotalAmount,
			Currency:       prepare.Currency,
			CountryCode:    prepare.VatCountryCode,
			MerchantId:     prepare.MerchantInfo.Id,
			CompanyId:      prepare.MerchantInfo.CompanyId,
			BillingReason:  prepare.Invoice.InvoiceName,
			ReturnUrl:      req.ReturnUrl,
		},
		Platform:      "WEB",
		DeviceType:    "Web",
		ShopperUserId: strconv.FormatInt(one.UserId, 10),
		ShopperEmail:  prepare.Email,
		ShopperLocale: "en",
		Mobile:        mobile,
		Invoice:       invoice_compute.ConvertInvoiceToSimplify(invoice),
		ShopperName: &v1.OutShopperName{
			FirstName: firstName,
			LastName:  lastName,
			Gender:    gender,
		},
		MediaData:              map[string]string{"BillingReason": prepare.Invoice.InvoiceName},
		MerchantOrderReference: one.SubscriptionId,
	})
	if err != nil {
		return nil, err
	}
	createRes = &ro.GatewayCreateSubscriptionInternalResp{
		GatewaySubscriptionId: createPaymentResult.PaymentId,
		Data:                  utility.MarshalToJsonString(createPaymentResult),
		Link:                  createPaymentResult.Link,
		Paid:                  createPaymentResult.Status == consts.PaymentSuccess,
	}

	//Update Subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().GatewaySubscriptionId: createRes.GatewaySubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
		dao.Subscription.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	one.GatewaySubscriptionId = createRes.GatewaySubscriptionId
	one.Status = consts.GatewayPlanStatusCreate
	one.Link = createRes.Link

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCreate.Topic,
		Tag:   redismq2.TopicSubscriptionCreate.Tag,
		Body:  one.SubscriptionId,
	})
	return &subscription.SubscriptionCreateRes{
		Subscription: one,
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}

type SubscriptionUpdatePrepareInternalRes struct {
	Subscription      *entity.Subscription               `json:"subscription"`
	Plan              *entity.SubscriptionPlan           `json:"plan"`
	Quantity          int64                              `json:"quantity"`
	Gateway           *entity.MerchantGateway            `json:"gateway"`
	MerchantInfo      *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.PlanAddonVo                  `json:"addons"`
	TotalAmount       int64                              `json:"totalAmount"                `
	Currency          string                             `json:"currency"              `
	UserId            int64                              `json:"userId" `
	OldPlan           *entity.SubscriptionPlan           `json:"oldPlan"`
	Invoice           *ro.InvoiceDetailSimplify          `json:"invoice"`
	NextPeriodInvoice *ro.InvoiceDetailSimplify          `json:"nextPeriodInvoice"`
	ProrationDate     int64                              `json:"prorationDate"`
	EffectImmediate   bool                               `json:"EffectImmediate"`
}

// SubscriptionUpdatePreview Default行为，升级订阅主方案不管总金额是否比之前高，都将按比例计算发票立即生效；降级订阅方案，次月生效；问题点，降级方案如果 addon 多可能的总金额可能比之前高
func SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq, prorationDate int64, merchantUserId int64) (res *SubscriptionUpdatePrepareInternalRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	// todo mark addon binding check

	plan := query.GetPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	//试用期内不允许修改计划
	//utility.Assert(sub.TrialEnd < gtime.Now().Timestamp(), "subscription is in trial period ,should end trial before update")
	//设置了下周期取消不允许修改计划
	utility.Assert(sub.CancelAtPeriodEnd == 0, "subscription cannot be update as it will cancel at period end, should resume subscription first")
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)

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
	//oldPlanChannel := query.GetGatewayPlan(ctx, int64(oldPlan.Id), sub.Id)
	//utility.Assert(oldPlanChannel != nil, "oldPlangateway not found")

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
		var oldAddonMap = make(map[uint64]int64)
		for _, oldAddon := range oldAddonParams {
			if _, ok := oldAddonMap[oldAddon.AddonPlanId]; ok {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddonMap[oldAddon.AddonPlanId] + oldAddon.Quantity
			} else {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddon.Quantity
			}
		}
		var newAddonMap = make(map[uint64]int64)
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
					effectImmediate = true
					break
				}
			} else {
				effectImmediate = true
				break
			}
		}
		var changed = false
		if len(oldAddonMap) != len(newAddonMap) {
			changed = true
		} else {
			for newAddonPlanId, newAddonQuantity := range newAddonMap {
				if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
					if oldAddonQuantity != newAddonQuantity {
						changed = true
						break
					}
				} else {
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
	var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, uint64(req.NewPlanId))

	var totalAmount int64
	var prorationInvoice *ro.InvoiceDetailSimplify
	//var nextPeriodTotalAmount int64
	var nextPeriodInvoice *ro.InvoiceDetailSimplify
	if effectImmediate {
		// Generate Proration Invoice Preview
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
			// after period end before trial end, also or sub data not sync or use testClock in stage env
			prorationInvoice = &ro.InvoiceDetailSimplify{
				InvoiceName:                    "SubscriptionUpgrade",
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
				InvoiceName:       "SubscriptionUpgrade",
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
			InvoiceName:                    "SubscriptionCycle",
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
			InvoiceName:                    "SubscriptionUpgrade",
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

	if prorationInvoice.TotalAmount <= 0 {
		effectImmediate = false // todo mark effectImmediate = true with negative proration invoice should not allowed
	}

	return &SubscriptionUpdatePrepareInternalRes{
		Subscription:      sub,
		Plan:              plan,
		Quantity:          req.Quantity,
		Gateway:           gateway,
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

type UpdateSubscriptionInternalResp struct {
	GatewayUpdateId string          `json:"gatewayUpdateId" description:""`
	Data            string          `json:"data"`
	Link            string          `json:"link" description:""`
	Paid            bool            `json:"paid" description:""`
	Invoice         *entity.Invoice `json:"invoice" description:""`
}

func SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq, merchantUserId int64) (*subscription.SubscriptionUpdateRes, error) {
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

	//subscription prepare
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	if prepare.Invoice.TotalAmount <= 0 {
		utility.Assert(prepare.EffectImmediate == false, "System Error, Cannot Effect Immediate With Negative Amount")
	}

	var effectImmediate = 0
	var effectTime = prepare.Subscription.CurrentPeriodEnd
	if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
		effectImmediate = 1
		effectTime = gtime.Now().Timestamp()
	}

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:           prepare.MerchantInfo.Id,
		GatewayId:            prepare.Subscription.GatewayId,
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
		UpdatePlanId:         prepare.Plan.Id,
		UpdateQuantity:       prepare.Quantity,
		UpdateAddonData:      utility.MarshalToJsonString(prepare.AddonParams),
		Status:               consts.PendingSubStatusInit,
		Data:                 "",
		MerchantUserId:       merchantUserId,
		ProrationDate:        req.ProrationDate,
		EffectImmediate:      effectImmediate,
		EffectTime:           effectTime,
		CreateTime:           gtime.Now().Timestamp(),
	}

	result, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	var subUpdateRes *UpdateSubscriptionInternalResp
	if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
		// createAndPayNewProrationInvoice
		merchantInfo := query.GetMerchantInfoById(ctx, one.MerchantId)
		utility.Assert(merchantInfo != nil, "merchantInfo not found")
		//utility.Assert(user != nil, "user not found")
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, prepare.Subscription)
		utility.AssertError(err, "System Error")
		createRes, err := service.CreateSubInvoiceAutomaticPayment(ctx, prepare.Subscription, invoice)
		if err != nil {
			return nil, err
		}
		// Upgrade
		subUpdateRes = &UpdateSubscriptionInternalResp{
			GatewayUpdateId: createRes.Invoice.InvoiceId,
			Data:            utility.MarshalToJsonString(createRes),
			Link:            createRes.Link,
			Paid:            createRes.Status == consts.PaymentSuccess,
			Invoice:         createRes.Invoice,
		}
	} else {
		prepare.EffectImmediate = false
		subUpdateRes = &UpdateSubscriptionInternalResp{
			Paid: false,
			Link: "",
		}
	}

	if err != nil {
		return nil, err
	}
	// bing to subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PendingUpdateId: one.UpdateSubscriptionId,
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	// only one need, cancel others
	// need cancel payment、 invoice and send invoice email
	CancelOtherUnfinishedPendingUpdatesBackground(prepare.Subscription.SubscriptionId, one.UpdateSubscriptionId, "CancelByNewUpdate-"+one.UpdateSubscriptionId)

	var PaidInt = 0
	if subUpdateRes.Paid {
		PaidInt = 1
	}
	var note = "Success"
	if effectImmediate == 1 && !subUpdateRes.Paid {
		note = "Payment Action Required"
	} else {
		note = "Will Effect At Period End"
	}

	one.Link = subUpdateRes.Link
	one.Status = consts.PendingSubStatusCreate
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:       consts.PendingSubStatusCreate,
		dao.SubscriptionPendingUpdate.Columns().ResponseData: subUpdateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:    gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Paid:         PaidInt,
		dao.SubscriptionPendingUpdate.Columns().Link:         subUpdateRes.Link,
		dao.SubscriptionPendingUpdate.Columns().InvoiceId:    subUpdateRes.GatewayUpdateId,
		dao.SubscriptionPendingUpdate.Columns().Note:         note,
	}).Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, one.UpdateSubscriptionId).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	if prepare.EffectImmediate && subUpdateRes.Paid {
		_, err = handler.HandlePendingUpdatePaymentSuccess(ctx, prepare.Subscription, one.UpdateSubscriptionId, subUpdateRes.Invoice)
		if err != nil {
			return nil, err
		}
		one.Status = consts.PendingSubStatusFinished
	}

	return &subscription.SubscriptionUpdateRes{
		SubscriptionPendingUpdate: one,
		Paid:                      len(subUpdateRes.Link) == 0 || subUpdateRes.Paid, // link is blank or paid is true, portal will not redirect
		Link:                      subUpdateRes.Link,
		Note:                      note,
	}, nil
}

func SubscriptionCancel(ctx context.Context, subscriptionId string, proration bool, invoiceNow bool, reason string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusCancelled, "subscription already cancelled")
	utility.Assert(sub.Status != consts.SubStatusExpired, "subscription already expired")
	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
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
	var nextStatus = consts.SubStatusCancelled
	// cancel will generate proration invoice need compute todo mark
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:       nextStatus,
		dao.Subscription.Columns().CancelReason: reason,
		dao.Subscription.Columns().TrialEnd:     sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	if user != nil {
		merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
		if merchant != nil {
			err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionImmediateCancel, "", &email.TemplateVariable{
				UserName:            user.FirstName + " " + user.LastName,
				MerchantProductName: plan.PlanName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        merchant.Name,
				PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
			})
			if err != nil {
				fmt.Printf("SendTemplateEmail err:%s", err.Error())
			}
		}
	}

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCancel.Topic,
		Tag:   redismq2.TopicSubscriptionCancel.Tag,
		Body:  sub.SubscriptionId,
	})
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
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1,
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
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin, "", &email.TemplateVariable{
			UserName:            user.FirstName + " " + user.LastName,
			MerchantProductName: plan.PlanName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
		})
		if err != nil {
			fmt.Printf("SendTemplateEmail err:%s", err.Error())
		}
	} else {
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelledAtPeriodEndByUser, "", &email.TemplateVariable{
			UserName:            user.FirstName + " " + user.LastName,
			MerchantProductName: plan.PlanName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
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
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
	err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelLastCancelledAtPeriodEnd, "", &email.TemplateVariable{
		UserName:            user.FirstName + " " + user.LastName,
		MerchantProductName: plan.PlanName,
		MerchantCustomEmail: merchant.Email,
		MerchantName:        merchant.Name,
		PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
	})
	if err != nil {
		fmt.Printf("SendTemplateEmail err:%s", err.Error())
	}
	return nil
}

func SubscriptionAddNewTrialEnd(ctx context.Context, subscriptionId string, AppendNewTrialEndByHour int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "sub cancelled or sub expired")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

	utility.Assert(AppendNewTrialEndByHour > 0, "invalid AppendNewTrialEndByHour , should > 0")
	newTrialEnd := sub.CurrentPeriodEnd + AppendNewTrialEndByHour*3600

	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, uint64(sub.PlanId))
	newStatus := sub.Status
	if newTrialEnd > gtime.Now().Timestamp() {
		//automatic change sub status to active
		newStatus = consts.SubStatusActive
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:             newStatus,
		dao.Subscription.Columns().TrialEnd:           newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
		dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
		dao.Subscription.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
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
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not trialed")
	err := EndTrialManual(ctx, sub.SubscriptionId)
	if err != nil {
		return err
	}

	return nil
}

func EndTrialManual(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not in trial period")
	newTrialEnd := sub.CurrentPeriodStart - 1
	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, uint64(sub.PlanId))
	newStatus := sub.Status
	if gtime.Now().Timestamp() > sub.CurrentPeriodEnd {
		// todo mark has unfinished pending update
		newStatus = consts.SubStatusIncomplete
		// Payment Pending Enter Incomplete
		plan := query.GetPlanById(ctx, sub.PlanId)

		var nextPeriodStart = gtime.Now().Timestamp()
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			Currency:      sub.Currency,
			PlanId:        sub.PlanId,
			Quantity:      sub.Quantity,
			AddonJsonData: sub.AddonData,
			TaxScale:      sub.TaxScale,
			PeriodStart:   nextPeriodStart,
			PeriodEnd:     nextPeriodEnd,
			InvoiceName:   "SubscriptionCycle",
		})
		one, err := handler2.CreateProcessingInvoiceForSub(ctx, invoice, sub)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateProcessingInvoiceForSub err:", err.Error())
			return err
		}
		createRes, err := service.CreateSubInvoiceAutomaticPayment(ctx, sub, one)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateSubInvoiceAutomaticPayment err:", err.Error())
			return err
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().CurrentPeriodStart: invoice.PeriodStart,
			dao.Subscription.Columns().CurrentPeriodEnd:   invoice.PeriodEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
		g.Log().Print(ctx, "EndTrialManual CreateSubInvoiceAutomaticPayment:", utility.MarshalToJsonString(createRes))
		err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, gtime.Now().Timestamp())
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual HandleSubscriptionIncomplete err:", err.Error())
			return err
		}
	} else {
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Status:             newStatus,
			dao.Subscription.Columns().TrialEnd:           newTrialEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
	}
	return nil
}
