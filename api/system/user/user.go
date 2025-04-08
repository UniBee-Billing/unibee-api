package user

import "github.com/gogf/gf/v2/frame/g"

type InternalWebhookSyncReq struct {
	g.Meta        `path:"/user_internal_webhook_sync" tags:"System-Admin" method:"post" summary:"Admin Sync User Internal Webhook (Analysis)" description:"Giving startId or startTime which data from, and endId or endTime which data to, unibee-api will create thread to sync data paging by count 100 to the giving end"`
	StartId       *string `json:"startId" dc:"The start Id of User to sync data" `
	StartTime     *int64  `json:"startTime" dc:"The start time to sync data, ignore if StartId provided" `
	EndId         *string `json:"endId" dc:"The end Id of User to sync data" `
	EndTime       *int64  `json:"endTime" dc:"The end time to sync data, ignore if EndId provided" `
	IsSynchronous bool    `json:"isSynchronous" dc:"Synchronous or not, default false" `
}

type InternalWebhookSyncRes struct {
	Total   int    `json:"total" dc:"The total of sync, only output when isSynchronous=true" `
	FirstId string `json:"firstId" dc:"The first userId of sync data, only output when isSynchronous=true" `
	LastId  string `json:"lastId" dc:"The last userId of sync data, only output when isSynchronous=true" `
}
