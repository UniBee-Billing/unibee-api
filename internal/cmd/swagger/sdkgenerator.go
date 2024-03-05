package swagger

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
	"unibee/internal/consts"
	"unibee/utility"
)

func SDKGeneratorJson(r *ghttp.Request) {
	url := fmt.Sprintf("http://127.0.0.1%s/%s", consts.GetConfigInstance().Server.Address, consts.GetConfigInstance().Server.OpenApiPath)
	request, err := utility.SendRequest(url, "GET", nil, nil)
	if err != nil {
		r.Exit()
		return
	} else {
		api := gjson.New(string(request))
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
		r.Response.WriteJson(api.String())
		r.Exit()
	}
}
