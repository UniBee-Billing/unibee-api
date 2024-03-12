package bean

type EmailTemplateVo struct {
	Id                  int64  `json:"id"                 description:""`                //
	MerchantId          uint64 `json:"merchantId"         description:""`                //
	TemplateName        string `json:"templateName"       description:""`                //
	TemplateDescription string `json:"templateDescription" description:""`               //
	TemplateTitle       string `json:"templateTitle"      description:""`                //
	TemplateContent     string `json:"templateContent"    description:""`                //
	TemplateAttachName  string `json:"templateAttachName" description:""`                //
	CreateTime          int64  `json:"createTime"         description:"create utc time"` // create utc time
	UpdateTime          int64  `json:"updateTime"         description:"update utc time"` // create utc time
	Status              string `json:"status"             description:""`                //
}
