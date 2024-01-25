package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetEmailTemplateByTemplateName(ctx context.Context, templateName string) (one *entity.EmailTemplate) {
	err := dao.EmailTemplate.Ctx(ctx).Where(entity.EmailTemplate{TemplateName: templateName}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
