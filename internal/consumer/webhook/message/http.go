package message

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"time"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func ResentWebhook(ctx context.Context, logId uint64) bool {
	var one *entity.MerchantWebhookLog
	err := dao.MerchantWebhookLog.Ctx(ctx).Where(dao.MerchantWebhookLog.Columns().Id, logId).Scan(&one)
	if err != nil {
		g.Log().Errorf(ctx, "ResentWebhook error:", err.Error())
		return false
	}
	utility.Assert(one != nil, "webhook log not found")
	merchant := query.GetMerchantById(ctx, one.MerchantId)
	if merchant == nil {
		g.Log().Errorf(ctx, "Webhook_Resend %s %s merchant not found\n", "POST", one.WebhookUrl)
		return false
	}
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	g.Log().Debugf(ctx, "Webhook_Start %s %s %s\n", "POST", one.WebhookUrl, one.Body)
	headers := map[string]string{
		"Content-Gateway": "application/json",
		"Msg-id":          msgId,
		"Datetime":        datetime,
		"Authorization":   fmt.Sprintf("Bearer %s", merchant.ApiKey),
	}
	body := []byte(one.Body)
	res, err := utility.SendRequest(one.WebhookUrl, "POST", body, headers)
	var response = string(res)
	if err != nil {
		response = utility.MarshalToJsonString(err)
		g.Log().Debugf(ctx, "Webhook_End %s %s response: %s error %s\n", "POST", one.WebhookUrl, response, err.Error())
	} else {
		g.Log().Debugf(ctx, "Webhook_End %s %s response: %s \n", "POST", one.WebhookUrl, response)
	}
	return true
}

func SendWebhookRequest(ctx context.Context, webhookMessage *WebhookMessage, reconsumeTimes int) bool {
	utility.Assert(webhookMessage.Data != nil, "param is nil")
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	err := webhookMessage.Data.Set("eventType", webhookMessage.Event)
	if err != nil {
		g.Log().Errorf(ctx, "Webhook_Send %s %s error %s\n", "POST", webhookMessage.Url, err.Error())
		return false
	}
	merchant := query.GetMerchantById(ctx, webhookMessage.MerchantId)
	if merchant == nil {
		g.Log().Errorf(ctx, "Webhook_Send %s %s merchant not found\n", "POST", webhookMessage.Url)
		return false
	}
	jsonString, err := webhookMessage.Data.ToJsonString()
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, webhookMessage.Data))
	g.Log().Debugf(ctx, "Webhook_Start %s %s %s\n", "POST", webhookMessage.Url, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Gateway": "application/json",
		"Msg-id":          msgId,
		"Datetime":        datetime,
		"Authorization":   fmt.Sprintf("Bearer %s", merchant.ApiKey),
	}
	res, err := utility.SendRequest(webhookMessage.Url, "POST", body, headers)
	var response = string(res)
	var responseMessage = "not success"
	if strings.Compare(response, "success") == 0 {
		responseMessage = response
	}
	if err != nil {
		response = utility.MarshalToJsonString(err)
		g.Log().Debugf(ctx, "Webhook_End %s %s response: %s error\n", "POST", webhookMessage.Url, responseMessage)
	} else {
		g.Log().Debugf(ctx, "Webhook_End %s %s response: %s \n", "POST", webhookMessage.Url, responseMessage)
	}
	one := &entity.MerchantWebhookLog{
		MerchantId:     webhookMessage.MerchantId,
		EndpointId:     int64(webhookMessage.EndpointId),
		WebhookUrl:     webhookMessage.Url,
		WebhookEvent:   string(webhookMessage.Event),
		RequestId:      msgId,
		Body:           jsonString,
		ReconsumeCount: reconsumeTimes,
		Response:       response,
		Mamo:           utility.MarshalToJsonString(err),
		CreateTime:     gtime.Now().Timestamp(),
	}
	_, saveErr := dao.MerchantWebhookLog.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if saveErr != nil {
		g.Log().Errorf(ctx, "Webhook_SaveLog error %s\n", saveErr.Error())
	}
	response = strings.Trim(response, " ")
	return err == nil && strings.Compare(response, "success") == 0
}

func generateMsgId() (msgId string) {
	return fmt.Sprintf("%s%s%d", utility.JodaTimePrefix(), utility.GenerateRandomAlphanumeric(5), utility.CurrentTimeMillis())
}

func getCurrentDateTime() (datetime string) {
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}
