package cloud

import (
	"context"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/email"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"
)

func MerchantSetupForCloudMode(ctx context.Context, merchantId uint64) error {
	if config.GetConfigInstance().Mode == "cloud" {
		//if cloud version setup default sendgrid and vat
		{
			name, data := email.GetDefaultMerchantEmailConfig(ctx, 15621)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := email.SetupMerchantEmailConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
		}
		{
			name, data := vat_gateway.GetDefaultMerchantVatConfig(ctx, 15621)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := vat_gateway.SetupMerchantVatConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SendMerchantRegisterEmail(ctx context.Context, to string, verificationCode string) {
	if config.GetConfigInstance().Mode == "cloud" {
		err := email.SendTemplateEmail(ctx, 15621, to, "", email.TemplateMerchantRegistrationCodeVerify, "", &email.TemplateVariable{
			CodeExpireMinute: "3",
			Code:             verificationCode,
		})
		utility.AssertError(err, "Server Error")
	} else {
		// deploy version todo mark send to unibee api, cloud version use cloud merchantId
		utility.Assert(true, "not support")
	}
}
