package utility

import (
	"fmt"
	"strings"
)

// ConvertFenToYuanMinUnitStr 将分转换为元的字符串
func ConvertFenToYuanMinUnitStr(cents int64) string {
	dollars := float64(cents) / 100.0
	return strings.Replace(fmt.Sprintf("%.2f", dollars), ".00", "", -1)
}
