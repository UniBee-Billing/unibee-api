package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type SubscriptionDetail struct {
	DayLeft                             int                              `json:"dayLeft" dc:"DayLeft util the period end, only available for webhook"`
	User                                *bean.UserAccount                `json:"user" dc:"user"`
	Subscription                        *bean.Subscription               `json:"subscription" dc:"Subscription"`
	Plan                                *bean.Plan                       `json:"plan" dc:"Plan"`
	Gateway                             *Gateway                         `json:"gateway" dc:"Gateway"`
	AddonParams                         []*bean.PlanAddonParam           `json:"addonParams" dc:"AddonParams"`
	Addons                              []*bean.PlanAddonDetail          `json:"addons" dc:"Addon"`
	LatestInvoice                       *bean.Invoice                    `json:"latestInvoice" dc:"LatestInvoice"`
	Discount                            *bean.MerchantDiscountCode       `json:"discount" dc:"Discount"`
	UnfinishedSubscriptionPendingUpdate *SubscriptionPendingUpdateDetail `json:"unfinishedSubscriptionPendingUpdate" dc:"processing pending update"`
}

type SubscriptionPendingUpdateEvent struct {
	MerchantId      uint64                  `json:"merchantId"           description:"MerchantId"`
	User            *bean.UserAccount       `json:"user" dc:"user"`
	Subscription    *bean.Subscription      `json:"subscription" dc:"Subscription"`
	Invoice         *bean.Invoice           `json:"invoice"            description:"invoice"`
	PendingUpdateId string                  `json:"pendingUpdateId"      description:"PendingUpdateId"`
	GmtCreate       *gtime.Time             `json:"gmtCreate"            description:"GmtCreate"`
	Amount          int64                   `json:"amount"               description:"CaptureAmount, Cent"`
	Status          int                     `json:"status"               description:"Status，1-Pending｜2-Finished｜3-Cancelled"`
	UpdateAmount    int64                   `json:"updateAmount"         description:"UpdateAmount, Cents"`
	ProrationAmount int64                   `json:"prorationAmount"      description:"ProrationAmount,Cents"`
	Currency        string                  `json:"currency"             description:"Currency"`
	UpdateCurrency  string                  `json:"updateCurrency"       description:"UpdateCurrency"`
	PlanId          uint64                  `json:"planId"               description:"PlanId"`
	UpdatePlanId    uint64                  `json:"updatePlanId"         description:"UpdatePlanId"`
	Quantity        int64                   `json:"quantity"             description:"quantity"`
	UpdateQuantity  int64                   `json:"updateQuantity"       description:"UpdateQuantity"`
	AddonData       string                  `json:"addonData"            description:"plan addon json data"`
	UpdateAddonData string                  `json:"updateAddonData"      description:"UpdateAddonData"`
	GatewayId       uint64                  `json:"gatewayId"            description:"Id"`
	GmtModify       *gtime.Time             `json:"gmtModify"            description:"GmtModify"`
	Paid            int                     `json:"paid"                 description:"Paid"`
	Link            string                  `json:"link"                 description:"Link"`
	MerchantMember  *MerchantMemberDetail   `json:"merchantMember"       description:"Merchant Member"`
	EffectImmediate int                     `json:"effectImmediate"      description:"EffectImmediate"`
	EffectTime      int64                   `json:"effectTime"           description:"effect_immediate=0, EffectTime unit_time"`
	Note            string                  `json:"note"                 description:"Update Note"`
	Plan            *bean.Plan              `json:"plan" dc:"Plan"`
	Addons          []*bean.PlanAddonDetail `json:"addons" dc:"Addons"`
	UpdatePlan      *bean.Plan              `json:"updatePlan" dc:"UpdatePlan"`
	UpdateAddons    []*bean.PlanAddonDetail `json:"updateAddons" dc:"UpdateAddons"`
	Metadata        map[string]interface{}  `json:"metadata" description:""`
}

type SubscriptionTimeLineDetail struct {
	MerchantId      uint64                  `json:"merchantId"      description:"merchant id"`             // merchant id
	UserId          uint64                  `json:"userId"          description:"userId"`                  // userId
	SubscriptionId  string                  `json:"subscriptionId"  description:"subscription id"`         // subscription id
	PeriodStart     int64                   `json:"periodStart"     description:"period_start"`            // period_start
	PeriodEnd       int64                   `json:"periodEnd"       description:"period_end"`              // period_end
	PeriodStartTime *gtime.Time             `json:"periodStartTime" description:"period start (datetime)"` // period start (datetime)
	PeriodEndTime   *gtime.Time             `json:"periodEndTime"   description:"period end (datatime)"`   // period end (datatime)
	InvoiceId       string                  `json:"invoiceId"       description:"invoice id"`              // invoice id
	PaymentId       string                  `json:"paymentId"       description:"payment id"`              // invoice id
	UniqueId        string                  `json:"uniqueId"        description:"unique id"`               // unique id
	Currency        string                  `json:"currency"        description:"currency"`                // currency
	PlanId          uint64                  `json:"planId"          description:"PlanId"`                  // PlanId
	Plan            *bean.Plan              `json:"plan" description:"Plan"`
	Quantity        int64                   `json:"quantity"        description:"quantity"` // quantity
	Addons          []*bean.PlanAddonDetail `json:"addons" description:"Addon"`
	GatewayId       uint64                  `json:"gatewayId"       description:"gateway_id"`                                  // gateway_id
	CreateTime      int64                   `json:"createTime"      description:"create utc time"`                             // create utc time
	Status          int                     `json:"status"          description:"1-processing,2-finish,3-cancelled,4-expired"` // 1-processing,2-finish
}

