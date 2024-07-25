package bean

import (
	entity "unibee/internal/model/entity/default"
)

type MerchantBatchExportTemplateSimplify struct {
	TemplateId    uint64                 `json:"templateId"            description:"templateId"`
	Name          string                 `json:"name"          description:"name"`        // name
	MerchantId    uint64                 `json:"merchantId"    description:"merchant_id"` // merchant_id
	MemberId      uint64                 `json:"memberId"      description:"member_id"`   // member_id
	Task          string                 `json:"task" dc:"Task,InvoiceExport|UserExport|SubscriptionExport|TransactionExport|DiscountExport|UserDiscountExport"`
	Payload       map[string]interface{} `json:"payload" dc:"Payload"`
	ExportColumns []string               `json:"exportColumns" dc:"ExportColumns, the export file column list, will export all columns if not specified"`
	Format        string                 `json:"format" dc:"The format of export file, xlsx|csv, will be xlsx if not specified"`
	CreateTime    int64                  `json:"createTime"    description:"create utc time"` // create utc time
}

func SimplifyMerchantBatchExportTemplate(one *entity.MerchantBatchExportTemplate) *MerchantBatchExportTemplateSimplify {
	if one == nil {
		return nil
	}
	var payload map[string]interface{}
	_ = UnmarshalFromJsonString(one.Payload, &payload)

	var exportColumns []string
	_ = UnmarshalFromJsonString(one.ExportColumns, &exportColumns)
	return &MerchantBatchExportTemplateSimplify{
		TemplateId:    one.Id,
		Name:          one.Name,
		MerchantId:    one.MerchantId,
		MemberId:      one.MemberId,
		Task:          one.Task,
		Payload:       payload,
		ExportColumns: exportColumns,
		Format:        one.Format,
		CreateTime:    one.CreateTime,
	}
}
