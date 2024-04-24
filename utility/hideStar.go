package utility

import (
	"regexp"
	"strings"
)

func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 && lens <= 10 {
				var replaceA = make([]byte, lens-4+3)
				for i, _ := range replaceA {
					replaceA[i] = '*'
				}
				result = string(nameRune[:2]) + string(replaceA) + string(nameRune[lens-2:lens])
			} else if lens > 10 {
				var replaceA = make([]byte, lens-4+3)
				for i, _ := range replaceA {
					replaceA[i] = '*'
				}
				result = string(nameRune[:4]) + string(replaceA) + string(nameRune[lens-4:lens])
			}
		}
		return
	}
}
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}
