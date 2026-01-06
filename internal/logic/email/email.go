package email

import (
	"context"
	"reflect"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	log2 "unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/email/gateway"
	"unibee/internal/logic/email/sender"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/middleware/rate_limit"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"

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
	KeyMerchantEmailName   = "KEY_MERCHANT_DEFAULT_EMAIL_NAME"
	IMPLEMENT_NAMES        = "sendgrid"
	KeyMerchantEmailSender = "KEY_MERCHANT_EMAIL_SENDER"
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

func GetDefaultMerchantEmailConfigWithClusterCloud(ctx context.Context, merchantId uint64) (name string, data string) {
	nameConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantEmailName)
	if nameConfig != nil {
		name = nameConfig.ConfigValue
	}
	valueConfig := merchant_config.GetMerchantConfig(ctx, merchantId, name)
	if valueConfig != nil {
		data = valueConfig.ConfigValue
	}
	if config.GetConfigInstance().Mode == "cloud" && len(data) == 0 {
		data, _ = getDefaultMerchantEmailConfigFromClusterCloud(ctx, merchantId)
	}
	return
}

func getDefaultMerchantEmailConfigFromClusterCloud(ctx context.Context, merchantId uint64) (string, error) {
	maxHourly := 1000
	if !config.GetConfigInstance().IsProd() {
		maxHourly = 50
	}
	checked, current := rate_limit.CheckRateLimit(ctx, fmt.Sprintf("UniBee#Cloud#MerchantDefaultMerchantEmailConfigFromClusterCloudHourlyLimitCheck#%d", merchantId), maxHourly, 3600)
	g.Log().Infof(ctx, "MerchantDefaultMerchantEmailConfigFromClusterCloudHourlyLimitCheck merchantId:%d currentQps:%d maxHourly:%d", merchantId, current, maxHourly)
	utility.Assert(checked, fmt.Sprintf("Reached max hourly email limitation, please upgrade your plan, current called:%d", current))

	sendgridRes := redismq.Invoke(ctx, &redismq.InvoiceRequest{
		Group:   "GID_UniBee_Cloud",
		Method:  "GetSendgridKey",
		Request: merchantId,
	}, 0)
	if sendgridRes == nil {
		return "", gerror.New("Server Error")
	}
	if !sendgridRes.Status {
		return "", gerror.New(fmt.Sprintf("%v", sendgridRes.Response))
	}
	if sendgridRes.Response == nil {
		return "", gerror.New("sendgrid key not found")
	}
	if key, ok := sendgridRes.Response.(string); ok {
		if len(key) == 0 {
			return "", gerror.New("sendgrid invalid")
		}
		return key, nil
	}
	return "", gerror.New("Get Sendgrid Key Error")
}

func GetMerchantEmailSender(ctx context.Context, merchantId uint64) *sender.Sender {
	config := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantEmailSender)
	var one *sender.Sender
	if config == nil {
		return nil
	} else {
		err := utility.UnmarshalFromJsonString(config.ConfigValue, &one)
		if err == nil && one != nil {
			return one
		} else {
			return nil
		}
	}
}

func SetupMerchantEmailSender(ctx context.Context, merchantId uint64, sender *sender.Sender) error {
	if merchantId > 0 && sender != nil && len(sender.Address) > 0 && len(sender.Name) > 0 {
		err := update.SetMerchantConfig(ctx, merchantId, KeyMerchantEmailSender, utility.MarshalToJsonString(sender))
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     merchantId,
			Target:         fmt.Sprintf("Name(%s)-Address(%v)", sender.Name, sender.Address),
			Content:        "SetupEmailSenderConfig",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
		return nil
	} else {
		return gerror.New("invalid data")
	}
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

func getEmailTemplateGroupVariables() []*bean.TemplateVariableGroup {
	// Create a sample instance to get field information
	sample := &bean.EmailTemplateVariable{}

	// Use reflection to get field information
	v := reflect.ValueOf(sample).Elem()
	t := v.Type()

	// Map to store groups by group name
	groupMap := make(map[string]*bean.TemplateVariableGroup)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		keyTag := field.Tag.Get("key")
		groupTag := field.Tag.Get("group")

		// Determine the variable name - use key if available, otherwise use json
		varName := jsonTag
		if keyTag != "" {
			varName = keyTag
		}

		// Create template variable
		templateVar := &bean.TemplateVariable{
			VariableName: varName,
		}

		// Get or create group based on group tag
		groupName := groupTag
		if groupName == "" {
			groupName = "Other Information" // Default group for fields without group tag
		}

		group, exists := groupMap[groupName]
		if !exists {
			group = &bean.TemplateVariableGroup{
				GroupName: groupName,
				Variables: []*bean.TemplateVariable{},
			}
			groupMap[groupName] = group
		}

		// Add variable to the group
		group.Variables = append(group.Variables, templateVar)
	}

	// Convert map to slice
	var groups []*bean.TemplateVariableGroup
	for _, group := range groupMap {
		groups = append(groups, group)
	}

	return groups
}

