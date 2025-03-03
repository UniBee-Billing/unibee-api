package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/balance"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/invoice"
	"github.com/stripe/stripe-go/v78/invoiceitem"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/paymentmethod"
	"github.com/stripe/stripe-go/v78/product"
	"github.com/stripe/stripe-go/v78/refund"
	"strconv"
	"strings"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Stripe struct {
}

func (s Stripe) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:               "Stripe",
		Description:        "Use public and private keys to secure the bank card payment.",
		DisplayName:        "Bank Cards",
		GatewayWebsiteLink: "https://stripe.com",
		GatewayLogo:        "https://api.unibee.top/oss/file/d76q2e3zyv4ylc6vyh.png",
		GatewayIcons:       []string{"https://api.unibee.top/oss/file/d6yhl1qz7qmcg6zafr.svg", "https://api.unibee.top/oss/file/d6yhlf1t8n3ev3ueii.svg", "https://api.unibee.top/oss/file/d6yhlpshof3muufphd.svg"},
		GatewayType:        consts.GatewayTypeCard,
		Sort:               100,
		AutoChargeEnabled:  true,
	}
}

func (s Stripe) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (s Stripe) GatewayTest(ctx context.Context, key string, secret string, subGateway string) (icon string, gatewayType int64, err error) {
	stripe.Key = secret
	s.setUnibeeAppInfo()
	utility.Assert(len(secret) > 0, "invalid gatewaySecret")
	utility.Assert(strings.HasPrefix(secret, "sk_"), "invalid gatewaySecret, should start with 'sk_'")

	params := &stripe.ProductListParams{}
	params.Limit = stripe.Int64(3)
	result := product.List(params)
	return "http://unibee.top/files/invoice/stripe.png", consts.GatewayTypeCard, result.Err()
}

func (s Stripe) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, userId)
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(gatewayUser.GatewayUserId),
	}
	_, err = paymentmethod.Attach(gatewayPaymentMethod, params)
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayUserAttachPaymentMethodResp{}, nil
}

func (s Stripe) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.PaymentMethodDetachParams{}
	_, err = paymentmethod.Detach(gatewayPaymentMethod, params)
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayUserDeAttachPaymentMethodResp{}, nil
}

func (s Stripe) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	utility.Assert(userId > 0, "userId is nil")
	gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, userId)
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSetup)),
		//Currency:           stripe.String(strings.ToUpper(currency)),
		Customer:           stripe.String(gatewayUser.GatewayUserId),
		Metadata:           utility.ConvertToStringMetadata(metadata),
		PaymentMethodTypes: []*string{stripe.String(string(stripe.PaymentMethodTypeCard))},
		PaymentMethodData:  &stripe.CheckoutSessionPaymentMethodDataParams{AllowRedisplay: stripe.String(string(stripe.PaymentMethodAllowRedisplayAlways))},
		SuccessURL:         stripe.String(webhook2.GetPaymentMethodRedirectEntranceUrlCheckout(gateway.Id, true, fmt.Sprintf("%s", metadata["SubscriptionId"]), fmt.Sprintf("%s", metadata["RedirectUrl"]))),
		CancelURL:          stripe.String(webhook2.GetPaymentMethodRedirectEntranceUrlCheckout(gateway.Id, false, fmt.Sprintf("%s", metadata["SubscriptionId"]), fmt.Sprintf("%s", metadata["RedirectUrl"]))),
	}
	if len(currency) > 0 {
		params.Currency = stripe.String(strings.ToUpper(currency))
	}
	result, err := session.New(params)
	log.SaveChannelHttpLog("GatewayUserCreateAndBindPaymentMethod", params, result, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayUserPaymentMethodCreateAndBindResp{
		PaymentMethod: nil,
		Url:           result.URL,
	}, nil
}

