package query

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetEmailDefaultTemplateByTemplateName(ctx context.Context, templateName string) *bean.MerchantEmailTemplateSimplify {
	if len(templateName) == 0 {
		return nil
	}
	var one *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	if err != nil {
		one = nil
	}
	if one != nil {
		return &bean.MerchantEmailTemplateSimplify{
			Id:                  one.Id,
			MerchantId:          0,
			TemplateName:        one.TemplateName,
			TemplateDescription: one.TemplateDescription,
			TemplateTitle:       one.TemplateTitle,
			TemplateContent:     one.TemplateContent,
			TemplateAttachName:  one.TemplateAttachName,
			CreateTime:          one.CreateTime,
			UpdateTime:          one.GmtModify.Timestamp(),
			Status:              "Active",
		}
	} else {
		return nil
	}
}

func GetMerchantEmailTemplateByTemplateName(ctx context.Context, merchantId uint64, templateName string) *bean.MerchantEmailTemplateSimplify {
	if len(templateName) == 0 || merchantId <= 0 {
		return nil
	}
	var one *entity.MerchantEmailTemplate
	err := dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	if one == nil {
		// No Setting, Get Default
		return GetEmailDefaultTemplateByTemplateName(ctx, templateName)
	}
	return bean.SimplifyMerchantEmailTemplate(one)
}

func GetEmailDefaultTemplateList(ctx context.Context) (mainlist []*bean.EmailDefaultTemplateSimplify, err error) {
	list := make([]*entity.EmailDefaultTemplate, 0)
	err = dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	var templates = make([]*bean.EmailDefaultTemplateSimplify, 0)
	for _, one := range list {
		templates = append(templates, bean.SimplifyEmailDefaultTemplate(one))
	}
	return templates, nil
}
