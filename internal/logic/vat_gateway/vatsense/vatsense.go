package vatsense

import (
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"io"
	"net/http"
)

type VatSense struct {
	Password string
	Name     string
}

func fetchVatSense(url string, passwd string) (*gjson.Json, error) {
	// username and password
	username := "user"
	password := passwd

	// Create HTTP Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("FetchVatSense Error creating request:", err)
		return nil, err
	}

	// Add HTTP Header
	auth := username + ":" + password
	authEncoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", authEncoded)

	// send HTTP
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

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FetchVatSense Error reading response body:", err)
		return nil, err
	}

	// print response
	fmt.Printf("FetchVatSense Response:%s\n", string(body))

	return gjson.LoadJson(string(body))
}

func (c VatSense) GetGatewayName() string {
	return c.Name
}

func (c VatSense) ListAllCountries() ([]*entity.CountryRate, error) {
	data, err := fetchVatSense("https://api.vatsense.com/1.0/countries", c.Password)
	if err != nil {
		return nil, err
	}
	if data.Get("code").Int() == 200 {
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
				Gateway:     c.GetGatewayName(),
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
				Gateway:               c.GetGatewayName(),
				CountryCode:           item.Get("country_code").String(),
				CountryName:           item.Get("country_name").String(),
				Latitude:              item.Get("latitude").String(),
				Longitude:             item.Get("longitude").String(),
				Vat:                   vat,
				Eu:                    eu,
				Provinces:             item.Get("provinces").String(),
				StandardTypes:         standard.Get("types").String(),
				StandardDescription:   standard.Get("description").String(),
				StandardTaxPercentage: int64(standard.Get("rate").Float64() * 100),
				Other:                 item.GetJson("other").String(),
			})
		}
		return countryRates, nil
	} else {
		return nil, gerror.New(response.String())
	}
}

func (c VatSense) ValidateVatNumber(vatNumber string, requesterVatNumber string) (*ro.ValidResult, error) {
	var response *gjson.Json
	var err error
	if len(requesterVatNumber) > 0 {
		response, err = fetchVatSense("https://api.vatsense.com/1.0/validate?vat_number="+vatNumber+"&requester_vat_number="+requesterVatNumber, c.Password)
	} else {
		response, err = fetchVatSense("https://api.vatsense.com/1.0/validate?vat_number="+vatNumber, c.Password)
	}
	if err != nil {
		return nil, err
	}
	if response.Get("code").Int() == 200 {
		data := response.GetJson("data")
		valid := data.Get("valid").Bool()
		if valid {
			company := data.GetJson("company")
			return &ro.ValidResult{
				Valid:          valid,
				VatNumber:      company.Get("vat_number").String(),
				CountryCode:    company.Get("country_code").String(),
				CompanyName:    company.Get("company_name").String(),
				CompanyAddress: company.Get("company_address").String(),
			}, nil
		} else {
			return &ro.ValidResult{
				Valid: false,
			}, nil
		}
	} else if response.Contains("error") && response.GetJson("error").Contains("detail") {
		return &ro.ValidResult{
			Valid:           false,
			ValidateMessage: fmt.Sprintf("%s-%s", response.GetJson("error").Get("title").String(), response.GetJson("error").Get("detail").String()),
		}, nil
	} else {
		return nil, gerror.New(response.String())
	}
}

func (c VatSense) ValidateEoriNumber(number string) (*ro.ValidResult, error) {
	response, err := fetchVatSense("https://api.vatsense.com/1.0/validate?eori_number="+number, c.Password)
	if err != nil {
		return nil, err
	}
	if response.Get("code").Int() == 200 {
		data := response.GetJson("data")
		valid := data.Get("valid").Bool()
		if valid {
			company := data.GetJson("company")
			return &ro.ValidResult{
				Valid:          valid,
				VatNumber:      company.Get("vat_number").String(),
				CountryCode:    company.Get("country_code").String(),
				CompanyName:    company.Get("company_name").String(),
				CompanyAddress: company.Get("company_address").String(),
			}, nil
		} else {
			return &ro.ValidResult{
				Valid: false,
			}, nil
		}
	} else {
		return nil, gerror.New(response.String())
	}
}
