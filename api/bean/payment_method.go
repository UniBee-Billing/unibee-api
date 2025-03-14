package bean

import "github.com/gogf/gf/v2/encoding/gjson"

type PaymentMethod struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	IsDefault bool        `json:"isDefault"`
	Data      *gjson.Json `json:"data"`
}
