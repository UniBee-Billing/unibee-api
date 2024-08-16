package email

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	log2 "unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/email/gateway"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "os"
	"fmt"
)

const (
	TemplateInvoiceAutomaticPaid                            = "InvoiceAutomaticPaid"
	TemplateInvoiceManualPaid                               = "InvoiceManualPaid"
	TemplateNewProcessingInvoice                            = "NewProcessingInvoice"
	TemplateNewProcessingInvoiceForPaidTrial                = "NewProcessingInvoiceForPaidTrial"
	TemplateNewProcessingInvoiceAfterTrial                  = "NewProcessingInvoiceAfterTrial"
	TemplateNewProcessingInvoiceForWireTransfer             = "NewProcessingInvoiceForWireTransfer"
	TemplateInvoiceCancel                                   = "InvoiceCancel"
	TemplateMerchantRegistrationCodeVerify                  = "MerchantRegistrationCodeVerify"
	TemplateMerchantOTPLogin                                = "MerchantOTPLogin"
	TemplateUserRegistrationCodeVerify                      = "UserRegistrationCodeVerify"
	TemplateUserOTPLogin                                    = "UserOTPLogin"
	TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin = "SubscriptionCancelledAtPeriodEndByMerchantAdmin"
	TemplateSubscriptionCancelledAtPeriodEndByUser          = "SubscriptionCancelledAtPeriodEndByUser"
	TemplateSubscriptionCancelledByTrialEnd                 = "SubscriptionCancelledByTrialEnd"
	TemplateSubscriptionCancelLastCancelledAtPeriodEnd      = "SubscriptionCancelLastCancelledAtPeriodEnd"
	TemplateSubscriptionImmediateCancel                     = "SubscriptionImmediateCancel"
	TemplateSubscriptionUpdate                              = "SubscriptionUpdate"
	TemplateSubscriptionNeedAuthorized                      = "SubscriptionNeedAuthorized"
	TemplateSubscriptionTrialStart                          = "SubscriptionTrialStart"
	TemplateInvoiceRefundCreated                            = "InvoiceRefundCreated"
	TemplateInvoiceRefundPaid                               = "InvoiceRefundPaid"
	TemplateMerchantMemberInvite                            = "MerchantMemberInvite"
)

const (
	KeyMerchantEmailName = "KEY_MERCHANT_DEFAULT_EMAIL_NAME"
	IMPLEMENT_NAMES      = "sendgrid"
)

func GetDefaultMerchantEmailConfig(ctx context.Context, merchantId uint64) (name string, data string) {
	nameConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantEmailName)
	if nameConfig != nil {
		name = nameConfig.ConfigValue
	}
	valueConfig := merchant_config.GetMerchantConfig(ctx, merchantId, name)
	if valueConfig != nil {
		data = valueConfig.ConfigValue
	}
	return
}

