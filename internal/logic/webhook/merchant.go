package webhook

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func SetupMerchantWebhook(ctx context.Context, merchantId int64, url string, events []string) error {
	utility.Assert(len(url) > 0, "url is nil")
	utility.Assert(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"), "Invalid Url")
	// events check todo mark
	one := query.GetMerchantWebhookByUrl(ctx, url)
	if one == nil {
		one = &entity.MerchantWebhook{
			MerchantId:    merchantId,
			WebhookUrl:    url,
			WebhookEvents: strings.Join(events, ","),
			CreateTime:    gtime.Now().Timestamp(),
		}
		result, err := dao.MerchantWebhook.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			g.Log().Errorf(ctx, "SetupMerchantWebhook Insert err:%s", err.Error())
			return gerror.NewCode(gcode.New(500, "server error", nil))
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(id)
	} else {
		_, err := dao.MerchantWebhook.Ctx(ctx).Data(g.Map{
			dao.MerchantWebhook.Columns().MerchantId:    merchantId,
			dao.MerchantWebhook.Columns().WebhookUrl:    url,
			dao.MerchantWebhook.Columns().WebhookEvents: strings.Join(events, ","),
			dao.MerchantWebhook.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.MerchantWebhook.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "SetupMerchantWebhook Update err:%s", err.Error())
			return gerror.NewCode(gcode.New(500, "server error", nil))
		}
	}
	return nil
}
