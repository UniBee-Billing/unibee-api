package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	defaultAlipayClient "unibee/internal/logic/gateway/api/alipay/api"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request/pay"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
	"unibee/internal/logic/gateway/api/alipay/api/tools"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

//https://global.alipay.com/platform/site/ihome
//https://global.alipay.com/docs/
//https://github.com/alipay/global-open-sdk-go
//https://docs.antom.com/ac/ams_zh-cn/api

type Alipay struct {
}

func (c Alipay) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Alipay",
		Description:                   "Antom Online Card Payments, Use public and private keys to secure the Alipay payment.",
		DisplayName:                   "Alipay",
		GatewayWebsiteLink:            "https://global.alipay.com/platform/site/ihome",
		GatewayWebhookIntegrationLink: "",
		Sort:                          91,
		//GatewayLogo:                   "https://api.unibee.top/oss/file/d76q5bxsotbt0uzajb.png",
		//GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d76q5bxsotbt0uzajb.png"},
		GatewayLogo:       "https://api.unibee.top/oss/file/d7xy50zqf0iae7q9s6.png",
		GatewayIcons:      []string{"https://api.unibee.top/oss/file/d7xy50zqf0iae7q9s6.png"},
		GatewayType:       consts.GatewayTypeCard,
		SubGatewayName:    "Client Id",
		PublicKeyName:     "Alipay Public Key",
		PrivateSecretName: "Merchant Private Key",
		Host:              "https://open-de-global.alipay.com",
	}
}

func (c Alipay) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   0,
		CryptoCurrency: "USDT",
		Rate:           0,
	}, nil
}

func (c Alipay) GatewayTest(ctx context.Context, key string, secret string, subGateway string) (icon string, gatewayType int64, err error) {
	var alipayClientId = subGateway
	client := defaultAlipayClient.NewDefaultAlipayClient(
		"https://open-de-global.alipay.com",
		alipayClientId,
		secret,
		key, false)

	payRequest, request := pay.NewAlipayPayRequest()
	request.PaymentRequestId = "paymentRequestId01"
	order := &model.Order{}
	order.OrderDescription = "antom test order"
	order.ReferenceOrderId = "3232db07-91f7-4364-85bc-829a4c1c653f"
	order.OrderAmount = model.NewAmount("4200", "EUR")
	order.Buyer = &model.Buyer{
		BuyerEmail: "mail@hotmail.com",
	}
	request.Order = order
	request.PaymentMethod = &model.PaymentMethod{PaymentMethodType: "CARD"}
	request.PaymentAmount = model.NewAmount("4200", "EUR")
	request.PaymentNotifyUrl = "https://www.gaga.com/notify"
	request.PaymentRedirectUrl = "https://www.alipay.com"
	request.PaymentFactor = &model.PaymentFactor{
		IsAuthorization: true,
	}
	request.Env = &model.Env{ClientIp: utility.GetPublicIP(), TerminalType: model.WEB}
	request.ProductCode = model.CASHIER_PAYMENT

	execute, err := client.Execute(payRequest)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	response := execute.(*responsePay.AlipayPayResponse)
	g.Log().Debugf(ctx, "responseJson :%s", utility.MarshalToJsonString(response))
	utility.Assert(len(response.NormalUrl) > 0, "invalid keys, NormalUrl is nil")
	g.Log().Infof(ctx, "Redirect Url:%s", tools.Decode(response.NormalUrl))
	return "https://api.unibee.top/oss/file/d76q5bxsotbt0uzajb.png", consts.GatewayTypeCard, nil
}

