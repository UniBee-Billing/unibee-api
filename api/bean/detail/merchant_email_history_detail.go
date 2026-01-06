package detail

import (
	"context"
	entity "unibee/internal/model/entity/default"
)

type MerchantEmailHistoryStatistics struct {
	TotalSend    int64 `json:"totalSend" dc:"Total Send"`
	TotalSuccess int64 `json:"totalSuccess" dc:"Total Success"`
	TotalFail    int64 `json:"totalFail" dc:"Total Fail"`
}

type MerchantEmailHistoryDetail struct {
	Id         uint64 `json:"id"         description:"Id"`                            // Id
	MerchantId uint64 `json:"merchantId" description:"merchantId"`                    // merchantId
	Email      string `json:"email"      description:"Email address"`                 // Email address
	Title      string `json:"title"      description:"Email title"`                   // Email title
	Content    string `json:"content"    description:"Email content"`                 // Email content
	AttachFile string `json:"attachFile" description:"Attachment file"`               // Attachment file
	Response   string `json:"response"   description:"Email response"`                // Email response
	CreateTime int64  `json:"createTime" description:"create utc time"`               // create utc time
	Status     int    `json:"status"     description:"0-pending,1-success,2-failure"` // 0-pending,1-success,2-failure
}

func ConvertMerchantEmailHistoryDetail(ctx context.Context, one *entity.MerchantEmailHistory) *MerchantEmailHistoryDetail {
	if one == nil {
		return nil
	}

	return &MerchantEmailHistoryDetail{
		Id:         one.Id,
		MerchantId: one.MerchantId,
		Email:      one.Email,
		Title:      one.Title,
		Content:    one.Content,
		AttachFile: one.AttachFile,
		Response:   one.Response,
		CreateTime: one.CreateTime,
		Status:     one.Status,
	}
}
