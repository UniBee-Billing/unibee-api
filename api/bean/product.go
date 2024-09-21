package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
)

type Product struct {
	Id          uint64 `json:"id"          description:""`
	MerchantId  uint64 `json:"merchantId"  description:"merchant id"`                                // merchant id
	ProductName string `json:"productName" description:"PlanName"`                                   // PlanName
	Description string `json:"description" description:"description"`                                // description
	ImageUrl    string `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     string `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      int    `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	IsDeleted   int    `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"`                      // 0-UnDeleted，1-Deleted
	CreateTime  int64  `json:"createTime"  description:"create utc time"`                            // create utc time
	MetaData    string `json:"metaData"    description:"meta_data(json)"`                            // meta_data(json)
}

func SimplifyProduct(one *entity.Product) *Product {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyPlan Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &Product{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		ProductName: one.ProductName,
		Description: one.Description,
		ImageUrl:    one.ImageUrl,
		HomeUrl:     one.HomeUrl,
		Status:      one.Status,
		IsDeleted:   one.IsDeleted,
		CreateTime:  one.CreateTime,
		MetaData:    one.MetaData,
	}
}

func SimplifyProductList(ones []*entity.Product) (list []*Product) {
	if len(ones) == 0 {
		return make([]*Product, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyProduct(one))
	}
	return list
}
