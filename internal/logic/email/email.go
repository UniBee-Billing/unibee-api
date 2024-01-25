package email

import (
	"context"
	"encoding/base64"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"os"
	"reflect"
	"strings"

	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "os"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	TemplateInvoiceAutomaticPaid                            = "InvoiceAutomaticPaid"
	TemplateUserRegistrationCodeVerify                      = "UserRegistrationCodeVerify"
	TemplateUserOTPLogin                                    = "UserOTPLogin"
	TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin = "SubscriptionCancelledAtPeriodEndByMerchantAdmin"
	TemplateSubscriptionCancelledAtPeriodEndByUser          = "SubscriptionCancelledAtPeriodEndByUser"
)

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

func SendPdfAttachEmailToUser(mailTo string, subject string, body string, pdfFilePath string, pdfFileName string) error {
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
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

type TemplateVariable struct {
	InvoiceId           string `json:"InvoiceId"`
	UserName            string `json:"User name"`
	MerchantProductName string `json:"Merchant Product Name"`
	MerchantCustomEmail string `json:"Merchant’s customer support email address"`
	MerchantName        string `json:"Merchant Name"`
	DateNow             string `json:"DateNow"`
	PaymentAmount       string `json:"Payment Amount"`
	TokenExpireMinute   string `json:"TokenExpireMinute"`
	CodeExpireMinute    string `json:"CodeExpireMinute"`
	Code                string `json:"Code"`
	PeriodEnd           string `json:"PeriodEnd"`
}

func ToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// SendTemplateEmail template should convert by html tools like https://www.iamwawa.cn/text2html.html
func SendTemplateEmail(ctx context.Context, merchantId int64, mailTo string, templateName string, pdfFilePath string, templateVariables *TemplateVariable) error {
	var template *entity.EmailTemplate
	if merchantId > 0 {
		template = query.GetMerchantEmailTemplateByTemplateName(ctx, merchantId, templateName)
	} else {
		template = query.GetEmailTemplateByTemplateName(ctx, templateName)
	}
	utility.Assert(template != nil, "template not found")
	utility.Assert(templateVariables != nil, "templateVariables not found")
	variableMap, err := ToMap(templateVariables)
	if err != nil {
		return err
	}
	var title = template.TemplateTitle
	var content = template.TemplateContent
	var attachName = template.TemplateAttachName
	utility.Assert(variableMap != nil, "template parse error")
	for key, value := range variableMap {
		mapKey := "{" + key + "}"
		mapKey = strings.Replace(mapKey, " ", "&nbsp;", 10)
		mapValue := "<strong>" + value.(string) + "</strong>"
		if len(title) > 0 {
			title = strings.Replace(title, mapKey, mapValue, 1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, mapKey, mapValue, 1)
		}
		if len(attachName) > 0 {
			attachName = strings.Replace(attachName, mapKey, mapValue, 1)
		}
	}
	if len(pdfFilePath) > 0 && len(attachName) == 0 {
		attachName = "attach"
	}
	if len(pdfFilePath) > 0 {
		return SendPdfAttachEmailToUser(mailTo, title, content, pdfFilePath, attachName+".pdf")
	} else {
		return SendEmailToUser(mailTo, title, content)
	}
}