func SetupMerchantEmailConfig(ctx context.Context, merchantId uint64, name string, data string, isDefault bool) error {
	utility.Assert(strings.Contains(IMPLEMENT_NAMES, name), "gateway not support, should be "+IMPLEMENT_NAMES)
	err := update.SetMerchantConfig(ctx, merchantId, name, data)
	if err != nil {
		return err
	}
	if isDefault {
		err = update.SetMerchantConfig(ctx, merchantId, KeyMerchantEmailName, name)
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     merchantId,
		Target:         fmt.Sprintf("EmailGateway(%s)-SetDefault(%v)", name, isDefault),
		Content:        "SetupEmailGateway",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

type TemplateVariable struct {
	InvoiceId             string      `json:"InvoiceId"`
	UserName              string      `json:"User name" key:"UserName"`
	MerchantProductName   string      `json:"Merchant Product Name" key:"ProductName"`
	MerchantCustomerEmail string      `json:"Merchant’s customer support email address" key:"SupportEmail"`
	MerchantName          string      `json:"Merchant Name" key:"MerchantName"`
	DateNow               *gtime.Time `json:"DateNow" layout:"2006-01-02"`
	PeriodEnd             *gtime.Time `json:"PeriodEnd" layout:"2006-01-02"`
	PaymentAmount         string      `json:"Payment Amount" key:"PaymentAmount"`
	RefundAmount          string      `json:"Refund Amount" key:"RefundAmount"`
	Currency              string      `json:"Currency" key:"Currency"`
	TokenExpireMinute     string      `json:"TokenExpireMinute"`
	CodeExpireMinute      string      `json:"CodeExpireMinute"`
	Code                  string      `json:"Code"`
	Link                  string      `json:"Link"`
	AccountHolder         string      `json:"Account Holder" key:"WireTransferAccountHolder"`
	Address               string      `json:"Address" key:"WireTransferAddress"`
	BIC                   string      `json:"BIC" key:"WireTransferBIC"`
	IBAN                  string      `json:"IBAN" key:"WireTransferIBAN"`
}

func SendTemplateEmailByOpenApi(ctx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, variableMap map[string]interface{}) (err error) {
	_, emailGatewayKey := GetDefaultMerchantEmailConfig(ctx, merchantId)
	if len(emailGatewayKey) == 0 {
		if strings.Compare(templateName, TemplateUserOTPLogin) == 0 || strings.Compare(templateName, TemplateUserRegistrationCodeVerify) == 0 {
			utility.Assert(false, "Default Email Gateway Need Setup")
		} else {
			return gerror.New("Default Email Gateway Need Setup")
		}
	}
	var template *bean.MerchantEmailTemplate
	if merchantId > 0 {
		template = query.GetMerchantEmailTemplateByTemplateName(ctx, merchantId, templateName)
	} else {
		template = query.GetEmailDefaultTemplateByTemplateName(ctx, templateName)
	}
	utility.Assert(template != nil, "template not found:"+templateName)
	utility.Assert(strings.Compare(template.Status, "Active") == 0, "template not active status")
	utility.Assert(template != nil, "template not found")
	utility.Assert(variableMap != nil, "variableMap not found")
	var title = toLocalizationTitle(template.LanguageData, template.TemplateTitle, language)
	var content = template.TemplateContent
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
	for key, value := range variableMap {
		mapKey := "{{" + key + "}}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(title) > 0 {
			title = strings.Replace(title, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	for key, value := range variableMap {
		mapKey := "{" + key + "}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(title) > 0 {
			title = strings.Replace(title, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = "attach"
	}

	var response string
	if len(pdfFilePath) > 0 {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s%s", mailTo, title, content, attachName))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(template.GatewayTemplateId) > 0 {
			response, err = gateway.SendDynamicPdfAttachEmailToUser(emailGatewayKey, mailTo, title, template.GatewayTemplateId, variableMap, language, pdfFilePath, attachName+".pdf")
		} else {
			response, err = gateway.SendPdfAttachEmailToUser(emailGatewayKey, mailTo, title, content, pdfFilePath, attachName+".pdf")
		}
		if err != nil {
			SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", err.Error())
		} else {
			SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", response)
		}
		return err
	} else {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s", mailTo, title, content))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(template.GatewayTemplateId) > 0 {
			response, err = gateway.SendDynamicTemplateEmailToUser(emailGatewayKey, mailTo, title, template.GatewayTemplateId, variableMap, language)
		} else {
			response, err = gateway.SendEmailToUser(emailGatewayKey, mailTo, title, content)
		}
		if err != nil {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", err.Error())
		} else {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", response)
		}
		return err
	}
}

// SendTemplateEmail template should convert by html tools like https://www.iamwawa.cn/text2html.html
func SendTemplateEmail(superCtx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, templateVariables *TemplateVariable) error {
	_, emailGatewayKey := GetDefaultMerchantEmailConfig(superCtx, merchantId)
	if len(emailGatewayKey) == 0 {
		if strings.Compare(templateName, TemplateUserOTPLogin) == 0 || strings.Compare(templateName, TemplateUserRegistrationCodeVerify) == 0 {
			utility.Assert(false, "Default Email Gateway Need Setup")
		} else {
			return gerror.New("Default Email Gateway Need Setup")
		}
	}
	go func() {
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log2.PrintPanic(backgroundCtx, err)
				return
			}
		}()
		err = sendTemplateEmailInternal(backgroundCtx, merchantId, mailTo, timezone, language, templateName, pdfFilePath, templateVariables, emailGatewayKey)
		utility.AssertError(err, "sendTemplateEmailInternal")
	}()
	return nil
}

func sendTemplateEmailInternal(ctx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, templateVariables *TemplateVariable, emailGatewayKey string) error {
	var template *bean.MerchantEmailTemplate
	if merchantId > 0 {
		template = query.GetMerchantEmailTemplateByTemplateName(ctx, merchantId, templateName)
	} else {
		template = query.GetEmailDefaultTemplateByTemplateName(ctx, templateName)
	}
	utility.Assert(template != nil, "template not found:"+templateName)
	utility.Assert(strings.Compare(template.Status, "Active") == 0, "template not active status")
	utility.Assert(template != nil, "template not found")
	utility.Assert(templateVariables != nil, "templateVariables not found")
	variableMap, err := utility.ReflectTemplateStructToMap(templateVariables, timezone)
	if err != nil {
		return err
	}
	var title = toLocalizationTitle(template.LanguageData, template.TemplateTitle, language)
	var content = template.TemplateContent
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
	for key, value := range variableMap {
		mapKey := "{{" + key + "}}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(title) > 0 {
			title = strings.Replace(title, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	for key, value := range variableMap {
		mapKey := "{" + key + "}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(title) > 0 {
			title = strings.Replace(title, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = "attach"
	}

	var response string
	if len(pdfFilePath) > 0 {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s%s", mailTo, title, content, attachName))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(template.GatewayTemplateId) > 0 {
			response, err = gateway.SendDynamicPdfAttachEmailToUser(emailGatewayKey, mailTo, title, template.GatewayTemplateId, variableMap, language, pdfFilePath, attachName+".pdf")
		} else {
			response, err = gateway.SendPdfAttachEmailToUser(emailGatewayKey, mailTo, title, content, pdfFilePath, attachName+".pdf")
		}
		if err != nil {
			if len(template.GatewayTemplateId) > 0 {
				SaveHistory(ctx, merchantId, mailTo, title, utility.MarshalToJsonString(variableMap), attachName+".pdf", err.Error())
			} else {
				SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", err.Error())
			}
		} else {
			if len(template.GatewayTemplateId) > 0 {
				SaveHistory(ctx, merchantId, mailTo, title, utility.MarshalToJsonString(variableMap), attachName+".pdf", response)
			} else {
				SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", response)
			}
		}
		return err
	} else {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s", mailTo, title, content))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(template.GatewayTemplateId) > 0 {
			response, err = gateway.SendDynamicTemplateEmailToUser(emailGatewayKey, mailTo, title, template.GatewayTemplateId, variableMap, language)
		} else {
			response, err = gateway.SendEmailToUser(emailGatewayKey, mailTo, title, content)
		}
		if err != nil {
			if len(template.GatewayTemplateId) > 0 {
				SaveHistory(ctx, merchantId, mailTo, title, utility.MarshalToJsonString(variableMap), "", err.Error())
			} else {
				SaveHistory(ctx, merchantId, mailTo, title, content, "", err.Error())
			}
		} else {
			if len(template.GatewayTemplateId) > 0 {
				SaveHistory(ctx, merchantId, mailTo, title, utility.MarshalToJsonString(variableMap), "", response)
			} else {
				SaveHistory(ctx, merchantId, mailTo, title, content, "", response)
			}
		}
		return err
	}
}

func SaveHistory(ctx context.Context, merchantId uint64, mailTo string, title string, content string, attachFilePath string, response string) {
	var err error
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			g.Log().Errorf(ctx, "SaveEmailHistory panic error:%s", err.Error())
			return
		}
	}()
	one := &entity.MerchantEmailHistory{
		MerchantId: merchantId,
		Email:      mailTo,
		Title:      title,
		Content:    content,
		AttachFile: attachFilePath,
		Response:   response,
		CreateTime: gtime.Now().Timestamp(),
	}
	_, _ = dao.MerchantEmailHistory.Ctx(ctx).Data(one).OmitNil().Insert(one)
}

func toLocalizationTitle(languageData string, defaultTitle string, lang string) (title string) {
	title = defaultTitle
	if len(languageData) == 0 || len(lang) == 0 {
		return title
	}
	var list []*bean.EmailLocalizationTemplate
	err := bean.UnmarshalFromJsonString(languageData, &list)
	if err == nil {
		for _, one := range list {
			if one.Language == lang {
				title = one.Title
			}
		}
	}
	return title
}
