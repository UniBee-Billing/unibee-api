package gateway_bean

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/oversea_pay"
)

type GatewayNewPaymentReq struct {
	CheckoutMode         bool                    `json:"checkoutMode"`
	Pay                  *entity.Payment         `json:"pay"`
	Gateway              *entity.MerchantGateway `json:"gateway"`
	ExternalUserId       string                  `json:"externalUserId"`
	Email                string                  `json:"email"`
	Metadata             map[string]string       `json:"metadata"`
	Invoice              *bean.InvoiceSimplify   `json:"invoice"`
	DaysUtilDue          int                     `json:"daysUtilDue"`
	GatewayPaymentMethod string                  `json:"gatewayPaymentMethod"`
	PayImmediate         bool                    `json:"payImmediate"`
}

type GatewayNewPaymentResp struct {
	Status                 consts.PaymentStatusEnum `json:"status"`
	PaymentId              string                   `json:"paymentId"`
	GatewayPaymentId       string                   `json:"gatewayPaymentId"`
	GatewayPaymentIntentId string                   `json:"gatewayPaymentIntentId"`
	GatewayPaymentMethod   string                   `json:"gatewayPaymentMethod"`
	Link                   string                   `json:"link"`
	Action                 *gjson.Json              `json:"action"`
	Invoice                *entity.Invoice          `json:"invoice"`
}

// GatewayPaymentCaptureResp is the golang structure for table oversea_pay.
type GatewayPaymentCaptureResp struct {
	MerchantId       string `json:"merchantId"         `
	GatewayCaptureId string `json:"gatewayCaptureId"            `
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Status           string `json:"status"`
}

// GatewayPaymentCancelResp is the golang structure for table oversea_pay.
type GatewayPaymentCancelResp struct {
	MerchantId      string                   `json:"merchantId"         `
	GatewayCancelId string                   `json:"gatewayCancelId"            `
	PaymentId       string                   `json:"paymentId"              `
	Status          consts.PaymentStatusEnum `json:"status"`
}

// GatewayPaymentRefundResp is the golang structure for table oversea_pay.
type GatewayPaymentRefundResp struct {
	MerchantId       string                  `json:"merchantId"         `
	GatewayRefundId  string                  `json:"gatewayRefundId"            `
	GatewayPaymentId string                  `json:"gatewayPaymentId"            `
	Status           consts.RefundStatusEnum `json:"status"`
	Reason           string                  `json:"reason"              `
	RefundAmount     int64                   `json:"refundFee"              `
	Currency         string                  `json:"currency"              `
	RefundTime       *gtime.Time             `json:"refundTime" `
}

type GatewayPaymentListReq struct {
	UserId int64 `json:"userId"         `
}

// GatewayPaymentRo is the golang structure for table oversea_pay.
type GatewayPaymentRo struct {
	MerchantId           uint64      `json:"merchantId"         `
	Status               int         `json:"status"`
	AuthorizeStatus      int         `json:"captureStatus"`
	AuthorizeReason      string      `json:"authorizeReason" `
	Currency             string      `json:"currency"              `
	TotalAmount          int64       `json:"totalAmount"              `
	PaymentAmount        int64       `json:"paymentAmount"              `
	BalanceAmount        int64       `json:"balanceAmount"              `
	RefundAmount         int64       `json:"refundAmount"              `
	BalanceStart         int64       `json:"balanceStart"              `
	BalanceEnd           int64       `json:"balanceEnd"              `
	Reason               string      `json:"reason"              `
	PayTime              *gtime.Time `json:"payTime" `
	CreateTime           *gtime.Time `json:"createTime" `
	CancelTime           *gtime.Time `json:"cancelTime" `
	CancelReason         string      `json:"cancelReason" `
	PaymentData          string      `json:"paymentData" `
	GatewayId            uint64      `json:"gatewayId"         `
	GatewayPaymentId     string      `json:"gatewayPaymentId"              `
	GatewayPaymentMethod string      `json:"gatewayPaymentMethod"              `
}

type GatewayCreateSubscriptionResp struct {
	GatewayUserId         string                                   `json:"gatewayUserId"`
	GatewaySubscriptionId string                                   `json:"gatewaySubscriptionId"`
	Data                  string                                   `json:"data"`
	Link                  string                                   `json:"link"`
	Status                consts.SubscriptionGatewayPlanStatusEnum `json:"status"`
	Paid                  bool                                     `json:"paid"`
}

type GatewayBalance struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type GatewayUserCreateResp struct {
	GatewayUserId string `json:"gatewayUserId"`
}

type GatewayUserDetailQueryResp struct {
	GatewayUserId        string            `json:"gatewayUserId"`
	DefaultPaymentMethod string            `json:"defaultPaymentMethod"`
	Balance              *GatewayBalance   `json:"balance"`
	CashBalance          []*GatewayBalance `json:"cashBalance"`
	InvoiceCreditBalance []*GatewayBalance `json:"invoiceCreditBalance"`
	Email                string            `json:"email"`
	Description          string            `json:"description"`
}

type GatewayUserAttachPaymentMethodResp struct {
}

type GatewayUserDeAttachPaymentMethodResp struct {
}

type GatewayUserPaymentMethodReq struct {
	UserId           int64  `json:"userId"`
	GatewayPaymentId string `json:"gatewayPaymentId"`
}

type GatewayUserPaymentMethodListResp struct {
	PaymentMethods []*bean.PaymentMethod `json:"paymentMethods"`
}

type GatewayUserPaymentMethodCreateAndBindResp struct {
	PaymentMethod *bean.PaymentMethod `json:"paymentMethod"`
}

type GatewayMerchantBalanceQueryResp struct {
	AvailableBalance       []*GatewayBalance `json:"available"`
	ConnectReservedBalance []*GatewayBalance `json:"connectReserved"`
	PendingBalance         []*GatewayBalance `json:"pending"`
}

type GatewayRedirectResp struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ReturnUrl string `json:"returnUrl"`
	QueryPath string `json:"queryPath"`
}