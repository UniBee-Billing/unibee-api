package email

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
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
	list := FetchDefaultEmailTemplateFromCloudApi()
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

func FetchDefaultEmailTemplateFromCloudApi() []*entity.EmailDefaultTemplate {
	var list = make([]*entity.EmailDefaultTemplate, 0)
	response, err := utility.SendRequest("http://api.cloud.unibee.top/cloud/email/default_template_list", "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("emailTemplateList") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("emailTemplateList").String()), &list)
	}
	return list
}