func (s Stripe) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	var paymentMethods = make([]*gateway_bean.PaymentMethod, 0)
	if len(req.GatewayPaymentMethodId) > 0 {
		params := &stripe.PaymentMethodParams{}
		paymentMethod, err := paymentmethod.Get(req.GatewayPaymentMethodId, params)
		if err == nil {
			if paymentMethod.Type == stripe.PaymentMethodTypeCard {
				data := gjson.New(nil)
				_ = data.Set("brand", paymentMethod.Card.Brand)
				_ = data.Set("checks", paymentMethod.Card.Checks)
				_ = data.Set("country", paymentMethod.Card.Country)
				_ = data.Set("last4", paymentMethod.Card.Last4)
				_ = data.Set("expMonth", paymentMethod.Card.ExpMonth)
				_ = data.Set("expYear", paymentMethod.Card.ExpYear)
				_ = data.Set("fingerprint", paymentMethod.Card.Fingerprint)
				_ = data.Set("description", paymentMethod.Card.Description)
				paymentMethods = append(paymentMethods, &gateway_bean.PaymentMethod{
					Id:   paymentMethod.ID,
					Type: "card",
					Data: data,
				})
			}
		}
	} else {
		utility.Assert(req.UserId > 0, "userId is nil")
		gatewayUserId := req.GatewayUserId
		if len(gatewayUserId) == 0 {
			gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, req.UserId)
			utility.Assert(gatewayUser != nil, "stripe create customer error")
			gatewayUserId = gatewayUser.GatewayUserId
		}

		params := &stripe.CustomerListPaymentMethodsParams{
			Customer: stripe.String(gatewayUserId),
		}
		params.Limit = stripe.Int64(10)
		result := customer.ListPaymentMethods(params)

		for _, paymentMethod := range result.PaymentMethodList().Data {
			// only append card type
			if paymentMethod.Type == stripe.PaymentMethodTypeCard {
				data := gjson.New(nil)
				_ = data.Set("brand", paymentMethod.Card.Brand)
				_ = data.Set("checks", paymentMethod.Card.Checks)
				_ = data.Set("country", paymentMethod.Card.Country)
				_ = data.Set("last4", paymentMethod.Card.Last4)
				_ = data.Set("expMonth", paymentMethod.Card.ExpMonth)
				_ = data.Set("expYear", paymentMethod.Card.ExpYear)
				_ = data.Set("fingerprint", paymentMethod.Card.Fingerprint)
				_ = data.Set("description", paymentMethod.Card.Description)
				paymentMethods = append(paymentMethods, &gateway_bean.PaymentMethod{
					Id:   paymentMethod.ID,
					Type: "card",
					Data: data,
				})
			}
		}
	}

	return &gateway_bean.GatewayUserPaymentMethodListResp{
		PaymentMethods: paymentMethods,
	}, nil
}

func (s Stripe) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.CustomerParams{
		//Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
		Email: stripe.String(user.Email),
	}

	createCustomResult, err := customer.New(params)
	log.SaveChannelHttpLog("GatewayUserCreate", params, createCustomResult, err, "", nil, gateway)
	if err != nil {
		g.Log().Printf(ctx, "customer.New: %v", err.Error())
		return nil, err
	}
	return &gateway_bean.GatewayUserCreateResp{GatewayUserId: createCustomResult.ID}, nil

}

func (s Stripe) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, listReq.UserId)

	params := &stripe.PaymentIntentListParams{}
	params.Customer = stripe.String(gatewayUser.GatewayUserId)
	params.Limit = stripe.Int64(200)
	paymentList := paymentintent.List(params)
	log.SaveChannelHttpLog("GatewayPaymentList", params, paymentList, err, "", nil, gateway)
	var list []*gateway_bean.GatewayPaymentRo
	for _, item := range paymentList.PaymentIntentList().Data {
		list = append(list, s.parseStripePayment(ctx, gateway, item))
	}

	return list, nil
}

func (s Stripe) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.RefundListParams{}
	params.PaymentIntent = stripe.String(gatewayPaymentId)
	params.Limit = stripe.Int64(100)
	refundList := refund.List(params)
	log.SaveChannelHttpLog("GatewayRefundList", params, refundList, err, "", nil, gateway)
	var list []*gateway_bean.GatewayPaymentRefundResp
	for _, item := range refundList.RefundList().Data {
		list = append(list, parseStripeRefund(item))
	}

	return list, nil
}

