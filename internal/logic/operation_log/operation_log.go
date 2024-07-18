package operation_log

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/analysis/segment"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type OptLogRequest struct {
	MerchantId     uint64
	Target         string
	Content        string
	UserId         uint64
	SubscriptionId string
	InvoiceId      string
	PlanId         uint64
	DiscountCode   string
}

func AppendOptLog(superCtx context.Context, req *OptLogRequest, optError error) {
	var merchantId = req.MerchantId
	if merchantId <= 0 {
		g.Log().Errorf(superCtx, "AppendOptLog error invalid merchantId:%v", merchantId)
		return
	}
	if optError != nil {
		g.Log().Infof(superCtx, "AppendOptLog hasError skip")
		return
	}
	var memberId uint64 = 0
	var optAccount = ""
	if _interface.Context().Get(superCtx) != nil && _interface.Context().Get(superCtx).MerchantMember != nil {
		memberId = _interface.Context().Get(superCtx).MerchantMember.Id
		optAccount = fmt.Sprintf("Member(%v)", memberId)
	} else if _interface.Context().Get(superCtx) != nil && _interface.Context().Get(superCtx).IsOpenApiCall {
		memberId = 0
		optAccount = fmt.Sprintf("OpenApi(%v)", _interface.Context().Get(superCtx).OpenApiKey)
		var targetUserId uint64 = 0
		if req.UserId > 0 {
			targetUserId = req.UserId
		} else if len(req.SubscriptionId) > 0 {
			sub := query.GetSubscriptionBySubscriptionId(superCtx, req.SubscriptionId)
			if sub != nil {
				targetUserId = sub.UserId
			}
		} else if len(req.InvoiceId) > 0 {
			in := query.GetInvoiceByInvoiceId(superCtx, req.InvoiceId)
			if in != nil {
				targetUserId = in.UserId
			}
		}
		if targetUserId > 0 {
			userAccount := query.GetUserAccountById(superCtx, targetUserId)
			if userAccount != nil {
				segment.TrackSegmentEventBackground(superCtx, userAccount.MerchantId, userAccount, req.Target, map[string]interface{}{
					"OptAccount": optAccount,
					"OptTarget":  req.Target,
					"OptContent": req.Content,
				})
			}
		}
	} else {
		memberId = 0
		optAccount = fmt.Sprintf("System")
	}
	operationLog := &entity.MerchantOperationLog{
		CompanyId:      0,
		MerchantId:     merchantId,
		MemberId:       memberId,
		OptAccount:     optAccount,
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
	if memberId <= 0 {
		if optAccount == "System" {
			g.Log().Debugf(superCtx, "Receive SystemOperationLog:%s", utility.MarshalToJsonString(operationLog))
			return
		} else {
			g.Log().Infof(superCtx, "Receive OpenApiOperation:%s", utility.MarshalToJsonString(operationLog))
		}
	}
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

		_, err = dao.MerchantOperationLog.Ctx(ctx).Data(operationLog).OmitNil().Insert(operationLog)
		if err != nil {
			g.Log().Errorf(ctx, "AppendOptLog Error %s", err.Error())
		}
	}()
}
