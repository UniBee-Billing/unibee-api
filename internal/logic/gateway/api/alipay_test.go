package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"testing"
	_interface "unibee/internal/interface"
	defaultAlipayClient "unibee/internal/logic/gateway/api/alipay/api"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request/auth"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

var key = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmWsFjKU/GXuzW/1g3B/JTy0gvMpEQv8BODuiD1SxMUoUximCixWD951gw2NnTqh0NnaJwnJ7kManNJJ8gvIUzxgAeK9AIWYaHwzJ1pWVpJpkwsuE1KdduD2Ui1uROWgz1DNWILZqFYz3tUM3W2rxlZEwGqmZqWdNIfL7/1a5s0Pg1oQxz3Czj1MVD+sLN8YZ54TX4xD5hbAfMhtM56XBz2C7fZAgzneHO2wy2NMb8c2Cf8SNfGjcxbINbxDa+iOx3o+xket3DthQJcglKfiknXn+JaiwagOapjhkddwVutMPjReDZt/GESvWuHFaKtWxosvOIZD7I38CFHkHpYsHcwIDAQAB"
var secret = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCAE/NHRmtje56HLcgppwljqiLO5Kkh6lkXDOn4qqiZkSQap/u+EpROjV7M0EXMVwwVF59ZDevKCh+2nEEOup3gDJh25+2vs9JPZhv6/VqS6pw5NVXIrlhXWnm+AZ6zzZNiONT/lXSjpsxT35JvHvvRUGqXG6fHUTzeAY6lzg4vhP4qIw8d1UdbJnpthokwC6RMvHOfi6gYeUw8zoVzUTaM4LRp7UP0WYSYViWJb+bk5n99Ow5/ruZbSNvTJ4l5+3Q2GlBsjLdAftpXhBfrd2lMJF2znvkmWAFalD49d+Hd9Ia+6U8LsUuZdZjKcPO4GJ4pF1oFnrW3LE9h2AyxISMTAgMBAAECggEAXhjKMoJdIYDQDnanSVrMPingWuqKDD3VaGb3etc++Vw2D1N9U77osPGSRZ16uk71tIVfcBkXM5/OfuY7seuPU+1NEocBDIZrrCPTyMncgnXVgv5ZYRAeHUd+jAc6ptURRCeG7aPLRvSjx7dJKVS1I6oWNaB+2qQnuN+iAtTpfST5wCVke0y7s9tl3mZqp8EkHV2yXJqIYpKiUOYCuSMr99ARCLGmeRAD5w5cpTHeZoMzdLxoITFNvout/a4Is1NNGzIjMWQD9WOrhmE15cMukOW6jFKXkrQDI9cnUifkS/zmPFLrnHd4Zqo9SoGQIxEYZLQ1U12NbRtV1Ss9vZijOQKBgQD+l9uqUqT+fw11Mp9xaDgZzxt1rcd488qQtNZjbD0eoA2a1tn2g7YZeIA+Y9lPAXv99q7b6yFeeecJ/cOKOgD8wIjDigjPnL2wi9ZBVp1ZUXJSAQXaYZFgkiyOApoHdB62jdnuecP4G5dLs4+MSO/pciGk6/eoIhWrDJTdO3vH9wKBgQCAySBkNeK2SZ+2jtv259x+xmrR+7b2FSAal9wKclytoZHtUYlZ5hsURgkbBoZMpH5VJBoaveOOsocoKpOjn4LbYY9eTSCQN9yzF4JS5PXFoNmjQ6P2Ndcorogi2pPkyCj3nrwa+zc9zFzUzYOrlPyswEr+mNTtgfyBNRhhTOVOxQKBgQCoQVQ7TEMernkGa15UZLwu0mEjdKXPmc7Vs628J1x9UOms2zFRadp/GtQmZ3bGcASx4sXNMafr+ERopf0E7TCZ2eSI1kDcdIook0IWDFgRH3KeH27u1Gxvlis77xw8sNFbdIQCxxZsck+bCCBmZg2oCnWRuSEDTQNk9/up+hXkIQKBgG7QoZaY52ODJnKnqo5iJFDR2sikl2JX+y/my+gRT7338OEL7+vzHAnt2ZfvnVAFms8YKX4pNs1qwPHG8RMyBh9Pa1Xxd7ug1b8k03cQnIpZRew+H6+T1Hek9m9HNUr/EIFBjQqKb5Y1awuRa2MQ5/qd2+oHB/D2kJd9YGUZDZchAoGBAO4IaKpc/iiqNV8ZZuSVnZT4HMoyszJ0q86wmKcITgc5qhYbbgzXCChPWnKSLqBTeWrqmXKJkqhMT9TunDu0Xvu2OPg0xzOWl7GjGBKNRXElqgEHilGK99no/5cK/Vww1nC9x9hwpJDgSdTkvi6mTv89M2SbErJdBydXOHswGEKZ"
var subGateway = "SANDBOX_5YES442ZS5S203863"

