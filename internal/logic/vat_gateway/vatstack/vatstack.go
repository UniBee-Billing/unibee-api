package vatstack

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"io"
	"net/http"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type VatStack struct {
	ApiData string
	Name    string
}

func fetchVatSense(url string, publicKey string) (*gjson.Json, error) {

	// Create HTTP Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("FetchVatStack Error creating request:", err)
		return nil, err
	}

	// Add HTTP Header
	req.Header.Set("X-API-KEY", publicKey)

	// send HTTP
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("FetchVatStack Error making request:", err)
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
		fmt.Println("FetchVatStack Error reading response body:", err)
		return nil, err
	}

	// print response
	fmt.Printf("FetchVatStack Response:%s\n", string(body))

	return gjson.LoadJson(string(body))
}

func (v VatStack) GetGatewayName() string {
	return v.Name
}

func (v VatStack) ListAllCountries() ([]*entity.CountryRate, error) {
	return []*entity.CountryRate{}, nil
}

func (v VatStack) ListAllRates() ([]*entity.CountryRate, error) {
	response, err := fetchVatSense("https://api.vatstack.com/v1/rates?limit=100", v.ApiData)
	if err != nil {
		return nil, err
	}
	var countryRates []*entity.CountryRate
	for _, item := range response.GetJsons("rates") {
		var vat = 1
		var eu = 0
		if item.Get("member_state") != nil && item.Get("member_state").Bool() {
			eu = 1
		} else {
			eu = 2
		}
		countryRates = append(countryRates, &entity.CountryRate{
			Gateway:               v.GetGatewayName(),
			CountryCode:           item.Get("country_code").String(),
			CountryName:           item.Get("country_name").String(),
			Vat:                   vat,
			Eu:                    eu,
			StandardDescription:   item.Get("local_name").String(),
			StandardTaxPercentage: int64(item.Get("standard_rate").Float64() * 100),
			Other:                 item.GetJson("categories").String(),
		})
	}
	return countryRates, nil
}

func (v VatStack) ValidateVatNumber(vatNumber string, requesterVatNumber string) (*bean.ValidResult, error) {
	if len(requesterVatNumber) > 0 {
		return nil, gerror.New("not support")
	}
	if len(vatNumber) > 0 {
		return nil, gerror.New("invalid vatNumber")
	}
	bodyReader := bytes.NewReader([]byte(utility.MarshalToJsonString(map[string]string{
		"query": vatNumber,
	})))
	request, err := http.NewRequest("POST", "https://api.vatstack.com/v1/validations", bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-API-KEY", v.ApiData)
	client := &http.Client{}
	responseData, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(responseData.Body)

	responseBody, err := io.ReadAll(responseData.Body)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 201 && responseData.StatusCode != 202 {
		return nil, gerror.NewCode(gcode.New(responseData.StatusCode, responseData.Status, responseData.Status+" "+string(responseBody)), responseData.Status+" "+string(responseBody))
	}
	if responseData.StatusCode == 202 {
		// todo mark VatStack will send webhook valid result
		return &bean.ValidResult{
			Valid: false,
		}, nil
	}
	response, err := gjson.LoadJson(responseBody)
	if err != nil {
		return nil, err
	}
	valid := response.Get("active").Bool()
	if valid {
		return &bean.ValidResult{
			Valid:          valid,
			VatNumber:      response.Get("query").String(),
			CountryCode:    response.Get("country_code").String(),
			CompanyName:    response.Get("company_name").String(),
			CompanyAddress: response.Get("company_address").String(),
		}, nil
	} else {
		return &bean.ValidResult{
			Valid: false,
		}, nil
	}
}

func (v VatStack) ValidateEoriNumber(number string) (*bean.ValidResult, error) {
	return nil, gerror.New("not support")
}
