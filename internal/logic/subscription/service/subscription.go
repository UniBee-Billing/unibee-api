package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/api/user/vat"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/subscription/handler"
	"go-oversea-pay/internal/logic/vat_gateway"
	"go-oversea-pay/internal/logic/vat_gateway/base"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

type SubscriptionCreatePrepareInternalRes struct {
	Plan                  *entity.SubscriptionPlan           `json:"planId"`
	Quantity              int64                              `json:"quantity"`
	PlanChannel           *entity.SubscriptionPlanChannel    `json:"planChannel"`
	PayChannel            *entity.OverseaPayChannel          `json:"payChannel"`
	MerchantInfo          *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams           []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons                []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount           int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency              string                             `json:"currency"              `      // 货币
	VatCountryCode        string                             `json:"vatCountryCode"              `
	VatCountryName        string                             `json:"vatCountryName"              `
	VatNumber             string                             `json:"vatNumber"              `
	VatNumberValidate     *base.ValidResult                  `json:"vatNumberValidate"              `
	StandardTaxPercentage int64                              `json:"standardTaxPercentage"              `
	VatVerifyData         string                             `json:"vatVerifyData"              `
	Invoice               *ro.ChannelDetailInvoiceRo         `json:"invoice"`
	UserId                int64                              `json:"userId" `
	Email                 string                             `json:"email" `
	VatCountryRate        *vat_gateway.VatCountryRate        `json:"vatCountryRate" `
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
	utility.Assert(plan.Status == consts.PlanStatusPublished, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.PlanId, req.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//vat
	utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil, "Merchant Vat Gateway not setup")
	var vatCountryCode = req.VatCountryCode
	var standardTaxPercentage int64 = 0
	var vatCountryName = ""
	var vatCountryRate *vat_gateway.VatCountryRate
	var vatNumberValidate *base.ValidResult
	var err error
	if len(req.VatNumber) > 0 {
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.VatNumber, "")
		if err != nil {
			return nil, err
		}
		if !vatNumberValidate.Valid {
			return nil, gerror.New("vat number validate failure:" + req.VatNumber)
		}
		vatCountryCode = vatNumberValidate.CountryCode
	}

	if len(vatCountryCode) > 0 {
		vatCountryRate, err = vat_gateway.QueryVatCountryRateByMerchant(ctx, merchantInfo.Id, vatCountryCode)
		utility.Assert(err == nil, fmt.Sprintf("vat vatCountryCode check error:%s", err))
		utility.Assert(vatCountryRate != nil, fmt.Sprintf("vat not found for countryCode:%v", vatCountryCode))
		vatCountryName = vatCountryRate.CountryName
		if vatNumberValidate == nil || !vatNumberValidate.Valid {
			standardTaxPercentage = vatCountryRate.StandardTaxPercentage
		}
	}

	//设置默认值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	var currency = plan.Currency
	var TotalAmountExcludingTax = plan.Amount * req.Quantity

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
		TotalAmountExcludingTax = TotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	//生成临时账单
	var invoiceItems []*ro.ChannelDetailInvoiceItem
	invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
		Currency:               currency,
		Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxPercentageToInternalFloat(standardTaxPercentage)),
		AmountExcludingTax:     req.Quantity * plan.Amount,
		Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(standardTaxPercentage)),
		UnitAmountExcludingTax: plan.Amount,
		Description:            plan.PlanName,
		Quantity:               req.Quantity,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               currency,
			Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxPercentageToInternalFloat(standardTaxPercentage)),
			Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxPercentageToInternalFloat(standardTaxPercentage)),
			AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Description:            addon.AddonPlan.PlanName,
			Quantity:               addon.Quantity,
		})
	}
	var taxAmount = int64(float64(TotalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(standardTaxPercentage))
	var totalAmount = TotalAmountExcludingTax + taxAmount

	invoice := &ro.ChannelDetailInvoiceRo{
		TotalAmount:                    totalAmount,
		TotalAmountExcludingTax:        TotalAmountExcludingTax,
		Currency:                       currency,
		TaxAmount:                      taxAmount,
		SubscriptionAmount:             totalAmount,             // 在没有 discount 之前，保持于 Total 一致
		SubscriptionAmountExcludingTax: TotalAmountExcludingTax, // 在没有 discount 之前，保持于 Total 一致
		Lines:                          invoiceItems,
	}

	return &SubscriptionCreatePrepareInternalRes{
		Plan:                  plan,
		Quantity:              req.Quantity,
		PlanChannel:           planChannel,
		PayChannel:            payChannel,
		MerchantInfo:          merchantInfo,
		AddonParams:           req.AddonParams,
		Addons:                addons,
		TotalAmount:           totalAmount,
		Currency:              currency,
		VatCountryCode:        vatCountryCode,
		VatCountryName:        vatCountryName,
		VatNumber:             req.VatNumber,
		VatNumberValidate:     vatNumberValidate,
		VatVerifyData:         utility.FormatToJsonString(vatNumberValidate),
		StandardTaxPercentage: standardTaxPercentage,
		UserId:                req.UserId,
		Email:                 email,
		Invoice:               invoice,
		VatCountryRate:        vatCountryRate,
	}, nil
}

func SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (*subscription.SubscriptionCreateRes, error) {

	utility.Assert(len(req.VatCountryCode) > 0, "VatCountryCode invalid")

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

	//校验
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	//channelUserId 处理
	var channelUserId string
	channelUser := query.GetUserChannel(ctx, prepare.UserId, prepare.PlanChannel.ChannelId)
	if channelUser != nil {
		channelUserId = channelUser.ChannelUserId
	}

	one := &entity.Subscription{
		MerchantId:     prepare.MerchantInfo.Id,
		PlanId:         int64(prepare.Plan.Id),
		ChannelId:      prepare.PlanChannel.ChannelId,
		UserId:         prepare.UserId,
		Quantity:       prepare.Quantity,
		Amount:         prepare.TotalAmount,
		Currency:       prepare.Currency,
		AddonData:      utility.MarshalToJsonString(prepare.AddonParams),
		SubscriptionId: utility.CreateSubscriptionId(),
		Status:         consts.SubStatusCreate,
		CustomerEmail:  prepare.Email,
		ChannelUserId:  channelUserId,
		ReturnUrl:      req.ReturnUrl,
		Data:           "", //额外参数配置
		VatNumber:      prepare.VatNumber,
		VatVerifyData:  prepare.VatVerifyData,
		CountryCode:    prepare.VatCountryCode,
		TaxPercentage:  prepare.StandardTaxPercentage,
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	createRes, err := gateway.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionCreate(ctx, &ro.ChannelCreateSubscriptionInternalReq{
		Plan:           prepare.Plan,
		AddonPlans:     prepare.Addons,
		PlanChannel:    prepare.PlanChannel,
		VatCountryRate: prepare.VatCountryRate,
		Subscription:   one,
	})
	if err != nil {
		return nil, err
	}

	//更新 Subscription
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().ChannelUserId:         createRes.ChannelUserId,
		dao.Subscription.Columns().ChannelSubscriptionId: createRes.ChannelSubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusCreate,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
		dao.Subscription.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("SubscriptionCreate update err:%s", update)
	}
	one.ChannelSubscriptionId = createRes.ChannelSubscriptionId
	one.Status = consts.PlanChannelStatusCreate
	one.Link = createRes.Link
	one.ChannelUserId = createRes.ChannelUserId

	if channelUser == nil && len(createRes.ChannelUserId) > 0 {
		_, err := query.SaveUserChannel(ctx, prepare.UserId, prepare.PlanChannel.ChannelId, createRes.ChannelUserId)
		if err != nil {
			// ChannelUser 创建错误
			return nil, gerror.Newf("SubscriptionCreate ChannelUser save err:%s", err)
		}
	}

	return &subscription.SubscriptionCreateRes{
		Subscription: one,
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}

type SubscriptionUpdatePrepareInternalRes struct {
	Subscription      *entity.Subscription               `json:"subscription"`
	Plan              *entity.SubscriptionPlan           `json:"planId"`
	Quantity          int64                              `json:"quantity"`
	PlanChannel       *entity.SubscriptionPlanChannel    `json:"planChannel"`
	PayChannel        *entity.OverseaPayChannel          `json:"payChannel"`
	MerchantInfo      *entity.MerchantInfo               `json:"merchantInfo"`
	AddonParams       []*ro.SubscriptionPlanAddonParamRo `json:"addonParams"`
	Addons            []*ro.SubscriptionPlanAddonRo      `json:"addons"`
	TotalAmount       int64                              `json:"totalAmount"                ` // 金额,单位：分
	Currency          string                             `json:"currency"              `      // 货币
	UserId            int64                              `json:"userId" `
	OldPlan           *entity.SubscriptionPlan           `json:"oldPlan"`
	OldPlanChannel    *entity.SubscriptionPlanChannel    `json:"oldPlanChannel"`
	Invoice           *ro.ChannelDetailInvoiceRo         `json:"invoice"`
	NextPeriodInvoice *ro.ChannelDetailInvoiceRo         `json:"nextPeriodInvoice"`
	ProrationDate     int64                              `json:"prorationDate"`
	EffectImmediate   bool                               `json:"EffectImmediate"`
}

// SubscriptionUpdatePreview 默认行为，升级订阅主方案不管总金额是否比之前高，都将按比例计算发票立即生效；降级订阅方案，次月生效；问题点，降级方案如果 addon 多可能的总金额可能比之前高
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
	utility.Assert(plan.Status == consts.PlanStatusPublished, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, req.NewPlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	//设置默认值
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams, planChannel.ChannelId)

	var currency = sub.Currency
	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
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
	oldPlanChannel := query.GetPlanChannel(ctx, int64(oldPlan.Id), sub.ChannelId)
	utility.Assert(oldPlanChannel != nil, "oldPlanChannel not found")

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

	var totalAmount int64
	var invoice *ro.ChannelDetailInvoiceRo
	//var nextPeriodTotalAmount int64
	var nextPeroidInvoice *ro.ChannelDetailInvoiceRo
	if effectImmediate {
		updatePreviewRes, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionUpdateProrationPreview(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
			Plan:            plan,
			Quantity:        req.Quantity,
			OldPlan:         oldPlan,
			AddonPlans:      addons,
			PlanChannel:     planChannel,
			OldPlanChannel:  oldPlanChannel,
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
		invoice = &ro.ChannelDetailInvoiceRo{
			TotalAmount:                    updatePreviewRes.Invoice.TotalAmount,
			TotalAmountExcludingTax:        updatePreviewRes.Invoice.TotalAmountExcludingTax,
			Currency:                       updatePreviewRes.Invoice.Currency,
			TaxAmount:                      updatePreviewRes.Invoice.TaxAmount,
			SubscriptionAmount:             updatePreviewRes.Invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: updatePreviewRes.Invoice.SubscriptionAmountExcludingTax,
			Lines:                          updatePreviewRes.Invoice.Lines,
		}
		//nextPeriodTotalAmount = updatePreviewRes.NextPeriodInvoice.TotalAmount
		nextPeroidInvoice = &ro.ChannelDetailInvoiceRo{
			TotalAmount:                    updatePreviewRes.NextPeriodInvoice.TotalAmount,
			TotalAmountExcludingTax:        updatePreviewRes.NextPeriodInvoice.TotalAmountExcludingTax,
			Currency:                       updatePreviewRes.NextPeriodInvoice.Currency,
			TaxAmount:                      updatePreviewRes.NextPeriodInvoice.TaxAmount,
			SubscriptionAmount:             updatePreviewRes.NextPeriodInvoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: updatePreviewRes.NextPeriodInvoice.SubscriptionAmountExcludingTax,
			Lines:                          updatePreviewRes.NextPeriodInvoice.Lines,
		}
	} else {
		//下周期生效,输出Preview账单
		var nextPeriodTotalAmountExcludingTax = plan.Amount * req.Quantity
		for _, addon := range addons {
			utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
			utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
			utility.Assert(addon.AddonPlan.Status == consts.PlanStatusPublished, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
			nextPeriodTotalAmountExcludingTax = nextPeriodTotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
		}

		var nextPeriodInvoiceItems []*ro.ChannelDetailInvoiceItem
		nextPeriodInvoiceItems = append(nextPeriodInvoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               currency,
			Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
			AmountExcludingTax:     req.Quantity * plan.Amount,
			Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
			UnitAmountExcludingTax: plan.Amount,
			Description:            plan.PlanName,
			Quantity:               req.Quantity,
		})
		for _, addon := range addons {
			nextPeriodInvoiceItems = append(nextPeriodInvoiceItems, &ro.ChannelDetailInvoiceItem{
				Currency:               currency,
				Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
				Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
				AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
				UnitAmountExcludingTax: addon.AddonPlan.Amount,
				Description:            addon.AddonPlan.PlanName,
				Quantity:               addon.Quantity,
			})
		}
		var nextPeriodTaxAmount = int64(float64(nextPeriodTotalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage))
		nextPeroidInvoice = &ro.ChannelDetailInvoiceRo{
			TotalAmount:                    nextPeriodTotalAmountExcludingTax + nextPeriodTaxAmount,
			TotalAmountExcludingTax:        nextPeriodTotalAmountExcludingTax,
			Currency:                       currency,
			TaxAmount:                      nextPeriodTaxAmount,
			SubscriptionAmount:             nextPeriodTotalAmountExcludingTax + nextPeriodTaxAmount, // 在没有 discount 之前，保持于 Total 一致
			SubscriptionAmountExcludingTax: nextPeriodTotalAmountExcludingTax,                       // 在没有 discount 之前，保持于 Total 一致
			Lines:                          nextPeriodInvoiceItems,
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
		OldPlanChannel:    oldPlanChannel,
		TotalAmount:       totalAmount,
		Invoice:           invoice,
		NextPeriodInvoice: nextPeroidInvoice,
		ProrationDate:     prorationDate,
		EffectImmediate:   effectImmediate,
	}, nil

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

	//subscription prepare 检查
	utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch again")
	utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")

	// 需要取消其他 PendingUpdate，保证只有一个在 Create 状态
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusCancelled,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, prepare.Subscription.Id).WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).OmitNil().Update()
	if err != nil {
		return nil, err
	}

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
		UpdateAmount:         prepare.TotalAmount,
		UpdateCurrency:       prepare.Currency,
		UpdatePlanId:         int64(prepare.Plan.Id),
		UpdateQuantity:       prepare.Quantity,
		UpdatedAddonData:     utility.MarshalToJsonString(prepare.AddonParams), // addon 暂定不带上之前订阅
		Status:               consts.PendingSubStatusInit,
		Data:                 "", //额外参数配置
		MerchantUserId:       merchantUserId,
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

	updateRes, err := gateway.GetPayChannelServiceProvider(ctx, int64(prepare.PayChannel.Id)).DoRemoteChannelSubscriptionUpdate(ctx, &ro.ChannelUpdateSubscriptionInternalReq{
		Plan:            prepare.Plan,
		Quantity:        prepare.Quantity,
		OldPlan:         prepare.OldPlan,
		AddonPlans:      prepare.Addons,
		PlanChannel:     prepare.PlanChannel,
		OldPlanChannel:  prepare.OldPlanChannel,
		Subscription:    prepare.Subscription,
		ProrationDate:   req.ProrationDate,
		EffectImmediate: prepare.EffectImmediate,
	})
	if err != nil {
		return nil, err
	}
	var PaidInt = 0
	if updateRes.Paid {
		PaidInt = 1
	}

	one.Link = updateRes.Link
	one.Status = consts.PendingSubStatusCreate
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:          consts.PendingSubStatusCreate,
		dao.SubscriptionPendingUpdate.Columns().ResponseData:    updateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:       gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Paid:            PaidInt,
		dao.SubscriptionPendingUpdate.Columns().Link:            updateRes.Link,
		dao.SubscriptionPendingUpdate.Columns().ChannelUpdateId: updateRes.ChannelUpdateId,
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	if prepare.EffectImmediate && updateRes.Paid {
		//需要3DS校验的用户，在进行订阅更新，如果使用 PendingUpdate，经过验证也是需要 3DS 校验，如果不使用 PendingUpdate，下一周期再进行Invoice收款，可能面临发票自动收款失败，然后需要用户 3DS 校验的情况；使用了 PendingUpdate 提前收款只是把问题前置了
		one.Status = consts.PendingSubStatusFinished
		_, err = handler.FinishPendingUpdateForSubscription(ctx, one)
		if err != nil {
			return nil, err
		}
	}

	return &subscription.SubscriptionUpdateRes{
		SubscriptionPendingUpdate: one,
		Paid:                      updateRes.Paid,
		Link:                      updateRes.Link,
	}, nil
}

