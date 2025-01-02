package subscription

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ConfigReq struct {
	g.Meta `path:"/config" tags:"Subscription Config" method:"get" summary:"Subscription Config"`
}
type ConfigRes struct {
	Config *bean.SubscriptionConfig `json:"config" dc:"Config"`
}

type ConfigUpdateReq struct {
	g.Meta                             `path:"/config/update" tags:"Subscription Config" method:"post" summary:"Update Merchant Subscription Config"`
	DowngradeEffectImmediately         *bool                   `json:"downgradeEffectImmediately" dc:"DowngradeEffectImmediately, Immediate Downgrade (by default, the downgrades takes effect at the end of the period ）"`
	UpgradeProration                   *bool                   `json:"upgradeProration" dc:"UpgradeProration, Prorated Upgrade Invoices(Upgrades will generate prorated invoice by default)"`
	IncompleteExpireTime               *int64                  `json:"incompleteExpireTime" dc:"IncompleteExpireTime, seconds, Incomplete Status Duration(The period during which subscription remains in “incomplete”)"`
	InvoiceEmail                       *bool                   `json:"invoiceEmail" dc:"InvoiceEmail, Enable Invoice Email (Toggle to send invoice email to customers)"`
	TryAutomaticPaymentBeforePeriodEnd *int64                  `json:"tryAutomaticPaymentBeforePeriodEnd" dc:"TryAutomaticPaymentBeforePeriodEnd, Auto-charge Start Before Period End （Time Difference for Auto-Payment Activation Before Period End）"`
	GatewayVATRule                     []*bean.MerchantVatRule `json:"gatewayVATRule" dc:""`
	ShowZeroInvoice                    *bool                   `json:"showZeroInvoice" dc:"ShowZeroInvoice, Display Invoices With Zero Amount (Invoice With Zero Amount will hidden in list by default)"`
}

type ConfigUpdateRes struct {
	Config *bean.SubscriptionConfig `json:"config" dc:"Config"`
}