func (s Stripe) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.PaymentIntentParams{}
	response, err := paymentintent.Get(gatewayPaymentId, params)
	log.SaveChannelHttpLog("GatewayPaymentDetail", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}

	return s.parseStripePayment(ctx, gateway, response), nil
}

func (s Stripe) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.BalanceParams{}
	response, err := balance.Get(params)
	if err != nil {
		return nil, err
	}

	var availableBalances []*gateway_bean.GatewayBalance
	for _, item := range response.Available {
		availableBalances = append(availableBalances, &gateway_bean.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var connectReservedBalances []*gateway_bean.GatewayBalance
	for _, item := range response.ConnectReserved {
		connectReservedBalances = append(connectReservedBalances, &gateway_bean.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var pendingBalances []*gateway_bean.GatewayBalance
	for _, item := range response.ConnectReserved {
		pendingBalances = append(pendingBalances, &gateway_bean.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance:       availableBalances,
		ConnectReservedBalance: connectReservedBalances,
		PendingBalance:         pendingBalances,
	}, nil
}

func (s Stripe) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.CustomerParams{}
	response, err := customer.Get(gatewayUserId, params)
	if err != nil {
		return nil, err
	}
	var cashBalances []*gateway_bean.GatewayBalance
	if response.CashBalance != nil {
		for currency, amount := range response.CashBalance.Available {
			cashBalances = append(cashBalances, &gateway_bean.GatewayBalance{
				Amount:   amount,
				Currency: strings.ToUpper(currency),
			})
		}
	}

	var invoiceCreditBalances []*gateway_bean.GatewayBalance
	for currency, amount := range response.InvoiceCreditBalance {
		invoiceCreditBalances = append(invoiceCreditBalances, &gateway_bean.GatewayBalance{
			Amount:   amount,
			Currency: strings.ToUpper(currency),
		})
	}
	var defaultPaymentMethod string
	if response.InvoiceSettings != nil && response.InvoiceSettings.DefaultPaymentMethod != nil {
		defaultPaymentMethod = response.InvoiceSettings.DefaultPaymentMethod.ID
	}
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId:        response.ID,
		DefaultPaymentMethod: defaultPaymentMethod,
		Balance: &gateway_bean.GatewayBalance{
			Amount:   response.Balance,
			Currency: strings.ToUpper(string(response.Currency)),
		},
		CashBalance:          cashBalances,
		InvoiceCreditBalance: invoiceCreditBalances,
		Description:          response.Description,
		Email:                response.Email,
	}, nil
}

// Test Card Data
// Payment Success
// 4242 4242 4242 4242
// Need 3DS
// 4000 0025 0000 3155
// Reject
// 4000 0000 0000 9995
func (s Stripe) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "UniBee.api",
		Version: "1.0.0",
		URL:     "https://merchant.unibee.dev",
	})
}

