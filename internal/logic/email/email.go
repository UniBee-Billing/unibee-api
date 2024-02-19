package email

import (
	"context"
	"encoding/base64"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"os"
	"strings"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"

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
	TemplateUserRegistrationCodeVerify                      = "UserRegistrationCodeVerify"
	TemplateUserOTPLogin                                    = "UserOTPLogin"
	TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin = "SubscriptionCancelledAtPeriodEndByMerchantAdmin"
	TemplateSubscriptionCancelledAtPeriodEndByUser          = "SubscriptionCancelledAtPeriodEndByUser"
	TemplateSubscriptionCancelLastCancelledAtPeriodEnd      = "SubscriptionCancelLastCancelledAtPeriodEnd"
	TemplateSubscriptionImmediateCancel                     = "SubscriptionImmediateCancel"
	TemplateSubscriptionUpdate                              = "SubscriptionUpdate"
	TemplateSubscriptionNeedAuthorized                      = "SubscriptionNeedAuthorized"
	TemplateInvoiceRefundCreated                            = "InvoiceRefundCreated"
)

// const SG_KEY = "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"
const SG_KEY = "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"

func SendEmailToUser(mailTo string, subject string, body string) error {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	subject = subject
	to := mail.NewEmail(mailTo, mailTo)
	plainTextContent := body
	htmlContent := "<strong>" + body + " </strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client := sendgrid.NewSendClient(SG_KEY)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

func SendTemplateEmailToUser(mailTo string, subject string, body string) (result string, err error) {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	subject = subject
	to := mail.NewEmail(mailTo, mailTo)
	plainTextContent := body
	htmlContent := "<div>" + body + " </div>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client := sendgrid.NewSendClient(SG_KEY)
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

func SendPdfAttachEmailToUser(mailTo string, subject string, body string, pdfFilePath string, pdfFileName string) (result string, err error) {
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
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client := sendgrid.NewSendClient(SG_KEY)
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
	PaymentAmount       string      `json:"Payment Amount"`
	RefundAmount        string      `json:"Refund Amount"`
	Currency            string      `json:"Currency"`
	TokenExpireMinute   string      `json:"TokenExpireMinute"`
	CodeExpireMinute    string      `json:"CodeExpireMinute"`
	Code                string      `json:"Code"`
	Link                string      `json:"Link"`
}

// SendTemplateEmail template should convert by html tools like https://www.iamwawa.cn/text2html.html
func SendTemplateEmail(ctx context.Context, merchantId int64, mailTo string, timezone string, templateName string, pdfFilePath string, templateVariables *TemplateVariable) error {
	var template *query.EmailTemplateVo
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
			title = strings.Replace(title, mapKey, value.(string), 1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, htmlKey, htmlValue, 1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, value.(string), 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = "attach"
	}
	if len(pdfFilePath) > 0 {
		response, err := SendPdfAttachEmailToUser(mailTo, title, content, pdfFilePath, attachName+".pdf")
		if err != nil {
			SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", err.Error())
		} else {
			SaveHistory(ctx, merchantId, mailTo, title, content, attachName+".pdf", response)
		}
		return err
	} else {
		response, err := SendTemplateEmailToUser(mailTo, title, content)
		if err != nil {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", err.Error())
		} else {
			SaveHistory(ctx, merchantId, mailTo, title, content, "", response)
		}
		return err
	}
}

func SaveHistory(ctx context.Context, merchantId int64, mailTo string, title string, content string, attachFilePath string, response string) {
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
