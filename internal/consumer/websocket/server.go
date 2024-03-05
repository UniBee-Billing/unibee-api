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

func MerchantWebSocketMessageEntry(r *ghttp.Request) {
	merchantApiKey := r.Get("merchantApiKey").String()
	if len(merchantApiKey) == 0 {
		glog.Error(r.Context(), gerror.New("merchantApiKey invalid"))
		r.Exit()
	}
	merchant := query.GetMerchantByApiKey(r.Context(), merchantApiKey)
	if merchant == nil {
		glog.Error(r.Context(), gerror.New("merchantApiKey invalid"))
		r.Exit()
	}
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(r.Context(), err)
		r.Exit()
	}
	g.Log().Infof(r.Context(), "MerchantWebSocketMessageEntry:%d", merchant.Id)
	for {
		// todo mark use broadcast redis message is better
		var one *entity.MerchantWebhookMessage
		err := dao.MerchantWebhookMessage.Ctx(r.Context()).
			Where(dao.MerchantWebhookMessage.Columns().MerchantId, merchant.Id).
			Where(dao.MerchantWebhookMessage.Columns().WebsocketStatus, 10).
			WhereNotNull(dao.MerchantWebhookMessage.Columns().Data).
			OrderAsc(dao.MerchantWebhookMessage.Columns().CreateTime).
			Scan(&one)
		utility.AssertError(err, "merchant query MerchantWebhookMessage error")
		if one != nil {
			g.Log().Infof(r.Context(), "MerchantWebSocketMessage Start WriteMessage:%d", one.Id)
			if err = ws.WriteMessage(websocket.BinaryMessage, []byte(one.Data)); err != nil {
				g.Log().Errorf(r.Context(), "MerchantWebSocketMessage WriteMessage err:%s", err.Error())
				r.Exit()
				break
			}
			_, err = dao.MerchantWebhookMessage.Ctx(r.Context()).Data(g.Map{
				dao.MerchantWebhookMessage.Columns().WebsocketStatus: 20,
				dao.MerchantWebhookMessage.Columns().GmtModify:       gtime.Now(),
			}).Where(dao.MerchantWebhookMessage.Columns().Id, one.Id).OmitNil().Update()
			utility.AssertError(err, "merchant update websocket status error")
			g.Log().Infof(r.Context(), "MerchantWebSocketMessage Finish WriteMessage:%d", one.Id)
		}
		time.Sleep(100)
	}
	g.Log().Infof(r.Context(), "MerchantWebSocketMessageEntry Exit:%d", merchant.Id)
}
