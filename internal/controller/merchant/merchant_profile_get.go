package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/profile"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/quickbooks"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/email"
	member2 "unibee/internal/logic/member"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/multi_currencies"
	"unibee/internal/logic/multi_currencies/currency_exchange"
	"unibee/internal/logic/totp"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"
	"unibee/utility/unibee"
)

func (c *ControllerProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	var member *entity.MerchantMember
	var isOwner = false
	var memberRoles = make([]*bean.MerchantRole, 0)
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		member = query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
		if member != nil {
			isOwner, memberRoles = detail.ConvertMemberRole(ctx, member)
		}
	}

	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	vatGatewayName, vatGatewayKey := vat_gateway.GetDefaultMerchantVatConfig(ctx, merchant.Id)
	if vatGatewayName != "vatsense" {
		vatGatewayKey = ""
	}
	_, emailData := email.GetDefaultMerchantEmailConfig(ctx, merchant.Id)
	var emailSender *bean.Sender
	one := email.GetMerchantEmailSender(ctx, _interface.GetMerchantId(ctx))
	if one != nil {
		emailSender = &bean.Sender{
			Name:    one.Name,
			Address: one.Address,
		}
	}
	exchangeApiKey := ""
	exchangeApiKeyConfig := merchant_config.GetMerchantConfig(ctx, _interface.GetMerchantId(ctx), currency_exchange.FiatExchangeApiKey)
	if exchangeApiKeyConfig != nil {
		exchangeApiKey = exchangeApiKeyConfig.ConfigValue
	}
	apikey := merchant.ApiKey
	if config.GetConfigInstance().IsProd() {
		apikey = utility.HideStar(merchant.ApiKey)
	}
	session := ""
	if member != nil {
		session, _ = member2.NewMemberSession(ctx, int64(member.Id), "")
	}
	var defaultCurrency = "USD"
	onePlan := query.GetOneLatestMainPlanByMerchantId(ctx, merchant.Id)
	if onePlan != nil {
		defaultCurrency = onePlan.Currency
	} else if len(merchant.CountryCode) > 0 {
		if target, ok := countryCurrency[merchant.CountryCode]; ok {
			defaultCurrency = target
		}
	}
	var qbCompanyName = ""
	var qbLastSynchronized = "Please wait at least 2 hours for initial sync"
	var qbLastSyncError = "No errors"
	qbConfig := quickbooks.GetMerchantQuickBooksConfig(ctx, merchant.Id)
	if qbConfig != nil && len(qbConfig.CompanyName) > 0 && qbConfig.BearerToken != nil && qbConfig.BearerToken.AccessToken != "" {
		qbCompanyName = qbConfig.CompanyName
	}
	analyticsHost := fmt.Sprintf("%s/analytics?session=%s", config.GetConfigInstance().Server.GetServerPath(), session)
	if len(config.GetConfigInstance().Server.AnalyticsPath) > 0 {
		analyticsHost = fmt.Sprintf("%s?session=%s", config.GetConfigInstance().Server.AnalyticsPath, session)
	}
	cloudFeatureAnalyticsEnabled := false
	if config.GetConfigInstance().Mode == "cloud" {
		cloudFeatureAnalyticsEnabledConfig := merchant_config.GetMerchantConfig(ctx, _interface.GetMerchantId(ctx), "FeatureAnalyticsEnabled")
		if cloudFeatureAnalyticsEnabledConfig != nil && cloudFeatureAnalyticsEnabledConfig.ConfigValue == "true" {
			cloudFeatureAnalyticsEnabled = true
		}
	}
	return &profile.GetRes{
		Merchant:                     bean.SimplifyMerchant(merchant),
		MerchantMember:               detail.ConvertMemberToDetail(ctx, member),
		DefaultCurrency:              defaultCurrency,
		Currency:                     currency.GetMerchantCurrencies(),
		Env:                          config.GetConfigInstance().Env,
		IsProd:                       config.GetConfigInstance().IsProd(),
		TimeZone:                     time.GetTimeZoneList(),
		Gateways:                     detail.ConvertGatewayList(ctx, query.GetMerchantGatewayList(ctx, merchant.Id, unibee.Bool(false))),
		ExchangeRateApiKey:           utility.HideStar(exchangeApiKey),
		OpenAPIHost:                  config.GetConfigInstance().Server.GetServerPath(),
		OpenAPIKey:                   apikey,
		SendGridKey:                  utility.HideStar(emailData),
		EmailSender:                  emailSender,
		VatSenseKey:                  utility.HideStar(vatGatewayKey),
		SegmentServerSideKey:         segment.GetMerchantSegmentServerSideConfig(ctx, merchant.Id),
		SegmentUserPortalKey:         segment.GetMerchantSegmentUserPortalConfig(ctx, merchant.Id),
		GlobalTOPTEnabled:            totp.GetMerchantTotpGlobalConfig(ctx, merchant.Id),
		QuickBooksCompanyName:        qbCompanyName,
		QuickBooksLastSynchronized:   qbLastSynchronized,
		QuickBooksLastSyncError:      qbLastSyncError,
		IsOwner:                      isOwner,
		MemberRoles:                  memberRoles,
		AnalyticsHost:                analyticsHost,
		CloudFeatureAnalyticsEnabled: cloudFeatureAnalyticsEnabled,
		MultiCurrencies:              multi_currencies.GetMerchantMultiCurrenciesConfig(ctx, merchant.Id),
	}, nil
}