type SubscriptionPendingUpdateDetail struct {
	MerchantId      uint64                  `json:"merchantId"           description:"MerchantId"`
	UserId          uint64                  `json:"userId"               description:"UserId"`
	SubscriptionId  string                  `json:"subscriptionId"       description:"SubscriptionId"`
	PendingUpdateId string                  `json:"pendingUpdateId"      description:"PendingUpdateId"`
	InvoiceId       string                  `json:"invoiceId"            description:"invoice id"`
	GmtCreate       *gtime.Time             `json:"gmtCreate"            description:"GmtCreate"`
	Amount          int64                   `json:"amount"               description:"CaptureAmount, Cent"`
	Status          int                     `json:"status"               description:"Status，1-Pending｜2-Finished｜3-Cancelled"`
	UpdateAmount    int64                   `json:"updateAmount"         description:"UpdateAmount, Cents"`
	ProrationAmount int64                   `json:"prorationAmount"      description:"ProrationAmount,Cents"`
	Currency        string                  `json:"currency"             description:"Currency"`
	UpdateCurrency  string                  `json:"updateCurrency"       description:"UpdateCurrency"`
	PlanId          uint64                  `json:"planId"               description:"PlanId"`
	UpdatePlanId    uint64                  `json:"updatePlanId"         description:"UpdatePlanId"`
	Quantity        int64                   `json:"quantity"             description:"quantity"`
	UpdateQuantity  int64                   `json:"updateQuantity"       description:"UpdateQuantity"`
	AddonData       string                  `json:"addonData"            description:"plan addon json data"`
	UpdateAddonData string                  `json:"updateAddonData"      description:"UpdateAddonData"`
	GatewayId       uint64                  `json:"gatewayId"            description:"Id"`
	GmtModify       *gtime.Time             `json:"gmtModify"            description:"GmtModify"`
	Paid            int                     `json:"paid"                 description:"Paid"`
	Link            string                  `json:"link"                 description:"Link"`
	MerchantMember  *MerchantMemberDetail   `json:"merchantMember"       description:"Merchant Member"`
	EffectImmediate int                     `json:"effectImmediate"      description:"EffectImmediate"`
	EffectTime      int64                   `json:"effectTime"           description:"effect_immediate=0, EffectTime unit_time"`
	Note            string                  `json:"note"                 description:"Update Note"`
	Plan            *bean.Plan              `json:"plan" dc:"Plan"`
	Addons          []*bean.PlanAddonDetail `json:"addons" dc:"Addons"`
	UpdatePlan      *bean.Plan              `json:"updatePlan" dc:"UpdatePlan"`
	UpdateAddons    []*bean.PlanAddonDetail `json:"updateAddons" dc:"UpdateAddons"`
	Metadata        map[string]interface{}  `json:"metadata" description:""`
}

func ConvertSubscriptionPendingUpdateDetailByInvoiceId(ctx context.Context, invoiceId string) *SubscriptionPendingUpdateDetail {
	if len(invoiceId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().InvoiceId, invoiceId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("ConvertSubscriptionPendingUpdateDetailByInvoiceId Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &SubscriptionPendingUpdateDetail{
		MerchantId:      one.MerchantId,
		SubscriptionId:  one.SubscriptionId,
		PendingUpdateId: one.PendingUpdateId,
		GmtCreate:       one.GmtCreate,
		Amount:          one.Amount,
		Status:          one.Status,
		UpdateAmount:    one.UpdateAmount,
		Currency:        one.Currency,
		UpdateCurrency:  one.UpdateCurrency,
		PlanId:          one.PlanId,
		UpdatePlanId:    one.UpdatePlanId,
		Quantity:        one.Quantity,
		UpdateQuantity:  one.UpdateQuantity,
		AddonData:       one.AddonData,
		UpdateAddonData: one.UpdateAddonData,
		ProrationAmount: one.ProrationAmount,
		GatewayId:       one.GatewayId,
		UserId:          one.UserId,
		InvoiceId:       one.InvoiceId,
		GmtModify:       one.GmtModify,
		Paid:            one.Paid,
		Link:            one.Link,
		EffectImmediate: one.EffectImmediate,
		EffectTime:      one.EffectTime,
		Note:            one.Note,
		Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
		Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
		UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
		UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
		Metadata:        metadata,
	}
}

type SubscriptionOnetimeAddonDetail struct {
	Id             uint64                 `json:"id"             description:"id"`              // id
	SubscriptionId string                 `json:"subscriptionId" description:"subscription_id"` // subscription_id
	AddonId        uint64                 `json:"addonId"        description:"onetime addonId"` // onetime addonId
	Addon          *bean.Plan             `json:"addon"          description:"Addon"`
	Quantity       int64                  `json:"quantity"       description:"quantity"`                                      // quantity
	Status         int                    `json:"status"         description:"status, 1-create, 2-paid, 3-cancel, 4-expired"` // status, 1-create, 2-paid, 3-cancel, 4-expired
	CreateTime     int64                  `json:"createTime"     description:"create utc time"`                               // create utc time
	Payment        *bean.Payment          `json:"payment"        description:"Payment"`
	Metadata       map[string]interface{} `json:"metadata"       description:"Metadata"`
	User           *bean.UserAccount      `json:"user"           description:"User"`
}
