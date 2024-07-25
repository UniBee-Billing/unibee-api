package api

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"time"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

var GatewayNameMapping = map[string]_interface.GatewayInterface{
	"stripe":          &Stripe{},
	"changelly":       &Changelly{},
	"paypal":          &Paypal{},
	"invalid":         &Invalid{},
	"autotest_crypto": &AutoTestCrypto{},
	"autotest":        &AutoTest{},
	"coinbase":        &Coinbase{},
	"wire_transfer":   &Wire{},
}

var GatewayShortNameMapping = map[string]string{
	"stripe":          "ST",
	"changelly":       "CT",
	"paypal":          "PP",
	"invalid":         "IP",
	"autotest_crypto": "AC",
	"autotest":        "AP",
	"coinbase":        "CP",
	"wire_transfer":   "WT",
}

type GatewayProxy struct {
	Gateway     *entity.MerchantGateway
	GatewayName string
}

func (p GatewayProxy) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	to, err = p.getRemoteGateway().GatewayCryptoFiatTrans(ctx, from)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayCryptoFiatTrans cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return to, err
}

func (p GatewayProxy) getRemoteGateway() (one _interface.GatewayInterface) {
	utility.Assert(len(p.GatewayName) > 0, "gateway is not set")
	one = GatewayNameMapping[p.GatewayName]
	utility.Assert(one != nil, "gateway not support:"+p.GatewayName+" should be stripe|paypal|changelly|wire_transfer")
	return
}

func (p GatewayProxy) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayUserCreateAndBindPaymentMethod(ctx, gateway, userId, currency, metadata)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserCreateAndBindPaymentMethod cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	//defer func() {
	//	if exception := recover(); exception != nil {
	//		if v, ok := exception.(error); ok && gerror.HasStack(v) {
	//			err = v
	//		} else {
	//			err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
	//		}
	//		printChannelPanic(ctx, err)
	//		return
	//	}
	//}()
	//startTime := time.Now()
	icon, gatewayType, err = p.getRemoteGateway().GatewayTest(ctx, key, secret)

	//glog.Infof(ctx, "MeasureChannelFunction:GatewayTest cost：%s \n", time.Now().Sub(startTime))
	//if err != nil {
	//	err = gerror.NewCode(utility.GatewayError, err.Error())
	//}
	return icon, gatewayType, err
}

func (p GatewayProxy) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayUserAttachPaymentMethodQuery(ctx, gateway, userId, gatewayPaymentMethod)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayUserDeAttachPaymentMethodQuery(ctx, gateway, userId, gatewayPaymentMethod)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()

	res, err = p.getRemoteGateway().GatewayUserPaymentMethodListQuery(ctx, gateway, req)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserPaymentMethodListQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()

	res, err = p.getRemoteGateway().GatewayUserCreate(ctx, gateway, user)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserCreate cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()

	res, err = p.getRemoteGateway().GatewayMerchantBalancesQuery(ctx, gateway)

	glog.Infof(ctx, "MeasureChannelFunction:GatewayMerchantBalancesQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayUserDetailQuery(ctx, gateway, userId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayUserDetailQuery cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayNewPayment(ctx, createPayContext)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayNewPayment cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayCapture(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayCapture(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayCapture cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayCancel(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayCancel(ctx, pay)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayPaymentList(ctx, gateway, listReq)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPaymentList cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayPaymentDetail(ctx, gateway, gatewayPaymentId, payment)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayPaymentDetail cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayRefund(ctx, pay, refund)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefund cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayRefundDetail(ctx, gateway, gatewayRefundId, refund)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundDetail cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayRefundCancel(ctx, payment, refund)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundCancel cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func (p GatewayProxy) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			printChannelPanic(ctx, err)
			return
		}
	}()
	startTime := time.Now()
	res, err = p.getRemoteGateway().GatewayRefundList(ctx, gateway, gatewayPaymentId)
	glog.Infof(ctx, "MeasureChannelFunction:GatewayRefundList cost：%s \n", time.Now().Sub(startTime))
	if err != nil {
		err = gerror.NewCode(utility.GatewayError, err.Error())
	}
	return res, err
}

func printChannelPanic(ctx context.Context, err error) {
	var requestId = "init"
	if _interface.Context().Get(ctx) != nil {
		requestId = _interface.Context().Get(ctx).RequestId
	}
	g.Log().Errorf(ctx, "ChannelException panic requestId:%s error:%s", requestId, err.Error())
}
