package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"net/http"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/api/paypal"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

// https://developer.paypal.com/docs/checkout/save-payment-methods/during-purchase/js-sdk/paypal/
// link：https://developer.paypal.com/docs/api/payments/v1/#payment_create
// https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
//APIBaseSandBox = "https://api-m.sandbox.paypal.com"
//APIBaseLive = "https://api-m.paypal.com"

type Paypal struct {
}

func (p Paypal) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:               "PayPal",
		Description:        "Use ClientId and Secret to secure the payment.",
		DisplayName:        "PayPal",
		GatewayWebsiteLink: "https://www.paypal.com/c2/home",
		GatewayLogo:        "https://api.unibee.top/oss/file/d76q3sb18tx0d9uclg.png",
		GatewayIcons:       []string{"https://api.unibee.top/oss/file/d6yhmepg0oh4xwzzjb.svg"},
		GatewayType:        consts.GatewayTypePaypal,
		Sort:               7,
	}
}

func (p Paypal) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	utility.Assert(userId > 0, "userId is nil")
	var paymentSource = &paypal.PaymentSource{
		Paypal: &paypal.PaymentSourcePaypal{
			Description:  "Save Payment Method For Future Use",
			UsagePattern: "IMMEDIATE",
			UsageType:    "MERCHANT",
			CustomerType: "CONSUMER",
			ExperienceContext: &paypal.ExperienceContext{
				ReturnURL: webhook2.GetPaymentMethodRedirectEntranceUrlCheckout(gateway.Id, true, fmt.Sprintf("%s", metadata["SubscriptionId"]), fmt.Sprintf("%s", metadata["RedirectUrl"])),
				CancelURL: webhook2.GetPaymentMethodRedirectEntranceUrlCheckout(gateway.Id, false, fmt.Sprintf("%s", metadata["SubscriptionId"]), fmt.Sprintf("%s", metadata["RedirectUrl"])),
			},
		},
	}
	result, err := c.NewSetupTokens(ctx, nil, paymentSource, fmt.Sprintf("%v%v", userId, gtime.Now().Timestamp()))
	log.SaveChannelHttpLog("GatewayUserCreateAndBindPaymentMethod", paymentSource, result, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	var approveLink = ""
	for _, link := range result.Links {
		if strings.Compare(link.Rel, "approve") == 0 {
			approveLink = link.Href
			break
		}
		if strings.Compare(link.Rel, "payer-action") == 0 {
			approveLink = link.Href
			break
		}
	}
	return &gateway_bean.GatewayUserPaymentMethodCreateAndBindResp{
		PaymentMethod: nil,
		Url:           approveLink,
	}, nil
}

func (p Paypal) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	c, _ := NewClient(key, secret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	if err == nil {
		_, vaultErr := c.GetPaymentMethodTokens(ctx, "BEEB8ANDETATE")
		if re, ok := vaultErr.(*paypal.ErrorResponse); ok {
			utility.Assert(re.Response != nil && re.Response.StatusCode != 403, "Insufficient permissions to start automatic payment,see https://developer.paypal.com/docs/checkout/save-payment-methods/during-purchase/orders-api/paypal/")
		}
	}
	return "https://www.paypalobjects.com/webstatic/icon/favicon.ico", consts.GatewayTypePaypal, err
}

func (p Paypal) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	order, err := c.GetOrder(ctx, gatewayPaymentId)
	log.SaveChannelHttpLog("GatewayPaymentDetail", c.RequestBodyStr, c.ResponseStr, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return p.parsePaypalPayment(ctx, gateway, order, payment)
}

func (p Paypal) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	gatewayRefund, err := c.GetRefund(ctx, gatewayRefundId)
	log.SaveChannelHttpLog("GatewayRefundDetail", c.RequestBodyStr, c.ResponseStr, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return p.parsePaypalRefund(ctx, gateway, gatewayRefund)
}

func (p Paypal) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("not support")
}

