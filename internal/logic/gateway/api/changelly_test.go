package api

import (
	"context"
	"fmt"
	"testing"
	"time"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func TestForGetCurrencyProviderList(t *testing.T) {
	ctx := context.Background()
	changelly := Changelly{}
	trans, err := changelly.GatewayCryptoFiatTrans(ctx, &gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq{
		Amount:      106,
		Currency:    "USD",
		CountryCode: "FR",
		Gateway: &entity.MerchantGateway{
			Id:                    29,
			MerchantId:            0,
			EnumKey:               0,
			GatewayType:           0,
			GatewayName:           "",
			Name:                  "",
			SubGateway:            "",
			BrandData:             "",
			Logo:                  "",
			Host:                  "",
			GatewayAccountId:      "",
			GatewayKey:            "",
			GatewaySecret:         "",
			Custom:                "",
			GmtCreate:             nil,
			GmtModify:             nil,
			Description:           "",
			WebhookKey:            "",
			WebhookSecret:         "",
			UniqueProductId:       "",
			CreateTime:            0,
			IsDeleted:             0,
			CryptoReceiveCurrency: "",
		},
	})
	if err != nil {
		return
	}
	fmt.Println(utility.MarshalToJsonString(trans))
}

func TestForTimeFormat(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05.876Z"))
	fmt.Println(fmt.Sprintf("PDF Generated on %s", time.Now().Format(time.RFC850)))
	fmt.Println(utility.ConvertCentToDollarStr(108, "USDT"))
}