func SubscriptionCancel(ctx context.Context, subscriptionId string, proration bool, invoiceNow bool) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusCancelled, "subscription already cancelled")
	plan := query.GetPlanById(ctx, sub.PlanId)
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	utility.Assert(planChannel != nil && len(planChannel.ChannelProductId) > 0 && len(planChannel.ChannelPlanId) > 0, "internal error plan channel transfer not complete")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")
	merchantInfo := query.GetMerchantInfoById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	_, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancel(ctx, &ro.ChannelCancelSubscriptionInternalReq{
		Plan:         plan,
		PlanChannel:  planChannel,
		Subscription: sub,
		InvoiceNow:   invoiceNow,
		Prorate:      proration,
	})
	if err != nil {
		return err
	}
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:    consts.SubStatusPendingInActive,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionCancel update subscription err:%s", update)
	}
	return nil
}

// todo mark 在版本2018-02-28之前，发送到更新订阅 API 的任何参数都会停止挂起的取消，需验证
func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool) error {
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
	_, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx, plan, planChannel, sub)
	if err != nil {
		return err
	}
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1, // todo mark 如果您在计费周期结束时取消订阅（即设置cancel_at_period_end为true），customer.subscription.updated则会立即触发事件。该事件反映了订阅值的变化cancel_at_period_end。当订阅在期限结束时实际取消时，customer.subscription.deleted就会发生一个事件
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionCancel subscription err:%s", update)
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
	_, err := gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx, plan, planChannel, sub)
	if err != nil {
		return err
	}
	update, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0, // todo mark 如果您在计费周期结束时取消订阅（即设置cancel_at_period_end为true），customer.subscription.updated则会立即触发事件。该事件反映了订阅值的变化cancel_at_period_end。当订阅在期限结束时实际取消时，customer.subscription.deleted就会发生一个事件
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return gerror.Newf("SubscriptionCancel subscription err:%s", update)
	}
	return nil
}

