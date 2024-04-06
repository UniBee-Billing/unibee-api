package merchant

import (
	"context"
	"flag"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/logic/email"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func MerchantSetupForCloudMode(ctx context.Context, merchantId uint64) error {
	if config.GetConfigInstance().Mode == "cloud" {
		//if cloud version setup default sendgrid and vat
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
	if len(query.GetMerchantList(ctx)) == 0 {
		_, _, err := CreateMerchant(ctx, &CreateMerchantInternalReq{
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
