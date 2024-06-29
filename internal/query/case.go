package query

import (
	"strings"
)

//
//func Camel2Case(name string) string {
//	buffer := NewBuffer()
//	for i, r := range name {
//		if unicode.IsUpper(r) {
//			if i != 0 {
//				buffer.Append('_')
//			}
//			buffer.Append(unicode.ToLower(r))
//		} else {
//			buffer.Append(r)
//		}
//	}
//	return buffer.String()
//}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}
