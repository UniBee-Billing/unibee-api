package log

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"strconv"
)

func SaveChannelHttpLog(url string, request interface{}, response interface{}, err interface{}, memo string, requestId interface{}, channel *entity.OverseaPayChannel) {
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
		httpLog := &entity.ChannelHttpLog{
			Url:       url,
			Request:   utility.FormatToJsonString(request),
			Response:  utility.FormatToJsonString(utility.CheckReturn(err != nil, err, response)),
			RequestId: utility.FormatToJsonString(requestId),
			Mamo:      memo,
			ChannelId: strconv.FormatUint(channel.Id, 10),
		}
		_, _ = dao.ChannelHttpLog.Ctx(context.Background()).Data(httpLog).OmitNil().Insert(httpLog)
		g.Log().Infof(context.Background(), "SaveChannelHttpLog:%s", url)
	}()
}
