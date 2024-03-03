package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consumer/webhook/event"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

const SplitSep = ","

type MerchantWebhookEndpointVo struct {
	Id            uint64   `json:"id"            description:"id"`                       // id
	MerchantId    uint64   `json:"merchantId"    description:"webhook url"`              // webhook url
	WebhookUrl    string   `json:"webhookUrl"    description:"webhook url"`              // webhook url
	WebhookEvents []string `json:"webhookEvents" description:"webhook_events,split dot"` // webhook_events,split dot
	UpdateTime    int64    `json:"gmtModify"     description:"update time"`              // update time
	CreateTime    int64    `json:"createTime"    description:"create utc time"`          // create utc time
}

func MerchantWebhookEndpointList(ctx context.Context, merchantId uint64) []*MerchantWebhookEndpointVo {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var list = make([]*MerchantWebhookEndpointVo, 0)
	if merchantId > 0 {
		var entities []*entity.MerchantWebhook
		err := dao.MerchantWebhook.Ctx(ctx).
			Where(dao.MerchantWebhook.Columns().MerchantId, merchantId).
			Where(dao.MerchantWebhook.Columns().IsDeleted, 0).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, &MerchantWebhookEndpointVo{
					Id:            one.Id,
					MerchantId:    one.MerchantId,
					WebhookUrl:    one.WebhookUrl,
					WebhookEvents: strings.Split(one.WebhookEvents, SplitSep),
					UpdateTime:    one.GmtModify.Timestamp(),
					CreateTime:    one.CreateTime,
				})
			}
		}
	}
	return list
}

type EndpointLogListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	EndpointId int64  `json:"endpointId" dc:"EndpointId" v:"required"`
	Page       int    `json:"page" dc:"Page, Start WIth 0" `
	Count      int    `json:"count" dc:"Count Of Page" `
}

func MerchantWebhookEndpointLogList(ctx context.Context, req *EndpointLogListInternalReq) []*entity.MerchantWebhookLog {
	var mainList []*entity.MerchantWebhookLog
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	utility.Assert(req.MerchantId > 0, "merchantId not found")
	utility.Assert(req.EndpointId > 0, "endpointId not found")
	var sortKey = "create_time desc"
	_ = dao.MerchantWebhookLog.Ctx(ctx).
		Where(dao.MerchantWebhookLog.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantWebhookLog.Columns().EndpointId, req.EndpointId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().Scan(&mainList)
	return mainList
}

func NewMerchantWebhookEndpoint(ctx context.Context, merchantId uint64, url string, events []string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(url) > 0, "url is nil")
	utility.Assert(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"), "Invalid Url")
	// events valid check
	for _, e := range events {
		utility.Assert(event.WebhookEventInListeningEvents(event.MerchantWebhookEvent(e)), fmt.Sprintf("Event:%s Not In Event List", e))
	}
	one := query.GetMerchantWebhookByUrl(ctx, url)
	utility.Assert(one == nil, "endpoint already exist")
	one = &entity.MerchantWebhook{
		MerchantId:    merchantId,
		WebhookUrl:    url,
		WebhookEvents: strings.Join(events, SplitSep),
		CreateTime:    gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantWebhook.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "NewMerchantWebhookEndpoint Insert err:%s", err.Error())
		return gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	return nil
}

func UpdateMerchantWebhookEndpoint(ctx context.Context, merchantId uint64, endpointId int64, url string, events []string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(endpointId > 0, "invalid endpointId")
	utility.Assert(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"), "Invalid Url")
	// events valid check
	for _, e := range events {
		utility.Assert(event.WebhookEventInListeningEvents(event.MerchantWebhookEvent(e)), fmt.Sprintf("Event:%s Not In Event List", e))
	}
	one := query.GetMerchantWebhook(ctx, endpointId)
	utility.Assert(one != nil, "endpoint not found")
	_, err := dao.MerchantWebhook.Ctx(ctx).Data(g.Map{
		dao.MerchantWebhook.Columns().MerchantId:    merchantId,
		dao.MerchantWebhook.Columns().WebhookUrl:    url,
		dao.MerchantWebhook.Columns().WebhookEvents: strings.Join(events, SplitSep),
		dao.MerchantWebhook.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantWebhook.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "UpdateMerchantWebhookEndpoint Update err:%s", err.Error())
		return gerror.NewCode(gcode.New(500, "server error", nil))
	}

	return nil
}

func DeleteMerchantWebhookEndpoint(ctx context.Context, merchantId uint64, endpointId int64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(endpointId > 0, "invalid endpointId")
	one := query.GetMerchantWebhook(ctx, endpointId)
	utility.Assert(one != nil, "endpoint not found")
	_, err := dao.MerchantWebhook.Ctx(ctx).Data(g.Map{
		dao.MerchantWebhook.Columns().IsDeleted: 1,
		dao.MerchantWebhook.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantWebhook.Columns().Id, one.Id).OmitNil().Update()
	return err
}
