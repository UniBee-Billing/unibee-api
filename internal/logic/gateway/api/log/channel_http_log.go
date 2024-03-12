package log

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func SaveChannelHttpLog(url string, request interface{}, response interface{}, err interface{}, memo string, requestId interface{}, gateway *entity.MerchantGateway) {
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				var panicError error
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					panicError = v
				} else {
					panicError = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "SaveChannelHttpLog exception panic error:%s\n", panicError.Error())
				return
			}
		}()
		httpLog := &entity.GatewayHttpLog{
			Url:        url,
			Request:    utility.FormatToJsonString(request),
			Response:   utility.MarshalToJsonString(utility.CheckReturn(err != nil, err, response)),
			RequestId:  utility.FormatToJsonString(requestId),
			Mamo:       memo,
			GatewayId:  strconv.FormatUint(gateway.Id, 10),
			CreateTime: gtime.Now().Timestamp(),
		}
		_, _ = dao.GatewayHttpLog.Ctx(context.Background()).Data(httpLog).OmitNil().Insert(httpLog)
		//g.Log().Infof(context.Background(), "SaveChannelHttpLog:%s", url)
	}()
}
