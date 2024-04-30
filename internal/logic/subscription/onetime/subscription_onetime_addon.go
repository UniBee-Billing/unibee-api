package onetime

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
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionCreateOnetimeAddonInternalRes struct {
	MerchantId               uint64                                 `json:"merchantId" dc:"MerchantId"`
	SubscriptionOnetimeAddon *bean.SubscriptionOnetimeAddonSimplify `json:"subscriptionOnetimeAddon"  dc:"SubscriptionOnetimeAddon" `
	Paid                     bool                                   `json:"paid"`
	Link                     string                                 `json:"link"`
	Invoice                  *bean.InvoiceSimplify                  `json:"invoice"  dc:"Invoice" `
}

type SubscriptionCreateOnetimeAddonInternalReq struct {
	MerchantId         uint64                 `json:"merchantId" dc:"MerchantId"`
	SubscriptionId     string                 `json:"subscriptionId"  dc:"SubscriptionId" `
	AddonId            uint64                 `json:"addonId" dc:"addonId"`
	Quantity           int64                  `json:"quantity" dc:"Quantity"`
	RedirectUrl        string                 `json:"redirectUrl"  dc:"RedirectUrl" `
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                 `json:"discountCode"        dc:"DiscountCode, ignore if discountAmount or discountPercentage provide"`
	DiscountAmount     *int64                 `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage *int64                 `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount provide"`
	TaxPercentage      *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, use subscription's taxPercentage if not provide"`
	GatewayId          *uint64                `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
}

func CreateSubOneTimeAddon(ctx context.Context, req *SubscriptionCreateOnetimeAddonInternalReq) (*SubscriptionCreateOnetimeAddonInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.AddonId > 0, "AddonId invalid")
	req.Quantity = utility.MaxInt64(req.Quantity, 1)
	addon := query.GetPlanById(ctx, req.AddonId)
	utility.Assert(addon != nil, "addon not found")
	utility.Assert(addon.Status == consts.PlanStatusActive, "addon not active")
	utility.Assert(addon.Type == consts.PlanTypeOnetimeAddon, "addon not onetime type")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub.Currency == addon.Currency, "Server error: currency not match")
	utility.Assert(sub != nil, "sub not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "sub not active")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "sub plan not found")
	utility.Assert(plan.Status == consts.PlanStatusActive, "addon not active")
	utility.Assert(strings.Contains(plan.BindingOnetimeAddonIds, strconv.FormatUint(req.AddonId, 10)), "plan not contain this addon")
	var gatewayId = sub.GatewayId
	if req.GatewayId != nil {
		gatewayId = *req.GatewayId
	}
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	user := query.GetUserAccountById(ctx, sub.UserId)
	utility.Assert(user != nil, "user not found")
	var taxPercentage = sub.TaxPercentage
	if req.TaxPercentage != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "External TaxPercentage only available for api call")
		utility.Assert(*req.TaxPercentage > 0 && *req.TaxPercentage < 10000, "invalid taxPercentage")
		taxPercentage = *req.TaxPercentage
	}

	one := &entity.SubscriptionOnetimeAddon{
		UserId:         sub.UserId,
		SubscriptionId: req.SubscriptionId,
		AddonId:        req.AddonId,
		Quantity:       req.Quantity,
		Status:         1,
		CreateTime:     gtime.Now().Timestamp(),
		MetaData:       utility.MarshalToJsonString(req.Metadata),
	}

	result, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	totalAmountExcludingTax := addon.Amount * req.Quantity
	var discountAmount int64 = 0

	if req.DiscountAmount != nil && *req.DiscountAmount > 0 {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		discountAmount = utility.MinInt64(*req.DiscountAmount, totalAmountExcludingTax)
	} else if req.DiscountPercentage != nil && *req.DiscountPercentage > 0 {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		utility.Assert(*req.DiscountPercentage > 0 && *req.DiscountPercentage < 10000, "invalid discountPercentage")
		discountAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(*req.DiscountPercentage))
	} else if len(req.DiscountCode) > 0 {
		discountCode := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
		utility.Assert(discountCode.Type == 0, "invalid code, code is from external")
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:   req.MerchantId,
			UserId:       sub.UserId,
			DiscountCode: req.DiscountCode,
			Currency:     plan.Currency,
			PLanId:       plan.Id,
		})
		utility.Assert(canApply, message)
		utility.Assert(!isRecurring, "recurring discount code not available for one-time addon")
		discountAmount = utility.MinInt64(discount.ComputeDiscountAmount(ctx, plan.MerchantId, totalAmountExcludingTax, sub.Currency, req.DiscountCode, gtime.Now().Timestamp()), totalAmountExcludingTax)
	}

	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(taxPercentage))
	invoice := &bean.InvoiceSimplify{
		InvoiceName:             "OneTimeAddonPurchase-Subscription",
		OriginAmount:            totalAmountExcludingTax + taxAmount + discountAmount,
		TotalAmount:             totalAmountExcludingTax + taxAmount,
		DiscountCode:            req.DiscountCode,
		DiscountAmount:          discountAmount,
		TotalAmountExcludingTax: totalAmountExcludingTax,
		Currency:                sub.Currency,
		TaxPercentage:           taxPercentage,
		TaxAmount:               taxAmount,
		Lines: []*bean.InvoiceItemSimplify{{
			Currency:               sub.Currency,
			OriginAmount:           addon.Amount*req.Quantity + taxAmount,
			Amount:                 addon.Amount*req.Quantity + taxAmount - discountAmount,
			DiscountAmount:         discountAmount,
			Tax:                    taxAmount,
			AmountExcludingTax:     addon.Amount*req.Quantity - discountAmount,
			UnitAmountExcludingTax: addon.Amount,
			Description:            addon.Description,
			Quantity:               req.Quantity,
			Plan:                   bean.SimplifyPlan(addon),
		}},
	}

	createRes, err := service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		CheckoutMode: false,
		Gateway:      gateway,
		Pay: &entity.Payment{
			ExternalPaymentId: strconv.FormatUint(one.Id, 10),
			BizType:           consts.BizTypeOneTime,
			UserId:            sub.UserId,
			GatewayId:         gateway.Id,
			TotalAmount:       invoice.TotalAmount,
			Currency:          invoice.Currency,
			CountryCode:       sub.CountryCode,
			MerchantId:        sub.MerchantId,
			CompanyId:         0,
			ReturnUrl:         req.RedirectUrl,
		},
		Email:                user.Email,
		Metadata:             map[string]interface{}{"BillingReason": invoice.InvoiceName, "SubscriptionOnetimeAddonId": strconv.FormatUint(one.Id, 10)},
		Invoice:              invoice,
		PayImmediate:         true,
		GatewayPaymentMethod: sub.GatewayDefaultPaymentMethod,
	})
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	//update paymentId
	status := 1
	if createRes.Status == consts.PaymentSuccess {
		status = 2
	}
	_, err = dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().Status:    status,
		dao.SubscriptionOnetimeAddon.Columns().PaymentId: createRes.Payment.PaymentId,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	if len(req.DiscountCode) > 0 {
		_, err = discount.UserDiscountApply(ctx, &discount.UserDiscountApplyReq{
			MerchantId:     req.MerchantId,
			UserId:         sub.UserId,
			DiscountCode:   invoice.DiscountCode,
			SubscriptionId: one.SubscriptionId,
			PaymentId:      createRes.Payment.PaymentId,
			InvoiceId:      invoice.InvoiceId,
			ApplyAmount:    invoice.DiscountAmount,
			Currency:       invoice.Currency,
		})
		if err != nil {
			// todo mark success payment
			fmt.Printf("UserDiscountApply onetimeAddon Purchase createRes:%s err:%s", utility.MarshalToJsonString(createRes), err.Error())
		}
	}

	return &SubscriptionCreateOnetimeAddonInternalRes{
		SubscriptionOnetimeAddon: bean.SimplifySubscriptionOnetimeAddonSimplify(one),
		Link:                     createRes.Link,
		Paid:                     createRes.Status == consts.PaymentSuccess,
		Invoice:                  bean.SimplifyInvoice(createRes.Invoice),
	}, nil
}
