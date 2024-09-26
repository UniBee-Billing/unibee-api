package gateway

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"testing"
)

func TestSendgridTemplate(t *testing.T) {
	apiKey := "changeme"
	host := "https://api.sendgrid.com"
	t.Run("Test Templates", func(t *testing.T) {
		request := sendgrid.GetRequest(apiKey, "/v3/templates", host)
		request.Method = "GET"
		queryParams := make(map[string]string)
		queryParams["generations"] = "legacy,dynamic"
		queryParams["page_size"] = "200"
		request.QueryParams = queryParams
		response, err := sendgrid.API(request)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	})
	t.Run("Test Single Template", func(t *testing.T) {
		request := sendgrid.GetRequest(apiKey, "/v3/templates/d-42e8f092b8754b06bba426f4608aca8a", host)
		request.Method = "GET"
		response, err := sendgrid.API(request)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	})
	t.Run("Test Create Template", func(t *testing.T) {
		request := sendgrid.GetRequest(apiKey, "/v3/templates", host)
		request.Method = "POST"
		request.Body = []byte(` {
  "name": "[UniBee][MLX]example_template[Linked]",
  "generation":"dynamic"
}`)
		response, err := sendgrid.API(request)
		if err != nil {
			fmt.Printf("Create Sendgrid template error:%s\n", err.Error())
		}
		data := gjson.New(response)
		if data == nil || !data.Contains("id") || data.Get("id") == nil {
			fmt.Printf("Create template error,no templateId\n")
		}
		templateId := data.Get("id").String()
		fmt.Println(templateId)
	})

	t.Run("Test Create Template Version", func(t *testing.T) {
		request := sendgrid.GetRequest(apiKey, "/v3/templates/d-bc051c7cd0244fa9bf7c4391b72d0809/versions", host)
		request.Method = "POST"
		request.Body = []byte(` {
  "html_content": "<%html_content%>",
  "name": "example_version_name",
  "subject": "<%subject%>"
}`)
		response, err := sendgrid.API(request)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	})
	t.Run("Test Create Template By Api", func(t *testing.T) {
		templateId, err := SyncToGatewayTemplate(context.Background(), apiKey, "SubscriptionCancelledAtPeriodEndByUser", "Content", "", false)
		if err != nil {
			log.Println(err)
		}
		log.Println(templateId)
	})
}
