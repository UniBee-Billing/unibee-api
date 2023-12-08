package evonet

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/paychannel/ro"
	"go-oversea-pay/internal/logic/paychannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const ENDPOINT = "https://hkg-online-uat.everonet.com"

type Evonet struct{}

func (e Evonet) DoRemoteChannelPayment(ctx context.Context, createPayContext interface{}) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res ro.OutPayCaptureRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channel := util.GetOverseaPayChannel(ctx, uint64(pay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")
	urlPath := "/g2/v1/payment/mer/" + channel.ChannelAccountId + "/evo.e-commerce.capture" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreateMerchantOrderNo(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": pay.Currency,
			"value":    utility.ConvertFenToYuanMinUnitStr(pay.BuyerPayFee),
		},
		"webhook": fmt.Sprintf("%s/evonet/notify/webhooks/notifications?payId=%d", consts.GetNacosConfigInstance().HostPath, pay.Id),
	}
	data, err := sendEvonetRequest(ctx, "POST", urlPath, channel.ChannelKey, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay捕获失败 result is null")
	resultJson := responseJson.GetJson("result")
	captureJson := responseJson.GetJson("capture")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		captureJson != nil &&
		captureJson.Contains("evoTransInfo") &&
		captureJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay捕获失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := captureJson.Get("status").String()
	pspReference := captureJson.GetJson("evoTransInfo").Get("evoTransID").String()
	res = ro.OutPayCaptureRo{
		PspReference: pspReference,
		Status:       status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res ro.OutPayCancelRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channel := util.GetOverseaPayChannel(ctx, uint64(pay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")
	urlPath := "/g2/v1/payment/mer/" + channel.ChannelAccountId + "/evo.e-commerce.cancel" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreateMerchantOrderNo(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"webhook": fmt.Sprintf("%s/evonet/notify/webhooks/notifications?payId=%d", consts.GetNacosConfigInstance().HostPath, pay.Id),
	}
	data, err := sendEvonetRequest(ctx, "POST", urlPath, channel.ChannelKey, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay取消失败 result is null")
	resultJson := responseJson.GetJson("result")
	cancelJson := responseJson.GetJson("cancel")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		cancelJson != nil &&
		cancelJson.Contains("evoTransInfo") &&
		cancelJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay取消失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := cancelJson.Get("status").String()
	pspReference := cancelJson.GetJson("evoTransInfo").Get("evoTransID").String()
	res = ro.OutPayCancelRo{
		PspReference: pspReference,
		Status:       status,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res ro.OutPayRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channel := util.GetOverseaPayChannel(ctx, uint64(pay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")
	urlPath := "/g2/v1/payment/mer/" + channel.ChannelAccountId + "/evo.e-commerce.payment"
	param := map[string]interface{}{
		"merchantTransID": pay.MerchantOrderNo,
	}
	data, err := sendEvonetRequest(ctx, "GET", urlPath, channel.ChannelKey, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay支付查询失败 result is null")
	resultJson := responseJson.GetJson("result")
	payment := responseJson.GetJson("payment")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		payment != nil &&
		payment.Contains("status") &&
		payment.Contains("evoTransInfo") &&
		payment.GetJson("evoTransInfo").Contains("evoTransID") &&
		payment.GetJson("merchantTransInfo").Contains("merchantTransID"),
		fmt.Sprintf("Evonetpay支付查询失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := payment.Get("status").String()
	pspReference := payment.GetJson("evoTransInfo").Get("evoTransID").String()
	merchantPspReference := payment.GetJson("merchantTransInfo").Get("merchantTransID").String()
	utility.Assert(strings.Compare(merchantPspReference, pay.MerchantOrderNo) == 0, "merchantPspReference not match")
	res = ro.OutPayRo{
		PayFee:    pay.PaymentFee,
		PayStatus: consts.TO_BE_PAID,
	}
	if strings.Compare(status, "Failed") == 0 || strings.Compare(status, "Cancelled") == 0 {
		res.PayStatus = consts.PAY_FAILED
		res.Reason = "from_query:" + payment.Get("failureReason").String()
	} else if strings.Compare(status, "Captured") == 0 {
		res.PayStatus = consts.PAY_SUCCESS
		res.ChannelTradeNo = pspReference
		res.PayTime = gtime.Now()
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res ro.OutPayRefundRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channel := util.GetOverseaPayChannel(ctx, uint64(pay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")
	urlPath := "/g2/v1/payment/mer/" + channel.ChannelAccountId + "/evo.e-commerce.refund" + "?merchantTransID=" + pay.MerchantOrderNo
	param := map[string]interface{}{
		"merchantTransInfo": map[string]interface{}{
			"merchantTransID":   utility.CreateMerchantOrderNo(),
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": pay.Currency,
			"value":    utility.ConvertFenToYuanMinUnitStr(refund.RefundFee),
		},
		"webhook": fmt.Sprintf("%s/evonet/notify/webhooks/notifications?payId=%d", consts.GetNacosConfigInstance().HostPath, pay.Id),
	}
	data, err := sendEvonetRequest(ctx, "POST", urlPath, channel.ChannelKey, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay退款失败 result is null")
	resultJson := responseJson.GetJson("result")
	refundJson := responseJson.GetJson("refund")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		refundJson != nil &&
		refundJson.Contains("evoTransInfo") &&
		refundJson.GetJson("evoTransInfo").Contains("evoTransID"),
		fmt.Sprintf("Evonetpay取消失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	pspReference := refundJson.GetJson("evoTransInfo").Get("evoTransID").String()
	res = ro.OutPayRefundRo{
		ChannelRefundNo: pspReference,
		RefundStatus:    consts.REFUND_ING,
	}
	return res, nil
}

func (e Evonet) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res ro.OutPayRefundRo, err error) {
	utility.Assert(pay.ChannelId > 0, "支付渠道异常")
	channel := util.GetOverseaPayChannel(ctx, uint64(pay.ChannelId))
	utility.Assert(channel != nil, "支付渠道异常 channel not found")
	urlPath := "/g2/v1/payment/mer/" + channel.ChannelAccountId + "/evo.e-commerce.refund"
	param := map[string]interface{}{
		"merchantTransID": refund.OutRefundNo,
	}
	data, err := sendEvonetRequest(ctx, "GET", urlPath, channel.ChannelKey, param)
	utility.Assert(err == nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err == nil, fmt.Sprintf("json parse error %s", err))
	utility.Assert(responseJson.Contains("result"), "Evonetpay退款查询失败 result is null")
	resultJson := responseJson.GetJson("result")
	refundJson := responseJson.GetJson("refund")
	utility.Assert(resultJson.Contains("code") &&
		strings.Compare(resultJson.Get("code").String(), "S0000") == 0 &&
		refundJson != nil &&
		refundJson.Contains("status") &&
		refundJson.Contains("evoTransInfo") &&
		refundJson.GetJson("evoTransInfo").Contains("evoTransID") &&
		refundJson.GetJson("merchantTransInfo").Contains("merchantTransID"),
		fmt.Sprintf("Evonetpay退款查询失败:%s-%s", resultJson.Get("code").String(), resultJson.Get("message").String()))
	status := refundJson.Get("status").String()
	pspReference := refundJson.GetJson("evoTransInfo").Get("evoTransID").String()
	merchantPspReference := refundJson.GetJson("merchantTransInfo").Get("merchantTransID").String()
	utility.Assert(strings.Compare(merchantPspReference, refund.OutRefundNo) == 0, "merchantPspReference not match")
	res = ro.OutPayRefundRo{
		RefundFee:    refund.RefundFee,
		RefundStatus: consts.REFUND_ING,
	}
	if strings.Compare(status, "Failed") == 0 {
		res.RefundStatus = consts.REFUND_FAILED
		res.Reason = "from_query:" + refundJson.Get("failureReason").String()
	} else if strings.Compare(status, "Success") == 0 {
		res.RefundStatus = consts.REFUND_SUCCESS
		res.ChannelRefundNo = pspReference
		res.RefundTime = gtime.Now()
	}
	return res, nil
}

func sendEvonetRequest(ctx context.Context, method string, urlPath string, key string, param map[string]interface{}) (res []byte, err error) {
	utility.Assert(param != nil, "param is null")
	// 定义自定义的头部信息
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Infof(ctx, "\nEvonet_Start %s %s %s %s\n", method, urlPath, key, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Msgid":         msgId,
		"Datetime":      datetime,
		"Authorization": sign("POST", urlPath, msgId, datetime, key, body),
		"Signtype":      "SHA256",
	}
	response, err := sendRequest(ENDPOINT+urlPath, method, body, headers)
	g.Log().Infof(ctx, "\nEvonet_End %s %s response: %s error %s\n", method, urlPath, response, err)
	return response, nil
}

func sendRequest(url string, method string, data []byte, headers map[string]string) ([]byte, error) {
	// 创建一个字节数组读取器，用于将数据传递给请求体
	bodyReader := bytes.NewReader(data)

	// 创建一个POST请求
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置自定义头部信息
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// 关闭响应体
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func sign(method string, urlPath string, msgId string, dateTime string, key string, postJson []byte) (sign string) {
	var builder strings.Builder
	lineSeparator := lineSeparator()
	builder.WriteString(method)
	builder.WriteString(lineSeparator)
	builder.WriteString(urlPath)
	builder.WriteString(lineSeparator)
	builder.WriteString(dateTime)
	builder.WriteString(lineSeparator)
	builder.WriteString(key)
	builder.WriteString(lineSeparator)
	builder.WriteString(msgId)
	if postJson != nil {
		builder.WriteString(lineSeparator)
		builder.Write(postJson)
	}
	return sha256Encoding(builder.String())
}

func generateMsgId() (msgId string) {
	return fmt.Sprintf("%s%s%s", utility.JodaTimePrefix(), utility.GenerateRandomAlphanumeric(5), utility.CurrentTimeMillis())
}

func getCurrentDateTime() (datetime string) {
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}

func lineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func sha256Encoding(data string) (hash string) {
	hasher := sha256.New()

	// 添加数据到散列器
	hasher.Write([]byte(data))

	// 计算散列值
	hashValue := hasher.Sum(nil)

	// 将散列值转换为十六进制字符串
	return hex.EncodeToString(hashValue)
}
