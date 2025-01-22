package utility

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"io"
	"net/http"
)

func SendRequest(url string, method string, data []byte, headers map[string]string) ([]byte, error) {
	bodyReader := bytes.NewReader(data)
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	//client := &http.Client{}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, gerror.NewCode(gcode.New(response.StatusCode, response.Status, response.Status+" "+string(responseBody)), response.Status+" "+string(responseBody))
	}
	return responseBody, nil
}