func TestForAlipay(t *testing.T) {
	pay := &Alipay{}
	_, _, _ = pay.GatewayTest(context.Background(), &_interface.GatewayTestReq{
		Key:                 key,
		Secret:              secret,
		SubGateway:          subGateway,
		GatewayPaymentTypes: nil,
	})
}

func TestForAlipayPlus(t *testing.T) {
	pay := &AlipayPlus{}
	_, _, _ = pay.GatewayTest(context.Background(), &_interface.GatewayTestReq{
		Key:        key,
		Secret:     secret,
		SubGateway: subGateway,
		GatewayPaymentTypes: []*_interface.GatewayPaymentType{&_interface.GatewayPaymentType{
			Name:        "",
			PaymentType: "ALIPAY_CN",
			CountryName: "",
			AutoCharge:  false,
			Category:    "",
		}},
	})
}

func TestForAlipayGetPaymentDetail(t *testing.T) {
	pay := &Alipay{}
	res, err := pay.GatewayPaymentDetail(context.Background(), &entity.MerchantGateway{
		Host:          "https://open-de-global.alipay.com",
		GatewayKey:    key,
		GatewaySecret: secret,
		SubGateway:    subGateway,
	}, "202502151640109000001906B0208391400", &entity.Payment{
		PaymentId: "pay202502152PgZfEdH8wyayLu",
	})
	if err != nil {
		g.Log().Errorf(context.Background(), "error:%s", err.Error())
	} else {
		fmt.Printf("%s", utility.MarshalToJsonString(res))
	}
}

func TestForPaymentTypes(t *testing.T) {
	fmt.Printf("%s", utility.MarshalToJsonString(fetchAlipayPlusPaymentTypes(context.Background())))
}

func TestAutochargeAudit(t *testing.T) {
	client := defaultAlipayClient.NewDefaultAlipayClient(
		"https://open-de-global.alipay.com",
		subGateway,
		secret,
		key, false)
	request, authConsultRequest := auth.NewAlipayAuthConsultRequest()
	authConsultRequest.CustomerBelongsTo = model.ALIPAY_CN
	authConsultRequest.AuthRedirectUrl = "https://www.yourRedirectUrl.com"
	authConsultRequest.Scopes = []model.ScopeType{model.ScopeTypeAgreementPay}
	authConsultRequest.AuthState = uuid.NewString()
	authConsultRequest.TerminalType = model.WEB

	execute, err := client.Execute(request)
	if err != nil {
		print(err.Error())
		return
	}
	response := execute.(*responseAuth.AlipayAuthConsultResponse)
	fmt.Println("response: ", utility.MarshalToJsonString(response))
	fmt.Println("NormalUrl: ", response.NormalUrl)
	//https://www.yourredirecturl.com/?authCode=281001139369360449491196&authState=6f9c3c7b-8619-4a7a-aa63-96db0239c47d
}