func (c Alipay) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		gateway.Host,
		gateway.SubGateway,
		gateway.GatewaySecret,
		gateway.GatewayKey, false)
	payRequest, request := pay.NewAlipayPayRequest()
	request.PaymentRequestId = createPayContext.Pay.PaymentId
	{
		order := &model.Order{}
		{
			var name = ""
			var description = ""
			if len(createPayContext.Invoice.Lines) > 0 {
				var line = createPayContext.Invoice.Lines[0]
				if len(line.Name) == 0 {
					name = line.Description
				} else {
					name = line.Name
					description = line.Description
				}
			}
			order.OrderDescription = fmt.Sprintf("%s_%s", name, description)
		}
		order.ReferenceOrderId = createPayContext.Pay.PaymentId
		order.OrderAmount = model.NewAmount(fmt.Sprintf("%d", createPayContext.Pay.TotalAmount), createPayContext.Pay.Currency)
		order.Buyer = &model.Buyer{
			ReferenceBuyerId: fmt.Sprintf("%d", createPayContext.Pay.UserId),
			BuyerEmail:       createPayContext.Email,
		}
		var containNegative = false
		for _, line := range createPayContext.Invoice.Lines {
			if line.Amount <= 0 {
				containNegative = true
			}
		}
		var items []*model.Goods
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
				item := &model.Goods{
					GoodsName: fmt.Sprintf("%s", name),
					GoodsUnitAmount: &model.Amount{
						Currency: strings.ToLower(createPayContext.Pay.Currency),
						Value:    fmt.Sprintf("%d", line.Amount),
					},
					GoodsQuantity: fmt.Sprintf("%d", 1),
				}
				if len(description) > 0 {
					item.GoodsName = fmt.Sprintf("%s", description)
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
			item := &model.Goods{
				GoodsName: fmt.Sprintf("%s", productName),
				GoodsUnitAmount: &model.Amount{
					Currency: strings.ToLower(createPayContext.Pay.Currency),
					Value:    fmt.Sprintf("%d", createPayContext.Invoice.TotalAmount),
				},
				GoodsQuantity: fmt.Sprintf("%d", 1),
			}

			items = append(items, item)
		}

		request.Order = order
	}
	request.PaymentMethod = &model.PaymentMethod{PaymentMethodType: "CARD"}
	request.PaymentAmount = model.NewAmount(fmt.Sprintf("%d", createPayContext.Pay.TotalAmount), createPayContext.Pay.Currency)
	request.PaymentNotifyUrl = webhook2.GetPaymentWebhookEntranceUrl(createPayContext.Gateway.Id)
	request.PaymentRedirectUrl = webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true)
	request.PaymentFactor = &model.PaymentFactor{
		IsAuthorization: true,
	}
	request.ProductCode = model.CASHIER_PAYMENT
	request.Env = &model.Env{ClientIp: utility.GetPublicIP(), TerminalType: model.WEB}

	execute, err := client.Execute(payRequest)
	log.SaveChannelHttpLog("GatewayNewPayment", utility.MarshalToJsonString(payRequest), utility.MarshalToJsonString(execute), err, "AlipayNewPayment", nil, gateway)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	response := execute.(*responsePay.AlipayPayResponse)
	g.Log().Debugf(ctx, "responseJson :%s", utility.MarshalToJsonString(response))
	utility.Assert(len(response.NormalUrl) > 0, fmt.Sprintf("invalid keys, NormalUrl is nil,%s %s", response.Result.ResultCode, response.Result.ResultMessage))
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	gatewayPaymentId := response.PaymentId
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       gatewayPaymentId,
		GatewayPaymentIntentId: gatewayPaymentId,
		Link:                   tools.Decode(response.NormalUrl),
	}, nil
}

func (c Alipay) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		gateway.Host,
		gateway.SubGateway,
		gateway.GatewaySecret,
		gateway.GatewayKey, false)

	request, cancelRequest := pay.NewAlipayPayCancelRequest()
	cancelRequest.PaymentId = payment.GatewayPaymentId
	execute, err := client.Execute(request)
	log.SaveChannelHttpLog("GatewayPaymentCancel", utility.MarshalToJsonString(request), utility.MarshalToJsonString(execute), err, "AlipayPaymentCancel", nil, gateway)
	if err != nil {
		return nil, err
	}
	response := execute.(*responsePay.AlipayPayCancelResponse)
	utility.Assert(response != nil, "Alipay payment query failed, result is nil")
	utility.Assert(response != nil && response.Result.ResultCode == "SUCCESS", "invalid request, result not SUCCESS")
	detailRes, err := c.GatewayPaymentDetail(ctx, gateway, payment.GatewayPaymentId, payment)
	if err != nil {
		return nil, err
	}
	utility.Assert(detailRes != nil, "Alipay payment query failed, result is nil")
	return &gateway_bean.GatewayPaymentCancelResp{
		MerchantId:      strconv.FormatUint(payment.MerchantId, 10),
		GatewayCancelId: response.PaymentId,
		PaymentId:       payment.PaymentId,
		Status:          consts.PaymentStatusEnum(detailRes.Status),
	}, nil
}

func (c Alipay) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		gateway.Host,
		gateway.SubGateway,
		gateway.GatewaySecret,
		gateway.GatewayKey, false)
	queryRequest := pay.AlipayPayQueryRequest{}
	queryRequest.PaymentId = gatewayPaymentId
	request := queryRequest.NewRequest()
	execute, err := client.Execute(request)
	if err != nil {
		return nil, err
	}
	response := execute.(*responsePay.AlipayPayQueryResponse)
	log.SaveChannelHttpLog("GatewayPaymentDetail", utility.MarshalToJsonString(request), utility.MarshalToJsonString(execute), err, "AlipayPaymentDetail", nil, gateway)
	utility.Assert(response != nil, "Alipay payment query failed, result is nil")
	utility.Assert(response != nil && response.Result.ResultCode == "SUCCESS", "invalid keys, result not SUCCESS")

	return parseAlipayPayment(response), nil
}

