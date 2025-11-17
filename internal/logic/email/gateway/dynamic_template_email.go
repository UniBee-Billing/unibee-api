package gateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"unibee/internal/logic/email/sender"
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
func SendSandgridDynamicTemplateEmailToUser(f *sender.Sender, emailGatewayKey string, mailTo string, subject string, templateId string, variables map[string]interface{}, language string) (result string, err error) {
	if f == nil {
		f = sender.GetDefaultSender()
	}
	from := mail.NewEmail(f.Name, f.Address)
	to := mail.NewEmail(mailTo, mailTo)
	message := mail.NewV3Mail()
	message.SetFrom(from)
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("Subject", subject)
	p.Subject = subject
	lang := LangMap[language]
	if len(lang) > 0 {
		p.SetDynamicTemplateData(lang, true)
	}
	for key, value := range variables {
		p.SetDynamicTemplateData(key, fmt.Sprintf("%s", value))
	}
	message.AddPersonalizations(p)
	message.SetTemplateID(templateId)
	message.Subject = subject

	client := sendgrid.NewSendClient(emailGatewayKey)
	response, err := client.Send(message)
	if err != nil {
		fmt.Printf("SendSandgridDynamicTemplateEmailToUser error:%s\n", err.Error())
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}

func SendSendgridDynamicTemplateWithAttachFileEmailToUser(f *sender.Sender, emailGatewayKey string, mailTo string, subject string, templateId string, variables map[string]interface{}, language string, pdfFilePath string, pdfFileName string) (result string, err error) {
	if f == nil {
		f = sender.GetDefaultSender()
	}
	from := mail.NewEmail(f.Name, f.Address)
	to := mail.NewEmail(mailTo, mailTo)
	message := mail.NewV3Mail()
	message.SetFrom(from)
	p := mail.NewPersonalization()
	tos := []*mail.Email{
		to,
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("Subject", subject)
	p.Subject = subject
	lang := LangMap[language]
	if len(lang) > 0 {
		p.SetDynamicTemplateData(lang, true)
	}
	for key, value := range variables {
		p.SetDynamicTemplateData(key, fmt.Sprintf("%s", value))
	}
	message.AddPersonalizations(p)
	message.SetTemplateID(templateId)
	message.Subject = subject

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
		fmt.Printf("SendSendgridDynamicTemplateWithAttachFileEmailToUser error:%s\n", err.Error())
		return "", err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return utility.MarshalToJsonString(response), nil
}

var sendGridHost = "https://api.sendgrid.com"

func SyncToGatewayTemplate(ctx context.Context, apiKey string, templateName string, content string, oldTemplateId string, versionEnable bool) (templateId string, err error) {
	name := fmt.Sprintf("[UniBee]%s", templateName)
	if len(oldTemplateId) == 0 {
		request := sendgrid.GetRequest(apiKey, "/v3/templates", sendGridHost)
		request.Method = "POST"
		request.Body = []byte(utility.MarshalToJsonString(map[string]string{"name": name, "generation": "dynamic"}))
		response, err := sendgrid.API(request)
		if err != nil {
			g.Log().Error(ctx, "Create Sendgrid template error:%s", err.Error())
			return "", gerror.New(fmt.Sprintf("Create template error:%s", err.Error()))
		}
		data := gjson.New(response.Body)
		if data == nil || !data.Contains("id") || data.Get("id") == nil || len(data.Get("id").String()) == 0 || response.StatusCode != 201 {
			return "", gerror.Newf("Create template error,no templateId, code:%v", response.StatusCode)
		}
		templateId = data.Get("id").String()
		if len(templateId) == 0 {
			return "", gerror.Newf("Create template error,no templateId, code:%v", response.StatusCode)
		}
	} else {
		templateId = oldTemplateId
	}
	param := map[string]interface{}{"name": fmt.Sprintf("[UniBeeVersion]%d", gtime.Now().Timestamp()), "html_content": content, "subject": "{{Subject}}"}
	if versionEnable {
		param["active"] = 1
	}
	request := sendgrid.GetRequest(apiKey, fmt.Sprintf("/v3/templates/%s/versions", templateId), sendGridHost)
	request.Method = "POST"
	request.Body = []byte(utility.MarshalToJsonString(param))
	response, err := sendgrid.API(request)
	if err != nil {
		g.Log().Error(ctx, "Create Sendgrid template version error:%s", err.Error())
		return "", gerror.New(fmt.Sprintf("Create template version error:%s", err.Error()))
	}
	data := gjson.New(response.Body)
	if data == nil || !data.Contains("id") || data.Get("id") == nil || len(data.Get("id").String()) == 0 || response.StatusCode != 201 {
		return "", gerror.Newf("Create template version error,no templateVersionId, code:%v", response.StatusCode)
	}
	return templateId, nil
}
