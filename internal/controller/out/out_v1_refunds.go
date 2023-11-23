package out

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/out/v1"
	response "go-oversea-pay/utility"
)

func (c *ControllerV1) Refunds(ctx context.Context, req *v1.RefundsReq) (res *v1.RefundsRes, err error) {
	//g.RequestFromCtx(ctx).Response.WriteJson(response.JsonRes{
	//	Code:    200,
	//	Message: "success",
	//	Data:    "Hello Refund",
	//})
	//g.RequestFromCtx(ctx).Exit()
	//response.JsonExit(g.RequestFromCtx(ctx), 200, "success", "Hello Refund")
	response.SuccessJsonExit(g.RequestFromCtx(ctx), "Hello Refund")
	return
}
