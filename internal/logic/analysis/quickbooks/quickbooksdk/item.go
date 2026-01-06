// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooksdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
)

// Item represents a QuickBooks Item object (a product type).
type Item struct {
	Id                  string         `json:"Id,omitempty"`
	SyncToken           string         `json:"SyncToken,omitempty"`
	MetaData            *MetaData      `json:"MetaData,omitempty"`
	Name                string         `json:"Name"`
	SKU                 string         `json:"Sku,omitempty"`
	Description         string         `json:"Description,omitempty"`
	Active              bool           `json:"Active,omitempty"`
	Taxable             bool           `json:"Taxable,omitempty"`
	SalesTaxIncluded    bool           `json:"SalesTaxIncluded,omitempty"`
	UnitPrice           json.Number    `json:"UnitPrice,omitempty"`
	Type                string         `json:"Type"`
	IncomeAccountRef    *ReferenceType `json:"IncomeAccountRef,omitempty"`
	ExpenseAccountRef   *ReferenceType `json:"ExpenseAccountRef,omitempty"`
	PurchaseDesc        string         `json:"PurchaseDesc,omitempty"`
	PurchaseTaxIncluded bool           `json:"PurchaseTaxIncluded,omitempty"`
	PurchaseCost        json.Number    `json:"PurchaseCost,omitempty"`
	AssetAccountRef     *ReferenceType `json:"AssetAccountRef,omitempty"`
	TrackQtyOnHand      bool           `json:"TrackQtyOnHand,omitempty"`
	QtyOnHand           json.Number    `json:"QtyOnHand,omitempty"`
	SalesTaxCodeRef     *ReferenceType `json:"SalesTaxCodeRef,omitempty"`
	PurchaseTaxCodeRef  *ReferenceType `json:"PurchaseTaxCodeRef,omitempty"`
}

func (c *Client) CreateItem(item *Item) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.post("item", item, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// FindItems gets the full list of Items in the QuickBooks account.
func (c *Client) FindItems() ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Item", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no items could be found")
	}

	items := make([]Item, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Item ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Items == nil {
			return nil, errors.New("no items could be found")
		}

		items = append(items, resp.QueryResponse.Items...)
	}

	return items, nil
}

// FindItemById returns an item with a given Id.
func (c *Client) FindItemById(id string) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.get("item/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// QueryItems accepts an SQL query and returns all items found using it
func (c *Client) QueryItems(query string) ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Items == nil {
		return make([]Item, 0), nil
	}

	return resp.QueryResponse.Items, nil
}

// UpdateItem updates the item
func (c *Client) UpdateItem(item *Item) (*Item, error) {
	if item.Id == "" {
		return nil, errors.New("missing item id")
	}

	existingItem, err := c.FindItemById(item.Id)
	if err != nil {
		return nil, err
	}

	item.SyncToken = existingItem.SyncToken

	payload := struct {
		*Item
		Sparse bool `json:"sparse"`
	}{
		Item:   item,
		Sparse: true,
	}

	var itemData struct {
		Item Item
		Time Date
	}

	if err = c.post("item", payload, &itemData, nil); err != nil {
		return nil, err
	}

	return &itemData.Item, err
}

func (c *Client) FindOrCreateItem(ctx context.Context, name string, account *Account) (*Item, error) {
	var item *Item
	items, err := c.QueryItems(fmt.Sprintf("SELECT * FROM Item WHERE Name = '%s'", name))
	if err != nil {
		g.Log().Errorf(ctx, "Error querying QuickBooks item '%s': %s.", name, err.Error())
	} else if len(items) > 0 {
		item = &items[0]
		g.Log().Infof(ctx, "Found existing QuickBooks item ID: %s (%s) for '%s'", item.Id, item.Name, name)
	}

	if item == nil {
		newItem := &Item{
			Name: name,
			Type: "Service", // Fix to Service Type
			IncomeAccountRef: &ReferenceType{
				Value: account.Id,
				Name:  account.Name,
			},
		}
		item, err = c.CreateItem(newItem)
		if err != nil {
			g.Log().Errorf(ctx, "Failed to create item '%s' in QuickBooks: %s", name, err.Error())
			return nil, err
		}
		g.Log().Infof(ctx, "Created new QuickBooks item ID: %s (%s) for invoice item '%s'", item.Id, item.Name, name)
	}
	return item, nil
}
