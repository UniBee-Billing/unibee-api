package quickbooksdk

import (
	"errors"
)

type CustomerType struct {
	SyncToken string   `json:",omitempty"`
	Domain    string   `json:"domain,omitempty"`
	Name      string   `json:",omitempty"`
	Active    bool     `json:",omitempty"`
	Id        string   `json:",omitempty"`
	MetaData  MetaData `json:",omitempty"`
}

// FindCustomerTypeById returns a customerType with a given Id.
func (c *Client) FindCustomerTypeById(id string) (*CustomerType, error) {
	var r struct {
		CustomerType CustomerType
		Time         Date
	}

	if err := c.get("customertype/"+id, &r, nil); err != nil {
		return nil, err
	}

	return &r.CustomerType, nil
}

// QueryCustomerTypes accepts an SQL query and returns all customerTypes found using it
func (c *Client) QueryCustomerTypes(query string) ([]CustomerType, error) {
	var resp struct {
		QueryResponse struct {
			CustomerTypes []CustomerType `json:"CustomerType"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.CustomerTypes == nil {
		return nil, errors.New("could not find any customerTypes")
	}

	return resp.QueryResponse.CustomerTypes, nil
}
