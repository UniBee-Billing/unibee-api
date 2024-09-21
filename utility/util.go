package utility

func CheckReturn(check bool, a interface{}, b interface{}) interface{} {
	if check {
		return a
	} else {
		return b
	}
}
