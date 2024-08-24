package email

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/email/gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func SyncMerchantEmailTemplateToGateway(ctx context.Context, id int64, versionEnable bool) error {
	var one *entity.MerchantEmailTemplate
	err := dao.MerchantEmailTemplate.Ctx(ctx).
		Where(dao.MerchantEmailTemplate.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return err
	}
	_, emailGatewayKey := GetDefaultMerchantEmailConfig(ctx, one.MerchantId)
	if len(emailGatewayKey) == 0 {
		return gerror.New("Default Email Gateway Need Setup")
	}
	content := one.TemplateContent
	if len(one.LanguageData) == 0 {
		content = gateway.ConvertUniBeeTemplateToPlain(content)
		content = gateway.ConvertPainToHtmlContent(content)
		content = gateway.ConvertToHtmlPage(content)
	} else {
		content = ConvertToSendgridLocalizationHtmlContent(one.TemplateContent, one.LanguageData)
		content = gateway.ConvertToHtmlPage(content)
	}
	gatewayTemplateName := fmt.Sprintf("[%d][%s]", one.Id, one.TemplateName)
	if len(one.TemplateDescription) > 0 {
		gatewayTemplateName = TruncateWithEllipsis(fmt.Sprintf("[%d][%d][%s][%s]", one.Id, one.MerchantId, one.TemplateName, one.TemplateDescription), 90, "...")
	}
	templateId, err := gateway.SyncToGatewayTemplate(ctx, emailGatewayKey, gatewayTemplateName, content, one.GatewayTemplateId, versionEnable)
	if err != nil {
		return err
	}
	_, err = dao.MerchantEmailTemplate.Ctx(ctx).Data(g.Map{
		dao.MerchantEmailTemplate.Columns().GatewayTemplateId: templateId,
		dao.MerchantEmailTemplate.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).Update()
	return err
}

func TruncateWithEllipsis(s string, maxLength int, ellipsis string) string {
	if len(s) <= maxLength {
		return s
	}
	ellipsisLen := len(ellipsis)
	maxContentLength := maxLength - ellipsisLen
	if maxContentLength < 0 {
		maxContentLength = 0
	}
	return s[:maxContentLength] + ellipsis
}

func ConvertToSendgridLocalizationHtmlContent(enContent string, languageData string) string {
	var list []*bean.EmailLocalizationTemplate
	err := bean.UnmarshalFromJsonString(languageData, &list)
	utility.AssertError(err, "ConvertToSendgridLocalizationHtmlContent error")
	content := "{{#if english}}"
	content = fmt.Sprintf("%s\n%s", content, gateway.ConvertPainToHtmlContent(enContent))
	for _, one := range list {
		if one.Language == "en" {
			continue
		} else if _, ok := gateway.LangMap[one.Language]; ok {
			content = fmt.Sprintf("%s\n{{else if %s}}", content, gateway.LangMap[one.Language])
			content = fmt.Sprintf("%s\n%s", content, gateway.ConvertPainToHtmlContent(one.Content))
		}
	}
	content = fmt.Sprintf("%s\n{{else}}", content)
	content = fmt.Sprintf("%s\n%s", content, gateway.ConvertPainToHtmlContent(enContent))
	content = fmt.Sprintf("%s\n{{/if}}", content)
	return content
}
