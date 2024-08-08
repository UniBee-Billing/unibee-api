package swagger

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/utility"
)

func MerchantPortalAndSDKGeneratorSpecJson(r *ghttp.Request) {
	url := fmt.Sprintf("http://127.0.0.1%s/%s", config.GetConfigInstance().Server.Address, config.GetConfigInstance().Server.OpenApiPath)
	request, err := utility.SendRequest(url, "GET", nil, nil)
	if err != nil {
		r.Exit()
		return
	} else {
		json := strings.Replace(string(request), "uint64", "int64", -1)
		api := gjson.New(json)
		api.SetSplitChar('#')
		for key, path := range api.GetJsonMap("components#schemas") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if strings.HasPrefix(key, "unibee.api.user") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			//if path.Contains("post") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("post"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			//if path.Contains("get") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("get"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/merchant")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}
		// generator error to format type of map[string]interface {}
		response := api.String()
		//response := strings.Replace(api.String(), "map[string]interface {}", "interface {}", -1)
		//mapTarget := `"additionalProperties":{"$ref":"#/components/schemas/interface"},`
		//mapReplace := `"additionalProperties":{"format":"string","properties":{},"type":"string"},` // If generate map[string]interface{}, leave blank
		//response = strings.Replace(response, mapTarget, mapReplace, -1)
		r.Response.WriteJson(response)
		r.Exit()
	}
}

func MerchantPortalAndSDKGeneratorSpecYaml(r *ghttp.Request) {
	url := fmt.Sprintf("http://127.0.0.1%s/%s", config.GetConfigInstance().Server.Address, config.GetConfigInstance().Server.OpenApiPath)
	request, err := utility.SendRequest(url, "GET", nil, nil)
	if err != nil {
		r.Exit()
		return
	} else {
		json := strings.Replace(string(request), "uint64", "int64", -1)
		api := gjson.New(json)
		api.SetSplitChar('#')
		for key, path := range api.GetJsonMap("components#schemas") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if strings.HasPrefix(key, "unibee.api.user") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			//if path.Contains("post") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("post"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			//if path.Contains("get") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("get"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/merchant")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}
		err = api.Set("info#description#$ref", "description.md")
		if err != nil {
			r.Exit()
			return
		}
		apiYaml, err := yaml.JSONToYAML([]byte(api.String()))
		if err != nil {
			r.Exit()
			return
		}
		//// generator error to format type of map[string]interface {}
		//response := strings.Replace(api.String(), "map[string]interface {}", "map[string]string", -1)
		//mapTarget := `"additionalProperties":{"$ref":"#/components/schemas/interface"},`
		//mapReplace := `"additionalProperties":{"format":"string","properties":{},"type":"string"},`
		//response = strings.Replace(response, mapTarget, mapReplace, -1)
		//r.Response.WriteJson(response)
		r.Response.WriteJson(apiYaml)
		r.Exit()
	}
}

func UserPortalGeneratorSpecJson(r *ghttp.Request) {
	url := fmt.Sprintf("http://127.0.0.1%s/%s", config.GetConfigInstance().Server.Address, config.GetConfigInstance().Server.OpenApiPath)
	request, err := utility.SendRequest(url, "GET", nil, nil)
	if err != nil {
		r.Exit()
		return
	} else {
		json := strings.Replace(string(request), "uint64", "int64", -1)
		api := gjson.New(json)
		api.SetSplitChar('#')
		for key, path := range api.GetJsonMap("components#schemas") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if strings.HasPrefix(key, "unibee.api.merchant") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			//if path.Contains("post") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("post"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			//if path.Contains("get") {
			//	_ = api.Set("paths#"+key+"#api", path.GetJson("get"))
			//	_ = api.Remove("paths#" + key + "#post")
			//}
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/user")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}
		r.Response.WriteJson(api.String())
		r.Exit()
	}
}
