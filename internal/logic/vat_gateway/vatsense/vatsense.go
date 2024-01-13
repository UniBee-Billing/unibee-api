package vatsense

import (
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/internal/logic/vat_gateway/base"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"io"
	"net/http"
)

type VatSense struct {
	Name     string
	Password string
}

func fetchVatSense(url string, passwd string) (*gjson.Json, error) {
	// 用户名和密码
	username := "user"
	password := passwd

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("FetchVatSense Error creating request:", err)
		return nil, err
	}

	// 添加 HTTP 基本认证头部
	auth := username + ":" + password
	authEncoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", authEncoded)

	// 发起 HTTP 请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("FetchVatSense Error making request:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FetchVatSense Error reading response body:", err)
		return nil, err
	}

	// 打印响应内容
	fmt.Printf("FetchVatSense Response:%s\n", string(body))

	return gjson.LoadJson(string(body))
}

func (c VatSense) SetVatName(name string) {
	c.Name = name
}

func (c VatSense) GetVatName() string {
	return c.Name
}

func (c VatSense) SetVatSignData(data string) {
	c.Password = data
}

func (c VatSense) ListAllCountries() ([]*entity.CountryRate, error) {
	data, err := fetchVatSense("https://api.vatsense.com/1.0/countries", c.Password)
	if err != nil {
		return nil, err
	}
	if data.Get("code").Int() != 200 {
		var countryRates []*entity.CountryRate
		for _, item := range data.GetJsons("data") {
			var vat = 0
			if item.Get("vat") != nil && item.Get("vat").Bool() {
				vat = 1
			} else {
				vat = 2
			}
			var eu = 0
			if item.Get("eu") != nil && item.Get("eu").Bool() {
				eu = 1
			} else {
				eu = 2
			}
			countryRates = append(countryRates, &entity.CountryRate{
				VatName:     c.GetVatName(),
				CountryCode: item.Get("country_code").String(),
				CountryName: item.Get("country_name").String(),
				Latitude:    item.Get("latitude").String(),
				Longitude:   item.Get("longitude").String(),
				Vat:         vat,
				Eu:          eu,
				Provinces:   item.Get("provinces").String(),
			})
		}
		return countryRates, nil
	} else {
		return nil, gerror.New(data.String())
	}
}

func (c VatSense) ListAllRates() ([]*entity.CountryRate, error) {
	response, err := fetchVatSense("https://api.vatsense.com/1.0/rates", c.Password)
	if err != nil {
		return nil, err
	}
	if response.Get("code").Int() == 200 {
		var countryRates []*entity.CountryRate
		for _, item := range response.GetJsons("data") {
			standard := item.GetJson("standard")
			var eu = 0
			if item.Get("eu") != nil && item.Get("eu").Bool() {
				eu = 1
			} else {
				eu = 2
			}
			countryRates = append(countryRates, &entity.CountryRate{
				VatName:               c.GetVatName(),
				CountryCode:           item.Get("country_code").String(),
				CountryName:           item.Get("country_name").String(),
				Latitude:              item.Get("latitude").String(),
				Longitude:             item.Get("longitude").String(),
				Eu:                    eu,
				Provinces:             item.Get("provinces").String(),
				StandardTypes:         standard.Get("types").String(),
				StandardDescription:   standard.Get("description").String(),
				StandardTaxPencentage: standard.Get("rate").Int64(),
				Other:                 item.GetJson("other").String(),
			})
		}
		return countryRates, nil
	} else {
		return nil, gerror.New(response.String())
	}
}

func (c VatSense) ValidateVatNumber(varNumber string, requesterVatNumber string) (*base.ValidResult, error) {
	var response *gjson.Json
	var err error
	if len(requesterVatNumber) > 0 {
		response, err = fetchVatSense("https://api.vatsense.com/1.0/validate?vat_number="+varNumber+"&requester_vat_number="+requesterVatNumber, c.Password)
	} else {
		response, err = fetchVatSense("https://api.vatsense.com/1.0/validate?vat_number="+varNumber, c.Password)
	}
	if err != nil {
		return nil, err
	}
	if response.Get("code").Int() == 200 {
		data := response.GetJson("data")
		valid := data.Get("valid").Bool()
		if valid {
			company := data.GetJson("company")
			return &base.ValidResult{
				Valid:          valid,
				VatNumber:      company.Get("vat_number").String(),
				CountryCode:    company.Get("country_code").String(),
				CompanyName:    company.Get("company_name").String(),
				CompanyAddress: company.Get("company_address").String(),
			}, nil
		} else {
			return &base.ValidResult{
				Valid: false,
			}, nil
		}
	} else {
		return nil, gerror.New(response.String())
	}
}

func (c VatSense) ValidateEoriNumber(number string) (*base.ValidResult, error) {
	response, err := fetchVatSense("https://api.vatsense.com/1.0/validate?eori_number="+number, c.Password)
	if err != nil {
		return nil, err
	}
	if response.Get("code").Int() == 200 {
		data := response.GetJson("data")
		valid := data.Get("valid").Bool()
		if valid {
			company := data.GetJson("company")
			return &base.ValidResult{
				Valid:          valid,
				VatNumber:      company.Get("vat_number").String(),
				CountryCode:    company.Get("country_code").String(),
				CompanyName:    company.Get("company_name").String(),
				CompanyAddress: company.Get("company_address").String(),
			}, nil
		} else {
			return &base.ValidResult{
				Valid: false,
			}, nil
		}
	} else {
		return nil, gerror.New(response.String())
	}
}
