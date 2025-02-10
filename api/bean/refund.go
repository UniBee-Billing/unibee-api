package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Refund struct {
	MerchantId              uint64                   `json:"merchantId"           description:"merchant id"`                                        // merchant id
	UserId                  uint64                   `json:"userId"               description:"user_id"`                                            // user_id
	GatewayId               uint64                   `json:"gatewayId"            description:"gateway_id"`                                         // gateway_id
	ExternalRefundId        string                   `json:"externalRefundId"     description:"external_refund_id"`                                 // external_refund_id
	CountryCode             string                   `json:"countryCode"          description:"country code"`                                       // country code
	Currency                string                   `json:"currency"             description:"currency"`                                           // currency
	PaymentId               string                   `json:"paymentId"            description:"relative payment id"`                                // relative payment id
	RefundId                string                   `json:"refundId"             description:"refund id (system generate)"`                        // refund id (system generate)
	RefundAmount            int64                    `json:"refundAmount"         description:"refund amount, cent"`                                // refund amount, cent
	RefundComment           string                   `json:"refundComment"        description:"refund comment"`                                     // refund comment
	RefundCommentExplain    string                   `json:"refundCommentExplain"  description:"refund comment"`                                    // refund comment
	Status                  int                      `json:"status"               description:"status。10-pending，20-success，30-failure, 40-cancel"` // status。10-pending，20-success，30-failure, 40-cancel
	RefundTime              int64                    `json:"refundTime"           description:"refund success time"`                                // refund success time
	GatewayRefundId         string                   `json:"gatewayRefundId"      description:"gateway refund id"`                                  // gateway refund id
	ReturnUrl               string                   `json:"returnUrl"            description:"return url after refund success"`                    // return url after refund success
	SubscriptionId          string                   `json:"subscriptionId"       description:"subscription id"`                                    // subscription id
	CreateTime              int64                    `json:"createTime"           description:"create utc time"`                                    // create utc time
	Metadata                map[string]interface{}   `json:"metadata" description:""`
	Type                    int                      `json:"type"                 description:"1-gateway refund,2-mark refund"` // 1-gateway refund,2-mark refund
	InvoiceId               string                   `json:"invoiceId"            description:"invoice id"`                     // invoice id
	GatewayCurrencyExchange *GatewayCurrencyExchange `json:"gatewayCurrencyExchange"  description:"gateway currency exchange"`
}

func SimplifyRefund(one *entity.Refund) *Refund {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	var gatewayExchange *GatewayCurrencyExchange
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyRefund Unmarshal Metadata error:%s", err.Error())
		}
	}
	if s, ok := metadata[GatewayCurrencyExchangeKey]; ok {
		if value, ok2 := s.(string); ok2 {
			if len(value) > 0 {
				_ = utility.UnmarshalFromJsonString(value, &gatewayExchange)
			}
		}
	}
	if gatewayExchange != nil && len(gatewayExchange.FromCurrency) > 0 {
		gatewayExchange.ExchangeAmount = utility.ExchangeCurrencyConvert(one.RefundAmount, gatewayExchange.FromCurrency, gatewayExchange.ToCurrency, gatewayExchange.ExchangeRate)
	}
	return &Refund{
		MerchantId:              one.MerchantId,
		UserId:                  one.UserId,
		GatewayId:               one.GatewayId,
		ExternalRefundId:        one.ExternalRefundId,
		CountryCode:             one.CountryCode,
		Currency:                one.Currency,
		PaymentId:               one.PaymentId,
		RefundId:                one.RefundId,
		RefundAmount:            one.RefundAmount,
		RefundComment:           one.RefundComment,
		RefundCommentExplain:    one.RefundCommentExplain,
		Status:                  one.Status,
		RefundTime:              one.RefundTime,
		GatewayRefundId:         one.GatewayRefundId,
		ReturnUrl:               one.ReturnUrl,
		SubscriptionId:          one.SubscriptionId,
		CreateTime:              one.CreateTime,
		Metadata:                metadata,
		Type:                    one.Type,
		InvoiceId:               one.InvoiceId,
		GatewayCurrencyExchange: gatewayExchange,
	}
}

func SimplifyRefundList(ones []*entity.Refund) (list []*Refund) {
	if len(ones) == 0 {
		return make([]*Refund, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyRefund(one))
	}
	return list
}
