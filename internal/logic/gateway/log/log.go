package log

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func DoSaveChannelLog(ctx context.Context, request string, url string, response string, memo string, requestId string, gatewayId string) {
	log := &entity.GatewayHttpLog{
		Url:       url,
		Request:   request,
		Response:  response,
		RequestId: requestId,
		Mamo:      memo,
		GatewayId: gatewayId,
		CreateAt:  gtime.Now().Timestamp(),
	}
	_, err := dao.GatewayHttpLog.Ctx(ctx).Data(log).OmitNil().Insert(log)
	if err != nil {
		g.Log().Errorf(ctx, `record insert failure %s`, err)
	}
}

func ConvertToStringIgnoreErr(m interface{}) (response string) {
	jsonData, err := gjson.Marshal(m)
	if err != nil {
		fmt.Printf("RemoteChannel ConvertToStringIgnoreErr error %s", err)
	}
	response = string(jsonData)
	return
}
