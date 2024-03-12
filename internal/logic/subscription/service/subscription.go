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
	"unibee/api/bean"
	"unibee/api/user/subscription"
	"unibee/api/user/vat"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionCreatePrepareInternalRes struct {
	Plan              *entity.Plan            `json:"plan"`
	Quantity          int64                   `json:"quantity"`
	Gateway           *entity.MerchantGateway `json:"gateway"`
	Merchant          *entity.Merchant        `json:"merchantInfo"`
	AddonParams       []*bean.PlanAddonParam  `json:"addonParams"`
	Addons            []*bean.PlanAddonDetail `json:"addons"`
	TotalAmount       int64                   `json:"totalAmount"                `
	Currency          string                  `json:"currency"              `
	VatCountryCode    string                  `json:"vatCountryCode"              `
	VatCountryName    string                  `json:"vatCountryName"              `
	VatNumber         string                  `json:"vatNumber"              `
	VatNumberValidate *bean.ValidResult       `json:"vatNumberValidate"              `
	TaxScale          int64                   `json:"taxScale"              `
	VatVerifyData     string                  `json:"vatVerifyData"              `
	Invoice           *bean.InvoiceSimplify   `json:"invoice"`
	UserId            int64                   `json:"userId" `
	Email             string                  `json:"email" `
	VatCountryRate    *bean.VatCountryRate    `json:"vatCountryRate" `
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*bean.PlanAddonParam) []*bean.PlanAddonDetail {
	var addons []*bean.PlanAddonDetail
	var totalAddonIds []uint64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.Plan
	if len(totalAddonIds) > 0 {
		//query all plan
		err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//add to list
			mapPlans := make(map[uint64]*entity.Plan)
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
				addons = append(addons, &bean.PlanAddonDetail{
					Quantity:  param.Quantity,
					AddonPlan: bean.SimplifyPlan(mapPlans[param.AddonPlanId]),
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

func SubscriptionCreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (*SubscriptionCreatePrepareInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.GatewayId > 0, "Id invalid")
	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "User invalid")
	email := ""
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := MerchantGatewayCheck(ctx, plan.MerchantId, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	user := query.GetUserAccountById(ctx, _interface.BizCtx().Get(ctx).User.Id)
	utility.Assert(user != nil, "user not found")

	if gateway.GatewayType == consts.GatewayTypeCrypto {
		utility.Assert(len(plan.GasPayer) > 0, "gasPayer must set before crypto payment")
	}

	var err error
	utility.Assert(query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), merchantInfo.Id) == nil, "another active subscription find, only one subscription can create")

	//vat
	utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil, "Merchant Vat VATGateway not setup")
	var vatCountryCode = req.VatCountryCode
	var standardTaxScale int64 = 0
	var vatCountryName = ""
	var vatCountryRate *bean.VatCountryRate
	var vatNumberValidate *bean.ValidResult

	if len(req.VatNumber) > 0 {
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, int64(_interface.BizCtx().Get(ctx).User.Id), req.VatNumber, "")
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
		Merchant:          merchantInfo,
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
		UserId:            int64(_interface.BizCtx().Get(ctx).User.Id),
		Email:             email,
		Invoice:           invoice,
		VatCountryRate:    vatCountryRate,
	}, nil
}

