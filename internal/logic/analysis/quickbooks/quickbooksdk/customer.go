// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooksdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	//"gopkg.in/guregu/null.v4"
)

// Customer represents a QuickBooks Customer object.
type Customer struct {
	Id                   string           `json:"Id,omitempty"`
	SyncToken            string           `json:"SyncToken,omitempty"`
	MetaData             *MetaData        `json:"MetaData,omitempty"`
	Title                string           `json:"Title,omitempty"`
	GivenName            string           `json:"GivenName,omitempty"`
	MiddleName           string           `json:"MiddleName,omitempty"`
	FamilyName           string           `json:"FamilyName,omitempty"`
	Suffix               string           `json:"Suffix,omitempty"`
	DisplayName          string           `json:"DisplayName,omitempty"`
	FullyQualifiedName   string           `json:"FullyQualifiedName,omitempty"`
	CompanyName          string           `json:"CompanyName,omitempty"`
	PrintOnCheckName     string           `json:"PrintOnCheckName,omitempty"`
	Active               bool             `json:"Active,omitempty"`
	PrimaryPhone         *TelephoneNumber `json:"PrimaryPhone,omitempty"`
	AlternatePhone       *TelephoneNumber `json:"AlternatePhone,omitempty"`
	Mobile               *TelephoneNumber `json:"Mobile,omitempty"`
	Fax                  *TelephoneNumber `json:"Fax,omitempty"`
	CustomerTypeRef      *ReferenceType   `json:"CustomerTypeRef,omitempty"`
	PrimaryEmailAddr     *EmailAddress    `json:"PrimaryEmailAddr,omitempty"`
	WebAddr              *WebSiteAddress  `json:"WebAddr,omitempty"`
	Taxable              *bool            `json:"Taxable,omitempty"`
	TaxExemptionReasonId *string          `json:"TaxExemptionReasonId,omitempty"`
	BillAddr             *PhysicalAddress `json:"BillAddr,omitempty"`
	ShipAddr             *PhysicalAddress `json:"ShipAddr,omitempty"`
	Notes                string           `json:"Notes,omitempty"`
	BillWithParent       bool             `json:"BillWithParent,omitempty"`
	ParentRef            *ReferenceType   `json:"ParentRef,omitempty"`
	Level                int              `json:"Level,omitempty"`
	Balance              json.Number      `json:"Balance,omitempty"`
	OpenBalanceDate      Date             `json:"OpenBalanceDate,omitempty"`
	BalanceWithJobs      json.Number      `json:"BalanceWithJobs,omitempty"`
}

// GetAddress prioritizes the ship address, but falls back on bill address
func (c *Customer) GetAddress() PhysicalAddress {
	if c.ShipAddr != nil {
		return *c.ShipAddr
	}
	if c.BillAddr != nil {
		return *c.BillAddr
	}
	return PhysicalAddress{}
}

// GetWebsite de-nests the Website object
func (c *Customer) GetWebsite() string {
	if c.WebAddr != nil {
		return c.WebAddr.URI
	}
	return ""
}

// GetPrimaryEmail de-nests the PrimaryEmailAddr object
func (c *Customer) GetPrimaryEmail() string {
	if c.PrimaryEmailAddr != nil {
		return c.PrimaryEmailAddr.Address
	}
	return ""
}

