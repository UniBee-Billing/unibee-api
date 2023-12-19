package ro

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/out/v1"
	"go-oversea-pay/api/out/vo"
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type CreatePayContext struct {
	OpenApiId                int64                     `json:"openApiId"`
	AppId                    string                    `json:"appId"`
	Desc                     string                    `json:"desc"`
	Pay                      *entity.OverseaPay        `json:"pay"`
	PayChannel               *entity.OverseaPayChannel `json:"payChannel"`
	PaymentBrandAddition     *gjson.Json               `json:"paymentBrandAddition"`
	TerminalIp               string                    `json:"terminalIp"`
	UserId                   string                    `json:"userId"`
	ShopperEmail             string                    `json:"shopperEmail"`
	ShopperLocale            string                    `json:"shopperLocale"`
	Mobile                   string                    `json:"mobile"`
	MediaInfo                *gjson.Json               `json:"mediaInfo"`
	Items                    []*v1.OutLineItem         `json:"items"`
	BillingDetails           *v1.OutPayAddress         `json:"billingDetails"`
	ShippingDetails          *v1.OutPayAddress         `json:"shippingDetails"`
	ShopperName              *v1.OutShopperName        `json:"shopperName"`
	ShopperInteraction       string                    `json:"shopperInteraction"`
	RecurringProcessingModel string                    `json:"recurringProcessingModel"`
	StorePaymentMethod       bool                      `json:"storePaymentMethod"`
	TokenId                  string                    `json:"tokenId"`
	DeviceFingerprint        string                    `json:"deviceFingerprint"`
	MerchantOrderReference   string                    `json:"merchantOrderReference"`
	DateOfBirth              *gtime.Time               `json:"dateOfBirth"`
	Platform                 string                    `json:"platform"`
	DeviceType               string                    `json:"deviceType"`
}

type CreatePayInternalResp struct {
	AlipayOrderNo  string      `json:"alipayOrderNo"`
	PayOrderNo     string      `json:"payOrderNo"`
	AlreadyPaid    bool        `json:"alreadyPaid"`
	OrderString    string      `json:"orderString"`
	Message        string      `json:"message"`
	TppOrderNo     string      `json:"tppOrderNo"`
	TppPayId       string      `json:"tppPayId"`
	PayChannel     int64       `json:"payChannel"`
	PayChannelType string      `json:"payChannelType"`
	Action         *gjson.Json `json:"action"`
	AdditionalData *gjson.Json `json:"additionalData"`
}

// OutPayCaptureRo is the golang structure for table oversea_pay.
type OutPayCaptureRo struct {
	MerchantId   string          `json:"merchantId"         `      // 商户ID
	PspReference string          `json:"pspReference"            ` // 业务类型。1-订单
	Reference    string          `json:"reference"              `  // 业务id-即商户订单号
	Amount       *vo.PayAmountVo `json:"amount"`
	Status       string          `json:"status"`
}

// OutPayCancelRo is the golang structure for table oversea_pay.
type OutPayCancelRo struct {
	MerchantId   string `json:"merchantId"         `      // 商户ID
	PspReference string `json:"pspReference"            ` // 业务类型。1-订单
	Reference    string `json:"reference"              `  // 业务id-即商户订单号
	Status       string `json:"status"`
}

// OutPayRefundRo is the golang structure for table oversea_pay.
type OutPayRefundRo struct {
	MerchantId      string      `json:"merchantId"         `          // 商户ID
	ChannelRefundNo string      `json:"channelRefundNo"            `  // 业务类型。1-订单
	ChargeRefundNo  string      `json:"chargeRefundNo"              ` // 业务id-即商户订单号
	RefundStatus    int         `json:"refundStatus"`
	Reason          string      `json:"reason"              `    // 业务id-即商户订单号
	RefundFee       int64       `json:"refundFee"              ` // 业务id-即商户订单号
	RefundTime      *gtime.Time `json:"refundTime" `             // 创建时间
}

// OutPayRo is the golang structure for table oversea_pay.
type OutPayRo struct {
	MerchantId      string      `json:"merchantId"         `        // 商户ID
	MerchantOrderNo string      `json:"merchantOrderNo"         `   // 商户ID
	ChannelTradeNo  string      `json:"ChannelTradeNo"            ` // 业务类型。1-订单
	ChannelPayId    string      `json:"channelPayId"              ` // 业务id-即商户订单号
	PayStatus       int         `json:"payStatus"`
	Reason          string      `json:"reason"              ` // 业务id-即商户订单号
	PayFee          int64       `json:"PayFee"              ` // 业务id-即商户订单号
	PayTime         *gtime.Time `json:"PayTime" `             // 创建时间
}

type CreateProductInternalResp struct {
	ChannelProductId     string `json:"channelProductId"`
	ChannelProductStatus string `json:"channelProductStatus"`
}

type CreatePlanInternalResp struct {
	ChannelPlanId     string                            `json:"channelPlanId"`
	ChannelPlanStatus string                            `json:"channelPlanStatus"`
	Data              string                            `json:"data"`
	Status            consts.SubscriptionPlanStatusEnum `json:"status"`
}

type CreateSubscriptionInternalResp struct {
	ChannelUserId             string                            `json:"channelUserId"`
	ChannelSubscriptionId     string                            `json:"channelSubscriptionId"`
	ChannelSubscriptionStatus string                            `json:"channelSubscriptionStatus"`
	Data                      string                            `json:"data"`
	Status                    consts.SubscriptionPlanStatusEnum `json:"status"`
}

type CancelSubscriptionInternalResp struct {
}

type UpdateSubscriptionInternalResp struct {
}

type ListSubscriptionInternalResp struct {
}

type WebhookSubscriptionInternalResp struct {
}
