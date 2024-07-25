package bean

import entity "unibee/internal/model/entity/default"

type MerchantWebhookEndpoint struct {
	Id            uint64   `json:"id"            description:"id"`                       // id
	MerchantId    uint64   `json:"merchantId"    description:"webhook url"`              // webhook url
	WebhookUrl    string   `json:"webhookUrl"    description:"webhook url"`              // webhook url
	WebhookEvents []string `json:"webhookEvents" description:"webhook_events,split dot"` // webhook_events,split dot
	UpdateTime    int64    `json:"gmtModify"     description:"update time"`              // update time
	CreateTime    int64    `json:"createTime"    description:"create utc time"`          // create utc time
}

type MerchantWebhookLog struct {
	Id             uint64 `json:"id"             description:"id"`              // id
	MerchantId     uint64 `json:"merchantId"     description:"webhook url"`     // webhook url
	EndpointId     int64  `json:"endpointId"     description:""`                //
	ReconsumeCount int    `json:"reconsumeCount" description:""`                //
	WebhookUrl     string `json:"webhookUrl"     description:"webhook url"`     // webhook url
	WebhookEvent   string `json:"webhookEvent"   description:"webhook_event"`   // webhook_event
	RequestId      string `json:"requestId"      description:"request_id"`      // request_id
	Body           string `json:"body"           description:"body(json)"`      // body(json)
	Response       string `json:"response"       description:"response"`        // response
	Mamo           string `json:"mamo"           description:"mamo"`            // mamo
	CreateTime     int64  `json:"createTime"     description:"create utc time"` // create utc time
}

func SimplifyMerchantWebhookLog(one *entity.MerchantWebhookLog) *MerchantWebhookLog {
	if one == nil {
		return nil
	}
	return &MerchantWebhookLog{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		EndpointId:     one.EndpointId,
		ReconsumeCount: one.ReconsumeCount,
		WebhookUrl:     one.WebhookUrl,
		WebhookEvent:   one.WebhookEvent,
		RequestId:      one.RequestId,
		Body:           one.Body,
		Response:       one.Response,
		Mamo:           one.Mamo,
		CreateTime:     one.CreateTime,
	}
}
