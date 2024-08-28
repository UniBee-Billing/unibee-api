package generator

import "fmt"

// Address represent an address
type Address struct {
	//Address    string `json:"address,omitempty" validate:"required"`
	Address    string `json:"address,omitempty"`
	Address2   string `json:"address_2,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	RegNumber  string `json:"regNumber,omitempty"`
	VatNumber  string `json:"vatNumber,omitempty"`
}

// ToString output address as string
// Line break are added for new lines
func (a *Address) ToString() string {
	addrString := a.Address

	if len(a.Address2) > 0 {
		addrString += "\n"
		addrString += a.Address2
	}

	if len(a.PostalCode) > 0 {
		addrString += "\n"
		addrString += a.PostalCode
	}
	//else {
	//	addrString += "\n"
	//}

	if len(a.City) > 0 {
		addrString += " "
		addrString += a.City
	}

	if len(a.Country) > 0 {
		addrString += "\n"
		addrString += a.Country
	}

	if len(a.RegNumber) > 0 {
		addrString += "\n"
		addrString += fmt.Sprintf("%s", a.RegNumber)
	}

	if len(a.VatNumber) > 0 {
		addrString += "\n"
		addrString += fmt.Sprintf("VAT Number:%s", a.VatNumber)
	}

	return addrString
}
