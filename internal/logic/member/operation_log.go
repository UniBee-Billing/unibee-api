package member

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
)

type OptLogRequest struct {
	Target         string
	Content        string
	UserId         uint64
	SubscriptionId string
	InvoiceId      string
	PlanId         uint64
	DiscountCode   string
}

func AppendOptLog(superCtx context.Context, req *OptLogRequest) {
	member := _interface.Context().Get(superCtx).MerchantMember
	if member != nil {
		go func() {
			ctx := context.Background()
			var err error
			defer func() {
				if exception := recover(); exception != nil {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						err = v
					} else {
						err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
					}
					log.PrintPanic(ctx, err)
					return
				}
			}()
			operationLog := &entity.MerchantOperationLog{
				CompanyId:      0,
				MerchantId:     member.MerchantId,
				MemberId:       member.Id,
				OptAccount:     strconv.FormatUint(member.Id, 10),
				ClientType:     0,
				BizType:        0,
				OptTarget:      req.Target,
				OptContent:     req.Content,
				CreateTime:     gtime.Now().Timestamp(),
				GmtCreate:      gtime.Now(),
				GmtModify:      gtime.Now(),
				ServerType:     0,
				ServerTypeDesc: "",
				SubscriptionId: req.SubscriptionId,
				UserId:         req.UserId,
				InvoiceId:      req.InvoiceId,
				PlanId:         req.PlanId,
				DiscountCode:   req.DiscountCode,
			}
			_, err = dao.MerchantOperationLog.Ctx(ctx).Data(operationLog).OmitNil().Insert(operationLog)
			if err != nil {
				g.Log().Errorf(ctx, "AppendOptLog Error %s", err.Error())
			}
		}()
	}
}
