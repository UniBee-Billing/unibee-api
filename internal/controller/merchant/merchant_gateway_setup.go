package merchant

import (
	"context"
	gateway3 "unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error) {
	if req.CompanyIssuer != nil {
		if req.Metadata == nil {
			req.Metadata = map[string]interface{}{}
		}
		req.Metadata["IssueCompanyName"] = req.CompanyIssuer.IssueCompanyName
		req.Metadata["IssueAddress"] = req.CompanyIssuer.IssueAddress
		req.Metadata["IssueRegNumber"] = req.CompanyIssuer.IssueRegNumber
		req.Metadata["IssueVatNumber"] = req.CompanyIssuer.IssueVatNumber
		req.Metadata["IssueLogo"] = req.CompanyIssuer.IssueLogo
	}
	return &gateway.SetupRes{Gateway: gateway3.ConvertGatewayDetail(ctx, gateway2.SetupGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.GatewayKey, req.GatewaySecret, req.SubGateway, req.GatewayPaymentTypes, req.DisplayName, req.GatewayIcons, req.Sort, req.CurrencyExchange, req.Metadata))}, nil
}