func (s Stripe) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, createPayContext.Pay.UserId)

	if createPayContext.CheckoutMode {
		var items []*stripe.CheckoutSessionLineItemParams
		var containNegative = false
		for _, line := range createPayContext.Invoice.Lines {
			if line.Amount <= 0 {
				containNegative = true
			}
		}
		if !containNegative {
			for _, line := range createPayContext.Invoice.Lines {
				var name = ""
				var description = ""
				if len(line.Name) == 0 {
					name = line.Description
				} else {
					name = line.Name
					description = line.Description
				}
				item := &stripe.CheckoutSessionLineItemParams{
					PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
						Currency: stripe.String(strings.ToLower(createPayContext.Pay.Currency)),
						ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
							Name: stripe.String(fmt.Sprintf("%s", name)),
							//Description: stripe.String(fmt.Sprintf("%s", description)),
						},
						UnitAmount: stripe.Int64(line.Amount),
					},
					Quantity: stripe.Int64(1),
				}
				if len(description) > 0 {
					item.PriceData.ProductData.Description = stripe.String(fmt.Sprintf("%s", description))
				}
				items = append(items, item)
			}
		} else {
			var productName = createPayContext.Invoice.ProductName
			if len(productName) == 0 {
				productName = createPayContext.Invoice.InvoiceName
			}
			if len(productName) == 0 {
				productName = "DefaultProduct"
			}
			item := &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(strings.ToLower(createPayContext.Pay.Currency)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("%s", productName)),
					},
					UnitAmount: stripe.Int64(createPayContext.Invoice.TotalAmount),
				},
				Quantity: stripe.Int64(1),
			}

			items = append(items, item)
		}
		checkoutParams := &stripe.CheckoutSessionParams{
			Customer:          stripe.String(gatewayUser.GatewayUserId),
			Currency:          stripe.String(strings.ToLower(createPayContext.Pay.Currency)),
			LineItems:         items,
			PaymentMethodData: &stripe.CheckoutSessionPaymentMethodDataParams{AllowRedisplay: stripe.String(string(stripe.PaymentMethodAllowRedisplayAlways))},
			SuccessURL:        stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true)),
			CancelURL:         stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false)),
			Metadata:          utility.ConvertToStringMetadata(createPayContext.Metadata),
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				Metadata:         utility.ConvertToStringMetadata(createPayContext.Metadata),
				SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
			},
		}
		//if len(gatewayUser.GatewayDefaultPaymentMethod) > 0 {
		//	checkoutParams.PaymentMethodConfiguration = stripe.String(gatewayUser.GatewayDefaultPaymentMethod)
		//}
		checkoutParams.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
		checkoutParams.Metadata = utility.ConvertToStringMetadata(createPayContext.Metadata)
		//checkoutParams.ExpiresAt
		detail, err := session.New(checkoutParams)
		log.SaveChannelHttpLog("GatewayNewPayment", checkoutParams, detail, err, "CheckoutSession", nil, gateway)
		if err != nil {
			return nil, err
		}
		var status consts.PaymentStatusEnum = consts.PaymentCreated
		if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusOpen)) == 0 {
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusComplete)) == 0 {
			status = consts.PaymentSuccess
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusExpired)) == 0 {
			status = consts.PaymentFailed
		}
		var gatewayPaymentId string
		if detail.PaymentIntent != nil {
			gatewayPaymentId = detail.PaymentIntent.ID
		}
		return &gateway_bean.GatewayNewPaymentResp{
			Status:                 status,
			GatewayPaymentId:       gatewayPaymentId,
			GatewayPaymentIntentId: detail.ID,
			Link:                   detail.URL,
		}, nil
	} else {
		// need payment link
		params := &stripe.InvoiceParams{
			Metadata: utility.ConvertToStringMetadata(createPayContext.Metadata),
			Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
			Customer: stripe.String(gatewayUser.GatewayUserId)}

		if createPayContext.PayImmediate {
			params.CollectionMethod = stripe.String("charge_automatically")
			// check the gateway user contains the payment method now
			listQuery, err := s.GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
				UserId: gatewayUser.UserId,
			})
			var paymentMethodIds = make([]string, 0)
			for _, paymentMethod := range listQuery.PaymentMethods {
				paymentMethodIds = append(paymentMethodIds, paymentMethod.Id)
			}
			if err != nil {
				return nil, err
			}
			if len(createPayContext.GatewayPaymentMethod) > 0 && ContainString(paymentMethodIds, createPayContext.GatewayPaymentMethod) {
				params.DefaultPaymentMethod = stripe.String(createPayContext.GatewayPaymentMethod)
			} else if len(gatewayUser.GatewayDefaultPaymentMethod) > 0 && ContainString(paymentMethodIds, gatewayUser.GatewayDefaultPaymentMethod) {
				params.DefaultPaymentMethod = stripe.String(gatewayUser.GatewayDefaultPaymentMethod)
			} else if len(listQuery.PaymentMethods) > 0 {
				params.DefaultPaymentMethod = stripe.String(listQuery.PaymentMethods[0].Id)
			}
		} else {
			params.CollectionMethod = stripe.String("send_invoice")
			if createPayContext.DaysUtilDue > 0 {
				params.DaysUntilDue = stripe.Int64(int64(createPayContext.DaysUtilDue))
			} else {
				params.DaysUntilDue = stripe.Int64(5)
			}
		}
		result, err := invoice.New(params)
		log.SaveChannelHttpLog("GatewayNewPayment", params, result, err, "NewInvoice", nil, gateway)
		if err != nil {
			return nil, err
		}

		for _, line := range createPayContext.Invoice.Lines {
			var description = line.Description
			if len(description) == 0 {
				if line.Plan != nil {
					description = line.Plan.PlanName
				} else {
					description = "Default Product"
				}
			}
			ItemParams := &stripe.InvoiceItemParams{
				Invoice:     stripe.String(result.ID),
				Currency:    stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
				Amount:      stripe.Int64(line.Amount),
				Description: stripe.String(description),
				Customer:    stripe.String(gatewayUser.GatewayUserId)}
			itemResult, err := invoiceitem.New(ItemParams)
			log.SaveChannelHttpLog("GatewayNewPayment", ItemParams, itemResult, err, "NewInvoiceItem", nil, gateway)
			if err != nil {
				return nil, err
			}
		}
		finalizeInvoiceParam := &stripe.InvoiceFinalizeInvoiceParams{}
		if createPayContext.PayImmediate {
			finalizeInvoiceParam.AutoAdvance = stripe.Bool(true)
		} else {
			finalizeInvoiceParam.AutoAdvance = stripe.Bool(false)
		}
		detail, err := invoice.FinalizeInvoice(result.ID, finalizeInvoiceParam)
		log.SaveChannelHttpLog("GatewayNewPayment", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, gateway)
		if err != nil {
			return nil, err
		}
		if createPayContext.PayImmediate && strings.Compare(string(detail.Status), "paid") != 0 {
			paymentParam := &stripe.InvoicePayParams{}
			if len(createPayContext.GatewayPaymentMethod) > 0 {
				paymentParam.PaymentMethod = stripe.String(createPayContext.GatewayPaymentMethod)
			}
			response, payErr := invoice.Pay(result.ID, paymentParam)
			log.SaveChannelHttpLog("GatewayNewPayment", params, response, payErr, "PayInvoice", nil, gateway)
			if response != nil && payErr == nil {
				detail.Status = response.Status
			}
		}

		var status consts.PaymentStatusEnum = consts.PaymentCreated
		if strings.Compare(string(detail.Status), "draft") == 0 {
		} else if strings.Compare(string(detail.Status), "open") == 0 {
		} else if strings.Compare(string(detail.Status), "paid") == 0 {
			status = consts.PaymentSuccess
		} else if strings.Compare(string(detail.Status), "uncollectible") == 0 {
			status = consts.PaymentFailed
		} else if strings.Compare(string(detail.Status), "void") == 0 {
			status = consts.PaymentFailed
		}
		var gatewayPaymentId string
		if detail.PaymentIntent != nil {
			gatewayPaymentId = detail.PaymentIntent.ID
		}
		var gatewayPaymentMethod string
		if detail.PaymentIntent != nil && detail.PaymentIntent.PaymentMethod != nil {
			gatewayPaymentMethod = detail.PaymentIntent.PaymentMethod.ID
		}
		return &gateway_bean.GatewayNewPaymentResp{
			Status:                 status,
			GatewayPaymentId:       gatewayPaymentId,
			GatewayPaymentIntentId: detail.ID,
			Link:                   detail.HostedInvoiceURL,
			GatewayPaymentMethod:   gatewayPaymentMethod,
		}, nil
	}
}

