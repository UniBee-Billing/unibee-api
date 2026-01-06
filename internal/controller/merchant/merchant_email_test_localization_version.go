package merchant

import (
	"context"
	"fmt"
	"strings"
	"unibee/internal/consts"
	"unibee/utility"

	"github.com/gogf/gf/v2/os/gtime"

	"unibee/api/bean"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/query"
)

func (c *ControllerEmail) TestLocalizationVersion(ctx context.Context, req *email.TestLocalizationVersionReq) (res *email.TestLocalizationVersionRes, err error) {
	merchantId := _interface.GetMerchantId(ctx)

	// Validate language
	if !consts.IsSupportedLanguage(req.Language) {
		if req.Language == "" {
			req.Language = "en"
		}
		utility.Assert(consts.IsSupportedLanguage(req.Language), fmt.Sprintf("Unsupported language: %s. Supported languages are: %s", req.Language, strings.Join(consts.GetSupportedLanguagesList(), ", ")))
	}

	// Get the email template
	template := query.GetMerchantEmailTemplateByTemplateName(ctx, merchantId, req.TemplateName)
	utility.Assert(template != nil, "Template not found")

	// Find the specific version
	var targetVersion *bean.MerchantLocalizationVersion
	for _, version := range template.LocalizationVersions {
		if version.VersionId == req.VersionId {
			targetVersion = version
			break
		}
	}

	utility.Assert(targetVersion != nil, "localization version not found")

	// Check if the requested language exists in the version
	var targetLocalization *bean.EmailLocalizationTemplate
	if targetVersion != nil {
		for _, localization := range targetVersion.Localizations {
			if localization.Language == req.Language {
				targetLocalization = localization
				break
			}
		}
	}
	utility.Assert(targetLocalization != nil, fmt.Sprintf("Language %s not found in the specified version", req.Language))

	// Create test template variables using EmailTemplateVariable struct
	now := gtime.New(gtime.Now())
	periodEnd := gtime.New(gtime.Now().AddDate(0, 1, 0))

	templateVariables := &bean.EmailTemplateVariable{
		InvoiceId:             "INV-2024-001",
		UserName:              "John Doe",
		MerchantProductName:   "Premium Subscription",
		MerchantCustomerEmail: "support@unibee.dev",
		MerchantName:          "Example Company",
		DateNow:               now,
		PeriodEnd:             periodEnd,
		PaymentAmount:         "99.99",
		RefundAmount:          "49.99",
		Currency:              "USD",
		TokenExpireMinute:     "30",
		CodeExpireMinute:      "15",
		Code:                  "123456",
		Link:                  "https://unibee.dev",
		HttpLink:              "https://unibee.dev",
		AccountHolder:         "UniBee Company Ltd",
		Address:               "123 Business St, City, Country",
		BIC:                   "EXAMPLBIC",
		IBAN:                  "GB29NWBK60161331926819",
		BankData:              "Bank of Example, Account: 12345678",
	}

	// Send the test email using existing function
	if targetVersion != nil {
		err = email2.SendTemplateEmailByOpenApi(
			ctx,
			merchantId,
			req.Email,
			"UTC",
			req.Language,
			req.TemplateName,
			"", // No PDF attachment for test
			templateVariables,
			&targetVersion.Localizations,
		)
		utility.AssertError(err, "Failed to send test email")
	} else {
		utility.Assert(false, "localization version not found")
	}

	return &email.TestLocalizationVersionRes{}, nil
}
