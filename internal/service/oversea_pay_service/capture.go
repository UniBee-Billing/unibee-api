package oversea_pay_service

import (
	"go-oversea-pay/internal/consts"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func DoChannelCapture(overseaPay *entity.OverseaPay) (result interface{}, err error) {
	utility.Assert(overseaPay != nil, "entity not found")
	utility.Assert(overseaPay.PayStatus == consts.TO_BE_PAID, "payment not waiting for pay")
	utility.Assert(overseaPay.AuthorizeStatus != consts.WAITING_AUTHORIZED, "payment not authorised")
	utility.Assert(overseaPay.BuyerPayFee > 0, "capture value should > 0")
	utility.Assert(overseaPay.BuyerPayFee <= overseaPay.PaymentFee, "capture value should <= authorized value")

	// todo mark 实现 channel capture
	return overseaPay, nil
}