func SubscriptionAddNewTrialEnd(ctx context.Context, subscriptionId string, AppendNewTrialEnd int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusPublished, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	planChannel := query.GetPlanChannel(ctx, sub.PlanId, sub.ChannelId)
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId) //todo mark 改造成支持 Merchant 级别的 PayChannel
	utility.Assert(payChannel != nil, "payChannel not found")

	details, err := gateway.GetPayChannelServiceProvider(ctx, sub.ChannelId).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, sub)
	utility.Assert(err == nil, fmt.Sprintf("SubscriptionDetail Fetch error%s", err))
	err = handler.UpdateSubWithChannelDetailBack(ctx, sub, details)
	utility.Assert(err == nil, fmt.Sprintf("UpdateSubWithChannelDetailBack Fetch error%s", err))
	//utility.Assert(newTrialEnd > details.CurrentPeriodEnd, "newTrainEnd should > subscription's currentPeriodEnd")
	utility.Assert(AppendNewTrialEnd > 0, "invalid AppendNewTrialEnd , should > 0")
	newTrialEnd := details.CurrentPeriodEnd + AppendNewTrialEnd*3600
	_, err = gateway.GetPayChannelServiceProvider(ctx, int64(payChannel.Id)).DoRemoteChannelSubscriptionNewTrialEnd(ctx, plan, planChannel, sub, newTrialEnd)
	if err != nil {
		return err
	}
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().TrialEnd:  newTrialEnd,
		dao.Subscription.Columns().GmtModify: gtime.Now(), // todo 存在并发调用问题
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
