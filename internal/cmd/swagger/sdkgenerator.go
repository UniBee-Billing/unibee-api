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
		if r.Get("hideSecurity") != nil {
			_ = api.Remove("security")
			_ = api.Remove("components#securitySchemes")
		}
		for key, path := range api.GetJsonMap("components#schemas") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if strings.HasPrefix(key, "unibee.api.user") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
			if strings.HasPrefix(key, "unibee.api.merchant.oss") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
			if strings.HasPrefix(key, "unibee.api.merchant.task") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/merchant")) {
				_ = api.Remove("paths#" + key)
				continue
			}
			if strings.HasPrefix(key, fmt.Sprintf("/merchant/task")) {
				_ = api.Remove("paths#" + key)
				continue
			}
			if strings.HasPrefix(key, fmt.Sprintf("/merchant/oss")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}

		// generator error to format type of map[string]interface {}
		response := api.String()
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
		if r.Get("hideSecurity") != nil {
			_ = api.Remove("security")
			_ = api.Remove("components#securitySchemes")
		}
		for key, path := range api.GetJsonMap("components#schemas") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if strings.HasPrefix(key, "unibee.api.user") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
			if strings.HasPrefix(key, "unibee.api.merchant.oss") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
			if strings.HasPrefix(key, "unibee.api.merchant.task") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/merchant")) {
				_ = api.Remove("paths#" + key)
				continue
			}
			if strings.HasPrefix(key, fmt.Sprintf("/merchant/task")) {
				_ = api.Remove("paths#" + key)
				continue
			}
			if strings.HasPrefix(key, fmt.Sprintf("/merchant/oss")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}
		if r.Get("hideDescription") != nil {
			err = api.Set("info#description#$ref", "description.md")
			if err != nil {
				r.Exit()
				return
			}
		}
		apiYaml, err := yaml.JSONToYAML([]byte(api.String()))
		if err != nil {
			r.Exit()
			return
		}
		r.Response.WriteJson(apiYaml)
		r.Exit()
	}
}

func SystemGeneratorSpecJson(r *ghttp.Request) {
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
			if strings.HasPrefix(key, "unibee.api.user") {
				_ = api.Remove("components#schemas#" + key)
				continue
			}
		}
		for key, path := range api.GetJsonMap("paths") {
			utility.Assert(len(path.Array()) == 1, "error:"+key)
			if !strings.HasPrefix(key, fmt.Sprintf("/system")) {
				_ = api.Remove("paths#" + key)
				continue
			}
		}
		r.Response.WriteJson(api.String())
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
