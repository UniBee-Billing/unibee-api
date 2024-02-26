package onetime

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/onetime/mock"
)

func (c *ControllerMock) DetailPay(ctx context.Context, req *mock.DetailPayReq) (res *mock.DetailPayRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
