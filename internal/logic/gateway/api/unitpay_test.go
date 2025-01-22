package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
	"unibee/utility"
)

func TestForCreateNewUnitPay(t *testing.T) {
	unitpay := &UnitPay{}
	//_, _, _ = unitpay.GatewayTest(context.Background(), "423641", "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7") // indigo prod test key
	_, _, _ = unitpay.GatewayTest(context.Background(), "443597", "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7") // indigo unibee staging test key
	//_, _, _ = unitpay.GatewayTest(context.Background(), "443598", "b7bd68017bd0edb89258ef3b068f7771") // indigo x prod key
}

func TestForGetUnitPay(t *testing.T) {
	ctx := context.Background()
	//publickey := "423641-fae73"
	//secret := "0C9384AB748-A05AEE3555F-116BAA936E"
	secret := "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7"
	urlPath := "/api?method=getPayment"
	param := map[string]interface{}{}
	param["paymentId"] = 44359792787
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}

func TestForGetUnitPayForProd(t *testing.T) {
	ctx := context.Background()
	//publickey := "423641-fae73"
	//secret := "0C9384AB748-A05AEE3555F-116BAA936E"
	secret := "e7c14f504ada80c70f18f46f7ccf24c3"
	urlPath := "/api?method=getPayment"
	param := map[string]interface{}{}
	param["paymentId"] = 2183488712
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, true)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}

func TestForCreateRefundUnitPay(t *testing.T) {
	ctx := context.Background()
	secret := "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7"
	urlPath := "/api?method=refundPayment"
	param := map[string]interface{}{}
	param["paymentId"] = "383117770"
	param["sum"] = 1000
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}

func TestForRefundDetailUnitPay(t *testing.T) {
	ctx := context.Background()
	secret := "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7"
	urlPath := "/api?method=getPayment"
	param := map[string]interface{}{}
	param["paymentId"] = "12358132134"
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}

func TestForGetSubscriptionUnitPay(t *testing.T) {
	ctx := context.Background()
	secret := "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7"
	urlPath := "/api?method=getSubscription"
	param := map[string]interface{}{}
	param["subscriptionId"] = "1"
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}

func TestForGetSubscriptionListUnitPay(t *testing.T) {
	ctx := context.Background()
	secret := "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7"
	urlPath := "/api?method=listSubscriptions"
	param := map[string]interface{}{}
	param["projectId"] = "443597"
	responseJson, err := SendUnitPayPaymentRequest(ctx, secret, "GET", urlPath, param, false)
	if err != nil {
		t.Error(err)
	}
	g.Log().Debugf(ctx, "responseJson: %s", utility.MarshalToJsonString(responseJson))
}