func (p Paypal) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	var items = make([]paypal.Item, 0)
	for _, line := range createPayContext.Invoice.Lines {
		var name = ""
		var description = ""
		if len(line.Name) == 0 {
			name = line.Description
		} else {
			name = line.Name
			description = line.Description
		}
		item := paypal.Item{
			Name:        name,
			Description: description,
			UnitAmount: &paypal.Money{
				Value:    utility.ConvertCentToDollarStr(line.Amount, createPayContext.Pay.Currency),
				Currency: strings.ToUpper(createPayContext.Pay.Currency),
			},
			Quantity: strconv.FormatInt(line.Quantity, 10),
		}
		items = append(items, item)
	}

	var productName = createPayContext.Invoice.ProductName
	if len(productName) == 0 {
		productName = createPayContext.Invoice.InvoiceName
	}
	if len(productName) == 0 {
		productName = "DefaultProduct"
	}
	var paymentSource = &paypal.PaymentSource{
		Paypal: &paypal.PaymentSourcePaypal{},
	}
	if createPayContext.PayImmediate && !createPayContext.CheckoutMode && len(createPayContext.GatewayPaymentMethod) > 0 {
		paymentSource.Paypal.VaultId = createPayContext.GatewayPaymentMethod
	} else {
		paymentSource.Paypal.Attributes = &paypal.PaymentSourceAttributes{
			Vault: &paypal.PaymentSourceAttributesVault{
				StoreInVault: "ON_SUCCESS",
				UsageType:    "MERCHANT",
			},
		}
	}
	detail, err := c.CreateOrder(
		ctx,
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
					Currency: strings.ToUpper(createPayContext.Pay.Currency),
				},
			},
		},
		&paypal.CreateOrderPayer{},
		paymentSource,
		&paypal.ApplicationContext{
			BrandName:          "",
			Locale:             "",
			ShippingPreference: "",
			UserAction:         "",
			PaymentMethod:      paypal.PaymentMethod{},
			ReturnURL:          webhook2.GetPaypalPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
			CancelURL:          webhook2.GetPaypalPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		},
		createPayContext.Pay.PaymentId,
	)
	log.SaveChannelHttpLog("GatewayNewPayment", c.RequestBodyStr, c.ResponseStr, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}

	payment, err := p.parsePaypalPayment(ctx, gateway, detail, createPayContext.Pay)
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 consts.PaymentStatusEnum(payment.Status),
		GatewayPaymentId:       payment.GatewayPaymentId,
		GatewayPaymentIntentId: payment.GatewayPaymentId,
		GatewayPaymentMethod:   payment.GatewayPaymentMethod,
		Link:                   payment.Link,
	}, nil
}

func (p Paypal) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	utility.Assert(payment != nil, "payment not found")
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	captureRes, err := c.CaptureOrder(ctx, payment.GatewayPaymentId, paypal.CaptureOrderRequest{})
	log.SaveChannelHttpLog("GatewayCapture", c.RequestBodyStr, c.ResponseStr, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayPaymentCaptureResp{
		MerchantId:       gateway.MerchantId,
		GatewayCaptureId: captureRes.ID,
		Amount:           payment.PaymentAmount,
		Currency:         payment.Currency,
	}, nil
}

func (p Paypal) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (p Paypal) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(len(payment.PaymentData) > 0, "payment capture data not found")
	var availableCapture *paypal.CaptureAmount
	gatewayPaymentRo, err := p.GatewayPaymentDetail(ctx, gateway, payment.GatewayPaymentId, payment)
	utility.Assert(len(gatewayPaymentRo.PaymentData) > 0, "available capture not found")
	utility.AssertError(utility.UnmarshalFromJsonString(gatewayPaymentRo.PaymentData, &availableCapture), "parse capture data error")
	utility.Assert(availableCapture != nil, "available capture not found")
	utility.Assert(refund != nil, "refund not found")
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, p.GetPaypalHost())
	_, err = c.GetAccessToken(ctx)
	captureRefundRes, err := c.RefundCapture(ctx, availableCapture.ID, paypal.RefundCaptureRequest{
		Amount: &paypal.Money{
			Currency: strings.ToUpper(refund.Currency),
			Value:    utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency),
		},
		InvoiceID:   refund.RefundId,
		NoteToPayer: refund.RefundComment,
	})
	log.SaveChannelHttpLog("GatewayRefund", c.RequestBodyStr, c.ResponseStr, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}

	return p.GatewayRefundDetail(ctx, gateway, captureRefundRes.ID, refund)
}

func (p Paypal) parsePaypalRefund(ctx context.Context, gateway *entity.MerchantGateway, item *paypal.Refund) (*gateway_bean.GatewayPaymentRefundResp, error) {
	var status consts.RefundStatusEnum = consts.RefundCreated
	if strings.Compare(item.Status, "COMPLETED") == 0 {
		status = consts.RefundSuccess
	} else if strings.Compare(item.Status, "FAILED") == 0 {
		status = consts.RefundFailed
	} else if strings.Compare(item.Status, "CANCELLED") == 0 {
		status = consts.RefundCancelled
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:      "",
		GatewayRefundId: item.ID,
		Status:          status,
		Reason:          item.NoteToPayer,
		RefundAmount:    utility.ConvertDollarStrToCent(item.Amount.Value, item.Amount.Currency),
		Currency:        strings.ToUpper(item.Amount.Currency),
		RefundTime:      gtime.New(item.UpdateTime),
	}, nil
}

