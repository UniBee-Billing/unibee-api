package websocket

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"time"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantWebhookMessageVo struct {
	Id           uint64      `json:"id"              description:"id"`              // id
	MerchantId   uint64      `json:"merchantId"      description:"merchantId"`      // merchantId
	WebhookEvent string      `json:"webhookEvent"    description:"webhook_event"`   // webhook_event
	Data         interface{} `json:"data"            description:"data(json)"`      // data(json)
	CreateTime   int64       `json:"createTime"      description:"create utc time"` // create utc time
}

func MerchantWebSocketMessageEntry(r *ghttp.Request) {
	merchantApiKey := r.Get("merchantApiKey").String()
	if len(merchantApiKey) == 0 {
		glog.Error(r.Context(), gerror.New("MerchantWebSocketMessage merchantApiKey invalid"))
		r.Exit()
	}
	merchant := query.GetMerchantByApiKey(r.Context(), merchantApiKey)
	if merchant == nil {
		glog.Error(r.Context(), gerror.New("MerchantWebSocketMessage merchantApiKey invalid"))
		r.Exit()
	}
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(r.Context(), err)
		r.Exit()
	}
	g.Log().Infof(r.Context(), "MerchantWebSocketMessage Entry:%d", merchant.Id)
	for {
		// todo mark use broadcast redis message is better
		var one *entity.MerchantWebhookMessage
		err := dao.MerchantWebhookMessage.Ctx(r.Context()).
			Where(dao.MerchantWebhookMessage.Columns().MerchantId, merchant.Id).
			//Where(dao.MerchantWebhookMessage.Columns().WebhookEvent, event.MERCHANT_WEBHOOK_TAG_USER_METRIC_UPDATED).
			Where(dao.MerchantWebhookMessage.Columns().WebsocketStatus, 10).
			WhereNotNull(dao.MerchantWebhookMessage.Columns().Data).
			OrderAsc(dao.MerchantWebhookMessage.Columns().CreateTime).
			Scan(&one)
		utility.AssertError(err, "MerchantWebSocketMessage merchant query MerchantWebSocketMessage error")
		if one != nil {
			g.Log().Infof(r.Context(), "MerchantWebSocketMessage Start WriteMessage:%d", one.Id)
			if err = ws.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				g.Log().Errorf(r.Context(), "MerchantWebSocketMessage WritePingMessage err:%s", err.Error())
				break
			}
			if err = ws.WriteMessage(websocket.BinaryMessage, []byte(utility.FormatToGJson(&MerchantWebhookMessageVo{
				Id:           one.Id,
				MerchantId:   one.MerchantId,
				WebhookEvent: one.WebhookEvent,
				Data:         one.Data,
				CreateTime:   one.CreateTime,
			}).String())); err != nil {
				g.Log().Errorf(r.Context(), "MerchantWebSocketMessage WriteMessage err:%s", err.Error())
				break
			}
			_, err = dao.MerchantWebhookMessage.Ctx(r.Context()).Data(g.Map{
				dao.MerchantWebhookMessage.Columns().WebsocketStatus: 20,
				dao.MerchantWebhookMessage.Columns().GmtModify:       gtime.Now(),
			}).Where(dao.MerchantWebhookMessage.Columns().Id, one.Id).OmitNil().Update()
			utility.AssertError(err, "MerchantWebSocketMessage merchant update websocket status error")
			g.Log().Infof(r.Context(), "MerchantWebSocketMessage Finish WriteMessage:%d", one.Id)
		}
		time.Sleep(100)
	}
	g.Log().Infof(r.Context(), "MerchantWebSocketMessage Exit:%d", merchant.Id)

}
