package link

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/util"
	session2 "unibee/internal/logic/session"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

type LinkCheckRes struct {
	Message string
	Link    string
	Payment *entity.Payment
}

func LinkCheck(ctx context.Context, paymentId string, time int64) *LinkCheckRes {
	var res = &LinkCheckRes{
		Message: "",
		Link:    "",
		Payment: nil,
	}
	one := query.GetPaymentByPaymentId(ctx, paymentId)
	if one == nil {
		g.Log().Errorf(ctx, "LinkEntry payment not found paymentId: %s", paymentId)
		res.Message = "Payment Not Found"
		return res
	}
	res.Payment = one
	if one.Status == consts.PaymentCancelled {
		res.Message = "Payment Cancelled"
	} else if one.Status == consts.PaymentFailed {
		res.Message = "Payment Failure"
	} else if one.Status == consts.PaymentSuccess {
		res.Message = "Payment Already Success"
	} else if one.ExpireTime != 0 && one.ExpireTime < time {
		res.Message = "Payment Expired"
	} else if len(one.GatewayLink) > 0 {
		res.Link = one.GatewayLink
	} else if strings.Contains(one.Link, "unibee.top") || strings.Contains(one.Link, "unibee.dev") {
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		if gateway != nil && gateway.GatewayType == consts.GatewayTypeWireTransfer {
			if len(config.GetConfigInstance().Server.HostedPagePath) > 0 {
				targetUrl := util.GetPaymentRedirectUrl(ctx, one, "true")
				cancelUrl := util.GetPaymentRedirectUrl(ctx, one, "false")
				_, userSession, err := session2.NewUserSession(ctx, one.MerchantId, one.UserId, targetUrl, cancelUrl)
				if err == nil && len(userSession) > 0 {
					res.Link = fmt.Sprintf("%s/payment_checker?merchantId=%d&paymentId=%s&session=%s&env=%s", config.GetConfigInstance().Server.HostedPagePath, one.MerchantId, one.PaymentId, userSession, config.GetConfigInstance().Env)
				}
			}
			if len(res.Link) == 0 {
				res.Message = "Please finish your wire transfer payment"
			}
		} else {
			res.Message = "Payment Link Error"
		}
	} else {
		res.Link = one.Link // old version
	}
	return res
}
