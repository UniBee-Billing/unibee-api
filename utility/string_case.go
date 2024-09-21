package utility

import (
	"strings"
	"unicode"
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

func IsStartUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

func IsStartLower(s string) bool {
	return unicode.IsLower([]rune(s)[0])
}

func ToFirstCharLowerCase(s string) string {
	var result string
	if !IsStartUpper(s) {
		return s
	}
	for i, char := range s {
		if i == 0 {
			result += strings.ToLower(string(char))
		} else {
			result += string(char)
		}
	}
	return result
}

func ToFirstCharUpperCase(s string) string {
	var result string
	if !IsStartLower(s) {
		return s
	}
	for i, char := range s {
		if i == 0 {
			result += strings.ToUpper(string(char))
		} else {
			result += string(char)
		}
	}
	return result
}
