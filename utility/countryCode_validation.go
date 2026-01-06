package utility

import (
	"fmt"
	"strings"
)

// ValidateCountryCode validates if the given country code is a valid ISO 3166-1 alpha-2 country code
// Returns error if the country code is invalid, nil if valid
func ValidateCountryCode(countryCode string) error {
	if countryCode == "" {
		return fmt.Errorf("country code cannot be empty")
	}

	// Convert to uppercase for case-insensitive comparison
	countryCode = strings.ToUpper(countryCode)

	// Check if the country code is exactly 2 characters
	if len(countryCode) != 2 {
		return fmt.Errorf("country code must be exactly 2 characters, got: %s", countryCode)
	}

	// Check if the country code contains only letters
	for _, char := range countryCode {
		if char < 'A' || char > 'Z' {
			return fmt.Errorf("country code must contain only letters, got: %s", countryCode)
		}
	}

	// Validate against ISO 3166-1 alpha-2 standard country codes
	if !isValidISO3166CountryCode(countryCode) {
		return fmt.Errorf("invalid ISO 3166-1 alpha-2 country code: %s", countryCode)
	}

	return nil
}

// isValidISO3166CountryCode checks if the given country code is a valid ISO 3166-1 alpha-2 code
// This includes all officially assigned country codes according to ISO 3166-1 standard
func isValidISO3166CountryCode(countryCode string) bool {
	validCountryCodes := map[string]bool{
		"AF": true, // Afghanistan
		"AX": true, // Åland Islands
		"AL": true, // Albania
		"DZ": true, // Algeria
		"AS": true, // American Samoa
		"AD": true, // Andorra
		"AO": true, // Angola
		"AI": true, // Anguilla
		"AQ": true, // Antarctica
		"AG": true, // Antigua and Barbuda
		"AR": true, // Argentina
		"AM": true, // Armenia
		"AW": true, // Aruba
		"AU": true, // Australia
		"AT": true, // Austria
		"AZ": true, // Azerbaijan
		"BS": true, // Bahamas
		"BH": true, // Bahrain
		"BD": true, // Bangladesh
		"BB": true, // Barbados
		"BY": true, // Belarus
		"BE": true, // Belgium
		"BZ": true, // Belize
		"BJ": true, // Benin
		"BM": true, // Bermuda
		"BT": true, // Bhutan
		"BO": true, // Bolivia
		"BQ": true, // Bonaire, Sint Eustatius and Saba
		"BA": true, // Bosnia and Herzegovina
		"BW": true, // Botswana
		"BV": true, // Bouvet Island
		"BR": true, // Brazil
		"IO": true, // British Indian Ocean Territory
		"BN": true, // Brunei Darussalam
		"BG": true, // Bulgaria
		"BF": true, // Burkina Faso
		"BI": true, // Burundi
		"CV": true, // Cabo Verde
		"KH": true, // Cambodia
		"CM": true, // Cameroon
		"CA": true, // Canada
		"KY": true, // Cayman Islands
		"CF": true, // Central African Republic
		"TD": true, // Chad
		"CL": true, // Chile
		"CN": true, // China
		"CX": true, // Christmas Island
		"CC": true, // Cocos (Keeling) Islands
		"CO": true, // Colombia
		"KM": true, // Comoros
		"CG": true, // Congo
		"CD": true, // Congo, Democratic Republic of the
		"CK": true, // Cook Islands
		"CR": true, // Costa Rica
		"CI": true, // Côte d'Ivoire
		"HR": true, // Croatia
		"CU": true, // Cuba
		"CW": true, // Curaçao
		"CY": true, // Cyprus
		"CZ": true, // Czech Republic
		"DK": true, // Denmark
		"DJ": true, // Djibouti
		"DM": true, // Dominica
		"DO": true, // Dominican Republic
		"EC": true, // Ecuador
		"EG": true, // Egypt
		"SV": true, // El Salvador
		"GQ": true, // Equatorial Guinea
		"ER": true, // Eritrea
		"EE": true, // Estonia
		"SZ": true, // Eswatini
		"ET": true, // Ethiopia
		"FK": true, // Falkland Islands
		"FO": true, // Faroe Islands
		"FJ": true, // Fiji
		"FI": true, // Finland
		"FR": true, // France
		"GF": true, // French Guiana
		"PF": true, // French Polynesia
		"TF": true, // French Southern Territories
		"GA": true, // Gabon
		"GM": true, // Gambia
		"GE": true, // Georgia
		"DE": true, // Germany
		"GH": true, // Ghana
		"GI": true, // Gibraltar
		"GR": true, // Greece
		"GL": true, // Greenland
		"GD": true, // Grenada
		"GP": true, // Guadeloupe
		"GU": true, // Guam
		"GT": true, // Guatemala
		"GG": true, // Guernsey
		"GN": true, // Guinea
		"GW": true, // Guinea-Bissau
		"GY": true, // Guyana
		"HT": true, // Haiti
		"HM": true, // Heard Island and McDonald Islands
		"VA": true, // Holy See
		"HN": true, // Honduras
		"HK": true, // Hong Kong
		"HU": true, // Hungary
		"IS": true, // Iceland
		"IN": true, // India
		"ID": true, // Indonesia
		"IR": true, // Iran
		"IQ": true, // Iraq
		"IE": true, // Ireland
		"IM": true, // Isle of Man
		"IL": true, // Israel
		"IT": true, // Italy
		"JM": true, // Jamaica
		"JP": true, // Japan
		"JE": true, // Jersey
		"JO": true, // Jordan
		"KZ": true, // Kazakhstan
		"KE": true, // Kenya
		"KI": true, // Kiribati
		"KP": true, // Korea, Democratic People's Republic of
		"KR": true, // Korea, Republic of
		"KW": true, // Kuwait
		"KG": true, // Kyrgyzstan
		"LA": true, // Lao People's Democratic Republic
		"LV": true, // Latvia
		"LB": true, // Lebanon
		"LS": true, // Lesotho
		"LR": true, // Liberia
		"LY": true, // Libya
		"LI": true, // Liechtenstein
		"LT": true, // Lithuania
		"LU": true, // Luxembourg
		"MO": true, // Macao
		"MG": true, // Madagascar
		"MW": true, // Malawi
		"MY": true, // Malaysia
		"MV": true, // Maldives
		"ML": true, // Mali
		"MT": true, // Malta
		"MH": true, // Marshall Islands
		"MQ": true, // Martinique
		"MR": true, // Mauritania
		"MU": true, // Mauritius
		"YT": true, // Mayotte
		"MX": true, // Mexico
		"FM": true, // Micronesia
		"MD": true, // Moldova
		"MC": true, // Monaco
		"MN": true, // Mongolia
		"ME": true, // Montenegro
		"MS": true, // Montserrat
		"MA": true, // Morocco
		"MZ": true, // Mozambique
		"MM": true, // Myanmar
		"NA": true, // Namibia
		"NR": true, // Nauru
		"NP": true, // Nepal
		"NL": true, // Netherlands
		"NC": true, // New Caledonia
		"NZ": true, // New Zealand
		"NI": true, // Nicaragua
		"NE": true, // Niger
		"NG": true, // Nigeria
		"NU": true, // Niue
		"NF": true, // Norfolk Island
		"MK": true, // North Macedonia
		"MP": true, // Northern Mariana Islands
		"NO": true, // Norway
		"OM": true, // Oman
		"PK": true, // Pakistan
		"PW": true, // Palau
		"PS": true, // Palestine, State of
		"PA": true, // Panama
		"PG": true, // Papua New Guinea
		"PY": true, // Paraguay
		"PE": true, // Peru
		"PH": true, // Philippines
		"PN": true, // Pitcairn
		"PL": true, // Poland
		"PT": true, // Portugal
		"PR": true, // Puerto Rico
		"QA": true, // Qatar
		"RE": true, // Réunion
		"RO": true, // Romania
		"RU": true, // Russian Federation
		"RW": true, // Rwanda
		"BL": true, // Saint Barthélemy
		"SH": true, // Saint Helena, Ascension and Tristan da Cunha
		"KN": true, // Saint Kitts and Nevis
		"LC": true, // Saint Lucia
		"MF": true, // Saint Martin (French part)
		"PM": true, // Saint Pierre and Miquelon
		"VC": true, // Saint Vincent and the Grenadines
		"WS": true, // Samoa
		"SM": true, // San Marino
		"ST": true, // Sao Tome and Principe
		"SA": true, // Saudi Arabia
		"SN": true, // Senegal
		"RS": true, // Serbia
		"SC": true, // Seychelles
		"SL": true, // Sierra Leone
		"SG": true, // Singapore
		"SX": true, // Sint Maarten (Dutch part)
		"SK": true, // Slovakia
		"SI": true, // Slovenia
		"SB": true, // Solomon Islands
		"SO": true, // Somalia
		"ZA": true, // South Africa
		"GS": true, // South Georgia and the South Sandwich Islands
		"SS": true, // South Sudan
		"ES": true, // Spain
		"LK": true, // Sri Lanka
		"SD": true, // Sudan
		"SR": true, // Suriname
		"SJ": true, // Svalbard and Jan Mayen
		"SE": true, // Sweden
		"CH": true, // Switzerland
		"SY": true, // Syrian Arab Republic
		"TW": true, // Taiwan
		"TJ": true, // Tajikistan
		"TZ": true, // Tanzania
		"TH": true, // Thailand
		"TL": true, // Timor-Leste
		"TG": true, // Togo
		"TK": true, // Tokelau
		"TO": true, // Tonga
		"TT": true, // Trinidad and Tobago
		"TN": true, // Tunisia
		"TR": true, // Turkey
		"TM": true, // Turkmenistan
		"TC": true, // Turks and Caicos Islands
		"TV": true, // Tuvalu
		"UG": true, // Uganda
		"UA": true, // Ukraine
		"AE": true, // United Arab Emirates
		"GB": true, // United Kingdom
		"US": true, // United States
		"UM": true, // United States Minor Outlying Islands
		"UY": true, // Uruguay
		"UZ": true, // Uzbekistan
		"VU": true, // Vanuatu
		"VE": true, // Venezuela
		"VN": true, // Vietnam
		"VG": true, // Virgin Islands, British
		"VI": true, // Virgin Islands, U.S.
		"WF": true, // Wallis and Futuna
		"EH": true, // Western Sahara
		"YE": true, // Yemen
		"ZM": true, // Zambia
		"ZW": true, // Zimbabwe
	}

	return validCountryCodes[countryCode]
}