func ContainString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s Stripe) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	utility.Assert(payment.GatewayId > 0, "invalid payment gatewayId")
	utility.Assert(len(payment.GatewayPaymentIntentId) > 0, "invalid payment GatewayPaymentIntentId")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	var status = consts.PaymentCreated
	var gatewayCancelId string
	if strings.HasPrefix(payment.GatewayPaymentIntentId, "in_") {
		params := &stripe.InvoiceVoidInvoiceParams{}
		result, err := invoice.VoidInvoice(payment.GatewayPaymentIntentId, params)
		log.SaveChannelHttpLog("GatewayCancel", payment.GatewayPaymentIntentId, result, err, "VoidInvoice", nil, gateway)
		if err != nil {
			result, err = invoice.Get(payment.GatewayPaymentIntentId, &stripe.InvoiceParams{})
			log.SaveChannelHttpLog("GatewayCancel", payment.GatewayPaymentIntentId, result, err, "GetInvoice", nil, gateway)
			if err != nil {
				return nil, err
			}
		}
		gatewayCancelId = result.ID
		if strings.Compare(string(result.Status), string(stripe.InvoiceStatusPaid)) == 0 {
			status = consts.PaymentSuccess
		} else if strings.Compare(string(result.Status), string(stripe.InvoiceStatusUncollectible)) == 0 || strings.Compare(string(result.Status), string(stripe.InvoiceStatusVoid)) == 0 {
			status = consts.PaymentFailed
		}
	} else if strings.HasPrefix(payment.GatewayPaymentIntentId, "cs_") {
		params := &stripe.CheckoutSessionExpireParams{}
		result, err := session.Expire(
			payment.GatewayPaymentIntentId,
			params,
		)
		log.SaveChannelHttpLog("GatewayCancel", payment.GatewayPaymentIntentId, result, err, "ExpireSession", nil, gateway)
		if err != nil {
			result, err = session.Get(payment.GatewayPaymentIntentId, &stripe.CheckoutSessionParams{})
			log.SaveChannelHttpLog("GatewayCancel", payment.GatewayPaymentIntentId, result, err, "GetSession", nil, gateway)
			if err != nil {
				return nil, err
			}
		}
		if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusOpen)) == 0 {
		} else if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusComplete)) == 0 {
			status = consts.PaymentSuccess
		} else if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusExpired)) == 0 {
			status = consts.PaymentFailed
		}
		gatewayCancelId = result.ID
	} else {
		params := &stripe.PaymentIntentCancelParams{}
		result, err := paymentintent.Cancel(payment.GatewayPaymentIntentId, params)
		log.SaveChannelHttpLog("GatewayCancel", params, result, err, "CancelPaymentIntent", nil, gateway)
		if err != nil {
			result, err = paymentintent.Get(payment.GatewayPaymentIntentId, &stripe.PaymentIntentParams{})
			if err != nil {
				return nil, err
			}
		}
		paymentDetails := s.parseStripePayment(ctx, gateway, result)
		status = paymentDetails.Status
		gatewayCancelId = paymentDetails.GatewayPaymentId
	}

	return &gateway_bean.GatewayPaymentCancelResp{
		MerchantId:      strconv.FormatUint(payment.MerchantId, 10),
		GatewayCancelId: gatewayCancelId,
		PaymentId:       payment.PaymentId,
		Status:          consts.PaymentStatusEnum(status),
	}, nil
}