func (c Alipay) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Alipay) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		gateway.Host,
		gateway.SubGateway,
		gateway.GatewaySecret,
		gateway.GatewayKey, false)
	queryRefundRequest := pay.AlipayInquiryRefundRequest{}
	queryRefundRequest.RefundId = gatewayRefundId
	request := queryRefundRequest.NewRequest()
	execute, err := client.Execute(request)
	log.SaveChannelHttpLog("GatewayRefundDetail", utility.MarshalToJsonString(request), utility.MarshalToJsonString(execute), err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	response := execute.(*responsePay.AlipayInquiryRefundResponse)
	utility.Assert(response != nil, "Alipay refund query failed, result is nil")
	utility.Assert(response != nil && response.RefundId != "", "invalid keys, resultId not found")
	var status consts.RefundStatusEnum = consts.RefundCreated
	if response.RefundStatus == model.TransactionStatusType_SUCCESS {
		status = consts.RefundSuccess
	} else if response.RefundStatus == model.TransactionStatusType_FAIL {
		status = consts.RefundFailed
	} else if response.RefundStatus == model.TransactionStatusType_CANCELLED {
		status = consts.RefundCancelled
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:      "",
		GatewayRefundId: response.RefundId,
		//GatewayPaymentId: response.AcquirerInfo.,
		Status:       status,
		Reason:       refund.RefundComment,
		RefundAmount: utility.ConvertCentStrToCent(response.RefundAmount.Value, response.RefundAmount.Currency),
		Currency:     strings.ToUpper(response.RefundAmount.Currency),
		RefundTime:   gtime.Now(),
	}, nil
}

func (c Alipay) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		gateway.Host,
		gateway.SubGateway,
		gateway.GatewaySecret,
		gateway.GatewayKey, false)
	refundRequest := pay.AlipayRefundRequest{}
	refundRequest.RefundRequestId = createPaymentRefundContext.Refund.RefundId
	refundRequest.PaymentId = createPaymentRefundContext.Payment.GatewayPaymentId
	refundRequest.RefundAmount = model.NewAmount(fmt.Sprintf("%d", createPaymentRefundContext.Refund.RefundAmount), createPaymentRefundContext.Refund.Currency)
	refundRequest.RefundReason = createPaymentRefundContext.Refund.RefundComment
	request := refundRequest.NewRequest()
	execute, err := client.Execute(request)
	log.SaveChannelHttpLog("GatewayRefund", utility.MarshalToJsonString(request), utility.MarshalToJsonString(execute), err, "refund", nil, gateway)
	utility.Assert(err == nil, fmt.Sprintf("call Alipay refund error %s", err))
	utility.Assert(execute != nil, "Alipay refund failed, result is nil")
	response := execute.(*responsePay.AlipayRefundResponse)
	utility.Assert(response != nil, "Alipay refund failed, result is nil")
	if response.RefundId == "" {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("invalid keys, resultId not found,%s %s", response.Result.ResultCode, response.Result.ResultMessage),
		}, nil
	}
	utility.Assert(response != nil && response.RefundId != "", fmt.Sprintf("invalid keys, resultId not found,%s %s", response.Result.ResultCode, response.Result.ResultMessage))

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: response.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (c Alipay) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("not support")
}

func parseAlipayPayment(item *responsePay.AlipayPayQueryResponse) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if item.PaymentStatus == model.TransactionStatusType_SUCCESS {
		status = consts.PaymentSuccess
	} else if item.PaymentStatus == model.TransactionStatusType_CANCELLED {
		status = consts.PaymentCancelled
	} else if item.PaymentStatus == model.TransactionStatusType_FAIL {
		status = consts.PaymentFailed
	}
	var authorizeReason = ""
	var paymentAmount = utility.ConvertCentStrToCent(item.PaymentAmount.Value, item.PaymentAmount.Currency)
	if item.ActualPaymentAmount.Currency != "" {
		paymentAmount = utility.ConvertCentStrToCent(item.ActualPaymentAmount.Value, item.ActualPaymentAmount.Currency)
	}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: item.PaymentId,
		Status:           status,
		AuthorizeStatus:  authorizeStatus,
		AuthorizeReason:  authorizeReason,
		CancelReason:     "",
		PaymentData:      utility.MarshalToJsonString(item),
		TotalAmount:      utility.ConvertCentStrToCent(item.PaymentAmount.Value, item.PaymentAmount.Currency),
		PaymentAmount:    paymentAmount,
		PaidTime:         gtime.Now(),
	}
}
