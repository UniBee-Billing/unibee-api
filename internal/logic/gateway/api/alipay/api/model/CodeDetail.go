package model

type CodeDetail struct {
	CodeValueType CodeValueType `json:"codeValueType,omitempty"`
	CodeValue     string        `json:"codeValue,omitempty"`
	DisplayType   DisplayType   `json:"displayType,omitempty"`
}

type CodeValueType string

const (
	CodeValueType_BARCODE CodeValueType = "BARCODE"
	CodeValueType_QRCODE  CodeValueType = "QRCODE"
)

type DisplayType string

const (
	DisplayType_TEXT        DisplayType = "TEXT"
	DisplayType_MIDDLEIMAGE DisplayType = "MIDDLEIMAGE"
	DisplayType_BIGIMAGE    DisplayType = "BIGIMAGE"
	DisplayType_SMALLIMAGE  DisplayType = "SMALLIMAGE"
	DisplayType_IMAGE       DisplayType = "IMAGE"
)