func (s Stripe) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(createPaymentRefundContext.Payment.GatewayId > 0, "Gateway Not Found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{PaymentIntent: stripe.String(createPaymentRefundContext.Payment.GatewayPaymentId)}
	params.Reason = stripe.String("requested_by_customer")
	params.Amount = stripe.Int64(createPaymentRefundContext.Refund.RefundAmount)
	var metadata = make(map[string]interface{})
	if len(createPaymentRefundContext.Refund.MetaData) > 0 {
		err := utility.UnmarshalFromJsonString(createPaymentRefundContext.Refund.MetaData, &metadata)
		if err != nil {
			g.Log().Errorf(ctx, "GatewayRefund Unmarshal Metadata error:%s", err.Error())
		}
	}
	params.Metadata = utility.ConvertToStringMetadata(metadata)
	result, err := refund.New(params)
	log.SaveChannelHttpLog("GatewayRefund", params, result, err, "refund", nil, gateway)
	utility.Assert(err == nil, fmt.Sprintf("call stripe refund error %s", err))
	utility.Assert(result != nil, "Stripe refund failed, result is nil")
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: result.ID,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (s Stripe) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, one *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{}
	response, err := refund.Get(gatewayRefundId, params)
	log.SaveChannelHttpLog("GatewayRefundDetail", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return parseStripeRefund(response), nil
}

func (s Stripe) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, one *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	utility.Assert(payment.GatewayId > 0, "Gateway Not Found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundCancelParams{}
	response, err := refund.Cancel(one.GatewayRefundId, params)
	log.SaveChannelHttpLog("GatewayRefundCancel", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return parseStripeRefund(response), nil
}

func parseStripeRefund(item *stripe.Refund) *gateway_bean.GatewayPaymentRefundResp {
	var gatewayPaymentId string
	if item.PaymentIntent != nil {
		gatewayPaymentId = item.PaymentIntent.ID
	}
	var status consts.RefundStatusEnum = consts.RefundCreated
	if strings.Compare(string(item.Status), string(stripe.RefundStatusSucceeded)) == 0 {
		status = consts.RefundSuccess
	} else if strings.Compare(string(item.Status), string(stripe.RefundStatusFailed)) == 0 {
		status = consts.RefundFailed
	} else if strings.Compare(string(item.Status), string(stripe.RefundStatusCanceled)) == 0 {
		status = consts.RefundCancelled
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:       "",
		GatewayRefundId:  item.ID,
		GatewayPaymentId: gatewayPaymentId,
		Status:           status,
		Reason:           string(item.Reason),
		RefundAmount:     item.Amount,
		Currency:         strings.ToUpper(string(item.Currency)),
		RefundTime:       gtime.NewFromTimeStamp(item.Created),
	}
}

func (s Stripe) parseStripePayment(ctx context.Context, gateway *entity.MerchantGateway, item *stripe.PaymentIntent) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	if strings.Compare(string(item.Status), string(stripe.PaymentIntentStatusSucceeded)) == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(string(item.Status), string(stripe.PaymentIntentStatusCanceled)) == 0 {
		status = consts.PaymentCancelled
	}
	var captureStatus = consts.Authorized
	var authorizeReason = ""
	var paymentData = ""
	if strings.Compare(string(item.Status), string(stripe.PaymentIntentStatusRequiresPaymentMethod)) == 0 || strings.Compare(string(item.Status), string(stripe.PaymentIntentStatusRequiresAction)) == 0 {
		captureStatus = consts.WaitingAuthorized
		if item.LastPaymentError != nil {
			authorizeReason = item.LastPaymentError.Msg
		}
	} else if strings.Compare(string(item.Status), string(stripe.PaymentIntentStatusRequiresConfirmation)) == 0 {
		captureStatus = consts.CaptureRequest
	}
	if item.NextAction != nil {
		paymentData = utility.MarshalToJsonString(item.NextAction)
	}
	var gatewayPaymentMethod string
	var paymentCode string
	if item.PaymentMethod != nil {
		gatewayPaymentMethod = item.PaymentMethod.ID
		if len(gatewayPaymentMethod) > 0 {
			query, _ := s.GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{GatewayPaymentMethodId: gatewayPaymentMethod})
			if query != nil {
				paymentCode = utility.MarshalToJsonString(query)
			}
		}
	}
	var lastError = ""
	if item.LastPaymentError != nil {
		lastError = item.LastPaymentError.Msg
	}
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.ID,
		Status:               status,
		AuthorizeStatus:      captureStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         string(item.CancellationReason),
		PaymentData:          paymentData,
		TotalAmount:          item.Amount,
		PaymentAmount:        item.AmountReceived,
		GatewayPaymentMethod: gatewayPaymentMethod,
		PaymentCode:          paymentCode,
		Currency:             strings.ToUpper(string(item.Currency)),
		PaidTime:             gtime.NewFromTimeStamp(item.Created),
		CreateTime:           gtime.NewFromTimeStamp(item.Created),
		CancelTime:           gtime.NewFromTimeStamp(item.CanceledAt),
		LastError:            lastError,
	}
}