func SendTemplateEmailByOpenApi(ctx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, templateVariables *bean.EmailTemplateVariable, languageData *[]*bean.EmailLocalizationTemplate) (err error) {
	mailTo = strings.ToLower(mailTo)
	_, emailGatewayKey := GetDefaultMerchantEmailConfigWithClusterCloud(ctx, merchantId)
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
	utility.Assert(templateVariables != nil, "templateVariables not found")
	variableMap, err := utility.ReflectTemplateStructToMap(templateVariables, timezone)
	if err != nil {
		return err
	}
	var subject = template.LocalizationSubject(language, languageData)
	var content = template.LocalizationContent(language, languageData)
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
	merchant := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchant != nil, "merchant not found")
	variableMap["CompanyName"] = merchant.CompanyName
	for key, value := range variableMap {
		mapKey := "{{" + key + "}}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, mapKey, htmlValue, -1)
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
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, mapKey, htmlValue, -1)
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = fmt.Sprintf("invoice_%s", time.Now().Format("20060102"))
	}
	return Send(ctx, &SendgridEmailReq{
		MerchantId:        merchantId,
		MailTo:            mailTo,
		Subject:           subject,
		Content:           content,
		LocalFilePath:     pdfFilePath,
		AttachName:        attachName + ".pdf",
		APIKey:            emailGatewayKey,
		VariableMap:       variableMap,
		Language:          language,
		GatewayTemplateId: template.GatewayTemplateId,
	})
}

