package response

type Result struct {
	ResultCode    string `json:"resultCode"`
	ResultStatus  string `json:"resultStatus"`
	ResultMessage string `json:"resultMessage"`
}

type AlipayResponse struct {
	Result Result `json:"result"`
}
