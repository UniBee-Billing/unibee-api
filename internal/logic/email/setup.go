package email

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func StandAloneInit(ctx context.Context) {
	list := query.GetEmailDefaultTemplateList(ctx)
	if len(list) == 0 {
		InitDefaultEmailTemplate(ctx)
	}
}

func InitDefaultEmailTemplate(ctx context.Context) {
	list := FetchDefaultEmailTemplateFromCloudApi()
	for _, one := range list {
		_, err := dao.EmailDefaultTemplate.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			glog.Errorf(ctx, "InitDefaultEmailTemplate error:%s", err.Error())
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
