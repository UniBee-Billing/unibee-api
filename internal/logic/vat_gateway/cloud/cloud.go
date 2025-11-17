package cloud

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	go_redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetCloudVatCountryList(ctx context.Context, merchantId uint64) []*entity.CountryRate {
	res := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:   "GID_UniBee_Cloud",
		Method:  "GetVatCountryList",
		Request: "",
	}, 0)
	if res == nil || !res.Status {
		return make([]*entity.CountryRate, 0)
	}
	var list = make([]*entity.CountryRate, 0)
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(res.Response), &list)
	if err != nil {
		return make([]*entity.CountryRate, 0)
	}
	return list
}

func GetCloudVatCountryListByCountryCode(ctx context.Context, merchantId uint64, countryCode string) *entity.CountryRate {
	res := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:   "GID_UniBee_Cloud",
		Method:  "GetVatCountryListByCountryCode",
		Request: countryCode,
	}, 0)
	if res == nil || !res.Status {
		return nil
	}
	var one *entity.CountryRate
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(res.Response), &one)
	if err != nil {
		return nil
	}
	return one
}

func ValidateVatNumberFromCloud(ctx context.Context, vatNumber string, requesterVatNumber string) (*bean.ValidResult, error) {
	res := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:  "GID_UniBee_Cloud",
		Method: "VatNumberValidate",
		Request: utility.MarshalToJsonString(&NumberValidateReq{
			VatNumber:        vatNumber,
			RequestVatNumber: requesterVatNumber,
		}),
	}, 0)
	if res == nil {
		return nil, gerror.New("System Error")
	}
	if !res.Status {
		return nil, gerror.New(fmt.Sprintf("%s", res.Response))
	}
	var one *bean.ValidResult
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(res.Response), &one)
	if err != nil {
		return nil, err
	}
	return one, nil
}

type NumberValidateReq struct {
	VatNumber        string `json:"vatNumber"`
	RequestVatNumber string `json:"requestVatNumber"`
}
