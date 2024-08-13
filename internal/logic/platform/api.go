package platform

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"net/http"
	"time"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

const ApiHost = "https://api.unibee.dev"

func SentPlatformMerchantOTP(param map[string]string) error {
	response, err := utility.SendRequest(fmt.Sprintf("%s/cloud/email/member_otp", ApiHost), "POST", []byte(utility.MarshalToJsonString(param)), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" {
		return nil
	}
	return gerror.New(fmt.Sprintf("Post Error:%s", response))
}

func SentPlatformMerchantInviteMember(param map[string]string) error {
	response, err := utility.SendRequest(fmt.Sprintf("%s/cloud/email/invite_member", ApiHost), "POST", []byte(utility.MarshalToJsonString(param)), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" {
		return nil
	}
	return gerror.New(fmt.Sprintf("Post Error:%s", response))
}

func SentPlatformMerchantRegisterEmail(param map[string]string) error {
	response, err := utility.SendRequest(fmt.Sprintf("%s/cloud/email/merchant_register", ApiHost), "POST", []byte(utility.MarshalToJsonString(param)), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" {
		return nil
	}
	return gerror.New(fmt.Sprintf("Post Error:%s", response))
}

func FetchDefaultEmailTemplateFromPlatformApi() []*entity.EmailDefaultTemplate {
	var list = make([]*entity.EmailDefaultTemplate, 0)
	response, err := utility.SendRequest(fmt.Sprintf("%s/cloud/email/default_template_list", ApiHost), "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("emailTemplateList") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("emailTemplateList").String()), &list)
	}
	return list
}

func FetchColumnAppendListFromPlatformApi() []*bean.TableUpgrade {
	var list = make([]*bean.TableUpgrade, 0)
	var env = 1
	if config.GetConfigInstance().IsProd() {
		env = 2
	}
	response, err := sendRequestInMainCtxStart(fmt.Sprintf("%s/cloud/table/column_append?databaseType=%s&env=%v", ApiHost, g.DB("default").GetConfig().Type, env), "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("tableUpgrades") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("tableUpgrades").String()), &list)
	}
	return list
}

func sendRequestInMainCtxStart(url string, method string, data []byte, headers map[string]string) ([]byte, error) {
	bodyReader := bytes.NewReader(data)
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, gerror.NewCode(gcode.New(response.StatusCode, response.Status, response.Status+" "+string(responseBody)), response.Status+" "+string(responseBody))
	}
	return responseBody, nil
}
