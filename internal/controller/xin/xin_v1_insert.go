package xin

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/xin/v1"
	"go-oversea-pay/internal/service/xin_service"
)

func (c *ControllerV1) Insert(ctx context.Context, req *v1.InsertReq) (res *v1.InsertRes, err error) {

	test, err := xin_service.InsertTest(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	g.RequestFromCtx(ctx).Response.Writeln(test)
	return
}
