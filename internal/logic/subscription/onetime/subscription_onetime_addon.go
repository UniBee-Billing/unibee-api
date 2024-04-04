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
	MerchantId     uint64            `json:"merchantId" dc:"MerchantId"`
	SubscriptionId string            `json:"subscriptionId"  dc:"SubscriptionId" `
	AddonId        uint64            `json:"addonId" dc:"addonId"`
	Quantity       int64             `json:"quantity" dc:"Quantity"`
	RedirectUrl    string            `json:"redirectUrl"  dc:"RedirectUrl" `
	Metadata       map[string]string `json:"metadata" dc:"Metadata，Map"`
}

func CreateSubscriptionOneTimeAddon(ctx context.Context, req *SubscriptionCreateOnetimeAddonInternalReq) (*SubscriptionCreateOnetimeAddonInternalRes, error) {
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
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	user := query.GetUserAccountById(ctx, sub.UserId)
	utility.Assert(user != nil, "user not found")

	one := &entity.SubscriptionOnetimeAddon{
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
	invoice := &bean.InvoiceSimplify{
		InvoiceName:                    "OneTimeAddonPurchase-Subscription",
		TotalAmount:                    addon.Amount * req.Quantity,
		TotalAmountExcludingTax:        addon.Amount * req.Quantity,
		Currency:                       sub.Currency,
		TaxAmount:                      0,
		SubscriptionAmount:             0,
		SubscriptionAmountExcludingTax: 0,
		Lines: []*bean.InvoiceItemSimplify{{
			Currency:               sub.Currency,
			Amount:                 addon.Amount * req.Quantity,
			Tax:                    0,
			AmountExcludingTax:     addon.Amount * req.Quantity,
			TaxScale:               0,
			UnitAmountExcludingTax: addon.Amount,
			Description:            addon.Description,
			Proration:              false,
			Quantity:               req.Quantity,
			PeriodEnd:              0,
			PeriodStart:            0,
			Plan:                   bean.SimplifyPlan(addon),
		}},
	}

	createRes, err := service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		CheckoutMode: false,
		Gateway:      gateway,
		Pay: &entity.Payment{
			SubscriptionId:    req.SubscriptionId,
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
		Metadata:             map[string]string{"BillingReason": invoice.InvoiceName, "SubscriptionOnetimeAddonId": strconv.FormatUint(one.Id, 10)},
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
		dao.SubscriptionOnetimeAddon.Columns().PaymentId: createRes.PaymentId,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	return &SubscriptionCreateOnetimeAddonInternalRes{
		SubscriptionOnetimeAddon: bean.SimplifySubscriptionOnetimeAddonSimplify(one),
		Link:                     createRes.Link,
		Paid:                     createRes.Status == consts.PaymentSuccess,
		Invoice:                  bean.SimplifyInvoice(createRes.Invoice),
	}, nil
}
