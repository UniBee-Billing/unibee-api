package evonet

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
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

func (e Evonet) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

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
			"merchantTransID":   "",
			"merchantTransTime": getCurrentDateTime(),
		},
		"transAmount": map[string]interface{}{
			"currency": pay.Currency,
			"value":    pay.BuyerPayFee,
		},
		"webhook": fmt.Sprintf("/evonet/notify/webhooks/notifications?payId=%d", pay.Id),
	}
	data, err := sendEvonetRequest(ctx, "POST", urlPath, channel.ChannelKey, param)
	utility.Assert(err != nil, fmt.Sprintf("call evonet error %s", err))
	responseJson, err := gjson.LoadJson(string(data))
	utility.Assert(err != nil, fmt.Sprintf("json parse error %s", err))
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
	res = ro.OutPayCaptureRo{}
	res.PspReference = pspReference
	res.Status = status
	return res, nil
}

func (e Evonet) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (e Evonet) DoRemoteChannelStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func sendEvonetRequest(ctx context.Context, method string, urlPath string, key string, param map[string]interface{}) (res []byte, err error) {
	g.Log().Infof(ctx, "Evonet start %s %s %s %s", method, urlPath, key, param)
	utility.Assert(param != nil, "param is null")
	// 定义自定义的头部信息
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	jsonData, err := gjson.Marshal(param)
	utility.Assert(err != nil, fmt.Sprintf("json format error %s", err))
	body := []byte(string(jsonData))
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Msgid":         msgId,
		"Datetime":      datetime,
		"Authorization": sign("POST", urlPath, msgId, datetime, key, body),
		"Signtype":      "SHA256",
	}
	response, err := sendRequest(ENDPOINT+urlPath, method, body, headers)
	g.Log().Infof(ctx, "Evonet end %s %s response: %s error %s", method, urlPath, response, err)
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
	return fmt.Sprintf("%s%s%s", utility.JodaTimePrefix(), utility.GenerateRandomString(5), utility.CurrentTimeMillis())
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
