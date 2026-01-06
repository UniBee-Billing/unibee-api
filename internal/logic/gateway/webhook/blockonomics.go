package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type BlockonomicsWebhook struct {
}

func (b BlockonomicsWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (b BlockonomicsWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	// Blockonomics does not need to set up webhooks, uses polling to check payment status
	return nil
}

func GetLatestPaymentByGatewayPaymentIntentId(ctx context.Context, gatewayPaymentId string) (one *entity.Payment) {
	if len(gatewayPaymentId) == 0 {
		return nil
	}
	err := dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().GatewayPaymentIntentId, gatewayPaymentId).
		Where(dao.Payment.Columns().Status, consts.PaymentCreated).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func (b BlockonomicsWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	// Blockonomics supports two webhook methods:
	// 1. Pass payment information through URL parameters
	// 2. Pass payment information through POST request body
	// sample data : https://api.unibee.top?addr=bc1qfwktz5zq06vum0hslqavnnlq6fd64k7ta5akz5&status=1&crypto=BTC&txid=WarningThisIsAGeneratedTestPaymentAndNotARealBitcoinTransaction&value=100000000 200

	var jsonData *gjson.Json
	var err error
	var gatewayPaymentId string = ""
	var txid string = ""
	var status int = 0
	var amount int64 = 0

	// Check if there is a POST request body
	if r.Method == "POST" && len(r.GetBody()) > 0 {
		jsonData, err = r.GetJson()
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Failed to parse webhook body:%s error: %s\n", gateway.GatewayName, r.GetBody(), err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
		} else {
			gatewayPaymentId = jsonData.Get("addr").String()
			txid = jsonData.Get("txid").String()
			status = jsonData.Get("status").Int()
			amount = jsonData.Get("value").Int64()
		}
	} else {
		// Build data from URL parameters
		jsonData = gjson.New("")
		gatewayPaymentId = r.Get("addr").String()
		_ = jsonData.Set("addr", gatewayPaymentId)
		txid = r.Get("txid").String()
		_ = jsonData.Set("txid", txid)
		status = r.Get("status").Int()
		_ = jsonData.Set("status", status)
		amount = r.Get("value").Int64()
		_ = jsonData.Set("value", amount)
	}

	g.Log().Info(r.Context(), "Receive_Webhook_Channel:", gateway.GatewayName, " hook:", jsonData.String())

	var payment *entity.Payment
	if gatewayPaymentId != "" {
		if strings.HasPrefix(gatewayPaymentId, "0x") {
			if txid != "" {
				payment = query.GetPaymentByGatewayPaymentId(r.Context(), txid)
			}
		} else {
			payment = GetLatestPaymentByGatewayPaymentIntentId(r.Context(), gatewayPaymentId)
		}
	}

	// Process payment notification
	var responseBack = http.StatusOK
	if payment != nil {
		if !config.GetConfigInstance().IsProd() && txid == "WarningThisIsAGeneratedTestPaymentAndNotARealBitcoinTransaction" {
			//sandbox test
			var lastErr string = ""
			if status == 0 {
				lastErr = "Unconfirmed"
			} else if status == 1 {
				lastErr = "Partially Confirmed"
			} else if status == 2 {
				lastErr = "Confirmed"
				err = handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
					PaymentId:              payment.PaymentId,
					GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
					GatewayPaymentId:       payment.GatewayPaymentId,
					PayStatusEnum:          consts.PaymentSuccess,
					PaymentAmount:          payment.TotalAmount,
					PaymentCode:            txid,
				})
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error HandlePaySuccess: %s\n", gateway.GatewayName, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
			handler2.UpdatePaymentLastGatewayError(r.Context(), payment.PaymentId, lastErr)
		} else {
			err = ProcessPaymentWebhook(r.Context(), payment.PaymentId, payment.GatewayPaymentId, gateway)
			if err != nil {
				g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Error ProcessPaymentWebhook: %s\n", gateway.GatewayName, err.Error())
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			}
		}
	} else {
		g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Unhandled paymentId\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		responseBack = http.StatusBadRequest
	}
	log.SaveChannelHttpLog("GatewayWebhook", jsonData, responseBack, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	r.Response.WriteHeader(responseBack)
	return
}

func (b BlockonomicsWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	// Handle payment redirect
	payIdStr := r.Get("paymentId").String()
	txHash := r.Get("txhash").String()
	var response string
	var status = false
	var returnUrl = ""
	var isSuccess = false
	var payment *entity.Payment
	if len(payIdStr) > 0 {
		response = ""
		//Payment Redirect
		payment = query.GetPaymentByPaymentId(r.Context(), payIdStr)
		if payment != nil {
			success := r.Get("success")
			if success != nil {
				if success.String() == "true" {
					isSuccess = true
				}
				returnUrl = util.GetPaymentRedirectUrl(r.Context(), payment, success.String())
			} else {
				returnUrl = util.GetPaymentRedirectUrl(r.Context(), payment, "")
			}
			if len(txHash) > 0 {
				payment.Code = txHash
				payment.GatewayPaymentId = txHash
				_, _ = dao.Payment.Ctx(r.Context()).Data(g.Map{
					dao.Payment.Columns().Code:             txHash,
					dao.Payment.Columns().GatewayPaymentId: txHash,
				}).Where(dao.Payment.Columns().PaymentId, payment.PaymentId).OmitNil().Update()
			}
		}
		if r.Get("success").Bool() {
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				response = "success"
				status = true
			} else {
				//find
				paymentIntentDetail, err := api.GetGatewayServiceProvider(r.Context(), gateway.Id).GatewayPaymentDetail(r.Context(), gateway, payment.GatewayPaymentId, payment)
				if err != nil {
					response = fmt.Sprintf("%v", err)
				} else {
					if paymentIntentDetail.Status == consts.PaymentSuccess {
						err := handler2.HandlePaySuccess(r.Context(), &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							GatewayUserId:          paymentIntentDetail.GatewayUserId,
							TotalAmount:            paymentIntentDetail.TotalAmount,
							PayStatusEnum:          consts.PaymentSuccess,
							PaidTime:               paymentIntentDetail.PaidTime,
							PaymentAmount:          paymentIntentDetail.PaymentAmount,
							Reason:                 paymentIntentDetail.Reason,
							GatewayPaymentMethod:   paymentIntentDetail.GatewayPaymentMethod,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							response = "payment success"
							status = true
						}
					} else if paymentIntentDetail.Status == consts.PaymentFailed {
						err := handler2.HandlePayFailure(r.Context(), &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							PayStatusEnum:          consts.PaymentFailed,
							Reason:                 paymentIntentDetail.Reason,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						}
					}
				}
			}
		} else {
			response = "user cancelled"
		}
	}
	log.SaveChannelHttpLog("GatewayRedirect", r.URL, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	return &gateway_bean.GatewayRedirectResp{
		Payment:   payment,
		Status:    status,
		Message:   response,
		Success:   isSuccess,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}
