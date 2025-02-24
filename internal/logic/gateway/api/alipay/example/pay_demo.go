package main

import (
	"fmt"
	"github.com/google/uuid"
	defaultAlipayClient "unibee/internal/logic/gateway/api/alipay/api"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request/pay"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

func main() {
	const alipayGatewayUrl = "https://open-de-global.alipay.com"
	const alipayClientId = "SANDBOX_5YES442ZS5S203863"
	const alipayMerchantPrivateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCAE/NHRmtje56HLcgppwljqiLO5Kkh6lkXDOn4qqiZkSQap/u+EpROjV7M0EXMVwwVF59ZDevKCh+2nEEOup3gDJh25+2vs9JPZhv6/VqS6pw5NVXIrlhXWnm+AZ6zzZNiONT/lXSjpsxT35JvHvvRUGqXG6fHUTzeAY6lzg4vhP4qIw8d1UdbJnpthokwC6RMvHOfi6gYeUw8zoVzUTaM4LRp7UP0WYSYViWJb+bk5n99Ow5/ruZbSNvTJ4l5+3Q2GlBsjLdAftpXhBfrd2lMJF2znvkmWAFalD49d+Hd9Ia+6U8LsUuZdZjKcPO4GJ4pF1oFnrW3LE9h2AyxISMTAgMBAAECggEAXhjKMoJdIYDQDnanSVrMPingWuqKDD3VaGb3etc++Vw2D1N9U77osPGSRZ16uk71tIVfcBkXM5/OfuY7seuPU+1NEocBDIZrrCPTyMncgnXVgv5ZYRAeHUd+jAc6ptURRCeG7aPLRvSjx7dJKVS1I6oWNaB+2qQnuN+iAtTpfST5wCVke0y7s9tl3mZqp8EkHV2yXJqIYpKiUOYCuSMr99ARCLGmeRAD5w5cpTHeZoMzdLxoITFNvout/a4Is1NNGzIjMWQD9WOrhmE15cMukOW6jFKXkrQDI9cnUifkS/zmPFLrnHd4Zqo9SoGQIxEYZLQ1U12NbRtV1Ss9vZijOQKBgQD+l9uqUqT+fw11Mp9xaDgZzxt1rcd488qQtNZjbD0eoA2a1tn2g7YZeIA+Y9lPAXv99q7b6yFeeecJ/cOKOgD8wIjDigjPnL2wi9ZBVp1ZUXJSAQXaYZFgkiyOApoHdB62jdnuecP4G5dLs4+MSO/pciGk6/eoIhWrDJTdO3vH9wKBgQCAySBkNeK2SZ+2jtv259x+xmrR+7b2FSAal9wKclytoZHtUYlZ5hsURgkbBoZMpH5VJBoaveOOsocoKpOjn4LbYY9eTSCQN9yzF4JS5PXFoNmjQ6P2Ndcorogi2pPkyCj3nrwa+zc9zFzUzYOrlPyswEr+mNTtgfyBNRhhTOVOxQKBgQCoQVQ7TEMernkGa15UZLwu0mEjdKXPmc7Vs628J1x9UOms2zFRadp/GtQmZ3bGcASx4sXNMafr+ERopf0E7TCZ2eSI1kDcdIook0IWDFgRH3KeH27u1Gxvlis77xw8sNFbdIQCxxZsck+bCCBmZg2oCnWRuSEDTQNk9/up+hXkIQKBgG7QoZaY52ODJnKnqo5iJFDR2sikl2JX+y/my+gRT7338OEL7+vzHAnt2ZfvnVAFms8YKX4pNs1qwPHG8RMyBh9Pa1Xxd7ug1b8k03cQnIpZRew+H6+T1Hek9m9HNUr/EIFBjQqKb5Y1awuRa2MQ5/qd2+oHB/D2kJd9YGUZDZchAoGBAO4IaKpc/iiqNV8ZZuSVnZT4HMoyszJ0q86wmKcITgc5qhYbbgzXCChPWnKSLqBTeWrqmXKJkqhMT9TunDu0Xvu2OPg0xzOWl7GjGBKNRXElqgEHilGK99no/5cK/Vww1nC9x9hwpJDgSdTkvi6mTv89M2SbErJdBydXOHswGEKZ"
	const alipayAlipayPublicKey = "MIIBIjANBgkqhkiG9...8CFHkHpYsHcwIDAQAB"

	client := defaultAlipayClient.NewDefaultAlipayClient(
		alipayGatewayUrl,
		alipayClientId,
		alipayMerchantPrivateKey,
		alipayAlipayPublicKey, false)

	doPay(client)
	//payQuery("1e6d724d-da95-407b-9802-6f2217c301d6", client)
	//refund("202408151940108001001886C0209996792", client)
	//queryRefund("ad53716a-81-4c4c-b151-216c5225e", client)
	//cancel("ad53716a-81-4c4c-b51-20916c5225e", client)
	//consult(client)
	//createSession(client)
}

func doPay(body *defaultAlipayClient.DefaultAlipayClient) {
	payRequest, request := pay.NewAlipayPayRequest()

	request.PaymentRequestId = uuid.NewString()
	order := &model.Order{}
	order.OrderDescription = "example order"
	order.ReferenceOrderId = uuid.NewString()
	order.OrderAmount = model.NewAmount("100", "HKD")
	merchant := &model.Merchant{}
	merchant.ReferenceMerchantId = "1238rye8yr8erwer"
	merchant.MerchantMCC = "7011"
	merchant.MerchantName = "example merchant"
	merchant.Store = &model.Store{StoreMCC: "7011", ReferenceStoreId: "289285674", StoreName: "store 1111"}
	order.Merchant = merchant
	order.Env = &model.Env{OsType: model.ANDROID, TerminalType: model.WEB}
	request.Order = order

	request.PaymentAmount = model.NewAmount("100", "HKD")

	request.PaymentNotifyUrl = "https://www.yourNotifyUrl.com"
	request.PaymentRedirectUrl = "https://www.yourRedirectUrl.com"

	request.PaymentMethod = &model.PaymentMethod{PaymentMethodType: model.ALIPAY_HK}

	request.ProductCode = model.CASHIER_PAYMENT

	execute, err := body.Execute(payRequest)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayPayResponse)
	fmt.Println("result: ", response.AlipayResponse.Result)
	fmt.Println("response: ", response)
}

