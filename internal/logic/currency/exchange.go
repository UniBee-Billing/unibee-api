package currency

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"strings"
	"unibee/utility"
)

// https://www.exchangerate-api.com/docs/overview
func GetExchangeConversionRates(ctx context.Context, key string, currency string, to string) (*float64, error) {
	redisCacheKey := fmt.Sprintf("FiatExchageConversionRate_%s", currency)
	value, err := g.Redis().Get(ctx, redisCacheKey)
	if err == nil && !value.IsNil() && len(value.String()) > 0 {
		var currencyRates = make(map[string]float64)
		err = gjson.Unmarshal([]byte(value.String()), &currencyRates)
		if err != nil {
			fmt.Printf("GetConversionRates Unmarshal error:%s", err.Error())
		}
		v, ok := currencyRates[strings.ToUpper(to)]
		if ok {
			return &v, nil
		}
	}
	response, err := utility.SendRequest(fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", key, strings.ToUpper(currency)), "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	if responseJson.Contains("result") && responseJson.Get("result").String() == "success" {
		conversionRates := responseJson.GetJsonMap("conversion_rates")
		_, _ = g.Redis().Set(ctx, redisCacheKey, utility.MarshalToJsonString(conversionRates))
		_, _ = g.Redis().Expire(ctx, redisCacheKey, 3600*12)
		v, ok := conversionRates[strings.ToUpper(to)]
		if ok {
			rateFloat, err := strconv.ParseFloat(v.String(), 64)
			if err == nil {
				return &rateFloat, nil
			} else {
				return nil, err
			}
		} else {
			return nil, gerror.New("not found")
		}
	} else {
		return nil, gerror.New(responseJson.String())
	}
}
