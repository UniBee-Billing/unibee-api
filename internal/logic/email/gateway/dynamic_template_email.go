package gateway

import (
	"encoding/base64"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"unibee/utility"
)

var LangMap = map[string]string{
	"cn": "chinese",
	"en": "english",
	"ru": "russian",
	"vi": "vietnamese",
	"pt": "portuguese",
}

// Sendgrid Template https://github.com/sendgrid/sendgrid-go/blob/main/use-cases/transactional-templates-with-mailer-helper.md
// https://www.twilio.com/docs/sendgrid/api-reference
func SendDynamicTemplateEmailToUser(emailGatewayKey string, mailTo string, subject string, templateId string, variables map[string]interface{}, language string) (result string, err error) {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	to := mail.NewEmail(mailTo, mailTo)
	message := mail.NewV3Mail()
	message.SetFrom(from)
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("subject", subject)
	lang := LangMap[language]
	if len(lang) > 0 {
		p.SetDynamicTemplateData(lang, true)
	}
	for key, value := range variables {
		p.SetDynamicTemplateData(key, fmt.Sprintf("%s", value))
	}
	message.AddPersonalizations(p)
	message.SetTemplateID(templateId)

	client := sendgrid.NewSendClient(emailGatewayKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("SendDynamicTemplateEmailToUser error:%s\n", err.Error())
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}

func SendDynamicPdfAttachEmailToUser(emailGatewayKey string, mailTo string, subject string, templateId string, variables map[string]interface{}, language string, pdfFilePath string, pdfFileName string) (result string, err error) {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	to := mail.NewEmail(mailTo, mailTo)
	message := mail.NewV3Mail()
	message.SetFrom(from)
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("subject", subject)
	lang := LangMap[language]
	if len(lang) > 0 {
		p.SetDynamicTemplateData(lang, true)
	}
	for key, value := range variables {
		p.SetDynamicTemplateData(key, fmt.Sprintf("%s", value))
	}
	message.AddPersonalizations(p)
	message.SetTemplateID(templateId)

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
	client := sendgrid.NewSendClient(emailGatewayKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("SendDynamicPdfAttachEmailToUser error:%s\n", err.Error())
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}