func payQuery(paymentRequestId string, body *defaultAlipayClient.DefaultAlipayClient) {
	queryRequest := pay.AlipayPayQueryRequest{}
	queryRequest.PaymentRequestId = paymentRequestId
	request := queryRequest.NewRequest()
	//1. 这里接收结果
	execute, err := body.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayPayQueryResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)
}

func refund(paymentId string, client *defaultAlipayClient.DefaultAlipayClient) {
	refundRequest := pay.AlipayRefundRequest{}
	refundRequest.RefundRequestId = uuid.NewString()
	refundRequest.PaymentId = paymentId
	refundRequest.RefundAmount = model.NewAmount("100", "HKD")
	refundRequest.RefundReason = "example refund"
	request := refundRequest.NewRequest()
	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayRefundResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)

}

func queryRefund(refundRequestId string, client *defaultAlipayClient.DefaultAlipayClient) {
	queryRefundRequest := pay.AlipayInquiryRefundRequest{}
	queryRefundRequest.RefundRequestId = refundRequestId
	request := queryRefundRequest.NewRequest()
	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayInquiryRefundResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)
}

func consult(client *defaultAlipayClient.DefaultAlipayClient) {
	consultRequest := &pay.AlipayPayConsultRequest{}
	request := consultRequest.NewRequest()
	consultRequest.SettlementStrategy = &model.SettlementStrategy{
		SettlementCurrency: "USD",
	}
	consultRequest.ProductCode = model.CASHIER_PAYMENT
	consultRequest.UserRegion = "SG"
	consultRequest.AllowedPaymentMethodRegions = []string{"SG", "HK", "CN"}
	consultRequest.Env = &model.Env{
		OsType:       model.IOS,
		TerminalType: model.APP,
	}
	consultRequest.PaymentAmount = model.NewAmount("1000", "USD")
	consultRequest.PaymentFactor = &model.PaymentFactor{
		PresentmentMode: model.BUNDLE,
	}

	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayPayConsultResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)

}

func cancel(paymentRequestId string, client *defaultAlipayClient.DefaultAlipayClient) {
	request, cancelRequest := pay.NewAlipayPayCancelRequest()
	cancelRequest.PaymentRequestId = paymentRequestId
	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayPayCancelResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)

}

func createSession(client *defaultAlipayClient.DefaultAlipayClient) {
	request, createSessionRequest := pay.NewAlipayPaymentSessionRequest()
	createSessionRequest.PaymentRequestId = uuid.NewString()
	createSessionRequest.Order = &model.Order{
		OrderDescription: "example order",
		ReferenceOrderId: "289473927358748",
		OrderAmount:      model.NewAmount("100", "HKD"),
		Buyer: &model.Buyer{
			ReferenceBuyerId: "111112132143434",
		},
	}
	createSessionRequest.PaymentAmount = model.NewAmount("100", "HKD")
	createSessionRequest.ProductCode = model.CASHIER_PAYMENT
	createSessionRequest.PaymentMethod = &model.PaymentMethod{
		PaymentMethodType: model.SHOPEEPAY_SG,
	}
	createSessionRequest.PaymentNotifyUrl = "https://www.yourNotifyUrl.com"
	createSessionRequest.PaymentRedirectUrl = "https://www.yourRedirectUrl.com"
	createSessionRequest.Env = &model.Env{
		OsType:       model.IOS,
		TerminalType: model.APP,
	}

	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responsePay.AlipayPaymentSessionResponse)
	fmt.Println(response.AlipayResponse.Result.ResultCode)
	fmt.Println(response)
}