var countryCurrency = map[string]string{"BD": "BDT", "BE": "EUR", "BF": "XOF", "BG": "BGN", "BA": "BAM", "BB": "BBD", "WF": "XPF", "BL": "EUR", "BM": "BMD", "BN": "BND", "BO": "BOB", "BH": "BHD", "BI": "BIF", "BJ": "XOF", "BT": "BTN", "JM": "JMD", "BV": "NOK", "BW": "BWP", "WS": "WST", "BQ": "USD", "BR": "BRL", "BS": "BSD", "JE": "GBP", "BY": "BYR", "BZ": "BZD", "RU": "RUB", "RW": "RWF", "RS": "RSD", "TL": "USD", "RE": "EUR", "TM": "TMT", "TJ": "TJS", "RO": "RON", "TK": "NZD", "GW": "XOF", "GU": "USD", "GT": "GTQ", "GS": "GBP", "GR": "EUR", "GQ": "XAF", "GP": "EUR", "JP": "JPY", "GY": "GYD", "GG": "GBP", "GF": "EUR", "GE": "GEL", "GD": "XCD", "GB": "GBP", "GA": "XAF", "SV": "USD", "GN": "GNF", "GM": "GMD", "GL": "DKK", "GI": "GIP", "GH": "GHS", "OM": "OMR", "TN": "TND", "JO": "JOD", "HR": "HRK", "HT": "HTG", "HU": "HUF", "HK": "HKD", "HN": "HNL", "HM": "AUD", "VE": "VEF", "PR": "USD", "PS": "ILS", "PW": "USD", "PT": "EUR", "SJ": "NOK", "PY": "PYG", "IQ": "IQD", "PA": "PAB", "PF": "XPF", "PG": "PGK", "PE": "PEN", "PK": "PKR", "PH": "PHP", "PN": "NZD", "PL": "PLN", "PM": "EUR", "ZM": "ZMK", "EH": "MAD", "EE": "EUR", "EG": "EGP", "ZA": "ZAR", "EC": "USD", "IT": "EUR", "VN": "VND", "SB": "SBD", "ET": "ETB", "SO": "SOS", "ZW": "ZWL", "SA": "SAR", "ES": "EUR", "ER": "ERN", "ME": "EUR", "MD": "MDL", "MG": "MGA", "MF": "EUR", "MA": "MAD", "MC": "EUR", "UZ": "UZS", "MM": "MMK", "ML": "XOF", "MO": "MOP", "MN": "MNT", "MH": "USD", "MK": "MKD", "MU": "MUR", "MT": "EUR", "MW": "MWK", "MV": "MVR", "MQ": "EUR", "MP": "USD", "MS": "XCD", "MR": "MRO", "IM": "GBP", "UG": "UGX", "TZ": "TZS", "MY": "MYR", "MX": "MXN", "IL": "ILS", "FR": "EUR", "IO": "USD", "SH": "SHP", "FI": "EUR", "FJ": "FJD", "FK": "FKP", "FM": "USD", "FO": "DKK", "NI": "NIO", "NL": "EUR", "NO": "NOK", "NA": "NAD", "VU": "VUV", "NC": "XPF", "NE": "XOF", "NF": "AUD", "NG": "NGN", "NZ": "NZD", "NP": "NPR", "NR": "AUD", "NU": "NZD", "CK": "NZD", "XK": "EUR", "CI": "XOF", "CH": "CHF", "CO": "COP", "CN": "CNY", "CM": "XAF", "CL": "CLP", "CC": "AUD", "CA": "CAD", "CG": "XAF", "CF": "XAF", "CD": "CDF", "CZ": "CZK", "CY": "EUR", "CX": "AUD", "CR": "CRC", "CW": "ANG", "CV": "CVE", "CU": "CUP", "SZ": "SZL", "SY": "SYP", "SX": "ANG", "KG": "KGS", "KE": "KES", "SS": "SSP", "SR": "SRD", "KI": "AUD", "KH": "KHR", "KN": "XCD", "KM": "KMF", "ST": "STD", "SK": "EUR", "KR": "KRW", "SI": "EUR", "KP": "KPW", "KW": "KWD", "SN": "XOF", "SM": "EUR", "SL": "SLL", "SC": "SCR", "KZ": "KZT", "KY": "KYD", "SG": "SGD", "SE": "SEK", "SD": "SDG", "DO": "DOP", "DM": "XCD", "DJ": "DJF", "DK": "DKK", "VG": "USD", "DE": "EUR", "YE": "YER", "DZ": "DZD", "US": "USD", "UY": "UYU", "YT": "EUR", "UM": "USD", "LB": "LBP", "LC": "XCD", "LA": "LAK", "TV": "AUD", "TW": "TWD", "TT": "TTD", "TR": "TRY", "LK": "LKR", "LI": "CHF", "LV": "EUR", "TO": "TOP", "LT": "LTL", "LU": "EUR", "LR": "LRD", "LS": "LSL", "TH": "THB", "TF": "EUR", "TG": "XOF", "TD": "XAF", "TC": "USD", "LY": "LYD", "VA": "EUR", "VC": "XCD", "AE": "AED", "AD": "EUR", "AG": "XCD", "AF": "AFN", "AI": "XCD", "VI": "USD", "IS": "ISK", "IR": "IRR", "AM": "AMD", "AL": "ALL", "AO": "AOA", "AQ": "", "AS": "USD", "AR": "ARS", "AU": "AUD", "AT": "EUR", "AW": "AWG", "IN": "INR", "AX": "EUR", "AZ": "AZN", "IE": "EUR", "ID": "IDR", "UA": "UAH", "QA": "QAR", "MZ": "MZN"}
