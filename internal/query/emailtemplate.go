package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetEmailTemplateByTemplateName(ctx context.Context, templateName string) (one *entity.EmailTemplate) {
	if len(templateName) == 0 {
		return nil
	}
	err := dao.EmailTemplate.Ctx(ctx).Where(entity.EmailTemplate{TemplateName: templateName}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetMerchantEmailTemplateByTemplateName(ctx context.Context, merchantId int64, templateName string) (one *entity.EmailTemplate) {
	if len(templateName) == 0 || merchantId <= 0 {
		return nil
	}
	err := dao.EmailTemplate.Ctx(ctx).
		Where(entity.EmailTemplate{TemplateName: templateName}).
		Where(entity.EmailTemplate{MerchantId: merchantId}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
