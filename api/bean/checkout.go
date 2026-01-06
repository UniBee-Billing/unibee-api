package bean

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
)

var DefaultCheckoutName = "DefaultTemplate"
var DefaultCheckoutDescription = "Merchant Default Checkout Template"

type CheckoutSignIn struct {
	Redirect              bool   `json:"redirect" description:"should redirect to sign in page"`
	Url                   string `json:"url" description:"redirect url"`
	DuplicateSubscription bool   `json:"DuplicateSubscription" description:"whether contain active or incomplete subscription or not"`
}

type MerchantCheckout struct {
	Id          int64                  `json:"id"          description:"ID"`          // ID
	MerchantId  uint64                 `json:"merchantId"  description:"merchantId"`  // merchantId
	Name        string                 `json:"name"        description:"name"`        // name
	Description string                 `json:"description" description:"description"` // description
	Data        map[string]interface{} `json:"data"        description:"checkout config data"`
	StagingData map[string]interface{} `json:"stagingData" description:"checkout staging config data"`
	CreateTime  int64                  `json:"createTime"  description:"create utc time"` // create utc time
	UpdateTime  int64                  `json:"updateTime"  description:"update utc time"` // update utc time
	IsDefault   bool                   `json:"isDefault"      description:"is default"`
	//CheckoutUrl string                 `json:"checkoutUrl" description:"checkout url"`
}

func SimplifyMerchantCheckout(ctx context.Context, one *entity.MerchantCheckout) *MerchantCheckout {
	if one == nil {
		return nil
	}
	var checkoutData = make(map[string]interface{})
	if len(one.Data) > 0 {
		err := gjson.Unmarshal([]byte(one.Data), &checkoutData)
		if err != nil {
			fmt.Printf("SimplifyPlan Unmarshal checkoutData error:%s", err.Error())
		}
	}

	var checkoutStagingData = make(map[string]interface{})
	if len(one.Data) > 0 {
		err := gjson.Unmarshal([]byte(one.Staging), &checkoutStagingData)
		if err != nil {
			fmt.Printf("SimplifyPlan Unmarshal checkoutstagingData error:%s", err.Error())
		}
	}
	updateTime := one.CreateTime
	if one.GmtModify != nil {
		updateTime = one.GmtModify.Timestamp()
	}
	isDefault := false
	if one.Name == DefaultCheckoutName && one.Description == DefaultCheckoutDescription {
		isDefault = true
	}
	return &MerchantCheckout{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		Name:        one.Name,
		Description: one.Description,
		Data:        checkoutData,
		StagingData: checkoutStagingData,
		CreateTime:  one.CreateTime,
		UpdateTime:  updateTime,
		IsDefault:   isDefault,
	}
}

func SimplifyMerchantCheckoutList(ctx context.Context, list []*entity.MerchantCheckout) []*MerchantCheckout {
	if list == nil || len(list) == 0 {
		return make([]*MerchantCheckout, 0)
	}
	result := make([]*MerchantCheckout, 0)
	for _, one := range list {
		result = append(result, SimplifyMerchantCheckout(ctx, one))
	}
	return result
}
