package email

import (
	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "os"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func SendPdfAttachEmailToUser(mailTo string, subject string, body string, pdfFileName string) error {
	from := mail.NewEmail("no-reply", "no-reply@unibee.dev")
	subject = subject
	to := mail.NewEmail(mailTo, mailTo)
	plainTextContent := body
	htmlContent := "<strong>" + body + " </strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	attach := mail.NewAttachment()
	attach.SetContent("TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gQ3JhcyBwdW12")
	attach.SetType("application/pdf")
	attach.SetFilename(pdfFileName)
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
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
