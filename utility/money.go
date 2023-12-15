package utility

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertFenToYuanMinUnitStr 将分转换为元的字符串
func ConvertFenToYuanMinUnitStr(cents int64) string {
	dollars := float64(cents) / 100.0
	return strings.Replace(fmt.Sprintf("%.2f", dollars), ".00", "", -1)
}

func ConvertYuanStrToFen(target string) int64 {
	// 将字符串解析为浮点数
	yuan, err := strconv.ParseFloat(target, 64)
	if err != nil {
		panic(fmt.Sprintf("ConvertYuanStrToFen panic:%s", target))
	}
	// 将浮点数表示的元转换为分的整数
	fen := int64(yuan * 100)
	return fen
}
