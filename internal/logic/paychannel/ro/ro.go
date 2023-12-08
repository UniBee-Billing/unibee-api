package ro

import "go-oversea-pay/api/out/vo"

// OverseaPay is the golang structure for table oversea_pay.
type OutPayCaptureRo struct {
	MerchantId   string         `json:"merchantId"         `      // 商户ID
	PspReference string         `json:"pspReference"            ` // 业务类型。1-订单
	Reference    string         `json:"reference"              `  // 业务id-即商户订单号
	Amount       vo.PayAmountVo `json:"amount"`
	Status       string         `json:"status"`
}