// CreateCustomer creates the given Customer on the QuickBooks server,
// returning the resulting Customer object.
func (c *Client) CreateCustomer(customer *Customer) (*Customer, error) {
	var resp struct {
		Customer Customer
		Time     Date
	}

	if err := c.post("customer", customer, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Customer, nil
}

// FindCustomers gets the full list of Customers in the QuickBooks account.
func (c *Client) FindCustomers() ([]Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customers     []Customer `json:"Customer"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Customer", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no customers could be found")
	}

	customers := make([]Customer, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Customer ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Customers == nil {
			return nil, errors.New("no customers could be found")
		}

		customers = append(customers, resp.QueryResponse.Customers...)
	}

	return customers, nil
}

// FindCustomerById returns a customer with a given Id.
func (c *Client) FindCustomerById(id string) (*Customer, error) {
	var r struct {
		Customer Customer
		Time     Date
	}

	if err := c.get("customer/"+id, &r, nil); err != nil {
		return nil, err
	}

	return &r.Customer, nil
}

// FindCustomerByName gets a customer with a given name.
func (c *Client) FindCustomerByName(name string) (*Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customer   []Customer
			TotalCount int
		}
	}

	query := "SELECT * FROM Customer WHERE DisplayName = '" + strings.Replace(name, "'", "''", -1) + "'"

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if len(resp.QueryResponse.Customer) == 0 {
		return nil, errors.New("no customers could be found")
	}

	return &resp.QueryResponse.Customer[0], nil
}

// QueryCustomers accepts an SQL query and returns all customers found using it
func (c *Client) QueryCustomers(query string) ([]Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customers     []Customer `json:"Customer"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Customers == nil {
		return make([]Customer, 0), nil
	}

	return resp.QueryResponse.Customers, nil
}

// UpdateCustomer updates the given Customer on the QuickBooks server,
// returning the resulting Customer object. It's a sparse update, as not all QB
// fields are present in our Customer object.
func (c *Client) UpdateCustomer(customer *Customer) (*Customer, error) {
	if customer.Id == "" {
		return nil, errors.New("missing customer id")
	}

	existingCustomer, err := c.FindCustomerById(customer.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing customer: %v", err)
	}

	customer.SyncToken = existingCustomer.SyncToken

	payload := struct {
		*Customer
		Sparse bool `json:"sparse"`
	}{
		Customer: customer,
		Sparse:   true,
	}

	var customerData struct {
		Customer Customer
		Time     Date
	}

	if err = c.post("customer", payload, &customerData, nil); err != nil {
		return nil, err
	}

	return &customerData.Customer, nil
}

func (c *Client) FindOrCreateCustomer(ctx context.Context, userAccount *bean.UserAccount) (*Customer, error) {
	customers, err := c.QueryCustomers("SELECT * FROM Customer WHERE DisplayName = '" + strings.Replace(userAccount.Email, "'", "''", -1) + "'")
	if err != nil {
		g.Log().Errorf(ctx, "FindOrCreateCustomer QueryCustomers err:%s", err.Error())
		return nil, err
	}
	var customer *Customer
	if customers == nil || len(customers) == 0 {
		// create customer
		customer, err = c.CreateCustomer(&Customer{
			MetaData: &MetaData{CreateTime: Date{
				Time: gtime.NewFromTimeStamp(userAccount.CreateTime).Time,
			}},
			GivenName:        userAccount.FirstName,
			FamilyName:       userAccount.LastName,
			DisplayName:      userAccount.Email,
			CompanyName:      userAccount.CompanyName,
			Active:           userAccount.Status == 0,
			PrimaryEmailAddr: &EmailAddress{Address: userAccount.Email},
			Notes:            fmt.Sprintf("UniBeeUserId:%d", userAccount.Id),
		})
		if err != nil {
			g.Log().Errorf(ctx, "FindOrCreateCustomer CreateCustomer err:%s", err.Error())
			return nil, err
		}
	} else {
		// update customer
		customer, err = c.UpdateCustomer(&Customer{
			Id: customers[0].Id,
			MetaData: &MetaData{CreateTime: Date{
				Time: gtime.NewFromTimeStamp(userAccount.CreateTime).Time,
			}},
			GivenName:        userAccount.FirstName,
			FamilyName:       userAccount.LastName,
			DisplayName:      userAccount.Email,
			CompanyName:      userAccount.CompanyName,
			Active:           userAccount.Status == 0,
			PrimaryEmailAddr: &EmailAddress{Address: userAccount.Email},
			Notes:            fmt.Sprintf("UniBeeUserId:%d", userAccount.Id),
		})
		if err != nil {
			g.Log().Errorf(ctx, "FindOrCreateCustomer UpdateCustomer err:%s", err.Error())
			return nil, err
		}
	}
	return customer, nil
}
