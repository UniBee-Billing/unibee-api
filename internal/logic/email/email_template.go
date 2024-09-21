package email

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetMerchantEmailTemplateList(ctx context.Context, merchantId uint64) ([]*bean.MerchantEmailTemplate, int) {
	var list = make([]*bean.MerchantEmailTemplate, 0)
	if merchantId > 0 {
		var defaultTemplateList []*entity.EmailDefaultTemplate
		err := dao.EmailDefaultTemplate.Ctx(ctx).
			Scan(&defaultTemplateList)
		if err == nil && len(defaultTemplateList) > 0 {
			for _, emailTemplate := range defaultTemplateList {
				var merchantEmailTemplate *entity.MerchantEmailTemplate
				err = dao.MerchantEmailTemplate.Ctx(ctx).
					Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
					Where(dao.MerchantEmailTemplate.Columns().TemplateName, emailTemplate.TemplateName).
					Scan(&merchantEmailTemplate)
				vo := &bean.MerchantEmailTemplate{
					Id:                  emailTemplate.Id,
					MerchantId:          0,
					TemplateName:        emailTemplate.TemplateName,
					TemplateDescription: emailTemplate.TemplateDescription,
					TemplateTitle:       emailTemplate.TemplateTitle,
					TemplateContent:     emailTemplate.TemplateContent,
					TemplateAttachName:  "", //pdf not customised here
					CreateTime:          emailTemplate.CreateTime,
					UpdateTime:          emailTemplate.GmtModify.Timestamp(),
					Status:              "Active", // default template status should be active
				}
				if err == nil && merchantEmailTemplate != nil {
					if merchantEmailTemplate.Status == 0 {
						vo.Status = "Active"
					} else {
						vo.Status = "InActive"
					}
					vo.TemplateDescription = merchantEmailTemplate.TemplateDescription
					vo.TemplateTitle = merchantEmailTemplate.TemplateTitle
					vo.TemplateContent = merchantEmailTemplate.TemplateContent
					vo.CreateTime = merchantEmailTemplate.CreateTime
					vo.UpdateTime = merchantEmailTemplate.GmtModify.Timestamp()
					vo.MerchantId = merchantEmailTemplate.MerchantId
				}
				list = append(list, vo)
			}
		}
	}
	return list, len(list)
}

func UpdateMerchantEmailTemplate(ctx context.Context, merchantId uint64, templateName string, templateTitle string, templateContent string) error {
	utility.Assert(merchantId > 0, "Invalid MerchantId")
	utility.Assert(len(templateName) > 0, "Invalid TemplateName")
	utility.Assert(len(templateTitle) > 0, "Invalid TemplateTitle")
	utility.Assert(len(templateContent) > 0, "Invalid TemplateContent")
	var defaultTemplate *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&defaultTemplate)
	utility.AssertError(err, "Server Error")
	utility.Assert(defaultTemplate != nil, "Default Template Not Found")
	var one *entity.MerchantEmailTemplate
	err = dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	utility.AssertError(err, "Server Error")
	if one == nil {
		//insert
		one = &entity.MerchantEmailTemplate{
			MerchantId:          merchantId,
			TemplateName:        defaultTemplate.TemplateName,
			TemplateDescription: defaultTemplate.TemplateDescription,
			TemplateTitle:       templateTitle,
			TemplateContent:     templateContent,
			TemplateAttachName:  defaultTemplate.TemplateAttachName,
			CreateTime:          gtime.Now().Timestamp(),
			Status:              0,
		}
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
		return err
	} else {
		//update
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailTemplate.Columns().MerchantId:          merchantId,
			dao.MerchantEmailTemplate.Columns().TemplateName:        defaultTemplate.TemplateName,
			dao.MerchantEmailTemplate.Columns().TemplateDescription: defaultTemplate.TemplateDescription,
			dao.MerchantEmailTemplate.Columns().TemplateTitle:       templateTitle,
			dao.MerchantEmailTemplate.Columns().TemplateContent:     templateContent,
			dao.MerchantEmailTemplate.Columns().TemplateAttachName:  defaultTemplate.TemplateAttachName,
			dao.MerchantEmailTemplate.Columns().GmtModify:           gtime.Now(),
			dao.MerchantEmailTemplate.Columns().Status:              0,
		}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		return err
	}
}

