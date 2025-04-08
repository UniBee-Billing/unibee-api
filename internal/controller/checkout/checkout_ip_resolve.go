package checkout

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/utility"

	"unibee/api/checkout/ip"
)

func (c *ControllerIp) Resolve(ctx context.Context, req *ip.ResolveReq) (res *ip.ResolveRes, err error) {
	body := fmt.Sprintf("{\"ip\":\"%s\"}", req.IP)
	response, err := utility.SendRequest("https://app.multiloginapp.com/resolve", "POST", []byte(body), map[string]string{"Content-Type": "application/json"})
	utility.AssertError(err, "Resolve Error")
	return &ip.ResolveRes{Location: gjson.New(string(response))}, nil
}
