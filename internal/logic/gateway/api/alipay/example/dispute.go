package main

import (
	defaultAlipayClient "unibee/internal/logic/gateway/api/alipay/api"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request/dispute"
	responseDispute "unibee/internal/logic/gateway/api/alipay/api/response/dispute"
)

func main() {
	const alipayGatewayUrl = ""
	const alipayClientId = ""
	const alipayMerchantPrivateKey = ""
	const alipayAlipayPublicKey = ""

	client := defaultAlipayClient.NewDefaultAlipayClient(
		alipayGatewayUrl,
		alipayClientId,
		alipayMerchantPrivateKey,
		alipayAlipayPublicKey, false)

	acceptDispute(client, "202409141900000000000001J0000009488")
}

func acceptDispute(client *defaultAlipayClient.DefaultAlipayClient, disputeId string) {
	request, acceptDisputeRequest := dispute.NewAlipayAcceptDisputeRequest()
	acceptDisputeRequest.DisputeId = disputeId

	response, err := client.Execute(request)
	if err != nil {
		panic(err)
	}
	println(response.(*responseDispute.AlipayAcceptDisputeResponse))
}

func supplyDefenseDocument(client *defaultAlipayClient.DefaultAlipayClient, disputeId string) {
	request, supplyDefenseDocumentRequest := dispute.NewAlipaySupplyDefenseDocumentRequest()
	supplyDefenseDocumentRequest.DisputeId = disputeId
	supplyDefenseDocumentRequest.DisputeEvidence = "test"

	response, err := client.Execute(request)
	if err != nil {
		panic(err)
	}
	println(response.(*responseDispute.AlipaySupplyDefenseDocumentResponse))
}

func downloadDisputeEvidence(client *defaultAlipayClient.DefaultAlipayClient, disputeId string) {
	request, downloadDisputeEvidenceRequest := dispute.NewAlipayDownloadDisputeEvidenceRequest()
	downloadDisputeEvidenceRequest.DisputeId = disputeId
	downloadDisputeEvidenceRequest.DisputeEvidenceType = model.DisputeEvidenceType_DISPUTE_EVIDENCE_TEMPLATE

	response, err := client.Execute(request)
	if err != nil {
		panic(err)
	}
	println(response.(*responseDispute.AlipayDownloadDisputeEvidenceResponse))
}