func SetMerchantEmailTemplateDefault(ctx context.Context, merchantId uint64, templateName string) error {
	utility.Assert(merchantId > 0, "Invalid MerchantId")
	utility.Assert(len(templateName) > 0, "Invalid TemplateName")
	var defaultTemplate *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&defaultTemplate)
	utility.AssertError(err, "Server Error")
	utility.Assert(defaultTemplate != nil, "Default Template Not Found")
	var one *entity.MerchantEmailTemplate
	err = dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	utility.AssertError(err, "Server Error")
	if one == nil {
		//insert
		one = &entity.MerchantEmailTemplate{
			MerchantId:          merchantId,
			TemplateName:        defaultTemplate.TemplateName,
			TemplateDescription: defaultTemplate.TemplateDescription,
			TemplateTitle:       defaultTemplate.TemplateTitle,
			TemplateContent:     defaultTemplate.TemplateContent,
			TemplateAttachName:  defaultTemplate.TemplateAttachName,
			CreateTime:          gtime.Now().Timestamp(),
			Status:              0,
		}
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
		return err
	} else {
		//update
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailTemplate.Columns().MerchantId:          merchantId,
			dao.MerchantEmailTemplate.Columns().TemplateName:        defaultTemplate.TemplateName,
			dao.MerchantEmailTemplate.Columns().TemplateDescription: defaultTemplate.TemplateDescription,
			dao.MerchantEmailTemplate.Columns().TemplateTitle:       defaultTemplate.TemplateTitle,
			dao.MerchantEmailTemplate.Columns().TemplateContent:     defaultTemplate.TemplateContent,
			dao.MerchantEmailTemplate.Columns().TemplateAttachName:  defaultTemplate.TemplateAttachName,
			dao.MerchantEmailTemplate.Columns().GmtModify:           gtime.Now(),
			dao.MerchantEmailTemplate.Columns().Status:              0,
		}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		return err
	}
}

func ActivateMerchantEmailTemplate(ctx context.Context, merchantId uint64, templateName string) error {
	utility.Assert(merchantId > 0, "Invalid MerchantId")
	utility.Assert(len(templateName) > 0, "Invalid TemplateName")
	var defaultTemplate *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&defaultTemplate)
	utility.AssertError(err, "Server Error")
	utility.Assert(defaultTemplate != nil, "Default Template Not Found")
	var one *entity.MerchantEmailTemplate
	err = dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	utility.AssertError(err, "Server Error")

	if one == nil {
		//insert
		one = &entity.MerchantEmailTemplate{
			MerchantId:          merchantId,
			TemplateName:        defaultTemplate.TemplateName,
			TemplateDescription: defaultTemplate.TemplateDescription,
			TemplateTitle:       defaultTemplate.TemplateTitle,
			TemplateContent:     defaultTemplate.TemplateContent,
			TemplateAttachName:  defaultTemplate.TemplateAttachName,
			CreateTime:          gtime.Now().Timestamp(),
			Status:              0,
		}
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
		return err
	} else {
		if one.Status == 0 {
			return nil
		}
		//update
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailTemplate.Columns().GmtModify: gtime.Now(),
			dao.MerchantEmailTemplate.Columns().Status:    0,
		}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		return err
	}
}

func DeactivateMerchantEmailTemplate(ctx context.Context, merchantId uint64, templateName string) error {
	utility.Assert(merchantId > 0, "Invalid MerchantId")
	utility.Assert(len(templateName) > 0, "Invalid TemplateName")
	var defaultTemplate *entity.EmailDefaultTemplate
	err := dao.EmailDefaultTemplate.Ctx(ctx).
		Where(dao.EmailDefaultTemplate.Columns().TemplateName, templateName).
		Scan(&defaultTemplate)
	utility.AssertError(err, "Server Error")
	utility.Assert(defaultTemplate != nil, "Default Template Not Found")
	var one *entity.MerchantEmailTemplate
	err = dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailTemplate.Columns().TemplateName, templateName).
		Scan(&one)
	utility.AssertError(err, "Server Error")
	if one == nil {
		//insert
		one = &entity.MerchantEmailTemplate{
			MerchantId:          merchantId,
			TemplateName:        defaultTemplate.TemplateName,
			TemplateDescription: defaultTemplate.TemplateDescription,
			TemplateTitle:       defaultTemplate.TemplateTitle,
			TemplateContent:     defaultTemplate.TemplateContent,
			TemplateAttachName:  defaultTemplate.TemplateAttachName,
			CreateTime:          gtime.Now().Timestamp(),
			Status:              1,
		}
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(one).Insert(one)
		return err
	} else {
		if one.Status == 1 {
			return nil
		}
		//update
		_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailTemplate.Columns().GmtModify: gtime.Now(),
			dao.MerchantEmailTemplate.Columns().Status:    1,
		}).Where(dao.Invoice.Columns().Id, one.Id).Update()
		return err
	}
}
