package util

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
)

func GetPaymentRedirectUrl(ctx context.Context, payment *entity.Payment, success string) string {
	if success == "false" {
		var metadata = make(map[string]string)
		if len(payment.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifyPayment Unmarshal Metadata error:%s", err.Error())
			}
		}
		cancelUrl := metadata["CancelUrl"]
		if cancelUrl != "" && len(cancelUrl) > 0 {
			return cancelUrl
		} else {
			return payment.ReturnUrl
		}
	} else {
		return payment.ReturnUrl
	}
}
