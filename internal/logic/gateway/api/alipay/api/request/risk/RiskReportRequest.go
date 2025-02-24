package risk

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseRisk "unibee/internal/logic/gateway/api/alipay/api/response/risk"
)

type RiskReportRequest struct {
	ReferenceTransactionId string `json:"referenceTransactionId,omitempty"`
	ReportReason           string `json:"reportReason,omitempty"`
	RiskType               string `json:"riskType,omitempty"`
	RiskOccurrenceTime     string `json:"riskOccurrenceTime,omitempty"`
}

func NewRiskReportRequest() (*request.AlipayRequest, *RiskReportRequest) {
	riskReportRequest := &RiskReportRequest{}
	alipayRequest := request.NewAlipayRequest(riskReportRequest, model.RISK_REPORT_PATH, &responseRisk.RiskReportResponse{})
	return alipayRequest, riskReportRequest
}
