package product

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type NewInternalReq struct {
	MerchantId  uint64                  `json:"merchantId" dc:"MerchantId" `
	ProductName string                  `json:"productName" description:"ProductName"`                                // ProductName
	Description string                  `json:"description" description:"description"`                                // description
	ImageUrl    string                  `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     string                  `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      int                     `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	Metadata    *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func ProductNew(ctx context.Context, req *NewInternalReq) (one *entity.Product, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(len(req.ProductName) > 0, "product name should not blank")

	merchant := query.GetMerchantById(ctx, req.MerchantId)
	if len(req.ImageUrl) == 0 {
		req.ImageUrl = merchant.CompanyLogo
	}
	if len(req.HomeUrl) == 0 {
		req.HomeUrl = merchant.HomeUrl
	}
	utility.Assert(merchant != nil, "merchant not found")
	if req.Status <= 0 {
		req.Status = 1
	}
	utility.Assert(req.Status == 1 || req.Status == 2, "status should be 1|2")
	if len(req.ProductName) == 0 {
		req.ProductName = req.ProductName
	}

	one = &entity.Product{
		CompanyId:   merchant.CompanyId,
		MerchantId:  req.MerchantId,
		ProductName: req.ProductName,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		HomeUrl:     req.HomeUrl,
		Status:      req.Status,
		CreateTime:  gtime.Now().Timestamp(),
		MetaData:    utility.MarshalToJsonString(req.Metadata),
	}
	result, err := dao.Product.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`ProductNew record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one, nil
}

type EditInternalReq struct {
	MerchantId  uint64                  `json:"merchantId" dc:"MerchantId" `
	ProductId   uint64                  `json:"productId" dc:"Id of product" v:"required"`
	ProductName *string                 `json:"productName" description:"ProductName"`                                // ProductName
	Description *string                 `json:"description" description:"description"`                                // description
	ImageUrl    *string                 `json:"imageUrl"    description:"image_url"`                                  // image_url
	HomeUrl     *string                 `json:"homeUrl"     description:"home_url"`                                   // home_url
	Status      *int                    `json:"status"      description:"status，1-active，2-inactive, default active"` // status，1-active，2-inactive, default active
	Metadata    *map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func ProductEdit(ctx context.Context, req *EditInternalReq) (one *entity.Product, err error) {
	utility.Assert(req != nil, "Req not found")
	utility.Assert(req.ProductId > 0, "ProductId should > 0")
	one = query.GetProductById(ctx, req.ProductId)
	utility.Assert(one != nil, fmt.Sprintf("product not found, id:%d", req.ProductId))
	utility.Assert(one.MerchantId == req.MerchantId, "Merchant not match")

	if req.ProductName != nil {
		utility.Assert(len(*req.ProductName) > 0, "Product name should not blank")
	}
	if req.Status != nil {
		utility.Assert(*req.Status == 1 || *req.Status == 2, "status should be 1|2")
	}

	_, err = dao.Product.Ctx(ctx).Data(g.Map{
		dao.Product.Columns().ProductName: req.ProductName,
		dao.Product.Columns().Description: req.Description,
		dao.Product.Columns().ImageUrl:    req.ImageUrl,
		dao.Product.Columns().HomeUrl:     req.HomeUrl,
		dao.Product.Columns().Status:      req.Status,
		dao.Product.Columns().IsDeleted:   0,
	}).Where(dao.Product.Columns().Id, req.ProductId).OmitNil().Update()
	if err != nil {
		return nil, gerror.Newf(`ProductEdit record insert failure %s`, err)
	}
	if req.Metadata != nil {
		_, _ = dao.Product.Ctx(ctx).Data(g.Map{
			dao.Product.Columns().MetaData: utility.MarshalToJsonString(req.Metadata),
		}).Where(dao.Product.Columns().Id, req.ProductId).OmitNil().Update()
	}

	one = query.GetProductById(ctx, req.ProductId)

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one, nil
}

func ProductActivate(ctx context.Context, productId uint64) error {
	utility.Assert(productId > 0, "invalid productId")
	one := query.GetProductById(ctx, productId)
	utility.Assert(one != nil, "product not found, invalid productId")
	if one.Status == 1 {
		return nil
	}
	_, err := dao.Product.Ctx(ctx).Data(g.Map{
		dao.Product.Columns().Status:    1,
		dao.Product.Columns().IsDeleted: 0,
		dao.Product.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Product.Columns().Id, productId).OmitNil().Update()
	if err != nil {
		return err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", one.Id),
		Content:        "Activate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func ProductInactivate(ctx context.Context, productId uint64) error {
	utility.Assert(productId > 0, "invalid productId")
	one := query.GetProductById(ctx, productId)
	utility.Assert(one != nil, "product not found, invalid productId")
	if one.Status == 2 {
		return nil
	}
	_, err := dao.Product.Ctx(ctx).Data(g.Map{
		dao.Product.Columns().Status:    2,
		dao.Product.Columns().IsDeleted: 0,
		dao.Product.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Product.Columns().Id, productId).OmitNil().Update()
	if err != nil {
		return err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", one.Id),
		Content:        "Inactivate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func ProductCopy(ctx context.Context, productId uint64) (one *entity.Product, err error) {
	utility.Assert(productId > 0, "ProductId should > 0")
	one = query.GetProductById(ctx, productId)
	utility.Assert(one != nil, fmt.Sprintf("product not found, id:%d", productId))
	one = &entity.Product{
		CompanyId:   one.CompanyId,
		MerchantId:  one.MerchantId,
		ProductName: one.ProductName + "(Copy)",
		Description: one.Description,
		ImageUrl:    one.ImageUrl,
		HomeUrl:     one.HomeUrl,
		Status:      one.Status,
		CreateTime:  gtime.Now().Timestamp(),
		MetaData:    one.MetaData,
	}
	result, err := dao.Product.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`ProductCopy record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", productId),
		Content:        fmt.Sprintf("CopyTo(%v)", one.Id),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one, nil
}

func ProductDelete(ctx context.Context, productId uint64) (err error) {
	utility.Assert(productId > 0, "productId invalid")
	one := query.GetProductById(ctx, productId)
	utility.Assert(one != nil, fmt.Sprintf("product not found, id:%d", productId))
	utility.Assert(one.Status == 2, fmt.Sprintf("product is not inactive status, id:%d", productId))
	list := query.GetPlansByProductId(ctx, int64(one.Id))
	utility.Assert(list == nil || len(list) == 0, "product can not delete while has plan linked")
	_, err = dao.Product.Ctx(ctx).Data(g.Map{
		dao.Product.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.Product.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Product.Columns().Id, one.Id).Update()
	if err != nil {
		return err
	}

	one.IsDeleted = int(gtime.Now().Timestamp())
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Product(%v)", one.Id),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}
