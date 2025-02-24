package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unibee/internal/logic/gateway/api/alipay/api/exception"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	"unibee/internal/logic/gateway/api/alipay/api/tools"
)

const readTimeout = 15 * time.Second
const connectTimeout = 15 * time.Second
const totalTimeout = 30 * time.Second

type DefaultAlipayClient struct {
	GatewayUrl         string
	ClientId           string
	MerchantPrivateKey string
	AlipayPublicKey    string
	IsSandboxMode      bool
	IsInStore          bool
}

/**
 * 构造方法
 * @param gatewayUrl 网关地址
 * @param clientId 应用ID
 * @param merchantPrivateKey 商户私钥
 * @param alipayPublicKey 支付宝公钥
 */
func NewDefaultAlipayClient(gatewayUrl string, clientId string, merchantPrivateKey string, alipayPublicKey string, isInStore bool) *DefaultAlipayClient {
	isSandboxMode := false
	if strings.HasPrefix(clientId, "SANDBOX_") {
		isSandboxMode = true
	}

	return &DefaultAlipayClient{
		GatewayUrl:         gatewayUrl,
		ClientId:           clientId,
		MerchantPrivateKey: merchantPrivateKey,
		AlipayPublicKey:    alipayPublicKey,
		IsSandboxMode:      isSandboxMode,
		IsInStore:          isInStore,
	}
}

// @Description 支持任意方式的http请求
// @Param url 请求地址
// @Param method 请求方式
// @Param headers 请求header头设置
// @Param params 请求地址栏后需要拼接参数操作
// @Param data 请求报文
// @Return 返回类型 "错误信息"
func (alipayClient *DefaultAlipayClient) httpDo(url, method string, params, headers map[string]string, data []byte, alipayResponse any) (any, error) {

	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: connectTimeout,
		}).DialContext,
		ResponseHeaderTimeout: readTimeout,
	}

	//自定义client
	client := &http.Client{
		Timeout:   totalTimeout,
		Transport: trans,
	}
	//http.post等方法只是在NewRequest上又封装来了一层而已
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, &exception.AlipayLibraryError{Message: "http.NewRequest is fail " + err.Error()}
	}
	req.Header.Set("Content-type", "application/json")

	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &exception.AlipayLibraryError{Message: "client.Do is fail " + err.Error()}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	responseTime := resp.Header.Get("response-time")
	sign, err := checkRspSign(resp.Request.Method, req.URL.Path, resp.Header.Get("Client-id"), responseTime, string(body), resp.Header.Get("Signature"), alipayClient.AlipayPublicKey)
	if err != nil {
		return nil, &exception.AlipayLibraryError{Message: "checkRspSign is fail " + err.Error()}
	}
	if !sign {
		return nil, &exception.AlipayLibraryError{Message: "check signature fail"}
	}
	//通过指针将request的response值赋值，这样虽然是any类型，
	//但是我们在上次必然已经定义了any的类型
	err = json.Unmarshal(body, alipayResponse)
	if err != nil {
		return nil, &exception.AlipayLibraryError{Message: "json.Unmarshal is fail " + err.Error()}
	}

	return alipayResponse, nil
}

func (alipayClient *DefaultAlipayClient) Execute(alipayRequest *request.AlipayRequest) (any, error) {
	reqPayload, err := json.Marshal(alipayRequest.Param)
	if err != nil {
		return nil, &exception.AlipayLibraryError{Message: "json.Marshal is fail " + err.Error()}
	}
	path := alipayRequest.Path
	if alipayClient.IsSandboxMode && !alipayClient.IsInStore {
		path = strings.ReplaceAll(path, "/ams", "/ams/sandbox")
	}
	httpMethod := alipayRequest.HttpMethod
	reqTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	sign, err := tools.GenSign(fmt.Sprintf("%s", httpMethod), path, alipayClient.ClientId, reqTime, string(reqPayload), alipayClient.MerchantPrivateKey)
	if err != nil {
		return nil, err
	}
	header := buildBaseHeader(reqTime, alipayClient.ClientId, alipayRequest.KeyVersion, sign)
	alipayResponse, err := alipayClient.httpDo(alipayClient.GatewayUrl+path, httpMethod, map[string]string{}, header, reqPayload, alipayRequest.AlipayResponse)
	if err != nil {
		return nil, err
	}
	return alipayResponse, nil
}

func checkRspSign(httpMethod string, path string, clientId string, responseTime string, rspBody string, rspSignValue string, alipayPublicKey string) (bool, error) {

	signature, err := tools.CheckSignature(path, httpMethod, clientId, responseTime, rspBody, rspSignValue, alipayPublicKey)
	if err != nil {
		return false, &exception.AlipayLibraryError{Message: "Failed to check signature  " + err.Error()}
	}
	if !signature {
		return false, &exception.AlipayLibraryError{Message: "check signature fail"}
	} else {
		return true, nil
	}
}

func buildBaseHeader(reqTime string, clientId string, keyVersion string, signatureValue string) map[string]string {
	if keyVersion == "" {
		keyVersion = "1"
	}
	signatureValue = "algorithm=RSA256,keyVersion=" + keyVersion + ",signature=" + signatureValue
	return map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Request-Time": reqTime,
		"Client-Id":    clientId,
		"Key-Version":  keyVersion,
		"Signature":    signatureValue,
	}
}
