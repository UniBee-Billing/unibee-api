package gateway

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"testing"
)

func TestSendgridTemplate(t *testing.T) {
	apiKey := "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"
	host := "https://api.sendgrid.com"
	t.Run("Test Templates", func(t *testing.T) {
		request := sendgrid.GetRequest(apiKey, "/v3/templates", host)
		request.Method = "GET"
		queryParams := make(map[string]string)
		queryParams["generations"] = "legacy,dynamic"
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
  "name": "example_template",
  "generation":"dynamic"
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
}