func SubscriptionCreate(ctx context.Context, req *subscription.CreateReq) (*subscription.CreateRes, error) {

	prepare, err := SubscriptionCreatePreview(ctx, &subscription.CreatePreviewReq{
		PlanId:         req.PlanId,
		Quantity:       req.Quantity,
		GatewayId:      req.GatewayId,
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
		MerchantId:                  prepare.Merchant.Id,
		Type:                        subType,
		PlanId:                      prepare.Plan.Id,
		GatewayId:                   prepare.Gateway.Id,
		UserId:                      prepare.UserId,
		Quantity:                    prepare.Quantity,
		Amount:                      prepare.TotalAmount,
		Currency:                    prepare.Currency,
		AddonData:                   utility.MarshalToJsonString(prepare.AddonParams),
		SubscriptionId:              utility.CreateSubscriptionId(),
		Status:                      consts.SubStatusCreate,
		CustomerEmail:               prepare.Email,
		ReturnUrl:                   req.ReturnUrl,
		Data:                        "", //额外参数配置
		VatNumber:                   prepare.VatNumber,
		VatVerifyData:               prepare.VatVerifyData,
		CountryCode:                 prepare.VatCountryCode,
		TaxScale:                    prepare.TaxScale,
		CurrentPeriodStart:          prepare.Invoice.PeriodStart,
		CurrentPeriodEnd:            prepare.Invoice.PeriodEnd,
		DunningTime:                 dunningTime,
		BillingCycleAnchor:          prepare.Invoice.PeriodStart,
		GatewayDefaultPaymentMethod: req.PaymentMethodId,
		CreateTime:                  gtime.Now().Timestamp(),
		MetaData:                    utility.MarshalToJsonString(req.Metadata),
		GasPayer:                    prepare.Plan.GasPayer,
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	var createRes *gateway_bean.GatewayCreateSubscriptionResp
	invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, one)
	utility.AssertError(err, "System Error")
	var createPaymentResult *gateway_bean.GatewayNewPaymentResp
	if len(req.PaymentMethodId) > 0 {
		// createAndPayNewProrationInvoice
		merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
		utility.Assert(merchantInfo != nil, "merchantInfo not found")
		//utility.Assert(user != nil, "user not found")
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, one)
		utility.AssertError(err, "System Error")
		createPaymentResult, err = service.CreateSubInvoiceAutomaticPayment(ctx, one, invoice)
		if err != nil {
			return nil, err
		}
	} else {
		createPaymentResult, err = service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
			CheckoutMode: true,
			Gateway:      prepare.Gateway,
			Pay: &entity.Payment{
				SubscriptionId:    one.SubscriptionId,
				ExternalPaymentId: one.SubscriptionId,
				BizType:           consts.BizTypeSubscription,
				UserId:            prepare.UserId,
				GatewayId:         prepare.Gateway.Id,
				TotalAmount:       prepare.Invoice.TotalAmount,
				Currency:          prepare.Currency,
				CountryCode:       prepare.VatCountryCode,
				MerchantId:        prepare.Merchant.Id,
				CompanyId:         prepare.Merchant.CompanyId,
				BillingReason:     prepare.Invoice.InvoiceName,
				ReturnUrl:         req.ReturnUrl,
				GasPayer:          prepare.Plan.GasPayer,
			},
			ExternalUserId: strconv.FormatInt(one.UserId, 10),
			Email:          prepare.Email,
			Invoice:        bean.SimplifyInvoice(invoice),
			Metadata:       map[string]string{"BillingReason": prepare.Invoice.InvoiceName},
		})
		if err != nil {
			return nil, err
		}
	}

	createRes = &gateway_bean.GatewayCreateSubscriptionResp{
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
	return &subscription.CreateRes{
		Subscription: bean.SimplifySubscription(one),
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}

type SubscriptionUpdatePrepareInternalRes struct {
	Subscription      *entity.Subscription    `json:"subscription"`
	Plan              *entity.Plan            `json:"plan"`
	Quantity          int64                   `json:"quantity"`
	Gateway           *entity.MerchantGateway `json:"gateway"`
	MerchantInfo      *entity.Merchant        `json:"merchantInfo"`
	AddonParams       []*bean.PlanAddonParam  `json:"addonParams"`
	Addons            []*bean.PlanAddonDetail `json:"addons"`
	TotalAmount       int64                   `json:"totalAmount"                `
	Currency          string                  `json:"currency"              `
	UserId            int64                   `json:"userId" `
	OldPlan           *entity.Plan            `json:"oldPlan"`
	Invoice           *bean.InvoiceSimplify   `json:"invoice"`
	NextPeriodInvoice *bean.InvoiceSimplify   `json:"nextPeriodInvoice"`
	ProrationDate     int64                   `json:"prorationDate"`
	EffectImmediate   bool                    `json:"EffectImmediate"`
}

// SubscriptionUpdatePreview Default行为，升级订阅主方案不管总金额是否比之前高，都将按比例计算发票立即生效；降级订阅方案，次月生效；问题点，降级方案如果 addon 多可能的总金额可能比之前高
func SubscriptionUpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq, prorationDate int64, merchantMemberId int64) (res *SubscriptionUpdatePrepareInternalRes, err error) {
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
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
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
		var oldAddonParams []*bean.PlanAddonParam
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

	var totalAmount int64
	var prorationInvoice *bean.InvoiceSimplify
	var nextPeriodInvoice *bean.InvoiceSimplify
	if effectImmediate {
		var oldAddonParams []*bean.PlanAddonParam
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
			if sub.TestClock > sub.CurrentPeriodStart && !consts.GetConfigInstance().IsProd() {
				prorationDate = sub.TestClock
			}
		}
		if prorationDate < sub.CurrentPeriodStart {
			// after period end before trial end, also or sub data not sync or use testClock in stage env
			prorationInvoice = &bean.InvoiceSimplify{
				InvoiceName:                    "SubscriptionUpgrade",
				TotalAmount:                    0,
				TotalAmountExcludingTax:        0,
				Currency:                       sub.Currency,
				TaxAmount:                      0,
				SubscriptionAmount:             0,
				SubscriptionAmountExcludingTax: 0,
				Lines:                          make([]*bean.InvoiceItemSimplify, 0),
				ProrationDate:                  prorationDate,
				PeriodStart:                    sub.CurrentPeriodStart,
				PeriodEnd:                      sub.CurrentPeriodEnd,
			}
		} else if prorationDate > sub.CurrentPeriodEnd {
			// after periodEnd, is not a prorationInvoice, just use it
			prorationInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:   "SubscriptionCycle",
				Currency:      sub.Currency,
				PlanId:        req.NewPlanId,
				Quantity:      req.Quantity,
				AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
				TaxScale:      sub.TaxScale,
				PeriodStart:   prorationDate,
				PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, prorationDate, req.NewPlanId),
			})
		} else {
			prorationInvoice = invoice_compute.ComputeSubscriptionProrationInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
				InvoiceName:       "SubscriptionUpgrade",
				Currency:          sub.Currency,
				TaxScale:          sub.TaxScale,
				ProrationDate:     prorationDate,
				OldProrationPlans: oldProrationPlanParams,
				NewProrationPlans: newProrationPlanParams,
				PeriodStart:       sub.CurrentPeriodStart,
				PeriodEnd:         sub.CurrentPeriodEnd,
			})
		}
		prorationDate = prorationInvoice.ProrationDate
		totalAmount = prorationInvoice.TotalAmount

		var nextPeriodStart = sub.CurrentPeriodEnd
		if sub.TrialEnd > sub.CurrentPeriodEnd {
			nextPeriodStart = sub.TrialEnd
		}
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, req.NewPlanId)
		// Generate Proration Invoice Preview
		nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			InvoiceName:   "SubscriptionCycle",
			Currency:      sub.Currency,
			PlanId:        req.NewPlanId,
			Quantity:      req.Quantity,
			AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
			TaxScale:      sub.TaxScale,
			PeriodStart:   nextPeriodStart,
			PeriodEnd:     nextPeriodEnd,
		})
	} else {
		prorationDate = sub.CurrentPeriodEnd
		prorationInvoice = &bean.InvoiceSimplify{
			InvoiceName:                    "SubscriptionUpgrade",
			TotalAmount:                    0,
			TotalAmountExcludingTax:        0,
			Currency:                       currency,
			TaxAmount:                      0,
			SubscriptionAmount:             0,
			SubscriptionAmountExcludingTax: 0,
			Lines:                          make([]*bean.InvoiceItemSimplify, 0),
			ProrationDate:                  prorationDate,
			PeriodStart:                    sub.CurrentPeriodStart,
			PeriodEnd:                      sub.CurrentPeriodEnd,
		}
		nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			InvoiceName:   "SubscriptionCycle",
			Currency:      sub.Currency,
			PlanId:        req.NewPlanId,
			Quantity:      req.Quantity,
			AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
			TaxScale:      sub.TaxScale,
			PeriodStart:   utility.MaxInt64(prorationInvoice.PeriodEnd, sub.TrialEnd),
			PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, utility.MaxInt64(prorationInvoice.PeriodEnd, sub.TrialEnd), req.NewPlanId),
		})
	}

	if prorationInvoice.TotalAmount <= 0 {
		effectImmediate = false
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

func SubscriptionUpdate(ctx context.Context, req *subscription.UpdateReq, merchantMemberId int64) (*subscription.UpdateRes, error) {
	prepare, err := SubscriptionUpdatePreview(ctx, &subscription.UpdatePreviewReq{
		SubscriptionId:      req.SubscriptionId,
		NewPlanId:           req.NewPlanId,
		Quantity:            req.Quantity,
		AddonParams:         req.AddonParams,
		WithImmediateEffect: req.WithImmediateEffect,
	}, req.ProrationDate, merchantMemberId)
	if err != nil {
		return nil, err
	}

	//subscription prepare
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	if prepare.Invoice.TotalAmount <= 0 {
		utility.Assert(prepare.EffectImmediate == false, "System Error, Cannot Effect Immediate With Negative CaptureAmount")
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
		MerchantMemberId:     merchantMemberId,
		ProrationDate:        req.ProrationDate,
		EffectImmediate:      effectImmediate,
		EffectTime:           effectTime,
		CreateTime:           gtime.Now().Timestamp(),
		MetaData:             utility.MarshalToJsonString(req.Metadata),
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
		merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
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
		dao.SubscriptionPendingUpdate.Columns().MetaData:     utility.MarshalToJsonString(req.Metadata),
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

	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: &bean.SubscriptionPendingUpdateDetail{
			MerchantId:           one.MerchantId,
			SubscriptionId:       one.SubscriptionId,
			UpdateSubscriptionId: one.UpdateSubscriptionId,
			GmtCreate:            one.GmtCreate,
			Amount:               one.Amount,
			Status:               one.Status,
			UpdateAmount:         one.UpdateAmount,
			Currency:             one.Currency,
			UpdateCurrency:       one.UpdateCurrency,
			PlanId:               one.PlanId,
			UpdatePlanId:         one.UpdatePlanId,
			Quantity:             one.Quantity,
			UpdateQuantity:       one.UpdateQuantity,
			AddonData:            one.AddonData,
			UpdateAddonData:      one.UpdateAddonData,
			ProrationAmount:      one.ProrationAmount,
			GatewayId:            one.GatewayId,
			UserId:               one.UserId,
			GmtModify:            one.GmtModify,
			Paid:                 one.Paid,
			Link:                 one.Link,
			MerchantMember:       bean.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate:      one.EffectImmediate,
			EffectTime:           one.EffectTime,
			Note:                 one.Note,
			Plan:                 bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:           bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:         addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:             req.Metadata,
		},
		Paid: len(subUpdateRes.Link) == 0 || subUpdateRes.Paid, // link is blank or paid is true, portal will not redirect
		Link: subUpdateRes.Link,
		Note: note,
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
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
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
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
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

func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool, merchantMemberId int64) error {
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
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	// SendEmail
	if merchantMemberId > 0 {
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
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
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
