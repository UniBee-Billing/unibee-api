package utility

func Assert(check bool, message string) {
	if !check {
		panic(message)
	}
}
