package email

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/platform"
	"unibee/internal/query"
)

func StandAloneInit(ctx context.Context) {
	list, err := query.GetEmailDefaultTemplateList(ctx)
	if err != nil {
		glog.Errorf(ctx, "StandAloneInit InitDefaultEmailTemplate error:%s", err.Error())
	}
	if err == nil && len(list) == 0 {
		glog.Infof(ctx, "StandAloneInit InitDefaultEmailTemplate start")
		InitDefaultEmailTemplate(ctx)
		glog.Infof(ctx, "StandAloneInit InitDefaultEmailTemplate end")
	}
}

func InitDefaultEmailTemplate(ctx context.Context) {
	list := platform.FetchDefaultEmailTemplateFromPlatformApi()
	glog.Infof(ctx, "StandAloneInit InitDefaultEmailTemplate cloud template count:%d", len(list))
	for _, one := range list {
		_, err := dao.EmailDefaultTemplate.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			glog.Errorf(ctx, "StandAloneInit InitDefaultEmailTemplate error:%s", err.Error())
		} else {
			glog.Infof(ctx, "StandAloneInit InitDefaultEmailTemplate template:%s", one.TemplateName)
		}
	}
}