// GetCountryCodeList returns a list of all valid ISO 3166-1 alpha-2 country codes
// This can be used for dropdown menus or validation purposes
func GetCountryCodeList() []string {
	return []string{
		"AF", "AX", "AL", "DZ", "AS", "AD", "AO", "AI", "AQ", "AG", "AR", "AM", "AW", "AU", "AT",
		"AZ", "BS", "BH", "BD", "BB", "BY", "BE", "BZ", "BJ", "BM", "BT", "BO", "BQ", "BA", "BW",
		"BV", "BR", "IO", "BN", "BG", "BF", "BI", "CV", "KH", "CM", "CA", "KY", "CF", "TD", "CL",
		"CN", "CX", "CC", "CO", "KM", "CG", "CD", "CK", "CR", "CI", "HR", "CU", "CW", "CY", "CZ",
		"DK", "DJ", "DM", "DO", "EC", "EG", "SV", "GQ", "ER", "EE", "SZ", "ET", "FK", "FO", "FJ",
		"FI", "FR", "GF", "PF", "TF", "GA", "GM", "GE", "DE", "GH", "GI", "GR", "GL", "GD", "GP",
		"GU", "GT", "GG", "GN", "GW", "GY", "HT", "HM", "VA", "HN", "HK", "HU", "IS", "IN", "ID",
		"IR", "IQ", "IE", "IM", "IL", "IT", "JM", "JP", "JE", "JO", "KZ", "KE", "KI", "KP", "KR",
		"KW", "KG", "LA", "LV", "LB", "LS", "LR", "LY", "LI", "LT", "LU", "MO", "MG", "MW", "MY",
		"MV", "ML", "MT", "MH", "MQ", "MR", "MU", "YT", "MX", "FM", "MD", "MC", "MN", "ME", "MS",
		"MA", "MZ", "MM", "NA", "NR", "NP", "NL", "NC", "NZ", "NI", "NE", "NG", "NU", "NF", "MK",
		"MP", "NO", "OM", "PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PN", "PL", "PT", "PR",
		"QA", "RE", "RO", "RU", "RW", "BL", "SH", "KN", "LC", "MF", "PM", "VC", "WS", "SM", "ST",
		"SA", "SN", "RS", "SC", "SL", "SG", "SX", "SK", "SI", "SB", "SO", "ZA", "GS", "SS", "ES",
		"LK", "SD", "SR", "SJ", "SE", "CH", "SY", "TW", "TJ", "TZ", "TH", "TL", "TG", "TK", "TO",
		"TT", "TN", "TR", "TM", "TC", "TV", "UG", "UA", "AE", "GB", "US", "UM", "UY", "UZ", "VU",
		"VE", "VN", "VG", "VI", "WF", "EH", "YE", "ZM", "ZW",
	}
}