// SendTemplateEmail template should convert by html tools like https://www.iamwawa.cn/text2html.html
func SendTemplateEmail(superCtx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, templateVariables *bean.EmailTemplateVariable) error {
	mailTo = strings.ToLower(mailTo)
	_, emailGatewayKey := GetDefaultMerchantEmailConfigWithClusterCloud(superCtx, merchantId)
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

func sendTemplateEmailInternal(ctx context.Context, merchantId uint64, mailTo string, timezone string, language string, templateName string, pdfFilePath string, templateVariables *bean.EmailTemplateVariable, emailGatewayKey string) error {
	mailTo = strings.ToLower(mailTo)
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
	var subject = template.LocalizationSubject(language, nil)
	var content = template.LocalizationContent(language, nil)
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
	merchant := query.GetMerchantById(ctx, merchantId)
	utility.Assert(merchant != nil, "merchant not found")
	variableMap["CompanyName"] = merchant.CompanyName
	variableMap["UserLanguage"] = language
	variableMap["UserTimezone"] = timezone
	for key, value := range variableMap {
		mapKey := "{{" + key + "}}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
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
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = fmt.Sprintf("invoice_%s", time.Now().Format("20060102"))
	}

	return Send(ctx, &SendgridEmailReq{
		MerchantId:        merchantId,
		MailTo:            mailTo,
		Subject:           subject,
		Content:           content,
		LocalFilePath:     pdfFilePath,
		AttachName:        attachName + ".pdf",
		APIKey:            emailGatewayKey,
		VariableMap:       variableMap,
		Language:          language,
		GatewayTemplateId: template.GatewayTemplateId,
	})
}

type SendgridEmailReq struct {
	MerchantId        uint64                 `json:"merchantId"`
	MailTo            string                 `json:"mailTo"`
	Subject           string                 `json:"subject"`
	Content           string                 `json:"content"`
	LocalFilePath     string                 `json:"localFilePath"`
	AttachName        string                 `json:"attachName"`
	APIKey            string                 `json:"apiKey"`
	VariableMap       map[string]interface{} `json:"variable_map"`
	Language          string                 `json:"language"`
	GatewayTemplateId string                 `json:"gatewayTemplateId"`
}

func Send(ctx context.Context, req *SendgridEmailReq) error {
	var err error
	var response string
	if len(req.LocalFilePath) > 0 {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s%s", req.MailTo, req.Subject, req.Content, req.AttachName))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(req.GatewayTemplateId) > 0 {
			response, err = gateway.SendSendgridDynamicTemplateWithAttachFileEmailToUser(GetMerchantEmailSender(ctx, req.MerchantId), req.APIKey, req.MailTo, req.Subject, req.GatewayTemplateId, req.VariableMap, req.Language, req.LocalFilePath, req.AttachName)
		} else {
			response, err = gateway.SendPdfAttachEmailToUser(GetMerchantEmailSender(ctx, req.MerchantId), req.APIKey, req.MailTo, req.Subject, req.Content, req.LocalFilePath, req.AttachName)
		}
		if err != nil {
			if len(req.GatewayTemplateId) > 0 {
				req.VariableMap["TemplateId"] = req.GatewayTemplateId
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, utility.MarshalToJsonString(req.VariableMap), req.AttachName, err.Error())
			} else {
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, req.Content, req.AttachName, err.Error())
			}
		} else {
			if len(req.GatewayTemplateId) > 0 {
				req.VariableMap["TemplateId"] = req.GatewayTemplateId
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, utility.MarshalToJsonString(req.VariableMap), req.AttachName, response)
			} else {
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, req.Content, req.AttachName, response)
			}
		}
		return err
	} else {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s", req.MailTo, req.Subject, req.Content))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		if len(req.GatewayTemplateId) > 0 {
			response, err = gateway.SendSendgridDynamicTemplateEmailToUser(GetMerchantEmailSender(ctx, req.MerchantId), req.APIKey, req.MailTo, req.Subject, req.GatewayTemplateId, req.VariableMap, req.Language)
		} else {
			response, err = gateway.SendEmailToUser(GetMerchantEmailSender(ctx, req.MerchantId), req.APIKey, req.MailTo, req.Subject, req.Content)
		}
		if err != nil {
			if len(req.GatewayTemplateId) > 0 {
				req.VariableMap["TemplateId"] = req.GatewayTemplateId
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, utility.MarshalToJsonString(req.VariableMap), "", err.Error())
			} else {
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, req.Content, "", err.Error())
			}
		} else {
			if len(req.GatewayTemplateId) > 0 {
				req.VariableMap["TemplateId"] = req.GatewayTemplateId
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, utility.MarshalToJsonString(req.VariableMap), "", response)
			} else {
				SaveHistory(ctx, req.MerchantId, req.MailTo, req.Subject, req.Content, "", response)
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
	status := 0
	if strings.Contains(response, "202") {
		status = 1
	} else {
		status = 2
	}
	one := &entity.MerchantEmailHistory{
		MerchantId: merchantId,
		Email:      mailTo,
		Title:      title,
		Content:    content,
		AttachFile: attachFilePath,
		Response:   response,
		Status:     status,
		CreateTime: gtime.Now().Timestamp(),
	}
	_, _ = dao.MerchantEmailHistory.Ctx(ctx).Data(one).OmitNil().Insert(one)
}

//
//func toLocalizationSubject(template *bean.MerchantEmailTemplate, lang string) (title string) {
//	if len(lang) == 0 {
//		lang = "en" // default language
//	}
//	if len(template.LanguageData) == 0 || len(lang) == 0 {
//		return template.TemplateTitle
//	}
//	for _, one := range template.LanguageData {
//		if one.Language == lang {
//			title = one.Title
//		}
//	}
//	return title
//}
//
//func toLocalizationContent(template *bean.MerchantEmailTemplate, lang string) (content string) {
//	if len(lang) == 0 {
//		lang = "en" // default language
//	}
//	if len(template.LanguageData) == 0 || len(lang) == 0 {
//		return template.TemplateContent
//	}
//	for _, one := range template.LanguageData {
//		if one.Language == lang {
//			content = one.Content
//		}
//	}
//	return content
//}
