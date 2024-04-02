package email

import (
	"context"
	"encoding/base64"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"os"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/merchant_config"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "os"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	TemplateInvoiceAutomaticPaid                            = "InvoiceAutomaticPaid"
	TemplateInvoiceManualPaid                               = "InvoiceManualPaid"
	TemplateNewProcessingInvoice                            = "NewProcessingInvoice "
	TemplateInvoiceCancel                                   = "InvoiceCancel "
	TemplateMerchantRegistrationCodeVerify                  = "MerchantRegistrationCodeVerify"
	TemplateMerchantOTPLogin                                = "MerchantOTPLogin"
	TemplateUserRegistrationCodeVerify                      = "UserRegistrationCodeVerify"
	TemplateUserOTPLogin                                    = "UserOTPLogin"
	TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin = "SubscriptionCancelledAtPeriodEndByMerchantAdmin"
	TemplateSubscriptionCancelledAtPeriodEndByUser          = "SubscriptionCancelledAtPeriodEndByUser"
	TemplateSubscriptionCancelLastCancelledAtPeriodEnd      = "SubscriptionCancelLastCancelledAtPeriodEnd"
	TemplateSubscriptionImmediateCancel                     = "SubscriptionImmediateCancel"
	TemplateSubscriptionUpdate                              = "SubscriptionUpdate"
	TemplateSubscriptionNeedAuthorized                      = "SubscriptionNeedAuthorized"
	TemplateInvoiceRefundCreated                            = "InvoiceRefundCreated"
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
	err := merchant_config.SetMerchantConfig(ctx, merchantId, name, data)
	if err != nil {
		return err
	}
	if isDefault {
		err = merchant_config.SetMerchantConfig(ctx, merchantId, KeyMerchantEmailName, name)
	}
	return err
}

func SendTemplateEmailToUser(key string, mailTo string, subject string, body string) (result string, err error) {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	to := mail.NewEmail(mailTo, mailTo)
	plainTextContent := body
	htmlContent := "<div>" + body + " </div>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}

func SendPdfAttachEmailToUser(key string, mailTo string, subject string, body string, pdfFilePath string, pdfFileName string) (result string, err error) {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	to := mail.NewEmail(mailTo, mailTo)
	plainTextContent := body
	htmlContent := "<div>" + body + " </div>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	attach := mail.NewAttachment()
	dat, err := os.ReadFile(pdfFilePath)
	if err != nil {
		fmt.Println(err)
	}
	encoded := base64.StdEncoding.EncodeToString(dat)
	attach.SetContent(encoded)
	attach.SetType("application/pdf")
	attach.SetFilename(pdfFileName)
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("SendPdfAttachEmailToUser error:%s\n", err.Error())
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}

type TemplateVariable struct {
	InvoiceId           string      `json:"InvoiceId"`
	UserName            string      `json:"User name"`
	MerchantProductName string      `json:"Merchant Product Name"`
	MerchantCustomEmail string      `json:"Merchant’s customer support email address"`
	MerchantName        string      `json:"Merchant Name"`
	DateNow             *gtime.Time `json:"DateNow" layout:"2006-01-02"`
	PeriodEnd           *gtime.Time `json:"PeriodEnd" layout:"2006-01-02"`
	PaymentAmount       string      `json:"Payment CaptureAmount"`
	RefundAmount        string      `json:"Refund CaptureAmount"`
	Currency            string      `json:"Currency"`
	TokenExpireMinute   string      `json:"TokenExpireMinute"`
	CodeExpireMinute    string      `json:"CodeExpireMinute"`
	Code                string      `json:"Code"`
	Link                string      `json:"Link"`
}

// SendTemplateEmail template should convert by html tools like https://www.iamwawa.cn/text2html.html
func SendTemplateEmail(ctx context.Context, merchantId uint64, mailTo string, timezone string, templateName string, pdfFilePath string, templateVariables *TemplateVariable) error {
	var template *bean.EmailTemplateVo
	if merchantId > 0 {
		template = query.GetMerchantEmailTemplateByTemplateName(ctx, merchantId, templateName)
	} else {
		template = query.GetEmailDefaultTemplateByTemplateName(ctx, templateName)
	}
	utility.Assert(strings.Compare(template.Status, "Active") == 0, "template not active status")
	utility.Assert(template != nil, "template not found")
	utility.Assert(templateVariables != nil, "templateVariables not found")
	variableMap, err := utility.ReflectTemplateStructToMap(templateVariables, timezone)
	if err != nil {
		return err
	}
	var title = template.TemplateTitle
	var content = template.TemplateContent
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
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
	_, key := GetDefaultMerchantEmailConfig(ctx, merchantId)
	if len(key) == 0 {
		if strings.Compare(templateName, TemplateUserOTPLogin) == 0 || strings.Compare(templateName, TemplateUserRegistrationCodeVerify) == 0 {
			utility.Assert(false, "Default Email Gateway Need Setup")
		} else {
			return gerror.New("Default Email Gateway Need Setup")
		}
	}

	if len(pdfFilePath) > 0 {
		md5 := utility.MD5(fmt.Sprintf("%s%s%s%s", mailTo, title, content, attachName))
		if !utility.TryLock(ctx, md5, 10) {
			utility.Assert(false, "duplicate email too fast")
		}
		response, err := SendPdfAttachEmailToUser(key, mailTo, title, content, pdfFilePath, attachName+".pdf")
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
		response, err := SendTemplateEmailToUser(key, mailTo, title, content)
		if err != nil {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", err.Error())
		} else {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", response)
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

func doubleRequestLimit(id string, r *ghttp.Request) {

}
