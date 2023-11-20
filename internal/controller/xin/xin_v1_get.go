package xin

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/xin/v1"
	"go-oversea-pay/internal/service/xin_service"
)

func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	test, err := xin_service.QueryTest(ctx)
	if err != nil {
		return nil, err
	}
	g.RequestFromCtx(ctx).Response.Writeln(test)
	return
}