func (p Paypal) parsePaypalPayment(ctx context.Context, gateway *entity.MerchantGateway, item *paypal.Order, payment *entity.Payment) (*gateway_bean.GatewayPaymentRo, error) {
	var availableCapture *paypal.CaptureAmount
	var paidTime *gtime.Time
	if item.PurchaseUnits != nil && len(item.PurchaseUnits) > 0 && item.PurchaseUnits[0].Payments != nil && len(item.PurchaseUnits[0].Payments.Captures) >= 1 {
		for _, one := range item.PurchaseUnits[0].Payments.Captures {
			if strings.Compare(item.Status, "COMPLETED") == 0 ||
				strings.Compare(item.Status, "REFUNDED") == 0 ||
				strings.Compare(item.Status, "PARTIALLY_REFUNDED") == 0 {
				availableCapture = &one
				break
			}
		}
		if availableCapture != nil && availableCapture.UpdateTime != nil {
			paidTime = gtime.New(availableCapture.UpdateTime)
		}
	}

	var cancelTime *gtime.Time
	var status = consts.PaymentCreated
	if strings.Compare(item.Status, "COMPLETED") == 0 && availableCapture != nil {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Status, "VOIDED") == 0 {
		status = consts.PaymentFailed
		cancelTime = gtime.New(item.UpdateTime)
	}
	var captureStatus = consts.Authorized
	var authorizeReason = ""
	if strings.Compare(item.Status, "CREATED") == 0 ||
		strings.Compare(item.Status, "SAVED") == 0 ||
		strings.Compare(item.Status, "PAYER_ACTION_REQUIRED") == 0 {
		captureStatus = consts.WaitingAuthorized
	} else if strings.Compare(item.Status, "APPROVED") == 0 {
		captureStatus = consts.Authorized
	} else if strings.Compare(item.Status, "COMPLETED") == 0 {
		captureStatus = consts.CaptureRequest
	}
	var gatewayPaymentMethod string
	var gatewayUserId string
	var paymentCode string
	if item.PaymentSource != nil &&
		item.PaymentSource.Paypal != nil &&
		item.PaymentSource.Paypal.Attributes != nil &&
		item.PaymentSource.Paypal.Attributes.Vault != nil &&
		len(item.PaymentSource.Paypal.Attributes.Vault.Id) > 0 && strings.Compare(item.PaymentSource.Paypal.Attributes.Vault.Status, "VAULTED") == 0 {
		gatewayPaymentMethod = item.PaymentSource.Paypal.Attributes.Vault.Id
		if item.PaymentSource.Paypal.Attributes.Vault.Customer != nil {
			gatewayUserId = item.PaymentSource.Paypal.Attributes.Vault.Customer.Id
		}
		if len(gatewayPaymentMethod) > 0 {
			paymentCode = utility.MarshalToJsonString(item.PaymentSource)
		}
	}
	var approveLink = ""
	for _, link := range item.Links {
		if strings.Compare(link.Rel, "approve") == 0 {
			approveLink = link.Href
			break
		}
		if strings.Compare(link.Rel, "payer-action") == 0 {
			approveLink = link.Href
			break
		}
	}
	var createTime *gtime.Time
	if item.CreateTime != nil {
		createTime = gtime.New(item.CreateTime)
	}
	var lastError = ""
	if len(item.FailureReason) > 0 {
		lastError = item.FailureReason
	}
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.ID,
		Status:               status,
		AuthorizeStatus:      captureStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          utility.MarshalToJsonString(availableCapture),
		TotalAmount:          payment.TotalAmount,
		PaymentAmount:        payment.TotalAmount,
		GatewayUserId:        gatewayUserId,
		GatewayPaymentMethod: gatewayPaymentMethod,
		PaymentCode:          paymentCode,
		Currency:             payment.Currency,
		PaidTime:             paidTime,
		CreateTime:           createTime,
		CancelTime:           cancelTime,
		Link:                 approveLink,
		LastError:            lastError,
	}, nil
}

// todo mark Verify that it is possible to transform to singletons without having to get them every time accessToken
func NewClient(clientID string, secret string, APIBase string) (*paypal.Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &paypal.Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

func init() {
	// gateway_webhook_entry
}

func (p Paypal) GetPaypalHost() string {
	var apiHost = "https://api-m.paypal.com"
	if !config.GetConfigInstance().IsProd() {
		apiHost = "https://api-m.sandbox.paypal.com"
	}
	return apiHost
}
