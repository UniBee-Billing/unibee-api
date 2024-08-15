package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) CustomizeLocalizationTemplateSync(ctx context.Context, req *email.CustomizeLocalizationTemplateSyncReq) (res *email.CustomizeLocalizationTemplateSyncRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
