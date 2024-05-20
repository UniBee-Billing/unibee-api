package merchant

import (
	"context"
	"flag"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/service"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func SetupForCloudMode(ctx context.Context, merchantId uint64) error {
	if config.GetConfigInstance().Mode == "cloud" {
		//if cloud version setup default sendgrid, vat, stripe gateway
		{
			name, data := email.GetDefaultMerchantEmailConfig(ctx, consts.CloudModeManagerMerchantId)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := email.SetupMerchantEmailConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
		}
		{
			name, data := vat_gateway.GetDefaultMerchantVatConfig(ctx, consts.CloudModeManagerMerchantId)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := vat_gateway.SetupMerchantVatConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
			err = vat_gateway.InitMerchantDefaultVatGateway(ctx, merchantId)
			if err != nil {
				return err
			}
		}
		{
			stripeGateway := query.GetGatewayByGatewayName(ctx, consts.CloudModeManagerMerchantId, "stripe")
			if stripeGateway != nil {
				service.SetupGateway(ctx, merchantId, stripeGateway.GatewayName, stripeGateway.GatewayKey, stripeGateway.GatewaySecret)
			}
		}
	}
	return nil
}

func StandAloneInit(ctx context.Context) {
	var merchantEmail string
	var password string
	flag.StringVar(&merchantEmail, "admin-email", utility.GetEnvParam("admin.email"), "admin email, default accounts.unibee@unibee.dev")
	flag.StringVar(&password, "admin-password", utility.GetEnvParam("admin.password"), "admin password, default changeme")
	if len(merchantEmail) == 0 {
		merchantEmail = "accounts.unibee@unibee.dev"
	}
	if len(password) == 0 {
		password = "changeme"
	}
	list, err := query.GetMerchantList(ctx)
	if err != nil {
		glog.Errorf(ctx, "StandAloneInit adminAccount error:%s", err.Error())
	}
	if err == nil && len(list) == 0 {
		_, _, err = CreateMerchant(ctx, &CreateMerchantInternalReq{
			FirstName: "unibee",
			LastName:  "unibee",
			Email:     merchantEmail,
			Password:  password,
			Phone:     "",
			UserName:  "",
		})
		if err != nil {
			g.Log().Errorf(ctx, "StandAloneInit adminAccount error:%s", err.Error())
			return
		} else {
			g.Log().Infof(ctx, "StandAloneInit adminAccount email:%s", merchantEmail)
		}
	}
}
